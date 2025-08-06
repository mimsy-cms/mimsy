package migrations

import (
	"context"
	"errors"
	"fmt"
	"os"

	migrations_interface "github.com/mimsy-cms/mimsy/internal/interfaces/migrations"
	"github.com/xataio/pgroll/pkg/backfill"
	"github.com/xataio/pgroll/pkg/migrations"
	"github.com/xataio/pgroll/pkg/roll"
	"github.com/xataio/pgroll/pkg/state"
)

type OptionFn func(*runConfig)

func NewRunConfig(opts ...OptionFn) *runConfig {
	config := &runConfig{
		Schema:      "public",
		StateSchema: "pgroll",
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}

// WithSchema sets the schema where migrations will be applied.
func WithMigrationsDir(dir string) OptionFn {
	return func(c *runConfig) {
		c.MigrationsDir = dir
	}
}

// WithSchema sets the schema where migrations will be applied.
func WithPgURL(url string) OptionFn {
	return func(c *runConfig) {
		c.PgURL = url
	}
}

// RunConfig holds the configuration for running migrations.
type runConfig struct {
	// MigrationsDir is the directory containing migration files.
	MigrationsDir string
	// PgURL is the PostgreSQL connection URL.
	PgURL string
	// Schema is the name of the schema where migrations will be applied.
	Schema string
	// StateSchema is the name of the schema where migration state will be stored.
	StateSchema string

	NewState    func(ctx context.Context, pgURL, schema string) (migrations_interface.State, error)
	NewMigrator func(ctx context.Context, pgURL, schema string, s migrations_interface.State) (migrations_interface.Migrator, error)
}

// Run executes the migrations defined in the migrations directory.
// It initializes the migration state, checks if a migration is already in progress,
// and applies all unapplied migrations in the specified directory.
// It returns the number of migrations applied or an error if something goes wrong.
func Run(ctx context.Context, config *runConfig) (int, error) {
	var (
		st  migrations_interface.State
		err error
	)

	if config.NewState != nil {
		st, err = config.NewState(ctx, config.PgURL, config.StateSchema)

	} else {
		st, err = state.New(ctx, config.PgURL, config.StateSchema)
		if err != nil {
			return 0, err
		}
	}
	if err != nil {
		return 0, err
	}

	var m migrations_interface.Migrator
	if config.NewMigrator != nil {
		m, err = config.NewMigrator(ctx, config.PgURL, config.Schema, st)
	} else {
		rollMigrator, rollErr := roll.New(ctx, config.PgURL, config.Schema, st.(*state.State))
		if rollErr != nil {
			err = rollErr
		} else {
			m = &migratorAdapter{rollMigrator}
		}
	}

	if err != nil {
		return 0, err
	}
	defer m.Close()

	ok, err := st.IsInitialized(ctx)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("migration state is not initialized")
	}

	latestMigration, err := m.State().LatestMigration(ctx, m.Schema())
	if err != nil {
		return 0, err
	}

	active, err := m.State().IsActiveMigrationPeriod(ctx, m.Schema())
	if err != nil {
		return 0, err
	}
	if active {
		name := "unknown"
		if latestMigration != nil {
			name = *latestMigration
		}
		return 0, fmt.Errorf("migration %q is active", name)
	}

	rawMigs, err := m.UnappliedMigrations(ctx, os.DirFS(config.MigrationsDir))
	if err != nil {
		return 0, err
	}

	migs := make([]*migrations.Migration, 0, len(rawMigs))
	for _, rawMig := range rawMigs {
		mig, err := migrations.ParseMigration(rawMig)
		if err != nil {
			return 0, err
		}
		migs = append(migs, mig)
	}

	backfillConfig := backfill.NewConfig()

	for _, mig := range migs {
		if err := m.Start(ctx, mig, backfillConfig); err != nil {
			return 0, err
		}

		if err := m.Complete(ctx); err != nil {
			return 0, err
		}
	}
	return len(migs), nil
}

// migratorAdapter adapts *roll.Roll to migrations_interface.Migrator
type migratorAdapter struct {
	*roll.Roll
}

// State adapts *state.State to migrations_interface.State
func (m *migratorAdapter) State() migrations_interface.State {
	return m.Roll.State()
}

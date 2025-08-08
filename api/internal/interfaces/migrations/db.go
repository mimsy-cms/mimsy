package migrations_interface

import (
	"context"
	"io/fs"

	"github.com/xataio/pgroll/pkg/backfill"
	"github.com/xataio/pgroll/pkg/migrations"
)

// State represents the pgroll state manager.
type State interface {
	IsInitialized(ctx context.Context) (bool, error)
	LatestMigration(ctx context.Context, schema string) (*string, error)
	IsActiveMigrationPeriod(ctx context.Context, schema string) (bool, error)
}

// Migrator represents a migrator capable of running migrations.
type Migrator interface {
	State() State
	Schema() string
	UnappliedMigrations(ctx context.Context, f fs.FS) ([]*migrations.RawMigration, error)
	Start(ctx context.Context, m *migrations.Migration, cfg *backfill.Config) error
	Complete(ctx context.Context) error
	Close() error
}

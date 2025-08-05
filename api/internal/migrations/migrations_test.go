package migrations

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	migrations_interface "github.com/mimsy-cms/mimsy/internal/interfaces/migrations"
	mocks "github.com/mimsy-cms/mimsy/internal/mocks/migrations"
	pgmigs "github.com/xataio/pgroll/pkg/migrations"
)

// TestRun_Success tests the Run function with a mock database.
func TestRun_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockState := mocks.NewMockState(ctrl)
	mockMigrator := mocks.NewMockMigrator(ctrl)

	// Setup expectations
	mockState.EXPECT().IsInitialized(ctx).Return(true, nil)
	mockState.EXPECT().LatestMigration(ctx, "public").Return(nil, nil)
	mockState.EXPECT().IsActiveMigrationPeriod(ctx, "public").Return(false, nil)

	mockMigrator.EXPECT().Close()

	mockMigrator.EXPECT().State().Return(mockState).AnyTimes()
	mockMigrator.EXPECT().Schema().Return("public").AnyTimes()

	mockMigrator.EXPECT().UnappliedMigrations(ctx, gomock.Any()).Return([]*pgmigs.RawMigration{
		{
			Name:       "001_init.up.sql",
			Operations: []byte(`[]`), // minimal valid empty operations array
		},
	}, nil)

	mockMigrator.EXPECT().Start(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockMigrator.EXPECT().Complete(gomock.Any()).Return(nil)

	config := runConfig{
		MigrationsDir: "testdata/migrations",
		PgURL:         "postgres://user:pass@localhost/db",
		Schema:        "public",
		StateSchema:   "pgroll",
		NewState: func(ctx context.Context, pgURL, schema string) (migrations_interface.State, error) {
			return mockState, nil
		},
		NewMigrator: func(ctx context.Context, pgURL, schema string, s migrations_interface.State) (migrations_interface.Migrator, error) {
			return mockMigrator, nil
		},
	}

	// Call the function
	n, err := Run(ctx, &config)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 migration to be applied, got %d", n)
	}
}

// TestRun_Failure_NotInitialized tests the Run function with a mock database that has not been initialized.

// TestRun_Failure_ActiveMigrationInProgress tests the Run function with a mock database that has an active migration in progress.

// TestRun_Failure_StateCreationError tests the Run function with a mock database that returns an error when checking the state.

// TestRun_Failure_RollCreationError tests the Run function with a mock database that returns an error when creating the roll.

// TestRun_Failure_LatestMigrationError tests the Run function with a mock database that returns an error when getting the latest migration.

// TestRun_Failure_UnappliedMigrationsError tests the Run function with a mock database that returns an error when getting unapplied migrations.

// TestRun_Failure_ParseMigrationError tests the Run function with a mock database that returns an error when parsing a migration.

// TestRun_Failure_StartMigrationError tests the Run function with a mock database that returns an error when starting a migration.

// TestRun_Failure_CompleteMigrationError tests the Run function with a mock database that returns an error when completing a migration.

package migrations

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	migrations_interface "github.com/mimsy-cms/mimsy/internal/interfaces/migrations"
	mocks "github.com/mimsy-cms/mimsy/internal/mocks/migrations"
	"github.com/xataio/pgroll/pkg/migrations"
	pgroll_migrations "github.com/xataio/pgroll/pkg/migrations"
)

// =================================================================================================
// Helper Functions
// =================================================================================================
type testDeps struct {
	ctrl         *gomock.Controller
	ctx          context.Context
	mockState    *mocks.MockState
	mockMigrator *mocks.MockMigrator
	config       *runConfig
}

func setupTest(t *testing.T) *testDeps {
	t.Helper()

	ctrl := gomock.NewController(t)
	ctx := context.Background()

	mockState := mocks.NewMockState(ctrl)
	mockMigrator := mocks.NewMockMigrator(ctrl)

	config := &runConfig{
		UnappliedMigrations: []*migrations.Migration{},
		PgURL:               "postgres://user:pass@localhost/db",
		Schema:              "public",
		StateSchema:         "pgroll",
		NewState: func(ctx context.Context, pgURL, schema string) (migrations_interface.State, error) {
			return mockState, nil
		},
		NewMigrator: func(ctx context.Context, pgURL, schema string, s migrations_interface.State) (migrations_interface.Migrator, error) {
			return mockMigrator, nil
		},
	}

	t.Cleanup(ctrl.Finish) // Automatically call Finish() at the end of the test

	return &testDeps{
		ctrl:         ctrl,
		ctx:          ctx,
		mockState:    mockState,
		mockMigrator: mockMigrator,
		config:       config,
	}
}

func setupCommonMigratorExpectations(m *mocks.MockMigrator, s *mocks.MockState) {
	m.EXPECT().State().Return(s).AnyTimes()
	m.EXPECT().Schema().Return("public").AnyTimes()
	m.EXPECT().Close()
}

func expectSuccessfulMigration(t *testing.T, deps *testDeps) {
	t.Helper()

	deps.mockState.EXPECT().IsInitialized(deps.ctx).Return(true, nil)
	deps.mockState.EXPECT().LatestMigration(deps.ctx, "public").Return(nil, nil)
	deps.mockState.EXPECT().IsActiveMigrationPeriod(deps.ctx, "public").Return(false, nil)

	// Create an empty unapplied migration
	deps.config.UnappliedMigrations = []*pgroll_migrations.Migration{
		{Name: "001_init.up.sql", Operations: []pgroll_migrations.Operation{}},
	}

	deps.mockMigrator.EXPECT().Start(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	deps.mockMigrator.EXPECT().Complete(gomock.Any()).Return(nil)
}

func expectNotInitialized(t *testing.T, deps *testDeps) {
	t.Helper()

	deps.mockState.EXPECT().IsInitialized(deps.ctx).Return(false, nil)
}

func expectActiveMigration(t *testing.T, deps *testDeps) {
	t.Helper()

	deps.mockState.EXPECT().IsInitialized(deps.ctx).Return(true, nil)
	deps.mockState.EXPECT().LatestMigration(deps.ctx, "public").Return(nil, nil)
	deps.mockState.EXPECT().IsActiveMigrationPeriod(deps.ctx, "public").Return(true, nil)
}

// =================================================================================================
// Run
// =================================================================================================

// TestRun_Success tests the Run function with a mock database.
func TestRun_Success(t *testing.T) {
	deps := setupTest(t)
	setupCommonMigratorExpectations(deps.mockMigrator, deps.mockState)
	expectSuccessfulMigration(t, deps)

	n, err := Run(deps.ctx, deps.config)
	if err != nil || n != 1 {
		t.Fatalf("unexpected result: n=%d, err=%v", n, err)
	}
}

// TestRun_Failure_NotInitialized tests the Run function with a mock database that has not been initialized.
func TestRun_Failure_NotInitialized(t *testing.T) {
	deps := setupTest(t)
	setupCommonMigratorExpectations(deps.mockMigrator, deps.mockState)
	expectNotInitialized(t, deps)

	n, err := Run(deps.ctx, deps.config)
	if err == nil {
		t.Fatal("expected an error due to uninitialized database, got nil")
	}
	if n != 0 {
		t.Fatalf("expected 0 migrations applied, got %d", n)
	}
}

// TestRun_Failure_ActiveMigrationInProgress tests the Run function with a mock database that has an active migration in progress.
func TestRun_Failure_ActiveMigrationInProgress(t *testing.T) {
	deps := setupTest(t)
	setupCommonMigratorExpectations(deps.mockMigrator, deps.mockState)
	expectActiveMigration(t, deps)

	n, err := Run(deps.ctx, deps.config)
	if err == nil {
		t.Fatal("expected an error due to active migration in progress, got nil")
	}
	if n != 0 {
		t.Fatalf("expected 0 migrations applied, got %d", n)
	}
}

// TestRun_Failure_StateCreationError tests the Run function with a mock database that returns an error when checking the state.
func TestRun_Failure_StateCreationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	expectedErr := context.DeadlineExceeded

	// Override NewState to return an error
	config := &runConfig{
		UnappliedMigrations: []*migrations.Migration{},
		PgURL:               "postgres://user:pass@localhost/db",
		Schema:              "public",
		StateSchema:         "pgroll",
		NewState: func(ctx context.Context, pgURL, schema string) (migrations_interface.State, error) {
			return nil, expectedErr
		},
		NewMigrator: func(ctx context.Context, pgURL, schema string, s migrations_interface.State) (migrations_interface.Migrator, error) {
			t.Fatal("NewMigrator should not be called if NewState fails")
			return nil, nil
		},
	}

	n, err := Run(ctx, config)
	if err == nil {
		t.Fatal("expected an error due to state creation failure, got nil")
	}
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if n != 0 {
		t.Fatalf("expected 0 migrations applied, got %d", n)
	}
}

// TestRun_Failure_RollCreationError tests the Run function with a mock database that returns an error when creating the roll.
func TestRun_Failure_RollCreationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockState := mocks.NewMockState(ctrl)
	expectedErr := context.Canceled

	config := &runConfig{
		UnappliedMigrations: []*migrations.Migration{},
		PgURL:               "postgres://user:pass@localhost/db",
		Schema:              "public",
		StateSchema:         "pgroll",
		NewState: func(ctx context.Context, pgURL, schema string) (migrations_interface.State, error) {
			return mockState, nil
		},
		NewMigrator: func(ctx context.Context, pgURL, schema string, s migrations_interface.State) (migrations_interface.Migrator, error) {
			return nil, expectedErr
		},
	}

	n, err := Run(ctx, config)
	if err == nil {
		t.Fatal("expected an error due to migrator creation failure, got nil")
	}
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if n != 0 {
		t.Fatalf("expected 0 migrations applied, got %d", n)
	}
}

// TestRun_Failure_LatestMigrationError tests the Run function with a mock database that returns an error when getting the latest migration.
func TestRun_Failure_LatestMigrationError(t *testing.T) {
	deps := setupTest(t)

	setupCommonMigratorExpectations(deps.mockMigrator, deps.mockState)

	deps.mockState.EXPECT().IsInitialized(deps.ctx).Return(true, nil)

	expectedErr := context.DeadlineExceeded
	deps.mockState.EXPECT().LatestMigration(deps.ctx, "public").Return(nil, expectedErr)

	n, err := Run(deps.ctx, deps.config)
	if err == nil {
		t.Fatal("expected an error due to latest migration retrieval failure, got nil")
	}
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if n != 0 {
		t.Fatalf("expected 0 migrations applied, got %d", n)
	}
}

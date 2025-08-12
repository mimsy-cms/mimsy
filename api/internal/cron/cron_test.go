package cron

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Allow overriding the test database URL via environment variable
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		// Use the main database for testing
		dbURL = "postgres://mimsy:mimsy@localhost:5432/mimsy?sslmode=disable"
	}
	
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Skipf("Failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Skipf("Failed to ping test database: %v", err)
	}

	// Clean up and recreate the test table
	_, err = db.Exec(`DROP TABLE IF EXISTS cron_locks`)
	if err != nil {
		t.Fatalf("Failed to clean up test database: %v", err)
	}

	// Create the cron_locks table for testing
	_, err = db.Exec(`
		CREATE TABLE cron_locks (
			key VARCHAR(255) PRIMARY KEY,
			locked_by VARCHAR(255) NOT NULL,
			locked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX idx_cron_locks_expires_at ON cron_locks(expires_at)`)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}

	return db
}

func TestPostgresLocker_SingleLock(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	locker, err := NewPostgresLocker(db)
	if err != nil {
		t.Fatalf("Failed to create locker: %v", err)
	}

	ctx := context.Background()
	lock, err := locker.Lock(ctx, "test-job")
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}

	err = lock.Unlock(ctx)
	if err != nil {
		t.Fatalf("Failed to release lock: %v", err)
	}
}

func TestPostgresLocker_ConcurrentLocks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	locker, err := NewPostgresLocker(db)
	if err != nil {
		t.Fatalf("Failed to create locker: %v", err)
	}

	ctx := context.Background()
	jobKey := "concurrent-test-job"
	
	var successCount int32
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			lock, err := locker.Lock(ctx, jobKey)
			if err == nil {
				atomic.AddInt32(&successCount, 1)
				time.Sleep(10 * time.Millisecond)
				lock.Unlock(ctx)
			}
		}(i)
	}

	wg.Wait()

	if successCount != 1 {
		t.Errorf("Expected only 1 goroutine to acquire lock, but %d succeeded", successCount)
	}
}

func TestPostgresLocker_LockExpiry(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	locker := &postgresLocker{
		db:            db,
		lockTableName: "cron_locks",
		lockTimeout:   100 * time.Millisecond,
	}

	ctx := context.Background()
	_, err := locker.Lock(ctx, "expiry-test")
	if err != nil {
		t.Fatalf("Failed to acquire first lock: %v", err)
	}

	// Check the lock in the database before sleep
	var expiresAt time.Time
	err = db.QueryRow("SELECT expires_at FROM cron_locks WHERE key = $1", "expiry-test").Scan(&expiresAt)
	if err != nil {
		t.Fatalf("Failed to query lock: %v", err)
	}
	t.Logf("Lock expires at: %v, current time: %v", expiresAt, time.Now())

	// Wait for the lock to expire
	time.Sleep(150 * time.Millisecond)

	// Check if lock is expired
	t.Logf("After sleep, current time: %v, lock expired: %v", time.Now(), time.Now().After(expiresAt))

	lock2, err := locker.Lock(ctx, "expiry-test")
	if err != nil {
		t.Fatalf("Failed to acquire lock after expiry: %v", err)
	}

	err = lock2.Unlock(ctx)
	if err != nil {
		t.Fatalf("Failed to release second lock: %v", err)
	}
}

func TestScheduler_RegisterAndRunJob(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	scheduler, err := NewScheduler(db)
	if err != nil {
		t.Fatalf("Failed to create scheduler: %v", err)
	}

	var counter int32
	job := Job{
		Name:     "test-job",
		Schedule: "1s",
		Function: func() {
			atomic.AddInt32(&counter, 1)
		},
		Params: []interface{}{},
	}

	err = scheduler.RegisterJob(job)
	if err != nil {
		t.Fatalf("Failed to register job: %v", err)
	}

	err = scheduler.Start()
	if err != nil {
		t.Fatalf("Failed to start scheduler: %v", err)
	}

	time.Sleep(2500 * time.Millisecond)

	scheduler.Stop()

	count := atomic.LoadInt32(&counter)
	if count < 2 || count > 3 {
		t.Errorf("Expected job to run 2-3 times, but ran %d times", count)
	}
}

func TestScheduler_RemoveJob(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	scheduler, err := NewScheduler(db)
	if err != nil {
		t.Fatalf("Failed to create scheduler: %v", err)
	}

	job := Job{
		Name:     "removable-job",
		Schedule: "5s",
		Function: func() {
			fmt.Println("Job executed")
		},
		Params: []interface{}{},
	}

	err = scheduler.RegisterJob(job)
	if err != nil {
		t.Fatalf("Failed to register job: %v", err)
	}

	jobs := scheduler.ListJobs()
	if len(jobs) != 1 {
		t.Errorf("Expected 1 job, found %d", len(jobs))
	}

	err = scheduler.RemoveJob("removable-job")
	if err != nil {
		t.Fatalf("Failed to remove job: %v", err)
	}

	jobs = scheduler.ListJobs()
	if len(jobs) != 0 {
		t.Errorf("Expected 0 jobs after removal, found %d", len(jobs))
	}
}

func TestScheduler_RunJobNow(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	scheduler, err := NewScheduler(db)
	if err != nil {
		t.Fatalf("Failed to create scheduler: %v", err)
	}

	var executed bool
	var mu sync.Mutex

	job := Job{
		Name:     "manual-job",
		Schedule: "1h",
		Function: func() {
			mu.Lock()
			executed = true
			mu.Unlock()
		},
		Params: []interface{}{},
	}

	err = scheduler.RegisterJob(job)
	if err != nil {
		t.Fatalf("Failed to register job: %v", err)
	}

	err = scheduler.Start()
	if err != nil {
		t.Fatalf("Failed to start scheduler: %v", err)
	}
	defer scheduler.Stop()

	err = scheduler.RunJobNow("manual-job")
	if err != nil {
		t.Fatalf("Failed to run job now: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if !executed {
		t.Error("Job was not executed when RunJobNow was called")
	}
	mu.Unlock()
}

func TestScheduler_CronExpression(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	scheduler, err := NewScheduler(db)
	if err != nil {
		t.Fatalf("Failed to create scheduler: %v", err)
	}

	job := Job{
		Name:     "cron-job",
		Schedule: "*/1 * * * *",
		Function: func() {
			fmt.Println("Cron job executed")
		},
		Params: []interface{}{},
	}

	err = scheduler.RegisterJob(job)
	if err != nil {
		t.Fatalf("Failed to register cron job: %v", err)
	}

	_, exists := scheduler.GetJob("cron-job")
	if !exists {
		t.Error("Cron job was not registered")
	}
}

func TestCronService(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service, err := NewCronService(db)
	if err != nil {
		t.Fatalf("Failed to create cron service: %v", err)
	}

	ctx := context.Background()

	var counter int32
	job := Job{
		Name:     "service-job",
		Schedule: "1s",
		Function: func() {
			atomic.AddInt32(&counter, 1)
		},
		Params: []interface{}{},
	}

	err = service.RegisterJob(ctx, job)
	if err != nil {
		t.Fatalf("Failed to register job via service: %v", err)
	}

	err = service.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start service: %v", err)
	}

	time.Sleep(2500 * time.Millisecond)

	jobs := service.ListJobs(ctx)
	if len(jobs) != 1 {
		t.Errorf("Expected 1 job in service, found %d", len(jobs))
	}

	err = service.Stop(ctx)
	if err != nil {
		t.Fatalf("Failed to stop service: %v", err)
	}

	count := atomic.LoadInt32(&counter)
	if count < 2 || count > 3 {
		t.Errorf("Expected job to run 2-3 times via service, but ran %d times", count)
	}
}
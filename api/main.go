package main

import (
	"context"
	"fmt"

	"github.com/mimsy-cms/mimsy/internal/migrations"
)

func main() {
	runConfig := migrations.NewRunConfig(
		migrations.WithMigrationsDir("./migrations"),
		migrations.WithPgURL("postgres://mimsy:mimsy@localhost?sslmode=disable"),
	)

	// NOTE: Migrations should not be run like this in production.
	migrationCount, err := migrations.Run(context.Background(), runConfig)
	if err != nil {
		fmt.Println("Failed to run migrations:", err)
	}

	fmt.Printf("Successfully applied %d migrations\n", migrationCount)
}

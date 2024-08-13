// cmd/main.go
package main

import (
	"flag"
	"log"

	"github.com/MXkodo/Management-of-School-materials/config"
	repo "github.com/MXkodo/Management-of-School-materials/internal/repo/db"

	"github.com/MXkodo/Management-of-School-materials/app"
)

func main() {
	applyMigrations := flag.Bool("am", false, "Apply database migrations")
	flag.Parse()
	cfg := config.LoadConfig()
	if *applyMigrations {
		if err := repo.ApplyMigrations(cfg); err != nil {
			log.Fatalf("error applying migrations: %v", err)
		}
		log.Println("Migrations applied successfully")
	}
	if err := app.RunApp(*cfg); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

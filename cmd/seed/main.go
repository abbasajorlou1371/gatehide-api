package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/database/seeders"
	"github.com/joho/godotenv"
)

func main() {
	var command = flag.String("command", "admin", "Seeder command to run (admin)")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	switch *command {
	case "admin":
		if err := seedAdmin(cfg); err != nil {
			log.Fatalf("Failed to seed admin: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", *command)
		fmt.Println("Available commands:")
		fmt.Println("  admin - Seed admin user")
		os.Exit(1)
	}
}

// seedAdmin seeds the default admin user
func seedAdmin(cfg *config.Config) error {
	return seeders.SeedAdmin(cfg)
}

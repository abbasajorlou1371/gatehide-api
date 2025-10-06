package seeders

import (
	"fmt"
	"log"

	"github.com/gatehide/gatehide-api/config"
)

// SeederFunc represents a seeder function
type SeederFunc func(cfg *config.Config) error

// Registry manages available seeders
type Registry struct {
	seeders map[string]SeederFunc
}

// NewRegistry creates a new seeder registry
func NewRegistry() *Registry {
	return &Registry{
		seeders: make(map[string]SeederFunc),
	}
}

// Register adds a seeder to the registry
func (r *Registry) Register(name string, seeder SeederFunc) {
	r.seeders[name] = seeder
}

// Get returns a seeder by name
func (r *Registry) Get(name string) (SeederFunc, bool) {
	seeder, exists := r.seeders[name]
	return seeder, exists
}

// List returns all available seeder names
func (r *Registry) List() []string {
	var names []string
	for name := range r.seeders {
		names = append(names, name)
	}
	return names
}

// Run executes a specific seeder
func (r *Registry) Run(name string, cfg *config.Config) error {
	seeder, exists := r.Get(name)
	if !exists {
		return fmt.Errorf("seeder '%s' not found. Available seeders: %v", name, r.List())
	}

	log.Printf("Running seeder: %s", name)
	return seeder(cfg)
}

// RunAll executes all registered seeders
func (r *Registry) RunAll(cfg *config.Config) error {
	log.Println("Running all seeders...")

	for name := range r.seeders {
		if err := r.Run(name, cfg); err != nil {
			return fmt.Errorf("failed to run seeder '%s': %w", name, err)
		}
	}

	log.Println("All seeders completed successfully")
	return nil
}

// Global registry instance
var globalRegistry = NewRegistry()

// RegisterSeeder registers a seeder in the global registry
func RegisterSeeder(name string, seeder SeederFunc) {
	globalRegistry.Register(name, seeder)
}

// GetSeeder returns a seeder from the global registry
func GetSeeder(name string) (SeederFunc, bool) {
	return globalRegistry.Get(name)
}

// ListSeeders returns all available seeder names from the global registry
func ListSeeders() []string {
	return globalRegistry.List()
}

// RunSeeder executes a specific seeder from the global registry
func RunSeeder(name string, cfg *config.Config) error {
	return globalRegistry.Run(name, cfg)
}

// RunAllSeeders executes all registered seeders from the global registry
func RunAllSeeders(cfg *config.Config) error {
	return globalRegistry.RunAll(cfg)
}

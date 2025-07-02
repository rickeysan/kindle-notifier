package seeder

// Seeder defines the interface for all seeders
type Seeder interface {
	// Seed runs the seeder
	Seed() error
	// Clear removes all seeded data
	Clear() error
}

// Registry holds all registered seeders
type Registry struct {
	seeders []Seeder
}

// NewRegistry creates a new seeder registry
func NewRegistry() *Registry {
	return &Registry{
		seeders: make([]Seeder, 0),
	}
}

// Register adds a seeder to the registry
func (r *Registry) Register(seeder Seeder) {
	r.seeders = append(r.seeders, seeder)
}

// SeedAll runs all registered seeders
func (r *Registry) SeedAll() error {
	for _, seeder := range r.seeders {
		if err := seeder.Seed(); err != nil {
			return err
		}
	}
	return nil
}

// ClearAll removes all seeded data
func (r *Registry) ClearAll() error {
	// Clear in reverse order to handle dependencies
	for i := len(r.seeders) - 1; i >= 0; i-- {
		if err := r.seeders[i].Clear(); err != nil {
			return err
		}
	}
	return nil
} 
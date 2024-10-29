package database

import (
	"github.com/jinzhu/gorm"
	"github.com/public-forge/go-gorm-unit-of-work/postgres"
	"go.uber.org/fx"
)

// DBModule is an Fx module that provides the configuration, connection, and holder for PostgreSQL database integration.
// It initializes the PostgreSQL configuration, establishes a connection, and verifies connectivity.
var DBModule = fx.Module("database",

	// Provides the PostgreSQL configuration using the GetPostgresConfig function.
	fx.Provide(GetPostgresConfig),

	// Provides a connection to the PostgreSQL database, initialized by NewConnect.
	fx.Provide(postgres.NewConnect),

	// Provides a holder for the database instance, facilitating dependency injection.
	fx.Provide(postgres.NewDBHolder),

	// Invokes a function to set the global database configuration in postgres.
	fx.Invoke(func(config *postgres.PgConfig) {
		postgres.DbConfig = config
	}),

	// Invokes a function to check the database connection health on startup.
	fx.Invoke(func(db *gorm.DB) {
		postgres.CheckConnection(db)
	}),
)

package database

import (
	"github.com/public-forge/go-gorm-unit-of-work/postgres"
	"github.com/urfave/cli/v2"
)

// Constants for PostgreSQL database configuration parameters.
const (
	postgresUser                  = "postgres-user"           // Database username
	postgresPassword              = "postgres-password"       // Database user password
	postgresDb                    = "postgres-db"             // Database name
	postgresSchema                = "postgres-schema"         // Database schema
	postgresLogMode               = "postgres-log-mode"       // Query logging mode
	postgresHost                  = "postgres-host"           // Database host
	postgresConnectionMaxLifeTime = "postgres-max-life-time"  // Maximum connection lifetime in milliseconds
	postgresMaxOpenConnection     = "postgres-max-connection" // Maximum number of open connections
)

// GetPostgresConfig creates and returns a PgConfig structure containing PostgreSQL
// configuration settings for establishing database connections.
//
// Parameters:
//   - c: The CLI context from which configuration values are read.
//
// Returns:
//
//	A pointer to a PgConfig struct with PostgreSQL settings.
func GetPostgresConfig(c *cli.Context) *postgres.PgConfig {
	return &postgres.PgConfig{
		User:                    c.String(postgresUser),
		Password:                c.String(postgresPassword),
		DBName:                  c.String(postgresDb),
		Schema:                  c.String(postgresSchema),
		LogMode:                 c.Bool(postgresLogMode),
		Host:                    c.String(postgresHost),
		ConnectionMaxLifetimeMS: c.Int(postgresConnectionMaxLifeTime),
		MaxOpenConnections:      c.Int(postgresMaxOpenConnection),
	}
}

// DatabaseFlags defines CLI flags for configuring PostgreSQL connections.
// These flags allow database connection settings to be specified via
// command-line arguments or environment variables.
var DatabaseFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    postgresUser,
		Value:   "test",
		Usage:   "PostgreSQL database username",
		EnvVars: []string{"POSTGRES_USER", "PG_USER"},
	},
	&cli.StringFlag{
		Name:    postgresPassword,
		Value:   "test",
		Usage:   "PostgreSQL database user password",
		EnvVars: []string{"POSTGRES_PASSWORD", "PG_PASSWORD"},
	},
	&cli.StringFlag{
		Name:    postgresDb,
		Value:   "node_art_slot_games",
		Usage:   "PostgreSQL database name",
		EnvVars: []string{"POSTGRES_DB", "PG_DB"},
	},
	&cli.StringFlag{
		Name:    postgresSchema,
		Value:   "public",
		Usage:   "PostgreSQL database schema",
		EnvVars: []string{"POSTGRES_SCHEMA", "PG_SCHEMA"},
	},
	&cli.StringFlag{
		Name:    postgresHost,
		Value:   "localhost:5432",
		Usage:   "PostgreSQL database host address",
		EnvVars: []string{"POSTGRES_HOST", "PG_HOST"},
	},
	&cli.StringFlag{
		Name:    postgresConnectionMaxLifeTime,
		Value:   "20",
		Usage:   "Maximum lifetime of a PostgreSQL connection in milliseconds",
		EnvVars: []string{"POSTGRES_CONNECTION_MAX_LIFE_TIME"},
	},
	&cli.StringFlag{
		Name:    postgresMaxOpenConnection,
		Value:   "300",
		Usage:   "Maximum number of open PostgreSQL connections",
		EnvVars: []string{"POSTGRES_MAX_OPEN_CONNECTION"},
	},
	&cli.BoolFlag{
		Name:    postgresLogMode,
		Value:   true,
		Usage:   "Enable or disable query logging in PostgreSQL",
		EnvVars: []string{"POSTGRES_QUERY_LOGGING"},
	},
}

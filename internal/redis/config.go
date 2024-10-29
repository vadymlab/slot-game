package redis

import "github.com/urfave/cli/v2"

// Constants defining the Redis configuration flags.
const (
	redisURL = "redis-url"
)

// Config represents the configuration settings required to connect to the Redis server.
// It includes a single field, URL, which specifies the Redis connection URL.
type Config struct {
	URL string // The Redis connection URL
}

// GetRedisConfig reads the Redis URL from the CLI context, allowing configuration via
// command-line arguments or environment variables.
//
// Parameters:
//   - c (*cli.Context): The CLI context containing flag and environment variable values.
//
// Returns:
//   - (*Config): A Config struct populated with the Redis URL.
func GetRedisConfig(c *cli.Context) *Config {
	return &Config{
		URL: c.String(redisURL),
	}
}

// Flags defines the CLI flags available for configuring the Redis connection.
// These flags enable the URL to be set via command-line arguments or environment variables.
var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    redisURL,                   // The flag name
		Value:   "redis://localhost:6379/0", // Default Redis URL for local connection
		Usage:   "Redis connection URL",     // Description for usage instructions
		EnvVars: []string{"REDIS_URL"},      // Environment variable to override the Redis URL
	},
}

package redis

import (
	libredis "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// Module provides the Redis client as an Fx module, enabling dependency injection
// for applications that require Redis as a data store.
var Module = fx.Options(
	fx.Provide(NewRedisClient),
)

// NewRedisClient initializes and returns a new Redis client instance configured with
// the provided Redis server URL from Config. This function parses the URL, creates
// a Redis client using the go-redis library, and ensures that the client can connect
// to the specified Redis server.
//
// Parameters:
//   - cfg (*Config): The configuration struct containing the Redis server URL.
//
// Returns:
//   - (*libredis.Client): The initialized Redis client instance.
//   - (error): An error if URL parsing or client creation fails.
func NewRedisClient(cfg *Config) (*libredis.Client, error) {
	option, err := libredis.ParseURL(cfg.Url)
	if err != nil {
		return nil, err
	}
	client := libredis.NewClient(option)
	return client, nil
}

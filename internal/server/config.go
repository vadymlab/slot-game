package server

import (
	"github.com/urfave/cli/v2"
)

// Constants defining CLI flags and environment variable names for API server configuration.
const (
	apiHost            = "server-host"                // API server host address
	apiPort            = "server-port"                // API server port
	apiMaxHeaderSize   = "server-max-header-size"     // Maximum size of request headers in bytes
	apiRequestTimeout  = "server-request-timeout"     // Maximum duration for reading request data
	apiResponseTimeout = "server-response-timeout"    // Maximum duration for writing response data
	jwtSecret          = "server-jwt-secret"          // JWT secret for authentication
	jwtSecretLifeTime  = "server-jwt-secret-lifetime" // JWT secret expiration time in minutes
	logRequest         = "server-log-request"         // Flag to enable or disable request logging
)

// APIConfig holds configuration settings for the API server.
type APIConfig struct {
	APIHost           string // Server host address
	APIPort           string // Server port number
	RequestTimeout    int    // Maximum request read duration in seconds
	ResponseTimeout   int    // Maximum response write duration in seconds
	MaxHeaderBytes    int    // Maximum size of request headers in bytes
	JWTSecret         string // JWT secret for signing tokens
	JWTSecretLifeTime int    // JWT token lifetime in minutes
	LogRequest        bool   // Enable request logging
}

// GetAPIConfig retrieves API server configuration from CLI flags or environment variables
// and initializes an APIConfig instance with these settings.
//
// Parameters:
//   - c: The CLI context containing parsed command-line arguments and environment variables.
//
// Returns:
//
//	A pointer to an ApiConfig instance populated with the specified configuration.
func GetAPIConfig(c *cli.Context) *APIConfig {
	return &APIConfig{
		APIHost:           c.String(apiHost),
		APIPort:           c.String(apiPort),
		RequestTimeout:    c.Int(apiRequestTimeout),
		ResponseTimeout:   c.Int(apiResponseTimeout),
		MaxHeaderBytes:    c.Int(apiMaxHeaderSize),
		LogRequest:        c.Bool(logRequest),
		JWTSecret:         c.String(jwtSecret),
		JWTSecretLifeTime: c.Int(jwtSecretLifeTime),
	}
}

// APIFlags defines a slice of CLI flags for configuring API server settings.
// These flags allow customization of server parameters through command-line arguments or environment variables.
var APIFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    apiHost,
		Value:   "0.0.0.0",
		Usage:   "API server host address",
		EnvVars: []string{"API_HOST"},
	},
	&cli.IntFlag{
		Name:    apiPort,
		Value:   8000,
		Usage:   "API server port",
		EnvVars: []string{"API_PORT"},
	},
	&cli.IntFlag{
		Name:    apiMaxHeaderSize,
		Value:   262144,
		Usage:   "Maximum size of request headers in bytes",
		EnvVars: []string{"API_MAX_HEADER_SIZE"},
	},
	&cli.IntFlag{
		Name:    apiRequestTimeout,
		Value:   5,
		Usage:   "Maximum duration for reading the entire request in seconds",
		EnvVars: []string{"API_REQUEST_TIMEOUT"},
	},
	&cli.IntFlag{
		Name:    apiResponseTimeout,
		Value:   5,
		Usage:   "Maximum duration before timing out writes of the response in seconds",
		EnvVars: []string{"API_RESPONSE_TIMEOUT"},
	},
	&cli.BoolFlag{
		Name:    logRequest,
		Value:   true,
		Usage:   "Enable or disable request logging",
		EnvVars: []string{"LOG_REQUEST"},
	},
	&cli.StringFlag{
		Name:    jwtSecret,
		Value:   "qi87x8Sd9KpQUuiOMP7gFMid3gRTQFjr",
		Usage:   "JWT secret used for signing authentication tokens",
		EnvVars: []string{"JWT_SECRET"},
	},
	&cli.IntFlag{
		Name:    jwtSecretLifeTime,
		Value:   60,
		Usage:   "JWT token lifetime in minutes",
		EnvVars: []string{"JWT_SECRET_LIFE_TIME"},
	},
}

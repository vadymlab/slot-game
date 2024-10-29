package config

import (
	log "github.com/public-forge/go-logger"
	"github.com/urfave/cli/v2"
)

// Constants for log configuration parameters.
const (
	loLevel = "log-level" // Log level setting for the application
	logJSON = "log-json"  // Determines if logs should be in JSON format
)

// GetLogConfig returns a configuration for logging based on CLI flags or environment variables.
// This includes the log level and whether logs should be formatted as JSON.
func GetLogConfig(c *cli.Context) *log.Config {
	conf := &log.Config{
		Level:  c.String(loLevel),
		IsJson: c.Bool(logJSON),
	}
	return conf
}

// LogFlags defines CLI flags for configuring the logging behavior.
// This includes flags for setting the log level and enabling JSON log formatting.
var LogFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    loLevel,
		Value:   "DEBUG",
		Usage:   "Application log level. Options: PANIC, FATAL, ERROR, WARNING, INFO, DEBUG, TRACE",
		EnvVars: []string{"LOG_LEVEL"},
	},
	&cli.BoolFlag{
		Name:    logJSON,
		Value:   false,
		Usage:   "Enable JSON format for logs. Defaults to false.",
		EnvVars: []string{"LOG_JSON"},
	},
}

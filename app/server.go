package app

import (
	"context"
	log "github.com/public-forge/go-logger"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"net/http"
)

// RunServer initializes and runs the server within an fx application lifecycle.
// It provides the CLI context, sets up logging, and manages the HTTP server lifecycle.
//
// Parameters:
//   - c: *cli.Context, a context object from the CLI, containing configuration and command line arguments.
//
// Returns:
//   - error: An error indicating the outcome of the server startup, or nil if the server runs successfully.
//
// Workflow:
//  1. Creates a new fx application (`newApp`) with `RootModule` as its primary dependency injection module.
//  2. Provides the CLI context as a dependency for use across the application.
//  3. Sets up logging by injecting the `log.Config` and creating a new logger.
//  4. Configures the HTTP server lifecycle, handling start and graceful shutdown operations.
//
// Lifecycle Management:
//   - OnStart: Launches the HTTP server in a separate goroutine to avoid blocking and logs the server start.
//   - OnStop: Gracefully shuts down the HTTP server by calling `srv.Shutdown`, waiting for ongoing requests to finish.
//
// Example usage:
//
//	err := RunServer(cliContext)
//	if err != nil {
//	    log.Fatalf("Failed to run server: %v", err)
//	}
func RunServer(c *cli.Context) error {

	newApp := fx.New(
		RootModule,
		fx.Provide(func() *cli.Context {
			return c
		}),
		// Sets up the logging configuration and initializes a logger instance.
		fx.Invoke(func(cfg *log.Config) (log.Logger, error) {
			return log.NewLogger(cfg)
		}),
		// Manages the HTTP server lifecycle using fx.Lifecycle hooks for OnStart and OnStop.
		fx.Invoke(func(lc fx.Lifecycle, srv *http.Server) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						err := srv.ListenAndServe()
						if err != nil {
							log.FromContext(ctx).Error(err)
						}
					}()
					log.FromContext(ctx).Info("Starting Auth API.")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return srv.Shutdown(ctx)
				},
			})
		}),
	)

	// Starts the fx application, blocking until the application exits or fails to start.
	newApp.Run()
	return nil
}

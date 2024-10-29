package app

import (
	"context"
	log "github.com/public-forge/go-logger"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"net/http"
)

func RunServer(c *cli.Context) error {

	newApp := fx.New(
		RootModule,
		fx.Provide(func() *cli.Context {
			return c
		}),
		fx.Invoke(func(cfg *log.Config) (log.Logger, error) {
			return log.NewLogger(cfg)
		}),
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
	newApp.Run()
	return nil
}

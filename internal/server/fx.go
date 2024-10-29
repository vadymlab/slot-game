package server

import "go.uber.org/fx"

// Module is an Fx module that provides dependencies for the server setup, including API configuration,
// the HTTP engine, and the server instance. These components are initialized using dependency injection
// and are required for running the API server.
var Module = fx.Module("server",

	// Provides API configuration settings, such as port and middleware options.
	fx.Provide(GetAPIConfig),

	// Provides a new HTTP engine, typically Gin, for handling HTTP requests and routing.
	fx.Provide(NewEngine),

	// Provides the server instance, which starts and runs the HTTP engine.
	fx.Provide(NewServer),
)

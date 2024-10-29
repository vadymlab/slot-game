package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/public-forge/go-logger"
	"github.com/vadymlab/slot-game/internal/middlewares"
	"net/http"
	"time"
)

// NewEngine creates and configures a new Gin engine instance.
// It applies middleware, including request logging (if enabled), request recovery, and CORS settings.
func NewEngine(config *ApiConfig) *gin.Engine {
	var router *gin.Engine
	if config.LogRequest {
		// Use the default Gin engine with logging and recovery middleware
		router = gin.Default()
	} else {
		// Create a new Gin engine without request logging
		router = gin.New()
		router.Use(gin.Recovery())
	}

	// Apply recovery middleware to handle panics gracefully
	router.Use(gin.Recovery())
	// Apply a trace middleware to manage request tracing IDs
	router.Use(middlewares.TraceMiddleware())

	// Configure CORS settings to allow all origins, methods, and headers,
	// with preflight requests cached for 12 hours
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))
	return router
}

// NewServer creates and configures a new HTTP server with a specified Gin router and API configuration.
// The server includes settings for address, timeouts, and max header bytes, with a timeout handler for request limits.
func NewServer(router *gin.Engine, config *ApiConfig) *http.Server {
	server := &http.Server{
		Addr:           config.ApiHost + ":" + config.ApiPort,                                                            // Server address
		Handler:        http.TimeoutHandler(router, time.Duration(config.RequestTimeout)*time.Second, "Request timeout"), // Timeout handler
		MaxHeaderBytes: config.MaxHeaderBytes,                                                                            // Maximum allowed header size
		ReadTimeout:    time.Duration(config.RequestTimeout) * time.Second,                                               // Timeout for reading request
		WriteTimeout:   time.Duration(config.ResponseTimeout) * time.Second,                                              // Timeout for writing response
	}

	// Log server startup details
	log.FromDefaultContext().Info("Starting server on " + config.ApiHost + ":" + config.ApiPort)
	return server
}

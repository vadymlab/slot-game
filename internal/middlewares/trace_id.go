package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/public-forge/go-logger"
	"github.com/vadymlab/slot-game/internal/constants"
)

// HeaderTraceID defines the HTTP header used for providing or propagating
// a trace ID, enabling consistent tracking of requests across services.
const HeaderTraceID = "X-Trace-ID"

// WithTraceID adds the trace ID to the context for request tracing purposes.
// This trace ID helps uniquely identify and track a request across services.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, constants.CtxFieldTraceID, traceID)
}

// WithLogger adds a logger instance to the context.
// The logger can then be retrieved from the context for structured logging with trace information.
func WithLogger(ctx context.Context, logger log.Logger) context.Context {
	return context.WithValue(ctx, constants.CtxFieldLogger, logger)
}

// TraceMiddleware is a middleware that attaches a trace ID and logger to both gin.Context and context.Context.
// If a trace ID is present in the request header, it is used; otherwise, a new trace ID is generated.
// The trace ID and logger are then added to the request context for consistent logging across services.
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or retrieve the trace ID from the request header.
		traceID := c.GetHeader(HeaderTraceID)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Attach the trace ID to gin.Context and context.Context.
		c.Set(string(constants.CtxFieldTraceID), traceID)
		ctx := WithTraceID(c.Request.Context(), traceID)

		// Initialize a logger with the trace ID and add it to the context for request logging.
		logger := log.FromContext(ctx).WithField(string(constants.CtxFieldTraceID), traceID)
		c.Set(string(constants.CtxFieldLogger), logger)
		ctx = WithLogger(ctx, logger)

		// Update the request with the new context containing trace ID and logger.
		c.Request = c.Request.WithContext(ctx)

		// Proceed to the next handler in the middleware chain.
		c.Next()
	}
}

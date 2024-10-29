package constants

type CtxKey string

// CtxFieldTraceID is the context key for storing the trace ID,
// which uniquely identifies a request for tracing across services.
const CtxFieldTraceID CtxKey = "trace_id"

// CtxFieldUserID is the context key for storing the user ID,
// allowing for user identification and access control throughout the request lifecycle.
const CtxFieldUserID CtxKey = "user_id"

// CtxFieldLogger is the context key for storing the logger instance,
// which facilitates structured and traceable logging within a request context.
const CtxFieldLogger CtxKey = "logger"

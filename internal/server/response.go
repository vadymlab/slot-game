package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/public-forge/go-logger"
	"net/http"
	"reflect"
)

// ErrorResponseMessage represents the structure of an error response with a list of error messages.
type ErrorResponseMessage struct {
	Errors []string `json:"errors"`
}

// SuccessResponse sends a successful HTTP response with status 200 and a response body.
func SuccessResponse(ctx *gin.Context, body interface{}) {
	response(ctx, http.StatusOK, body)
}

// UnauthorizedErrorResponse logs the error message and sends an unauthorized response with status 401.
// The function also aborts the current context.
func UnauthorizedErrorResponse(ctx *gin.Context, message string) {
	log.FromContext(ctx).Error(message)
	response(ctx, http.StatusUnauthorized, NewErrorMessage(message))
	ctx.Abort()
}

// ErrorsBadRequest logs the list of error messages and sends a bad request response with status 400.
// It uses a list of error messages and aborts the current context.
func ErrorsBadRequest(ctx *gin.Context, message []string) {
	log.FromContext(ctx).Error(message)
	response(ctx, http.StatusBadRequest, NewErrorMessages(message))
	ctx.Abort()
}

// ErrorBadRequest logs a single error message and sends a bad request response with status 400.
// The message can be of any type, and the context is aborted.
func ErrorBadRequest(ctx *gin.Context, message interface{}) {
	log.FromContext(ctx).Error(message)
	response(ctx, http.StatusBadRequest, NewErrorMessage(message))
	ctx.Abort()
}

// NewErrorMessage creates a new ErrorResponseMessage with a single error message.
// It accepts either an error object or a string and returns a pointer to ErrorResponseMessage.
func NewErrorMessage(err interface{}) *ErrorResponseMessage {
	var errorMessage string
	if e, ok := err.(error); ok {
		errorMessage = e.Error()
	} else if msg, ok := err.(string); ok {
		errorMessage = msg
	}
	return &ErrorResponseMessage{
		Errors: []string{errorMessage},
	}
}

// NewErrorMessages creates an ErrorResponseMessage with multiple error messages.
func NewErrorMessages(errors []string) *ErrorResponseMessage {
	return &ErrorResponseMessage{
		Errors: errors,
	}
}

// InternalErrorResponse logs the error message and sends an internal server error response with status 500.
// The function also aborts the current context.
func InternalErrorResponse(ctx *gin.Context, message string) {
	log.FromContext(ctx).Error(message)
	response(ctx, http.StatusInternalServerError, NewErrorMessage(message))
	ctx.Abort()
}

// ConflictErrorResponse logs the error message and sends a conflict response with status 409.
// The function also aborts the current context.
func ConflictErrorResponse(ctx *gin.Context, message string) {
	log.FromContext(ctx).Error(message)
	response(ctx, http.StatusConflict, NewErrorMessage(message))
	ctx.Abort()
}

// response sends an HTTP response based on the Accept header.
// Supports JSON and XML formats. Defaults to JSON if no specific format is requested.
// Handles nil and empty slice cases gracefully by setting appropriate HTTP status codes.
func response(ctx *gin.Context, code int, body interface{}) {
	accept := ctx.GetHeader("Accept")
	switch accept {
	case "application/json":
		ctx.JSON(code, body)
	case "application/xml":
		ctx.XML(code, body)
	default:
		if body != nil {
			v := reflect.ValueOf(body)
			if v.Kind() != reflect.Slice {
				ctx.JSON(code, body)
				return
			}
			if v.IsNil() || v.Len() == 0 {
				ctx.Status(code)
				return
			}
			ctx.JSON(code, body)
		} else {
			ctx.Status(code)
		}
	}
}

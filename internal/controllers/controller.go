package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vadymlab/slot-game/internal/constants"
	"github.com/vadymlab/slot-game/internal/server"
)

// BaseController defines a fundamental interface for controllers, offering
// essential methods for retrieving and initializing routes within a Gin router group.
type BaseController interface {

	// GetRoute provides the base route path for the controller.
	// This path serves as the primary endpoint for all routes within the controller.
	GetRoute() string

	// InitRoute sets up and initializes all routes for the controller within the specified
	// Gin router group. It returns the updated router group containing the controller’s configured routes.
	InitRoute(route *gin.RouterGroup) *gin.RouterGroup
}

// GetUserFromContext retrieves the user ID from the context, if available, and returns it as a UUID pointer.
// If the user ID is not present in the context, or if it is invalid, appropriate error responses
// are sent back to the client and nil is returned.
//
// Parameters:
//   - ctx: The Gin context from which to retrieve the user ID.
//
// Returns:
//
//	A pointer to the UUID representing the user's ID, or nil if the ID was not found
//	or could not be parsed.
func GetUserFromContext(ctx *gin.Context) *uuid.UUID {
	userID := ctx.GetString(string(constants.CtxFieldUserID))
	if userID == "" {
		server.UnauthorizedErrorResponse(ctx, "user not found")
		return nil
	}

	uUid, err := uuid.Parse(userID)
	if err != nil {
		server.InternalErrorResponse(ctx, err.Error())
		return nil
	}
	return &uUid
}

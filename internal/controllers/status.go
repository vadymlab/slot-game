package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vadymlab/slot-game/internal/server"
)

// StatusController is a simple controller that handles server status checks,
// providing a health check endpoint to verify that the server is up and running.
type StatusController struct{}

// NewStatusController creates a new instance of StatusController.
// This function provides a lightweight, easily initialized controller
// for handling server status requests.
//
// Returns:
//
//	A pointer to a StatusController instance.
func NewStatusController() *StatusController {
	return &StatusController{}
}

// InitRoute sets up the status check route under "/status", which responds with
// a success message. This route allows clients to confirm that the server is operational.
//
// Parameters:
//   - route: A Gin RouterGroup to which the status check route will be added.
//
// Returns:
//
//	An updated RouterGroup with the status route initialized.
func (c *StatusController) InitRoute(route *gin.RouterGroup) *gin.RouterGroup {
	route.GET("/status", c.onStatus)
	return route
}

// onStatus is the handler function for the "/status" endpoint.
// It responds with a success status in JSON format, indicating the server is running.
//
// @Summary Check server status
// @Description Returns a simple status message indicating the server is operational
// @Tags Status
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Server status message"
// @Router /api/status [get]
func (c *StatusController) onStatus(ctx *gin.Context) {
	server.SuccessResponse(ctx, gin.H{"status": "ok"})
}

// GetRoute returns the base route path for the StatusController.
// This path provides the main API route for server status checks.
func (c *StatusController) GetRoute() string {
	return "/api/status"
}

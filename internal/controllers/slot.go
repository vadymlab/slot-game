package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/public-forge/go-logger"
	libredis "github.com/redis/go-redis/v9"
	"github.com/vadymlab/slot-game/internal/config"
	"github.com/vadymlab/slot-game/internal/dto/request"
	"github.com/vadymlab/slot-game/internal/dto/response"
	serviceError "github.com/vadymlab/slot-game/internal/error"
	"github.com/vadymlab/slot-game/internal/interfaces"
	"github.com/vadymlab/slot-game/internal/middlewares"
	"github.com/vadymlab/slot-game/internal/server"
	"github.com/vadymlab/slot-game/internal/server/jwt"
	"github.com/vadymlab/slot-game/internal/validators"
)

// SlotController manages slot game operations, including processing spin requests
// and retrieving user spin history. It connects to slotService for core operations
// and applies JWT authentication for protected routes.
type SlotController struct {
	config      *server.ApiConfig       // API configuration, including JWT settings
	slotService interfaces.ISlotService // Service interface for slot game operations
	appConfig   *config.SlotConfig
	redisClient *libredis.Client
}

// NewSlotController initializes a new SlotController with the provided configuration
// and slotService, establishing the required dependencies for handling slot game operations.
//
// Parameters:
//   - config: A pointer to the API configuration struct.
//   - slotService: An implementation of the ISlotService interface for slot game functionality.
//
// Returns:
//
//	A pointer to a SlotController instance.
func NewSlotController(config *server.ApiConfig, appConfig *config.SlotConfig, redisClient *libredis.Client, slotService interfaces.ISlotService) *SlotController {
	return &SlotController{
		config:      config,
		slotService: slotService,
		appConfig:   appConfig,
		redisClient: redisClient,
	}
}

// InitRoute registers the slot game routes under the "/slot" endpoint, applying JWT
// middleware for authentication. Routes include "/spin" for spinning and "/history" for retrieving
// the user's spin history.
//
// Parameters:
//   - route: A Gin RouterGroup to which the slot game routes will be added.
//
// Returns:
//
//	An updated RouterGroup with initialized slot game routes.
func (c *SlotController) InitRoute(route *gin.RouterGroup) *gin.RouterGroup {
	g := route.Group("/slot", middlewares.NewRateLimiter(c.appConfig, c.redisClient), jwt.AuthMiddleware(c.config.JWTSecret))
	g.POST("/spin", c.spin)
	g.POST("/history", c.history)
	return route
}

// GetRoute returns the base route path for SlotController, serving as the primary
// API route for all slot game endpoints.
func (c *SlotController) GetRoute() string {
	return "/api"
}

// spin processes a slot spin request, validates input, retrieves the user ID from context,
// and invokes slotService.spin to perform the spin operation. If successful, it returns the spin result.
// In case of errors, it responds with appropriate error messages.
//
// @Summary spin the slot machine
// @Description Initiates a spin with the specified bet amount and returns the result.
// @Tags Slot
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param req body request.SpinRequest true "spin request body"
// @Success 200 {object} response.SpinResponse "spin result with win amount"
// @Failure 400 {string} string "Bad request due to invalid input or insufficient funds"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /api/slot/spin [post]
func (c *SlotController) spin(ctx *gin.Context) {
	req := request.SpinRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.FromContext(ctx).Error(err)
		server.ErrorBadRequest(ctx, err)
		return
	}
	if errs := validators.Validate(req); errs != nil {
		server.ErrorsBadRequest(ctx, errs)
		return
	}
	userId := GetUserFromContext(ctx)
	bit, err := c.slotService.RetrySpin(ctx.Request.Context(), userId, req.BetAmount)
	if err != nil {
		if errors.Is(err, serviceError.ErrInsufficientFunds) {
			server.ErrorBadRequest(ctx, err)
			return
		}
		server.InternalErrorResponse(ctx, err.Error())
		return
	}
	server.SuccessResponse(ctx, response.SpinFromModel(bit))
}

// history retrieves the user's spin history from slotService and returns it as a structured
// response. If an error occurs, it responds with an internal server error message.
//
// @Summary Get spin history
// @Description Retrieves the user's spin history, showing past spins with their results
// @Tags Slot
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} response.SpinHistoryResponse "List of past spin results"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /api/slot/history [post]
func (c *SlotController) history(ctx *gin.Context) {
	userId := GetUserFromContext(ctx)
	history, err := c.slotService.History(ctx.Request.Context(), userId)
	if err != nil {
		server.InternalErrorResponse(ctx, err.Error())
		return
	}
	server.SuccessResponse(ctx, response.SpinHistoryFromModels(history))
}

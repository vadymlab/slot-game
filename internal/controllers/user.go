package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/public-forge/go-logger"
	"github.com/vadymlab/slot-game/internal/dto/request"
	"github.com/vadymlab/slot-game/internal/dto/response"
	serviceError "github.com/vadymlab/slot-game/internal/error"
	"github.com/vadymlab/slot-game/internal/interfaces"
	"github.com/vadymlab/slot-game/internal/server"
	mw "github.com/vadymlab/slot-game/internal/server/jwt"
	"github.com/vadymlab/slot-game/internal/validators"
)

// UserController manages user-related actions, including registration, login, and profile retrieval.
// It connects to userService for core user operations and uses JWT authentication for protected routes.
type UserController struct {
	userService interfaces.IUserService // Service for managing user-related operations
	config      *server.APIConfig       // API configuration with JWT settings
}

// NewUserController creates a new instance of UserController with the given userService and config.
//
// Parameters:
//   - userService: Implementation of IUserService for user business logic.
//   - config: API configuration, including JWT settings.
//
// Returns:
//
//	A pointer to UserController.
func NewUserController(userService interfaces.IUserService, config *server.APIConfig) *UserController {
	return &UserController{
		userService: userService,
		config:      config,
	}
}

// InitRoute initializes routes for user-related endpoints, including registration, login, and profile retrieval.
// The profile endpoint is protected and requires JWT authentication.
//
// Parameters:
//   - route: A Gin RouterGroup to which user routes will be added.
//
// Returns:
//
//	An updated RouterGroup with initialized user routes.
func (c *UserController) InitRoute(route *gin.RouterGroup) *gin.RouterGroup {
	route.POST("/register", c.register)
	route.POST("/login", c.login)
	route.GET("/profile", mw.AuthMiddleware(c.config.JWTSecret), c.profile)
	return route
}

// GetRoute returns the base route for UserController, used for defining the primary API route.
func (c *UserController) GetRoute() string {
	return "/api"
}

// register handles user registration by validating input, checking for existing users,
// and creating a new user. If successful, it returns the user’s registration details.
//
// @Summary Register a new user
// @Description Allows a new user to register with their details
// @Tags User
// @Accept json
// @Produce json
// @Param req body request.RegisterRequest true "Registration request body"
// @Success 200 {object} response.RegisterResponse "User registered successfully"
// @Failure 400 {string} string "Bad request due to invalid input"
// @Failure 409 {string} string "Conflict - user already exists"
// @Failure 500 {string} string "Internal server error"
// @Router /api/register [post]
func (c *UserController) register(ctx *gin.Context) {
	req := request.RegisterRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.FromContext(ctx).Error(err)
		server.ErrorBadRequest(ctx, err)
		return
	}
	if errs := validators.Validate(req); errs != nil {
		server.ErrorsBadRequest(ctx, errs)
		return
	}
	user, err := c.userService.Register(ctx.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.As(err, &serviceError.UserAlreadyExists{}) {
			server.ConflictErrorResponse(ctx, err.Error())
			return
		}
		server.InternalErrorResponse(ctx, err.Error())
		return
	}
	server.SuccessResponse(ctx, response.RegisterFromModel(user))
}

// login authenticates a user by validating credentials and generating a JWT token if successful.
// Returns a token upon successful authentication; otherwise, returns an error.
//
// @Summary Login user
// @Description Authenticates a user and returns a JWT token
// @Tags User
// @Accept json
// @Produce json
// @Param req body request.LoginRequest true "Login request body"
// @Success 200 {object} map[string]string "Token for authenticated user"
// @Failure 400 {string} string "Bad request due to invalid input or incorrect login details"
// @Failure 500 {string} string "Internal server error"
// @Router /api/login [post]
func (c *UserController) login(ctx *gin.Context) {
	req := request.LoginRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.FromContext(ctx).Error(err)
		server.ErrorBadRequest(ctx, err)
		return
	}
	if errs := validators.Validate(req); errs != nil {
		server.ErrorsBadRequest(ctx, errs)
		return
	}
	usr, err := c.userService.Login(ctx.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, &serviceError.UserNotFound{}) {
			server.ErrorBadRequest(ctx, err)
			return
		}
		if errors.Is(err, &serviceError.InvalidPassword{}) {
			server.ErrorBadRequest(ctx, err)
			return
		}
	}
	token, err := mw.GenerateToken(usr.ExternalID, c.config.JWTSecret, c.config.JWTSecretLifeTime)
	if err != nil {
		server.InternalErrorResponse(ctx, err.Error())
		return
	}
	server.SuccessResponse(ctx, gin.H{"token": "Bearer " + token})
}

// profile retrieves the profile details of the authenticated user, including the user's ID, login, and balance.
// This endpoint requires JWT authentication.
//
// @Summary Get user profile
// @Description Retrieves the profile and balance of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.ProfileResponse "User profile information"
// @Failure 401 {string} string "Unauthorized - user not authenticated"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /api/profile [get]
func (c *UserController) profile(ctx *gin.Context) {
	uUID := GetUserFromContext(ctx)
	user, err := c.userService.GetByExternalID(ctx.Request.Context(), uUID)
	if err != nil {
		server.InternalErrorResponse(ctx, err.Error())
		return
	}
	responseDto := response.ProfileResponse{
		ID:      user.ExternalID,
		Login:   user.Login,
		Balance: user.Balance,
	}
	server.SuccessResponse(ctx, responseDto)
}

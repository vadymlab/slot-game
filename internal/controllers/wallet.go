package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/public-forge/go-logger"
	"github.com/vadymlab/slot-game/internal/dto/request"
	"github.com/vadymlab/slot-game/internal/dto/response"
	error2 "github.com/vadymlab/slot-game/internal/error"
	"github.com/vadymlab/slot-game/internal/interfaces"
	"github.com/vadymlab/slot-game/internal/server"
	"github.com/vadymlab/slot-game/internal/server/jwt"
	"github.com/vadymlab/slot-game/internal/validators"
)

// WalletController manages wallet-related operations, including depositing and withdrawing funds.
type WalletController struct {
	config      *server.ApiConfig       // API configuration settings, including JWT secret
	userService interfaces.IUserService // Service for user-related operations
}

// NewWalletController creates a new instance of WalletController with the provided API configuration and user service.
//
// Parameters:
//   - config: A pointer to the API configuration struct, including JWT settings.
//   - userService: Implementation of IUserService for managing user wallet operations.
//
// Returns:
//
//	A pointer to WalletController.
func NewWalletController(config *server.ApiConfig, userService interfaces.IUserService) *WalletController {
	return &WalletController{
		config:      config,
		userService: userService,
	}
}

// InitRoute initializes wallet-related routes within the provided router group,
// including deposit and withdraw endpoints, both protected by JWT authentication middleware.
//
// Parameters:
//   - route: A Gin RouterGroup to which wallet routes will be added.
//
// Returns:
//
//	An updated RouterGroup with initialized wallet routes.
func (c *WalletController) InitRoute(route *gin.RouterGroup) *gin.RouterGroup {
	g := route.Group("/wallet", jwt.AuthMiddleware(c.config.JWTSecret))
	g.POST("/deposit", c.deposit)
	g.POST("/withdraw", c.withdraw)
	return route
}

// GetRoute returns the base route path for WalletController.
func (c *WalletController) GetRoute() string {
	return "/api"
}

// deposit handles fund deposits to the user's wallet.
//
// @Summary      Deposit funds into wallet
// @Description  Allows the user to deposit funds into their wallet
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string              true  "JWT Token"                    format(bearer)
// @Param        data           body      request.DepositRequest true  "Deposit amount"
// @Success      200            {object}  response.DepositResponse "Updated wallet balance"
// @Failure      400            {string}  string "Invalid request payload"
// @Failure      401            {string}  string "Unauthorized - user not authenticated"
// @Failure      500            {string}  string "Internal server error"
// @Security     BearerAuth
// @Router       /api/wallet/deposit [post]
func (c *WalletController) deposit(ctx *gin.Context) {
	req := request.DepositRequest{}
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
	balance, err := c.userService.Deposit(ctx.Request.Context(), userId, req.Amount)
	if err != nil {
		if errors.Is(err, error2.ErrInvalidAmount) || errors.Is(err, error2.ErrInsufficientFunds) {
			server.ErrorBadRequest(ctx, err)
			return
		}
		server.InternalErrorResponse(ctx, err.Error())
		return
	}
	responseDto := response.DepositResponse{
		Balance: *balance,
	}
	server.SuccessResponse(ctx, responseDto)
}

// withdraw handles fund withdrawals from the user's wallet.
//
// @Summary      Withdraw funds from wallet
// @Description  Allows the user to withdraw funds from their wallet
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                true  "JWT Token"                    format(bearer)
// @Param        data           body      request.WithdrawRequest true  "Withdraw amount"
// @Success      200            {object}  response.WithdrawResponse "Updated wallet balance"
// @Failure      400            {string}  string "Invalid request payload"
// @Failure      401            {string}  string "Unauthorized - user not authenticated"
// @Failure      500            {string}  string "Internal server error"
// @Security     BearerAuth
// @Router       /api/wallet/withdraw [post]
func (c *WalletController) withdraw(ctx *gin.Context) {
	req := request.WithdrawRequest{}
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
	balance, err := c.userService.Withdraw(ctx.Request.Context(), userId, req.Amount)
	if err != nil {
		server.InternalErrorResponse(ctx, err.Error())
		return
	}
	responseDto := response.WithdrawResponse{
		Balance: *balance,
	}
	server.SuccessResponse(ctx, responseDto)
}

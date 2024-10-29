package app

import (
	"github.com/gin-gonic/gin"
	log "github.com/public-forge/go-logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/vadymlab/slot-game/docs" // Import for loading Swagger documentation
	"github.com/vadymlab/slot-game/internal/config"
	controller "github.com/vadymlab/slot-game/internal/controllers"
	"github.com/vadymlab/slot-game/internal/database"
	"github.com/vadymlab/slot-game/internal/redis"
	"github.com/vadymlab/slot-game/internal/repository"
	"github.com/vadymlab/slot-game/internal/server"
	"github.com/vadymlab/slot-game/internal/service"
	"go.uber.org/fx"
)

// initController initializes routes for a given controller.
// It sets up a handler group in the router based on the controller's base route
// and invokes InitRoute to register each endpoint.
func initController(router *gin.Engine, ctrl controller.BaseController) {
	handler := router.Group(ctrl.GetRoute())
	ctrl.InitRoute(handler)
}

// ConfigModule sets up the configuration dependencies for the application.
// It includes providers for logging, slot configuration, and Redis configuration.
var ConfigModule = fx.Module("config",
	fx.Provide(config.GetLogConfig),
	fx.Provide(config.GetSlotConfig),
	fx.Provide(redis.GetRedisConfig),
)

// Repositories defines providers for the repository layer, which is responsible
// for data persistence and retrieval logic. Includes providers for UserRepository
// and SlotRepository, which handle user data and slot game data, respectively.
var Repositories = fx.Provide(
	repository.NewUserRepository,
	repository.NewSlotRepository,
)

// Services defines providers for the service layer, which contains business logic.
// It includes UserService and SlotService, handling operations related to user
// management and slot game logic.
var Services = fx.Provide(
	service.NewUserService,
	service.NewSlotService,
)

// Controllers defines providers for HTTP controllers, responsible for handling
// HTTP requests and interacting with the service layer. This includes controllers
// for user management, system status, wallet operations, and slot game endpoints.
var Controllers = fx.Provide(
	controller.NewUserController,
	controller.NewStatusController,
	controller.NewWalletController,
	controller.NewSlotController,
)

// RootModule orchestrates the complete application setup, assembling repositories,
// services, controllers, and configurations into an fx.Module for dependency injection.
//
// Additionally, it sets up Swagger API documentation, initializes HTTP controllers,
// and enables logging capabilities.
var RootModule = fx.Module("server",
	Repositories,
	Services,
	Controllers,
	ConfigModule,
	database.DBModule,
	server.Module,
	redis.Module,
	fx.Provide(log.NewLogger),
	fx.Invoke(func(router *gin.Engine,

		userController *controller.UserController,
		statusController *controller.StatusController,
		walletController *controller.WalletController,
		slotController *controller.SlotController,
	) {
		// Registers Swagger API documentation handler on /swagger endpoint
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Initializes routes for each controller in the application
		initController(router, userController)
		initController(router, statusController)
		initController(router, walletController)
		initController(router, slotController)
	}),
)

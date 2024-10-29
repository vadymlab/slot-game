package app

import (
	"github.com/gin-gonic/gin"
	log "github.com/public-forge/go-logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/vadymlab/slot-game/docs"
	"github.com/vadymlab/slot-game/internal/config"
	controller "github.com/vadymlab/slot-game/internal/controllers"
	"github.com/vadymlab/slot-game/internal/database"
	"github.com/vadymlab/slot-game/internal/redis"
	"github.com/vadymlab/slot-game/internal/repository"
	"github.com/vadymlab/slot-game/internal/server"
	"github.com/vadymlab/slot-game/internal/service"
	"go.uber.org/fx"
)

func initController(router *gin.Engine, ctrl controller.BaseController) {
	handler := router.Group(ctrl.GetRoute())
	ctrl.InitRoute(handler)
}

var ConfigModule = fx.Module("config",
	fx.Provide(config.GetLogConfig),
	fx.Provide(config.GetSlotConfig),
	fx.Provide(redis.GetRedisConfig),
)
var Repositories = fx.Provide(
	repository.NewUserRepository,
	repository.NewSlotRepository,
)
var Services = fx.Provide(
	service.NewUserService,
	service.NewSlotService,
)
var Controllers = fx.Provide(

	controller.NewUserController,
	controller.NewStatusController,
	controller.NewWalletController,
	controller.NewSlotController,
)
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
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		initController(router, userController)
		initController(router, statusController)
		initController(router, walletController)
		initController(router, slotController)
	}),
)

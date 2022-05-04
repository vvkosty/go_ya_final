package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	config "github.com/vvkosty/go_ya_final/internal/app/config"
	handler "github.com/vvkosty/go_ya_final/internal/app/handlers"
	"github.com/vvkosty/go_ya_final/internal/app/helpers"
	"github.com/vvkosty/go_ya_final/internal/app/integrations"
	middleware "github.com/vvkosty/go_ya_final/internal/app/middlewares"
	storage "github.com/vvkosty/go_ya_final/internal/app/storage"
)

type App struct {
	Config           *config.Config
	Storage          storage.Repository
	Handler          *handler.Handler
	Middleware       *middleware.Middleware
	Encoder          *helpers.Encoder
	AccrualAPIClient *integrations.AccrualAPIClient
}

func (app *App) Init() {
	app.Handler.UserStorage = &storage.UserStorage{DB: app.Storage.Instance()}
	app.Handler.OrderStorage = &storage.OrderStorage{DB: app.Storage.Instance()}
	app.Handler.UserBalanceStorage = &storage.UserBalanceStorage{DB: app.Storage.Instance()}
	app.Handler.WithdrawHistoryStorage = &storage.WithdrawHistoryStorage{DB: app.Storage.Instance()}

	app.Handler.Config = app.Config
	app.Handler.Encoder = app.Encoder
	app.Handler.AccrualAPIClient = app.AccrualAPIClient

	app.Middleware.Config = app.Config
	app.Middleware.Encoder = app.Encoder
}

func (app *App) Start() {
	app.Migrate()
	r := app.SetupRouter()

	err := r.Run(app.Config.Address)
	if err != nil {
		fmt.Println(err)
	}
}

func (app *App) SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	r.POST("/api/user/register", app.Handler.RegisterUser)
	r.POST("/api/user/login", app.Handler.LoginUser)

	v1 := r.Group("/api/user").Use(app.Middleware.NeedAuth)
	{
		v1.GET("/orders", app.Handler.GetOrders)
		v1.GET("/balance", app.Handler.GetUserBalance)
		v1.GET("/balance/withdrawals", app.Handler.GetUserWithdrawals)

		v1.POST("/orders", app.Handler.SaveOrder)
		v1.POST("/balance/withdraw", app.Handler.Withdraw)
	}

	r.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})

	return r
}

func (app *App) Migrate() {
	driver, err := postgres.WithInstance(app.Storage.Instance(), &postgres.Config{})

	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/app/migrations", "postgres", driver)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		log.Print(err)
	}
}

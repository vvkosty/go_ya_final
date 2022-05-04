package main

import (
	"github.com/vvkosty/go_ya_final/internal/app"
	config "github.com/vvkosty/go_ya_final/internal/app/config"
	handler "github.com/vvkosty/go_ya_final/internal/app/handlers"
	"github.com/vvkosty/go_ya_final/internal/app/helpers"
	"github.com/vvkosty/go_ya_final/internal/app/integrations"
	middleware "github.com/vvkosty/go_ya_final/internal/app/middlewares"
	storage "github.com/vvkosty/go_ya_final/internal/app/storage"
)

func main() {
	var appConfig config.Config
	var appHandler handler.Handler
	var appMiddleware middleware.Middleware

	appConfig.LoadEnvs()
	appConfig.ParseCommandLine()

	application := app.App{
		Config:     &appConfig,
		Handler:    &appHandler,
		Middleware: &appMiddleware,
	}

	application.Encoder = &helpers.Encoder{SecretKey: appConfig.AppKey}
	application.AccrualAPIClient = &integrations.AccrualAPIClient{BaseAddress: appConfig.AccrualSystemAddress}

	application.Storage = storage.NewPostgresDatabase(appConfig.DatabaseDsn)
	defer application.Storage.Close()

	application.Init()
	application.Start()
}

package main

import (
	"github.com/vvkosty/go_ya_final/internal/app"
	config "github.com/vvkosty/go_ya_final/internal/app/config"
	handler "github.com/vvkosty/go_ya_final/internal/app/handlers"
	"github.com/vvkosty/go_ya_final/internal/app/helpers"
	middleware "github.com/vvkosty/go_ya_final/internal/app/middlewares"
	storage "github.com/vvkosty/go_ya_final/internal/app/storage"
)

func main() {
	var appConfig config.Config
	var appHandler handler.Handler
	var appMiddleware middleware.Middleware
	var appEncoder helpers.Encoder

	appConfig.LoadEnvs()
	appConfig.ParseCommandLine()

	application := app.App{
		Config:     &appConfig,
		Handler:    &appHandler,
		Middleware: &appMiddleware,
		Encoder:    &appEncoder,
	}

	application.Encoder = &helpers.Encoder{SecretKey: appConfig.AppKey}

	application.Storage = storage.NewPostgresDatabase(appConfig.DatabaseDsn)
	defer application.Storage.Close()

	application.Init()
	application.Start()
}

package main

import (
	"github.com/vvkosty/go_ya_final/internal/app"
	config "github.com/vvkosty/go_ya_final/internal/app/config"
	handler "github.com/vvkosty/go_ya_final/internal/app/handlers"
	middleware "github.com/vvkosty/go_ya_final/internal/app/middlewares"
	storage "github.com/vvkosty/go_ya_final/internal/app/storage"
)

func main() {
	var appConfig config.ServerConfig
	var appHandler handler.Handler
	var appMiddleware middleware.Middleware

	appConfig.LoadEnvs()
	appConfig.ParseCommandLine()

	application := app.App{
		Config:     &appConfig,
		Handler:    &appHandler,
		Middleware: &appMiddleware,
	}

	application.Storage = storage.NewPostgresDatabase(appConfig.DatabaseDsn)
	defer application.Storage.Close()

	application.Init()
	application.Start()
}

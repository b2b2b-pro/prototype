/*
Main инициализирует чтение конфига, запуск логгера и выполнение приложения.
Main initializes reading the config, launching the logger, and executing the application.
*/
package main

import (
	"log"

	"github.com/b2b2b-pro/prototype/http2repo/config"
	"github.com/b2b2b-pro/prototype/http2repo/internal/app"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Printf("Can't initialize logger, err: %v", err)
	}

	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	app.Run(cfg)
}

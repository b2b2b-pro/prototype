/*
Main инициализирует чтение конфига, запуск логгера и выполнение приложения.
*/
package main

import (
	"log"
	"os"

	"github.com/b2b2b-pro/prototype/xml2repo/config"
	"github.com/b2b2b-pro/prototype/xml2repo/internal/app"
	"go.uber.org/zap"
)

func main() {
	// временно, для дебага
	var f []string

	if len(os.Args) < 2 {
		// log.Fatal("Usage xml2repo file\n")
		f = append(f, "../../contract_1241101171822000013_74728698.xml")
	} else {
		f = os.Args[1:]
	}

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

	app.Run(cfg, f)
}

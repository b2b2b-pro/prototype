/*
App
*/
package app

import (
	"github.com/b2b2b-pro/prototype/xml2repo/config"
	"github.com/b2b2b-pro/prototype/xml2repo/internal/controller/fz44"

	"github.com/b2b2b-pro/lib/repo_client"
	"go.uber.org/zap"
)

// app.Run(*config.Config) - основная функция пакета, создаёт и запускает компоненты приложения
func Run(cfg *config.Config, files []string) {
	var err error

	zap.S().Debugf("Run. Config: %v\n", cfg)
	rp, err := repo_client.New(&cfg.ConfigRPC)
	if err != nil {
		zap.S().Fatal("GRPC to Repository initialization error: ", err, "\n")
	}
	defer rp.Stop()

	xp := fz44.NewParser(rp)

	for _, x := range files {
		err = xp.ParseXML(x)
		zap.S().Debugf("ParseXML %v error: %v\n", x, err)
	}

	rp.Stop()
}

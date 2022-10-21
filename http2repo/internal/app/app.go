/*
App запускает компоненты приложения:
- коннектор к сервису, предоставляющему доступ к репозиторию;
- http роутер;
- http сервер.
Слушает системные сигналы (SIGHUP, SIGINT, SIGTERM, SIGQUIT)
и при получении останавливает приложение.
*/
package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/b2b2b-pro/prototype/http2repo/config"
	"github.com/b2b2b-pro/prototype/http2repo/internal/controller/httprouter"
	"github.com/b2b2b-pro/prototype/http2repo/pkg/httpserver"

	"github.com/b2b2b-pro/lib/repo_client"
	"go.uber.org/zap"
)

// app.Run(*config.Config) - основная функция пакета, создаёт и запускает компоненты приложения
func Run(cfg *config.Config) {
	var err error

	zap.S().Debugf("Run. Config: %v\n", cfg)

	rp, err := repo_client.New(&cfg.ConfigRPC)
	if err != nil {
		zap.S().Debugf("GRPC to Repository initialization error: ", err, "\n")
		time.Sleep(time.Second * 10)
		rp, err = repo_client.New(&cfg.ConfigRPC)
		if err != nil {
			zap.S().Fatalf("GRPC to Repository initialization error: ", err, "\n")
		}
	}

	defer rp.Stop()

	// Создание http router'а
	wr := httprouter.New(rp)

	// Создание http server'а
	hs := httpserver.New(cfg, wr.R)

	// старт http server'а
	hs.Start()

	// обработка системных сигналов
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case err = <-hs.Notify(): // если http server завершился с ошибкой
		zap.S().Error("HTTP Server error: %v\n", err)
	case sig := <-interrupt: // если поступил сигнал от системы
		zap.S().Info("System call: %v\n", sig)
	}

	// остановить http server, repo
	hs.Stop()
	rp.Stop()
}

package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/b2b2b-pro/prototype/auth/config"
	"github.com/b2b2b-pro/prototype/auth/internal/controller/httprouter"
	"github.com/b2b2b-pro/prototype/auth/pkg/httpserver"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {
	zap.S().Debugf("Auth App.Run\n")

	var err error

	// Создание http router'а
	wr := httprouter.New()

	// Создание http server'а
	hs := httpserver.New(&cfg.ConfigHTTP, wr.R)

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

}

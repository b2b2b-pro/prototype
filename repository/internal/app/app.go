/*
App запускает компоненты приложения:
- реализацию репозитория (базу данных на postgres или простенькую эмуляцию в памяти);
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

	"github.com/b2b2b-pro/lib/repo_srv"
	"github.com/b2b2b-pro/prototype/repository/config"
	"github.com/b2b2b-pro/prototype/repository/internal/repo/repopg"

	"go.uber.org/zap"
)

// app.Run(*config.Config) - основная функция пакета, создаёт и запускает компоненты приложения
func Run(cfg *config.Config) {
	zap.S().Debugf("App.Run config: %v\n", cfg)
	var err error

	// создание базы данным
	rp, err := repopg.New(cfg) // база в postgresql

	if err != nil {
		zap.S().Fatal("Repository initialization error: ", err, "\n")
	}
	defer rp.Stop()

	gs, err := repo_srv.New(cfg.ConfigRPC, rp)
	if err != nil {
		zap.S().Fatalf("repo_srv create error: %v\n", err)
	}

	gs.Start()

	// обработка системных сигналов
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-interrupt: // если поступил сигнал от системы
		zap.S().Info("System call: %v\n", sig)
	case err = <-gs.Notify():
		zap.S().Errorf("gRPC server eroor: %v\n", err)
	}
	gs.Stop()
}

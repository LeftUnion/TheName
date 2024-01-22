package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/LeftUnion/theName/service/TheName/internal/autogen"
	configs "github.com/LeftUnion/theName/service/TheName/internal/config"
	"github.com/LeftUnion/theName/service/TheName/internal/database/postgres"
	"github.com/LeftUnion/theName/service/TheName/internal/service"
	"github.com/LeftUnion/theName/service/TheName/internal/transport"

	"github.com/sirupsen/logrus"
)

const (
	InfoCfgReaded         = "Конфигурация сервиса прочитана"
	InfoDataBaseConnected = "База данных подключена"
	InfoTransportCreated  = "Транспонртный слой проинициализирован"
	InfoServiceCreated    = "Сервисный слой проинициализирован"
)

func Run() {
	// Чтение конфигурации сервиса.
	cfg := configs.New()
	logrus.Info(InfoCfgReaded)

	// Создание слоя базы данных.
	db := postgres.NewPostgres(cfg.GetDatabaseCfg())
	logrus.Info(InfoDataBaseConnected)

	// Создание сервисного слоя.
	srv := service.NewService(db)
	logrus.Info(InfoServiceCreated)

	// Создание транспортного слоя.
	transport := transport.NewTransport(srv)
	logrus.Info(InfoTransportCreated)

	// Адрес сервера.
	serverCfg := cfg.GetServerCfg()
	addr := serverCfg.Host + ":" + serverCfg.Port

	// TODO: middleware - logger
	// router := http.NewServeMux()
	// router.Handle("/", middlewares.Use(autogen.Handler(transport), middlewares.Logger()))

	server := http.Server{
		Addr:    addr,
		Handler: autogen.Handler(transport),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Panicf("Ошибка старта сервиса: %s", err)
		}
		logrus.Info("Сервис остановлен")
	}()
	logrus.Infof("Сервис запущен %s", addr)

	channel := make(chan os.Signal, 1)
	signal.Notify(channel,
		syscall.SIGABRT,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-channel

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Panicf("Ошибка остановки сервиса: %s", err)
	}
}

package app

import (
	"avito-rest-api/config"
	v1 "avito-rest-api/internal/controller/http/v1"
	"avito-rest-api/internal/repository"
	"avito-rest-api/internal/service"
	"avito-rest-api/package/httpserver"
	"avito-rest-api/package/postgres"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Run(configPath string) {
	// Чтение файла конфигурации
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to intialize config: %s", err)
	}

	// Инициализация логгера
	err = SetupLogrus(cfg.Log.Level, cfg.Log.LogsPath)
	if err != nil {
		panic(fmt.Sprintf("failed to setup logger due to error: %s", err))
	}

	// Инициализация клиента базы данных
	log.Info("Initializing database...")
	pg, err := postgres.New(
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Database,
		cfg.PostgreSQL.Username,
		cfg.PostgreSQL.Password,
		postgres.MaxPoolSize(cfg.PostgreSQL.MaxPoolSize),
	)
	defer pg.Close()

	// Инициализация слоя-репозитория
	log.Info("Initializing repositories...")
	repositories := repository.NewRepositories(pg)

	// Инициализация сервисов
	log.Info("Initializing services...")
	dependencies := service.ServicesDependencies{Repositories: repositories}
	services := service.NewService(dependencies)

	// Echo-обработчик
	log.Info("Initializing echo...")
	handler := echo.New()
	v1.NewRouter(handler, services)

	// HTTP-сервер
	log.Info("Starting HTTP server...")
	log.Debugf("Startint HTTP server by %s:%s", cfg.HTTP.BindIP, cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Address(cfg.HTTP.BindIP, cfg.HTTP.Port))

	// Ожидание сигнала к выключению
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Infof("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		log.Errorf("app - Run - httpServer.Notify: %s", err)
	}

	// Плавное выключение
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Errorf("app - Run - httpServer.Shutdown: %s", err)
	}
}

// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"backend-test/config"
	"backend-test/internal/business/usecase"
	"backend-test/internal/business/usecase/repo"
	v1 "backend-test/internal/controller/http/v1"
	"backend-test/internal/db/gorm/mysql"
	"backend-test/pkg/httpserver"
	"backend-test/pkg/logger"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	db, err := mysql.NewGormMysql(cfg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - db.New: %w", err))
	}

	bussinessRepo := repo.NewBusinessRepo(db, l)

	if err := bussinessRepo.MigrateAndMockData(); err != nil {
		l.Fatal(fmt.Errorf("app - Run - businessRepo.MigrateAndMock: %w", err))
	}

	// Use case
	businessUseCase := usecase.NewBusinessUseCase(
		bussinessRepo,
		l,
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, businessUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}

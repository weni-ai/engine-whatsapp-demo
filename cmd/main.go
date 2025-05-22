package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/metric"
	"github.com/weni/whatsapp-router/servers/grpc"
	"github.com/weni/whatsapp-router/servers/http"
	"github.com/weni/whatsapp-router/storage"
)

func main() {
	db := storage.NewDB()
	defer storage.CloseDB(db)
	logger.Info("Starting application...")

	var err error
	metrics, err := metric.NewPrometheusService()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	httpServer := http.NewServer(db, metrics)
	if err := httpServer.Start(); err != nil {
		logger.Error(fmt.Sprintf("Server startup failed: %v", err))
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(db, metrics)
	if err := grpcServer.Start(); err != nil {
		logger.Error(fmt.Sprintf("grpc server startup failed: %v", err))
		os.Exit(1)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	signal := <-ch
	logger.Info(fmt.Sprintf("WHATSAPP ROUTER STOPING, signal %v", signal))

}

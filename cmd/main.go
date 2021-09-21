package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/servers/grpc"
	"github.com/weni/whatsapp-router/servers/http"
	"github.com/weni/whatsapp-router/storage"
)

func main() {
	db := storage.NewDB()
	logger.Info("Starting application...")

	httpServer := http.NewServer(db)
	if err := httpServer.Start(); err != nil {
		logger.Error(fmt.Sprintf("Server startup failed: %v", err))
		log.Fatal()
	}

	grpcServer := grpc.NewServer(db)
	if err := grpcServer.Start(); err != nil {
		logger.Error(fmt.Sprintf("grpc server startup failed: %v", err))
		log.Fatal()
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	signal := <-ch
	logger.Info(fmt.Sprintf("WHATSAPP ROUTER STOPING, signal %v", signal))

}

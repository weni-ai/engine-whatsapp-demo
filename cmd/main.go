package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/weni/whatsapp-router/servers/grpc"
	"github.com/weni/whatsapp-router/servers/http"
	"github.com/weni/whatsapp-router/storage"
)

func main() {
	db := storage.NewDB()

	httpServer := http.NewServer(db)
	if err := httpServer.Start(); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}

	grpcServer := grpc.NewServer(db)
	if err := grpcServer.Start(); err != nil {
		log.Fatalf("grpc server startup failed: %v", err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	signal := <-ch
	log.Printf("WHATSAPP ROUTER STOPING, signal %v", signal)

}

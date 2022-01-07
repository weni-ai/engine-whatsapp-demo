package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/repositories"
	"github.com/weni/whatsapp-router/servers/grpc"
	"github.com/weni/whatsapp-router/servers/http"
	"github.com/weni/whatsapp-router/services"
	"github.com/weni/whatsapp-router/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	db := storage.NewDB()
	defer storage.CloseDB(db)
	logger.Info("Starting application...")

	initAuthToken(db)

	httpServer := http.NewServer(db)
	if err := httpServer.Start(); err != nil {
		logger.Error(fmt.Sprintf("Server startup failed: %v", err))
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(db)
	if err := grpcServer.Start(); err != nil {
		logger.Error(fmt.Sprintf("grpc server startup failed: %v", err))
		os.Exit(1)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	signal := <-ch
	logger.Info(fmt.Sprintf("WHATSAPP ROUTER STOPING, signal %v", signal))

}

const tokenUpdateInterval = 12

func initAuthToken(db *mongo.Database) {
	configRepo := repositories.NewConfigRepository(db)
	configService := services.NewConfigService(configRepo)
	conf, err := configService.GetConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting config with whatsapp auth token: %s", err))
		os.Exit(1)
	}
	if conf == nil {
		whatsappService := services.NewWhatsappService()
		res, err := whatsappService.Login()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		var login services.LoginWhatsapp
		if err := json.NewDecoder(res.Body).Decode(&login); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		conf, err = configService.CreateOrUpdate(&models.Config{Token: login.Users[0].Token})
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}
	config.UpdateAuthToken(conf.Token)

	whatsappService := services.NewWhatsappService()

	s := gocron.NewScheduler(time.UTC)
	s.Every(tokenUpdateInterval).
		Hour().
		StartAt(time.Now().Add(time.Hour * tokenUpdateInterval)).
		Do(func() {
			res, err := whatsappService.Login()
			if err != nil {
				logger.Error(err.Error())
				return
			}

			var login services.LoginWhatsapp
			bdBytes, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				logger.Error(err.Error())
				return
			}
			if res.StatusCode != 200 {
				logger.Error(fmt.Sprintf("Couldn't update token: %s, %s", res.Status, string(bdBytes)))
				return
			}

			if err := json.Unmarshal(bdBytes, &login); err != nil {
				logger.Error(err.Error())
				return
			}
			newToken := login.Users[0].Token

			configService.CreateOrUpdate(&models.Config{Token: newToken})

			config.UpdateAuthToken(newToken)
			logger.Info("Whatsapp token updated")
		})

	s.StartAsync()
}

package config

import (
	"fmt"
	"log"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/weni/whatsapp-router/logger"
)

type Config struct {
	Server   ServerConfig
	DB       DbConfig
	Whatsapp WhatsappConfig
}

type ServerConfig struct {
	HttpPort int32 `env:"SERVER_HTTP_PORT,required"`
	GRPCPort int32 `env:"SERVER_GRPC_PORT,required"`
}

type DbConfig struct {
	Host     string `env:"DB_HOST,required"`
	Port     int32  `env:"DB_PORT,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Name     string `env:"DB_NAME,required"`
	AppName  string `env:"DB_APP_NAME,required"`
}

type WhatsappConfig struct {
	BaseURL   string `env:"WPP_BASEURL,required"`
	Username  string `env:"WPP_USERNAME,required"`
	Password  string `env:"WPP_PASSWORD,required"`
	AuthToken string `env:"WPP_AUTHTOKEN,required"`
}

var AppConf *Config

func GetConfig() *Config {
	if AppConf == nil {
		fmt.Println("Config is nil")
		AppConf = &Config{}
		if err := godotenv.Load("./config/.env"); err != nil {
			logger.Error("Error loading .env file")
		}

		if err := envdecode.StrictDecode(AppConf); err != nil {
			logger.Error(fmt.Sprintf("Failed to decode and load environment variables: %v", err))
			log.Fatal()
		}
	}
	return AppConf
}

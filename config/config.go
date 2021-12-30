package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joeshaw/envdecode"
)

type Config struct {
	Server   ServerConfig
	DB       DbConfig
	Whatsapp WhatsappConfig
}

type ServerConfig struct {
	HttpPort       int32  `env:"SERVER_HTTP_PORT,default=9000"`
	GRPCPort       int32  `env:"SERVER_GRPC_PORT,default=7000"`
	CourierBaseURL string `env:"SERVER_COURIER_BASE_URL,required"`
	SentryDSN      string `env:"SERVER_SENTRY_DSN"`
	LogLevel       string `env:"SERVER_LOG_LEVEL,default=debug"`
}

type DbConfig struct {
	Name string `env:"DB_NAME,required"`
	URI  string `env:"DB_URI,required"`
}

type WhatsappConfig struct {
	BaseURL  string `env:"WPP_BASEURL,required"`
	Username string `env:"WPP_USERNAME,required"`
	Password string `env:"WPP_PASSWORD,required"`
}

var appConf *Config

var authToken string

func GetConfig() *Config {
	if appConf == nil {
		log.Println("loading config")
		appConf = &Config{}
		if err := envdecode.Decode(appConf); err != nil {
			log.Println(fmt.Sprintf("Failed to decode and load environment variables: %v", err.Error()))
			os.Exit(1)
		}
	}
	return appConf
}

func GetAuthToken() string {
	return authToken
}

func UpdateAuthToken(token string) {
	authToken = token
}

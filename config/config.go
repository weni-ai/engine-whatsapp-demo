package config

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joeshaw/envdecode"
	"github.com/weni/whatsapp-router/utils"
)

type Config struct {
	App      App
	DB       DB
	Whatsapp Whatsapp
}

type App struct {
	HttpPort       int32  `env:"APP_HTTP_PORT,default=9000"`
	GRPCPort       int32  `env:"APP_GRPC_PORT,default=7000"`
	CourierBaseURL string `env:"APP_COURIER_BASE_URL,default=http://localhost:8000/c/wa"`
	SentryDSN      string `env:"APP_SENTRY_DSN"`
	LogLevel       string `env:"APP_LOG_LEVEL,default=debug"`
}

type DB struct {
	Name string `env:"DB_NAME,default=whatsapp-router"`
	URI  string `env:"DB_URI,default=mongodb://admin:admin@localhost:27017"`
}

type Whatsapp struct {
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
	if authToken == "" {
		req, err := http.NewRequest("GET", "v1/users/login", nil)
		if err != nil {
			log.Println(err.Error())
			return ""
		}
		client := utils.GetHTTPClient()
		_, err = client.Do(req)
		if err != nil {
			log.Println(err.Error())
			return ""
		}
	}
	return authToken
}

func UpdateAuthToken(token string) {
	authToken = token
}

package http

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/metric"
	"github.com/weni/whatsapp-router/repositories"
	"github.com/weni/whatsapp-router/servers/http/handlers"
	"github.com/weni/whatsapp-router/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	config     config.Config
	db         *mongo.Database
	httpServer *http.Server
	metrics    *metric.Service
}

func NewServer(db *mongo.Database, metrics *metric.Service) *Server {
	conf := config.GetConfig()
	return &Server{
		db:      db,
		config:  *conf,
		metrics: metrics,
	}
}

func (s *Server) Start() error {
	sRouter := NewRouter(s)
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.App.HttpPort),
		Handler:      sRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting http server :%v", s.config.App.HttpPort))
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	return nil
}

func NewRouter(s *Server) *chi.Mux {
	router := chi.NewRouter()

	contactRepoDb := repositories.NewContactRepositoryDb(s.db)
	channelRepoDb := repositories.NewChannelRepositoryDb(s.db)
	configRepoDb := repositories.NewConfigRepository(s.db)
	whatsappHandler := handlers.WhatsappHandler{
		ContactService:  services.NewContactService(contactRepoDb),
		ChannelService:  services.NewChannelService(channelRepoDb, s.metrics),
		CourierService:  services.NewCourierService(),
		WhatsappService: services.NewWhatsappService(),
		ConfigService:   services.NewConfigService(configRepoDb),
		Metrics:         s.metrics,
	}
	courierHandler := handlers.CourierHandler{
		WhatsappService: services.NewWhatsappService(),
	}

	router.Use(logger.MiddlewareLogger)

	router.Route("/wr/", func(r chi.Router) {
		r.Use(ContentTypeJson)
		r.Route("/receive", func(r chi.Router) {
			r.Post("/", whatsappHandler.HandleIncomingRequests)
		})
	})

	router.Route("/v1", func(r chi.Router) {
		r.Post("/messages", courierHandler.HandleSendMessage)
		r.Post("/users/login", whatsappHandler.RefreshToken)
		r.Get("/health", whatsappHandler.HandleHealth)
		r.Get("/media/{mediaID}", whatsappHandler.HandleGetMedia)
		r.Post("/media", whatsappHandler.HandlePostMedia)
		r.Patch("/settings/application", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Get("/metrics", promhttp.Handler().ServeHTTP)

	return router
}

func ContentTypeJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf8")
		next.ServeHTTP(w, r)
	})
}

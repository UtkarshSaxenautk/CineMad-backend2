package transport

import (
	"authentication-ms/pkg/config"
	"authentication-ms/pkg/svc"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type Http struct {
	HttpServer      *http.Server
	webServerConfig config.WebServer
	db              *mongo.Database
	svc             svc.SVC
}

func NewHttp(webServerConfig config.WebServer, db *mongo.Database, svc svc.SVC) *Http {
	return &Http{&http.Server{}, webServerConfig, db, svc}
}

func (h *Http) Initialize() {
	h.HttpServer.Handler = NewRouter(h.webServerConfig.RoutePrefix, h.svc).Initialize()
	h.HttpServer.Addr = ":" + h.webServerConfig.Port
	log.Println("http initialized")
}

func (h *Http) Run() {
	log.Println("HTTP server starting at #{h.HttpServer.Addr}")
	if err := h.HttpServer.ListenAndServe(); err != nil {
		log.Print("Error: #{err}")
	}
}

func (h *Http) Shutdown() {
	h.HttpServer.Close()
}

package transport

import (
	"authentication-ms/pkg/config"
	cache2 "authentication-ms/pkg/repo/cache"
	"authentication-ms/pkg/repo/mongo"
	sdk2 "authentication-ms/pkg/repo/sdk"
	svc2 "authentication-ms/pkg/svc"
	"context"
	"log"
)

type Server struct {
	worker    Booter
	http      Booter
	appConfig config.App
}

func NewServer(appConfig config.App) (*Server, error) {

	db, err := mongo.New(appConfig.Mongo.Config)
	if err != nil {
		log.Fatal("mongo connection failed : ", err)
		return &Server{}, err
	}
	log.Println("mongo connection established")

	dao := mongo.NewDal(db.Database)
	cache := cache2.NewCache(context.Background())
	sdk := sdk2.New()
	svc := svc2.New(dao, cache, sdk)
	worker, err := NewWorker(svc)
	if err != nil {
		log.Fatal("worker connection failed")
		return &Server{}, err
	}

	return &Server{
		worker:    worker,
		http:      NewHttp(appConfig.WebServer, db.Database, svc),
		appConfig: appConfig,
	}, nil
}

func (s *Server) Initialize() {
	s.worker.Initialize()
	s.http.Initialize()
}

func (s *Server) Run() {
	s.http.Run()
	go s.worker.Run()
}

func (s *Server) Shutdown() {
	s.worker.Shutdown()
}

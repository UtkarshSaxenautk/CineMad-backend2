package main

import (
	"authentication-ms/pkg/config"
	"authentication-ms/pkg/transport"
	"log"
)

func main() {
	appConfig, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}
	s, err := transport.NewServer(*appConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Shutdown()
	s.Initialize()
	s.Run()
}

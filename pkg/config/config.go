package config

import (
	"authentication-ms/pkg/repo/mongo"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
	"path/filepath"
	"runtime"
)

const (
	envPrefix = "MOVIE"
)

type Mongo struct {
	mongo.Config
}

type WebServer struct {
	Port        string `required:"true" split_words:"true"`
	RoutePrefix string `required:"true" split_words:"true"`
}

type App struct {
	Service   string `required:"true" split_words:"true"`
	Env       string `required:"true" split_words:"true"`
	WebServer WebServer
	Mongo     Mongo
}

func FromEnv() (*App, error) {
	fromFileToEnv()
	cfg := &App{}
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func fromFileToEnv() {
	cfgFileName := os.Getenv("ENV_FILE")
	if cfgFileName != "" {
		if err := godotenv.Load(cfgFileName); err != nil {
			fmt.Println("error: failure reading ENV_FILE file, ", err)
		} else {
			return
		}
	}

	_, b, _, _ := runtime.Caller(0)
	cfgFileName = filepath.Join(filepath.Dir(b), "../../etc/config.local.env")

	if err := godotenv.Load(cfgFileName); err != nil {
		fmt.Println("error: failure reading config file:, ", err)
	}
}

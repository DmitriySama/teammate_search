package main

import (
	"context"
	"log"

	"github.com/DmitriySama/teammate_search/config"
	"github.com/DmitriySama/teammate_search/internal/bootstrap"
)

func main() {
	cfg, _ := config.LoadConfig()

	log.Println(cfg != nil)
	producer := bootstrap.InitProducers(cfg)

	storage:= bootstrap.InitPGStorage(cfg, producer)
	cache := bootstrap.InitCache(cfg)
	service := bootstrap.InitTSService(storage, cache)
	api := bootstrap.InitRegistryAPI(service, cfg.ServiceName, storage)
	bootstrap.AppRun(context.Background(), cfg, api)
}

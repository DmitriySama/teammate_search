package main

import (
	"context"

	"github.com/DmitriySama/teammate_search/config"
	"github.com/DmitriySama/teammate_search/internal/bootstrap"
)

func main() {
	cfg, _ := config.LoadConfig()
	producer := bootstrap.InitProducers(cfg)

	storage:= bootstrap.InitPGStorage(cfg, producer)
	cache := bootstrap.InitCache(cfg)
	service := bootstrap.InitTSService(storage, cache)
	api := bootstrap.InitRegistryAPI(service, cfg.ServiceName, storage)

	bootstrap.AppRun(context.Background(), cfg, api)
}

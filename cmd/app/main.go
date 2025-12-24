package main

import (
	"context"

	"github.com/DmitriySama/teammate_search/config"
	"github.com/DmitriySama/teammate_search/internal/bootstrap"
	//"github.com/DmitriySama/teammate_search/internal/cache"
)

func main() {
	cfg, _ := config.LoadConfig()

	storage:= bootstrap.InitPGStorage(cfg)
	cache := bootstrap.InitCache(cfg)
	service := bootstrap.InitTSService(storage, cache)
	api := bootstrap.InitRegistryAPI(service, cfg.ServiceName, storage)
	// producers := bootstrap.InitProducers(cfg)
	// consumers := bootstrap.InitConsumers(cfg, service, producers)

	bootstrap.AppRun(context.Background(), cfg, api)
}

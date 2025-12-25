package bootstrap

import (
	"github.com/DmitriySama/teammate_search/config"
	"github.com/DmitriySama/teammate_search/internal/producer"
)

func InitProducers(cfg *config.Config) *producer.Manager {
	return producer.NewManager(cfg)
}


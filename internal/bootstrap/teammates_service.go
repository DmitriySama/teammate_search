package bootstrap

import (
	"github.com/DmitriySama/teammate_search/internal/cache"
	tsService "github.com/DmitriySama/teammate_search/internal/services/teammateSearchService"
	"github.com/DmitriySama/teammate_search/internal/storage/pgstorage"
)

func InitTSService(storage *pgstorage.PGstorage, cache *cache.Cache) *tsService.Service {
	return tsService.New(storage, cache)
}

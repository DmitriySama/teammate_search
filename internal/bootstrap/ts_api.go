package bootstrap

import (
	"github.com/DmitriySama/teammate_search/internal/api/ts_service_api"
	tsService "github.com/DmitriySama/teammate_search/internal/services/teammateSearchService"
	"github.com/DmitriySama/teammate_search/internal/storage/pgstorage"
)
func InitRegistryAPI(service *tsService.Service, serviceName string, pg *pgstorage.PGstorage) *ts_service_api.API {
	return ts_service_api.New(service, serviceName, pg)
}
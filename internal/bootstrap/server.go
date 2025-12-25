package bootstrap

import (
	"log"
	"net/http"
	"context"
	
	"github.com/DmitriySama/teammate_search/internal/api/ts_service_api"
	"github.com/DmitriySama/teammate_search/config"
)
//func AppRun(ctx context.Context, cfg *config.Config, api *ts_service_api.API, webService *http_logic.WebService) error   {
func AppRun(ctx context.Context, cfg *config.Config, api *ts_service_api.API) error   {
    r := api.Router()

    log.Println("Сервер запущен на http://localhost:3000")
    
    log.Fatal(http.ListenAndServe(":3000", r))
    return nil
}
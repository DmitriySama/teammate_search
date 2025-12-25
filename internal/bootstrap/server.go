package bootstrap

import (
	"context"
	"log"
	"net/http"

	"github.com/DmitriySama/teammate_search/config"
	"github.com/DmitriySama/teammate_search/internal/api/ts_service_api"
)

func AppRun(ctx context.Context, cfg *config.Config, api *ts_service_api.API) error   {
    r := api.Router()

    log.Println("Сервер запущен на http://localhost:3000")
    
    log.Fatal(http.ListenAndServe(":3000", r))
    return nil
}
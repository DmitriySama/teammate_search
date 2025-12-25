package bootstrap

import (
	"fmt"
	"log"

	"github.com/DmitriySama/teammate_search/config"
	"github.com/DmitriySama/teammate_search/internal/producer"
	"github.com/DmitriySama/teammate_search/internal/storage/pgstorage"
)

func InitPGStorage(cfg *config.Config, producer *producer.Manager) (*pgstorage.PGstorage) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.Username,
        cfg.Database.Password,
        cfg.Database.DBName)
		
	storage, err := pgstorage.InitDB(dsn, producer)
	if err != nil {
		log.Panic(fmt.Sprintf("ошибка инициализации БД, %v", err))
		panic(err)
	}
	return storage
}
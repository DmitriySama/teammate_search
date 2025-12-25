package bootstrap

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/DmitriySama/teammate_search/config"
	"github.com/DmitriySama/teammate_search/internal/cache"
)

func InitCache(cfg *config.Config) *cache.Cache {
	redisAddr := cfg.RedisAddr()
	log.Printf("Redis: инициализация подключения к Redis БД: %d", cfg.Redis.DB)
	
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Redis: ошибка подключения к Redis: %v", err)
		return nil
	}
	log.Printf("Redis: успешно подключено к Redis по адресу %s", redisAddr)

	c := cache.NewCache(client, cfg.Redis.TTL)
	return c
}

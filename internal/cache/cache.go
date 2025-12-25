package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"github.com/DmitriySama/teammate_search/internal/models"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

type LanguagesCache struct {
	Languages []models.Language `json:"languages"`
}

type GamesCache struct {
	Games     []models.Games `json:"games"`
}

type GenresCache struct {
	Genres    []models.Genres `json:"genres"`
}

type AppsCache struct {
	Apps    []models.Apps `json:"apps"`
}


func NewCache(client *redis.Client, ttlSeconds int) *Cache {
	return &Cache{client: client, ttl:time.Duration(ttlSeconds) * time.Second,}
}		

func (c *Cache) Key(prefix, id string) string {
	return fmt.Sprintf("%s:%s", prefix, id)
}


func (c *Cache) GetLanguages(ctx context.Context) ([]models.Language, bool) {
	if c == nil || c.client == nil {
		return nil, false
	}

	cacheKey := c.Key("languages", "all")
	log.Printf("Redis: попытка получить языки из кэша для ключа %s", cacheKey)

	data, err := c.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Redis: языки не найдены в кэше для ключа %s", cacheKey)
		} else {
			log.Printf("Redis: ошибка получения языков из кэша для ключа %s: %v", cacheKey, err)
		}
		return nil, false
	}

	var cached LanguagesCache
	if err := json.Unmarshal(data, &cached); err != nil {
		log.Printf("Redis: ошибка десериализации языков из кэша для ключа %s: %v", cacheKey, err)
		return nil, false
	}

	log.Printf("Redis: успешно получены языки из кэша для ключа %s", cacheKey)
	return cached.Languages, true
}


func (c *Cache) SetLanguages(ctx context.Context, languages []models.Language) error {
	if c == nil || c.client == nil {
		log.Printf("Redis: кэш не инициализирован, пропуск сохранения языков")
		return nil
	}

	cacheKey := c.Key("languages", "all")
	log.Printf("Redis: попытка сохранить языки в кэше для ключа %s (TTL: %v)", cacheKey, c.ttl)

	value, err := json.Marshal(LanguagesCache{
		Languages: languages,
	})
	if err != nil {
		log.Printf("Redis: ошибка сериализации языков для ключа %s: %v", cacheKey, err)
		return err
	}

	if err := c.client.Set(ctx, cacheKey, value, c.ttl).Err(); err != nil {
		log.Printf("Redis: ошибка сохранения языков в кэше для ключа %s: %v", cacheKey, err)
		return err
	} else {
		log.Printf("Redis: успешно сохранены языки в кэше для ключа %s", cacheKey)
	}
	return nil
}


func (c *Cache) GetGames(ctx context.Context) ([]models.Games, bool) {
	if c == nil || c.client == nil {
		return nil, false
	}

	cacheKey := c.Key("games", "all")
	log.Printf("Redis: попытка получить игры из кэша для ключа %s", cacheKey)

	data, err := c.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Redis: игры не найдены в кэше для ключа %s", cacheKey)
		} else {
			log.Printf("Redis: ошибка получения игр из кэша для ключа %s: %v", cacheKey, err)
		}
		return nil, false
	}

	var cached GamesCache
	if err := json.Unmarshal(data, &cached); err != nil {
		log.Printf("Redis: ошибка десериализации игр из кэша для ключа %s: %v", cacheKey, err)
		return nil, false
	}

	log.Printf("Redis: успешно получены игры из кэша для ключа %s", cacheKey)
	return cached.Games, true
}


func (c *Cache) SetGames(ctx context.Context, games []models.Games)  error {
	if c == nil || c.client == nil {
		log.Printf("Redis: кэш не инициализирован, пропуск сохранения игр")
		return nil
	}

	cacheKey := c.Key("games", "all")
	log.Printf("Redis: попытка сохранить игры в кэше для ключа %s (TTL: %v)", cacheKey, c.ttl)

	value, err := json.Marshal(GamesCache{
		Games:     games,
	})
	if err != nil {
		log.Printf("Redis: ошибка сериализации игр для ключа %s: %v", cacheKey, err)
		return err
	}

	if err := c.client.Set(ctx, cacheKey, value, c.ttl).Err(); err != nil {
		log.Printf("Redis: ошибка сохранения игр в кэше для ключа %s: %v", cacheKey, err)
	} else {
		log.Printf("Redis: успешно сохранены игры в кэше для ключа %s", cacheKey)
	}
	return nil
}


func (c *Cache) GetGenres(ctx context.Context) ([]models.Genres, bool) {
	if c == nil || c.client == nil {
		return nil, false
	}

	cacheKey := c.Key("genres", "all")
	log.Printf("Redis: попытка получить жанры игр из кэша для ключа %s", cacheKey)

	data, err := c.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Redis: жанры игр не найдены в кэше для ключа %s", cacheKey)
		} else {
			log.Printf("Redis: ошибка получения жанров игр из кэша для ключа %s: %v", cacheKey, err)
		}
		return nil, false
	}

	var cached GenresCache
	if err := json.Unmarshal(data, &cached); err != nil {
		log.Printf("Redis: ошибка десериализации жанров игр из кэша для ключа %s: %v", cacheKey, err)
		return nil, false
	}

	log.Printf("Redis: успешно получены жанры игр из кэша для ключа %s", cacheKey)
	return cached.Genres, true
}


func (c *Cache) SetGenres(ctx context.Context, genres []models.Genres) error {
	if c == nil || c.client == nil {
		log.Printf("Redis: кэш не инициализирован, пропуск сохранения жанров игр")
		return nil
	}

	cacheKey := c.Key("genres", "all")
	log.Printf("Redis: попытка сохранить жанры игр в кэше для ключа %s (TTL: %v)", cacheKey, c.ttl)

	value, err := json.Marshal(GenresCache{
		Genres:    genres,
	})
	if err != nil {
		log.Printf("Redis: ошибка сериализации жанров игр для ключа %s: %v", cacheKey, err)
		return err
	}

	if err := c.client.Set(ctx, cacheKey, value, c.ttl).Err(); err != nil {
		log.Printf("Redis: ошибка сохранения жанров игр в кэше для ключа %s: %v", cacheKey, err)
	} else {
		log.Printf("Redis: успешно сохранены жанры игр в кэше для ключа %s", cacheKey)
	}
	return nil
}

func (c *Cache) GetApps(ctx context.Context) ([]models.Apps, bool) {
	if c == nil || c.client == nil {
		return nil, false
	}

	cacheKey := c.Key("apps", "all")
	log.Printf("Redis: попытка получить apps из кэша для ключа %s", cacheKey)

	data, err := c.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Redis: apps не найдены в кэше для ключа %s", cacheKey)
		} else {
			log.Printf("Redis: ошибка получения apps из кэша для ключа %s: %v", cacheKey, err)
		}
		return nil, false
	}

	var cached AppsCache
	if err := json.Unmarshal(data, &cached); err != nil {
		log.Printf("Redis: ошибка десериализации apps из кэша для ключа %s: %v", cacheKey, err)
		return nil, false
	}

	log.Printf("Redis: успешно получены apps из кэша для ключа %s", cacheKey)
	return cached.Apps, true
}


func (c *Cache) SetApps(ctx context.Context, apps []models.Apps) error{
	if c == nil || c.client == nil {
		log.Printf("Redis: кэш не инициализирован, пропуск сохранения apps")
		return nil
	}

	cacheKey := c.Key("apps", "all")
	log.Printf("Redis: попытка сохранить apps в кэше для ключа %s (TTL: %v)", cacheKey, c.ttl)

	value, err := json.Marshal(AppsCache{
		Apps:    apps,
	})
	if err != nil {
		log.Printf("Redis: ошибка сериализации apps для ключа %s: %v", cacheKey, err)
		return err
	}

	if err := c.client.Set(ctx, cacheKey, value, c.ttl).Err(); err != nil {
		log.Printf("Redis: ошибка сохранения apps в кэше для ключа %s: %v", cacheKey, err)
	} else {
		log.Printf("Redis: успешно сохранены apps в кэше для ключа %s", cacheKey)
	}
	return nil
}
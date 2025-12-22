package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"github.com/redis/go-redis/v9"
	"github.com/DmitriySama/teammate_search/internal/models"
)

type UsersCache struct {
	client *redis.Client
	ttl    time.Duration
}


func NewUsersCache(client *redis.Client, ttlSeconds int) *UsersCache {
	return &UsersCache{client: client, ttl: time.Duration(ttlSeconds) * time.Second}
}


func (c *UsersCache) key(orderID, groupID string) string {
	return fmt.Sprintf("UsersCaches:%s:%s", orderID, groupID)
}

func (c *UsersCache) Get(ctx context.Context, orderID, groupID string) ([]models.User, string, bool) {
	if c == nil || c.client == nil {
		return nil, "", false
	}

	cacheKey := c.key(orderID, groupID)
	log.Printf("Redis: попытка получить данные из кэша для ключа %s", cacheKey)
	data, err := c.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Redis: данные не найдены в кэше для ключа %s", cacheKey)
		} else {
			log.Printf("Redis: ошибка получения данных из кэша для ключа %s: %v", cacheKey, err)
		}
		return nil, "", false
	}

	var cached cachedUsersCaches
	if err := json.Unmarshal(data, &cached); err != nil {
		log.Printf("Redis: ошибка десериализации данных из кэша для ключа %s: %v", cacheKey, err)
		return nil, "", false
	}

	log.Printf("Redis: успешно получены данные из кэша для ключа %s (приказ %s, группа %s)", cacheKey, orderID, groupID)
	return cached.UsersCaches, cached.GroupName, true
}

func (c *UsersCache) Set(ctx context.Context, orderID, groupID, groupName string, UsersCaches []models.User) {
	if c == nil || c.client == nil {
		log.Printf("Redis: кэш не инициализирован, пропуск сохранения данных")
		return
	}

	cacheKey := c.key(orderID, groupID)
	log.Printf("Redis: попытка сохранить данные в кэше для ключа %s (приказ %s, группа %s, TTL: %v)", cacheKey, orderID, groupID, c.ttl)
	value, err := json.Marshal(cachedUsersCaches{GroupName: groupName, UsersCaches: UsersCaches})
	if err != nil {
		log.Printf("Redis: ошибка сериализации данных для ключа %s: %v", cacheKey, err)
		return
	}

	if err := c.client.Set(ctx, cacheKey, value, c.ttl).Err(); err != nil {
		log.Printf("Redis: ошибка сохранения данных в кэше для ключа %s: %v", cacheKey, err)
	} else {
		log.Printf("Redis: успешно сохранены данные в кэше для ключа %s", cacheKey)
	}
}

func (c *UsersCache) InvalidateGroup(ctx context.Context, orderID, groupID string) {
	if c == nil || c.client == nil {
		log.Printf("Redis: кэш не инициализирован, пропуск инвалидации группы")
		return
	}
	cacheKey := c.key(orderID, groupID)
	log.Printf("Redis: инвалидация кэша для группы (ключ: %s, приказ: %s, группа: %s)", cacheKey, orderID, groupID)
	if err := c.client.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Redis: ошибка инвалидации кэша для ключа %s: %v", cacheKey, err)
	} else {
		log.Printf("Redis: успешно инвалидирован кэш для ключа %s", cacheKey)
	}
}

func (c *UsersCache) InvalidateOrder(ctx context.Context, orderID string) {
	if c == nil || c.client == nil {
		log.Printf("Redis: кэш не инициализирован, пропуск инвалидации приказа")
		return
	}

	pattern := fmt.Sprintf("UsersCaches:%s:*", orderID)
	log.Printf("Redis: инвалидация кэша для приказа %s (паттерн: %s)", orderID, pattern)
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	deletedCount := 0
	for iter.Next(ctx) {
		key := iter.Val()
		if err := c.client.Del(ctx, key).Err(); err != nil {
			log.Printf("Redis: ошибка удаления ключа %s: %v", key, err)
		} else {
			deletedCount++
		}
	}
	if err := iter.Err(); err != nil {
		log.Printf("Redis: ошибка при сканировании ключей для приказа %s: %v", orderID, err)
	} else {
		log.Printf("Redis: успешно инвалидирован кэш для приказа %s (удалено ключей: %d)", orderID, deletedCount)
	}
}

// Обновляет статус определенной задачи в кэше без инвалидации всего кэша
func (c *UsersCache) UpdateUsersCacheInCache(ctx context.Context, orderID, groupID string, UsersCacheID int64, status string) {
	if c == nil || c.client == nil {
		log.Printf("Redis: кэш не инициализирован, пропуск обновления задачи в кэше")
		return
	}

	cacheKey := c.key(orderID, groupID)
	log.Printf("Redis: попытка обновить задачу %d в кэше для ключа %s", UsersCacheID, cacheKey)

	data, err := c.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Redis: кэш для ключа %s не найден, пропуск обновления задачи", cacheKey)
		} else {
			log.Printf("Redis: ошибка получения данных из кэша для ключа %s: %v", cacheKey, err)
		}
		return
	}

	var cached cachedUsersCaches
	if err := json.Unmarshal(data, &cached); err != nil {
		log.Printf("Redis: ошибка десериализации данных из кэша для ключа %s: %v", cacheKey, err)
		return
	}

	found := false
	for i := range cached.UsersCaches {
		if cached.UsersCaches[i].ID == UsersCacheID {
			cached.UsersCaches[i].Status = status
			found = true
			break
		}
	}

	if !found {
		log.Printf("Redis: задача %d не найдена в кэше для ключа %s", UsersCacheID, cacheKey)
		return
	}

	value, err := json.Marshal(cached)
	if err != nil {
		log.Printf("Redis: ошибка сериализации данных для ключа %s: %v", cacheKey, err)
		return
	}

	if err := c.client.Set(ctx, cacheKey, value, c.ttl).Err(); err != nil {
		log.Printf("Redis: ошибка обновления данных в кэше для ключа %s: %v", cacheKey, err)
	} else {
		log.Printf("Redis: успешно обновлена задача %d в кэше для ключа %s", UsersCacheID, cacheKey)
	}
}

type cachedUsersCaches struct {
	GroupName string        `json:"group_name"`
	UsersCaches     []models.User `json:"UsersCaches"`
}

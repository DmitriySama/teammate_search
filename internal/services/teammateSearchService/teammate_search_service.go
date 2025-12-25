package teammateSearchService

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/DmitriySama/teammate_search/internal/models"
	"github.com/DmitriySama/teammate_search/internal/storage/pgstorage"
)


type UsersStorage interface {
	Register(username, password, description string, age int) (*pgstorage.AuthResult, error)
	UserExists(username string) (bool, error)
	Login(username, password string) (*pgstorage.AuthResult, error)
	UpdateUser(r *http.Request, user models.User) (error)
	FindUser(username, password string) (int, error)
	
	GetLanguages(ctx context.Context) ([]models.Language, error)
	GetGenres(ctx context.Context) ([]models.Genres, error)
	GetGames(ctx context.Context) ([]models.Games, error)
	GetUserByID(userID int) (*models.User, error)
	GetUserCount() (int, error)
	GetUsers(r *http.Request) (*sql.Rows, error)
	GetApps(ctx context.Context) ([]models.Apps, error)
}

type UsersCache interface {
	Key(prefix, id string) string
	
	GetGames(ctx context.Context) ([]models.Games, bool)
	SetGames(ctx context.Context, games []models.Games) error
	
	GetGenres(ctx context.Context) ([]models.Genres, bool)
	SetGenres(ctx context.Context, genres []models.Genres) error
	
	GetLanguages(ctx context.Context) ([]models.Language, bool)
	SetLanguages(ctx context.Context, languages []models.Language) error
	
	GetApps(ctx context.Context) ([]models.Apps, bool)
	SetApps(ctx context.Context, apps []models.Apps) error
}

type Service struct {
	storage UsersStorage
	cache   UsersCache
}

func New(storage UsersStorage, cache UsersCache) *Service {
	return &Service{storage: storage, cache: cache}
}


func (s *Service) GetLanguages(ctx context.Context) ([]models.Language, error) {
	cachedLanguages, ok := s.cache.GetLanguages(ctx)
	if ok {
		log.Println("Набор языков взят из REDIS")
		return cachedLanguages, nil
	}

	languages, err := s.storage.GetLanguages(ctx)
	if err != nil {
		return []models.Language{}, err
	}

	err = s.cache.SetLanguages(ctx, languages)
	if err == nil {
		log.Println("Набор языков установлен в REDIS")
	} else {
		log.Println("Ошибка при установке языков в REDIS:", err)
	}
	return languages, nil
}

func (s *Service) GetGenres(ctx context.Context) ([]models.Genres, error) {
	cachedGenres, ok := s.cache.GetGenres(ctx)
	if ok {
		log.Println("Набор жанров взят из REDIS")
		return cachedGenres, nil
	}

	genres, err := s.storage.GetGenres(ctx)
	if err != nil {
		return []models.Genres{}, err
	}

	err = s.cache.SetGenres(ctx, genres)
	if err == nil {
		log.Println("Жанры добавлены в REDIS")
	} else {
		log.Println("Ошибка при добавлении жанров в REDIS:", err)
	}
	return genres, nil

}

func (s *Service) GetGames(ctx context.Context) ([]models.Games, error) {
	cachedGames, ok := s.cache.GetGames(ctx)
	if ok {
		log.Println("Набор игр взят из REDIS")
		return cachedGames, nil
	}

	games, err := s.storage.GetGames(ctx)
	if err != nil {
		return []models.Games{}, err
	}

	err = s.cache.SetGames(ctx, games)
	if err == nil {
		log.Println("Набор игр добавлен в REDIS")
	} else {
		log.Println("Ошибка при добавлении набора игр в REDIS:", err)
	}
	return games, nil

}

func (s *Service) GetApps(ctx context.Context) ([]models.Apps, error) {
	cachedApps, ok := s.cache.GetApps(ctx)
	if ok {
		log.Println("Набор приложений взят из REDIS")
		return cachedApps, nil
	}

	apps, err := s.storage.GetApps(ctx)
	if err != nil {
		return []models.Apps{}, err
	}

	err = s.cache.SetApps(ctx, apps)
	if err == nil {
		log.Println("Набор приложений добавлен в REDIS")
	} else {
		log.Println("Ошибка при добавлении приложений в REDIS:", err)
	}
	return apps, nil

}
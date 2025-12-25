package teammateSearchService

import (
	"context"
	"net/http"
	"database/sql"

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
}

type UsersCache interface {
	Key(prefix, id string) string
	GetLanguages(ctx context.Context) ([]models.Language, bool)
	GetGenres(ctx context.Context) ([]models.Genres, bool)
	GetGames(ctx context.Context) ([]models.Games, bool)
	SetGames(ctx context.Context, games []models.Games)
	SetGenres(ctx context.Context, genres []models.Genres)
	SetLanguages(ctx context.Context, languages []models.Language)
}

type Service struct {
	storage UsersStorage
	cache   UsersCache
}

func New(storage UsersStorage, cache UsersCache) *Service {
	return &Service{storage: storage, cache: cache}
}

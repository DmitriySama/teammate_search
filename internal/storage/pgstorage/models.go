package pgstorage

import (
	"database/sql"
	"github.com/DmitriySama/teammate_search/internal/models"
	"log"
)


// AuthService содержит бизнес-логику авторизации
type AuthService struct {
    DB *sql.DB
}

// AuthResult содержит результат операции авторизации
type AuthResult struct {
    User    *models.User  `json:"user"`
    Success bool   `json:"success"`
    Message string `json:"message"`
}

func NewAuthService() *AuthService {
    db, err := InitDB()
    if err != nil {
        log.Printf("Ошибка получения дб: %v", err)
        return nil
    }   
    if db == nil {
        log.Fatal("Database connection is nil") // или вернуть ошибку
    }
    return &AuthService{DB: db}
}
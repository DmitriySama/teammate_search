package pgstorage

import (
	"github.com/DmitriySama/teammate_search/internal/models"
)

// AuthResult содержит результат операции авторизации
type AuthResult struct {
    User    *models.User  `json:"user"`
    Success bool   `json:"success"`
    Message string `json:"message"`
}

package models

import (
    "time"
)

type User struct {
    ID                int       `json:"id"`
    Username          string    `json:"username"`
    PasswordHash      string    `json:"-"`
    MainGame          string    `json:"main_game,omitempty"`
    Age               int       `json:"age,omitempty"`
    Country           string    `json:"country,omitempty"`
    Language          string    `json:"language,omitempty"`
    Communication     string    `json:"communication,omitempty"` 
    FavoriteGenres    string    `json:"favorite_genres,omitempty"`
    FavoriteGames     string    `json:"favorite_games,omitempty"`  
    Description       string    `json:"description,omitempty"`
    CreatedAt         time.Time `json:"created_at"`
}

type RegisterRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type RegisterResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    UserID  int    `json:"user_id,omitempty"`
    Token   string `json:"token,omitempty"`
}
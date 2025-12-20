package main

import (
    "time"
)

type User struct {
    ID                int       `json:"id"`
    Username          string    `json:"username"`
    Email             string    `json:"email"`
    PasswordHash      string    `json:"-"`
    MainGame          string    `json:"main_game,omitempty"`
    Rank              string    `json:"rank,omitempty"`
    Age               int       `json:"age,omitempty"`
    Country           string    `json:"country,omitempty"`
    Language          string    `json:"language,omitempty"`
    Communication     string    `json:"communication,omitempty"` 
    DiscordTag        string    `json:"discord_tag,omitempty"`
    FavoriteGenres    string    `json:"favorite_genres,omitempty"`
    FavoriteGames     string    `json:"favorite_games,omitempty"`  
    Playstyle         string    `json:"playstyle,omitempty"`
    Description       string    `json:"description,omitempty"`
    IsProfileComplete bool      `json:"is_profile_complete"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    MainGame string `json:"main_game,omitempty"`
    Rank     string `json:"rank,omitempty"`
}

type RegisterResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    UserID  int    `json:"user_id,omitempty"`
    Token   string `json:"token,omitempty"`
}
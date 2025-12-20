package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "regexp"
    "strings"
    "time"
    "github.com/teammate_search/models"
)


// In-memory хранилище (временное решение)
var users = make(map[string]User)
var sessions = make(map[string]Session)

type Session struct {
    UserID int
    Expiry time.Time
}

var db *sql.DB

// Инициализация БД (PostgreSQL)
func initDB() {
    var err error
    connStr := "user=postgres dbname=teamfind password=yourpassword host=localhost sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    
    // Создание таблиц
    createTables()
}

func createTables() {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        main_game VARCHAR(50),
        rank VARCHAR(50),
        age INTEGER,
        country VARCHAR(100),
        language VARCHAR(50),
        communication JSONB,
        discord_tag VARCHAR(100),
        favorite_genres JSONB,
        favorite_games JSONB,
        playstyle VARCHAR(50),
        description TEXT,
        is_profile_complete BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE TABLE IF NOT EXISTS sessions (
        token VARCHAR(255) PRIMARY KEY,
        user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
        expiry TIMESTAMP NOT NULL
    );
    `
    
    _, err := db.Exec(query)
    if err != nil {
        log.Fatal("Ошибка создания таблиц:", err)
    }
}

// registerHandler обрабатывает регистрацию пользователя
func registerHandler(w http.ResponseWriter, r *http.Request) {
    // Только POST запросы
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Парсим JSON тело запроса
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Неверный формат данных", http.StatusBadRequest)
        return
    }
    
    // Валидация данных
    if err := validateRegistration(req); err != nil {
        sendError(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Проверка существования пользователя
    if exists, err := userExists(req.Username, req.Email); err != nil {
        sendError(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
        log.Println("Ошибка проверки пользователя:", err)
        return
    } else if exists {
        sendError(w, "Пользователь с таким именем или email уже существует", http.StatusConflict)
        return
    }
    
    // Хеширование пароля
    passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        sendError(w, "Ошибка обработки пароля", http.StatusInternalServerError)
        log.Println("Ошибка хеширования пароля:", err)
        return
    }
    
    // Создание пользователя в БД
    userID, err := createUserInDB(req, string(passwordHash))
    if err != nil {
        sendError(w, "Ошибка создания пользователя", http.StatusInternalServerError)
        log.Println("Ошибка создания пользователя:", err)
        return
    }
    
    // Создание сессии (JWT токен)
    token, err := createSession(userID)
    if err != nil {
        sendError(w, "Ошибка создания сессии", http.StatusInternalServerError)
        log.Println("Ошибка создания сессии:", err)
        return
    }
    
    // Устанавливаем куки
    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    token,
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true,
        Path:     "/",
        SameSite: http.SameSiteStrictMode,
    })
    
    // Отправляем успешный ответ
    sendJSON(w, http.StatusCreated, RegisterResponse{
        Success: true,
        Message: "Регистрация успешна",
        UserID:  userID,
        Token:   token,
    })
    
    log.Printf("Новый пользователь зарегистрирован: %s (ID: %d)", req.Username, userID)
}

// Валидация данных регистрации
func validateRegistration(req RegisterRequest) error {
    // Проверка имени пользователя
    if len(req.Username) < 3 || len(req.Username) > 20 {
        return fmt.Errorf("имя пользователя должно быть от 3 до 20 символов")
    }
    
    usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    if !usernameRegex.MatchString(req.Username) {
        return fmt.Errorf("имя пользователя может содержать только буквы, цифры и подчеркивание")
    }
    
    // Проверка email
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(req.Email) {
        return fmt.Errorf("неверный формат email")
    }
    
    // Проверка пароля
    if len(req.Password) < 6 {
        return fmt.Errorf("пароль должен содержать не менее 6 символов")
    }
    
    // Дополнительные проверки пароля (опционально)
    if !strings.ContainsAny(req.Password, "0123456789") {
        // return fmt.Errorf("пароль должен содержать хотя бы одну цифру")
    }
    
    return nil
}

// Проверка существования пользователя
func userExists(username, email string) (bool, error) {
    // Если используем БД
    if db != nil {
        var count int
        err := db.QueryRow(
            "SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2",
            username, email,
        ).Scan(&count)
        
        if err != nil {
            return false, err
        }
        return count > 0, nil
    }
    
    // In-memory проверка
    for _, user := range users {
        if user.Username == username || user.Email == email {
            return true, nil
        }
    }
    return false, nil
}

// Создание пользователя в БД
func createUserInDB(req RegisterRequest, passwordHash string) (int, error) {
    if db != nil {
        var userID int
        err := db.QueryRow(`
            INSERT INTO users (username, email, password_hash, main_game, rank, is_profile_complete)
            VALUES ($1, $2, $3, $4, $5, $6)
            RETURNING id
        `, req.Username, req.Email, passwordHash, req.MainGame, req.Rank, false).Scan(&userID)
        
        return userID, err
    }
    
    // In-memory создание
    userID := len(users) + 1
    users[req.Email] = User{
        ID:                userID,
        Username:          req.Username,
        Email:             req.Email,
        PasswordHash:      passwordHash,
        MainGame:          req.MainGame,
        Rank:              req.Rank,
        IsProfileComplete: false,
        CreatedAt:         time.Now(),
        UpdatedAt:         time.Now(),
    }
    
    return userID, nil
}

// Создание сессии
func createSession(userID int) (string, error) {
    // Генерация токена (упрощенная версия)
    token := generateToken()
    expiry := time.Now().Add(24 * time.Hour)
    
    if db != nil {
        _, err := db.Exec(
            "INSERT INTO sessions (token, user_id, expiry) VALUES ($1, $2, $3)",
            token, userID, expiry,
        )
        if err != nil {
            return "", err
        }
    } else {
        // In-memory сессия
        sessions[token] = Session{
            UserID: userID,
            Expiry: expiry,
        }
    }
    
    return token, nil
}

// Генерация токена (упрощенная)
func generateToken() string {
    return fmt.Sprintf("%d_%d", time.Now().UnixNano(), len(sessions))
}

// Вспомогательные функции для ответов
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, message string, status int) {
    sendJSON(w, status, map[string]interface{}{
        "success": false,
        "message": message,
    })
}

// Auth middleware
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Получаем токен из куки
        cookie, err := r.Cookie("session_token")
        if err != nil {
            // Пробуем получить из заголовка
            token := r.Header.Get("Authorization")
            if token == "" {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }
            
            // Проверяем токен
            if !isValidToken(token) {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }
            
            // Устанавливаем токен в контекст
            // context.WithValue(r.Context(), "user_id", userID)
        } else {
            // Проверяем куки
            if !isValidToken(cookie.Value) {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }
        }
        
        next(w, r)
    }
}

func isValidToken(token string) bool {
    if db != nil {
        var count int
        err := db.QueryRow(
            "SELECT COUNT(*) FROM sessions WHERE token = $1 AND expiry > NOW()",
            token,
        ).Scan(&count)
        
        return err == nil && count > 0
    }
    
    // In-memory проверка
    session, exists := sessions[token]
    if !exists {
        return false
    }
    
    return session.Expiry.After(time.Now())
}
package pgstorage

import (
	"log"
	"strings"
	"strconv"
	"time"
	"database/sql"
	"net/http"
	"fmt"
	"github.com/DmitriySama/teammate_search/internal/models"
)


// Register регистрирует нового пользователя
func (aus *AuthService) Register(username, password, description string, age int) (*AuthResult, error) {
  
    // Проверка уникальности username
    exists, err := aus.userExists(username)
    if err != nil {
        log.Printf("Ошибка проверки пользователя: %v", err)
        return nil, err
    } 
    if exists {
        return &AuthResult{
            Success: false,
            Message: "Пользователь с таким именем или email уже существует",
        }, nil
    }
    
    // Создание пользователя
    user := &models.User{
        Username:     username,
        Password:     password,
        Description:     description,
        Age:     age,
        CreatedAt:    time.Now(),
    }
    
    // Сохранение в БД
    var userID int
    err = aus.DB.QueryRow(`
        INSERT INTO users (username, password, description, age, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`, user.Username, user.Password, user.Description, user.Age, user.CreatedAt).Scan(&userID)
    
    if err != nil {
        log.Printf("Ошибка сохранения пользователя в БД: %v", err)
        return nil, err
    }
    
    user.ID = userID
    
    log.Printf("Пользователь зарегистрирован: ID=%d, Username=%s", user.ID, user.Username)
    
    return &AuthResult{
        User:    user,
        Success: true,
        Message: "Регистрация успешно завершена",
    }, nil
}


// userExists проверяет существование пользователя
func (aus *AuthService) userExists(username string) (bool, error) {
    
	var count int
    err := aus.DB.QueryRow(`
        SELECT COUNT(*) 
        FROM users 
        WHERE username = $1
    `, username).Scan(&count)
    
    return count > 0, err
}

// Login выполняет вход пользователя
func (aus *AuthService) Login(username, password string) (*AuthResult, error) {    
    
    user_id, err := aus.findUser(username, password)
    if err != nil {
        if err == sql.ErrNoRows {
            return &AuthResult{
                Success: false,
                Message: "Неверное имя пользователя или пароль",
            }, nil
        }
        log.Printf("Ошибка поиска пользователя: %v", err)
        return nil, err
    }
    
    user, err := aus.GetUserByID(user_id)
    if err != nil {
        if err == sql.ErrNoRows {
            return &AuthResult{
                Success: false,
                Message: "",
            }, nil
        }
        log.Printf("Ошибка получения пользователя: %v", err)
        return nil, err
    }

    log.Printf("Пользователь вошел: ID=%d, Username=%s", user_id, user.Username)
    
    return &AuthResult{
        User:    user,
        Success: true,
        Message: "Вход выполнен успешно",
    }, nil
}

func (aus *AuthService) UpdateUser(r *http.Request, user models.User) (error) {

    // Динамически строим запрос
    var setClauses []string
    var args []interface{}
    argIndex := 1
    
    age, _ := strconv.Atoi(r.FormValue("age"))
    if age != user.Age {
        setClauses = append(setClauses, fmt.Sprintf("age = $%d", argIndex))
        args = append(args, age)
        argIndex++
    }
    if r.FormValue("description") != user.Description {
        setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIndex))
        args = append(args, r.FormValue("description"))
        argIndex++
    }
    MostLikeGame := strings.SplitN(r.FormValue("game"), " ", 2)[0]
    if MostLikeGame != user.MostLikeGame {
        setClauses = append(setClauses, fmt.Sprintf("most_like_game = $%d", argIndex))
        args = append(args, MostLikeGame)
        argIndex++
    }
    MostLikeGenre := strings.SplitN(r.FormValue("genre"), " ", 2)[0]
    if MostLikeGenre != user.MostLikeGenre {
        setClauses = append(setClauses, fmt.Sprintf("most_like_genre = $%d", argIndex))
        args = append(args, MostLikeGenre)
        argIndex++
    }
    if r.FormValue("app") != user.App {
        setClauses = append(setClauses, fmt.Sprintf("speaking_app = $%d", argIndex))
        args = append(args, r.FormValue("app"))
        argIndex++
    }
    Language := strings.SplitN(r.FormValue("language"), " ", 2)[0]
    if Language != user.Language {
        setClauses = append(setClauses, fmt.Sprintf("language = $%d", argIndex))
        args = append(args, Language)
        argIndex++
    }
    
    // Если ничего не изменилось
    log.Println(len(setClauses))
    log.Println(setClauses)
    log.Println(args)
    if len(setClauses) == 0 {
        log.Println("Ничего не изменилось")
        return nil
    }

    // Добавляем WHERE
    args = append(args, user.ID)
    query := fmt.Sprintf(
        "UPDATE users SET %s WHERE id = $%d",
        strings.Join(setClauses, ", "),
        argIndex,
    )

    _, err := aus.DB.Exec(query, args...)
    return err
}
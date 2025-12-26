package pgstorage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DmitriySama/teammate_search/internal/models"
)

// Register регистрирует нового пользователя
func (pg *PGstorage) Register(username, password, description string, age int) (*AuthResult, error) {
  
    // Проверка уникальности username
    exists, err := pg.UserExists(username)
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
    err = pg.DB.QueryRow(`
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

// UserExists проверяет существование пользователя
func (pg *PGstorage) UserExists(username string) (bool, error) {
    
	var count int
    err := pg.DB.QueryRow(`
        SELECT COUNT(*) 
        FROM users 
        WHERE username = $1
    `, username).Scan(&count)
    
    return count > 0, err
}


func (pg *PGstorage) SelectUser(username string) {
    pg.producer.SendUserPopularityData(context.Background(), username)
}   

func (pg *PGstorage) FilterData(fd models.FilterData) error {
    err := pg.producer.SendFilterData(context.Background(), fd)
    return err
}   


func (pg *PGstorage) Login(username, password string) (*AuthResult, error) {    
    
    user_id, err := pg.FindUser(username, password)
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
    
    user, err := pg.GetUserByID(user_id)
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

func (pg *PGstorage) UpdateUser(r *http.Request, user models.User) (error) {

    // Динамически строим запрос
    var setClpges []string
    var args []interface{}
    argIndex := 1
    
    age, _ := strconv.Atoi(r.FormValue("age"))
    if age != user.Age {
        setClpges = append(setClpges, fmt.Sprintf("age = $%d", argIndex))
        args = append(args, age)
        argIndex++
        log.Println("AGE:", age, user.Age)
    }
    desc := strings.TrimSpace(r.FormValue("description"))
    if desc != user.Description {
        setClpges = append(setClpges, fmt.Sprintf("description = $%d", argIndex))
        args = append(args, desc)
        argIndex++
        log.Println("description:", r.FormValue("description"), user.Description)
    }
    MostLikeGame := strings.SplitN(r.FormValue("game"), " ", 2)[1]
    if MostLikeGame != user.MostLikeGame {
        id_game := strings.SplitN(r.FormValue("game"), " ", 2)[0]
        setClpges = append(setClpges, fmt.Sprintf("most_like_game = $%d", argIndex))
        args = append(args, id_game)
        argIndex++
        log.Println("MostLikeGame:", MostLikeGame, user.MostLikeGame)
    }
    MostLikeGenre := strings.SplitN(r.FormValue("genre"), " ", 2)[1]
    if MostLikeGenre != user.MostLikeGenre {
        id_genre := strings.SplitN(r.FormValue("genre"), " ", 2)[0]
        setClpges = append(setClpges, fmt.Sprintf("most_like_genre = $%d", argIndex))
        args = append(args, id_genre)
        argIndex++
        log.Println("MostLikeGenre:", MostLikeGenre, user.MostLikeGenre)
    }
    App := strings.SplitN(r.FormValue("app"), " ", 2)[1]
    if App != user.App {
        id_app := strings.SplitN(r.FormValue("app"), " ", 2)[0]
        setClpges = append(setClpges, fmt.Sprintf("speaking_app = $%d", argIndex))
        args = append(args, id_app)
        argIndex++
        log.Println("App:", r.FormValue("app"), user.App)
    }
    Language := strings.SplitN(r.FormValue("language"), " ", 2)[1]
    if Language != user.Language {
        lang_id := strings.SplitN(r.FormValue("language"), " ", 2)[0]
        setClpges = append(setClpges, fmt.Sprintf("language = $%d", argIndex))
        args = append(args, lang_id)
        argIndex++
        log.Println("Language:", Language, user.Language)
    }
    
    // Если ничего не изменилось
    log.Println(len(setClpges))
    log.Println(setClpges)
    log.Println(args)
    if len(setClpges) == 0 {
        log.Println("Ничего не изменилось")
        return nil
    }

    // Добавляем WHERE
    args = append(args, user.ID)
    query := fmt.Sprintf(
        "UPDATE users SET %s WHERE id = $%d",
        strings.Join(setClpges, ", "),
        argIndex,
    )

    _, err := pg.DB.Exec(query, args...)
    return err
}
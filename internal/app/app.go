package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"strconv"

	auth "github.com/DmitriySama/teammate_search/internal/storage/pgstorage"
	"github.com/DmitriySama/teammate_search/internal/models"
)


var user = &models.User{}
var aus = auth.NewAuthService()

func main() {
    http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Println(filepath.Ext(r.URL.Path))
        switch filepath.Ext(r.URL.Path) {
        case ".css":
            w.Header().Set("Content-Type", "text/css; charset=utf-8")
        case ".js":
            w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
        case ".json":
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
        case ".png":
            w.Header().Set("Content-Type", "image/png")
        case ".jpg", ".jpeg":
            w.Header().Set("Content-Type", "image/jpeg")
        case ".ico":
            w.Header().Set("Content-Type", "image/x-icon")
        default:
            w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        }
    })

    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/main/home", mainMainHandler)
    http.HandleFunc("/main/search", mainSearchHandler)
    http.HandleFunc("/profile/look", handleGetProfile)
    http.HandleFunc("/profile/update", handleUpdateProfile)


    log.Println("Сервер запущен на http://localhost:3000")
    log.Fatal(http.ListenAndServe(":3000", nil))
}



func registerHandler(w http.ResponseWriter, r *http.Request) {
   
    if r.Method == "POST" {
        if err := r.ParseForm(); err != nil {
            log.Println("Ошибка при разборе формы")
        } else {
            age, _ := strconv.Atoi(r.FormValue("age"))

            result, err := aus.Register(r.FormValue("username"), r.FormValue("password"), r.FormValue("description"), age)
            if err != nil {
                log.Fatal("Ошибка регистрации:", err)
            }
            if result.Success {
                log.Printf("Пользователь зарегистрирован: ID=%d", result.User.ID)
                user = result.User
                http.Redirect(w, r, "/main/main", http.StatusSeeOther)
            } else {
                log.Printf("Ошибка регистрации: %s", result.Message)
            }
        }
    }
    if r.Method == "GET" {
        http.ServeFile(w, r, "../frontend/register.html")
    }
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
   
    if r.Method == "POST" {
        if err := r.ParseForm(); err != nil {
            log.Println("Ошибка при разборе формы")
        } else {
            
            aus := auth.NewAuthService()
            
            result, err := aus.Login(r.FormValue("username"), r.FormValue("password"))
            if err != nil {
                log.Fatal("Ошибка входа:", err)
            }
            if result.Success {
                log.Printf("Пользователь авторизован: ID=%d", result.User.ID)
                user = result.User
                http.Redirect(w, r, "/main/home", http.StatusSeeOther)
            } else {
                log.Printf("Ошибка авторизации: %s", result.Message)
            }
        }
    } 
    if r.Method == "GET" {
        http.ServeFile(w, r, "../frontend/login.html")
    }
}

func mainMainHandler(w http.ResponseWriter, r *http.Request) {    
    EmptyUserCheck(w, r)

    if r.Method == "GET" {
        profileData := getDataToShow(r, "main")
        tmpl, _ := template.ParseFiles("../frontend/main.html")
        tmpl.Execute(w, profileData)
    }
}

func mainSearchHandler(w http.ResponseWriter, r *http.Request) {    
    EmptyUserCheck(w, r)
    profileData := getDataToShow(r, "search")
        if r.Method == "GET" {
            tmpl, _ := template.ParseFiles("../frontend/main_search.html")
            profileData["User"] = []models.User{}
            log.Println(profileData)
            tmpl.Execute(w, profileData)
        }
    if r.Method == "POST" {
        if err := r.ParseForm(); err != nil {
            log.Println("Ошибка при разборе формы")
        } else {
            users, _ := aus.GetUsers(r)
            users_v2 := []models.UserListShow{}
            for users.Next() {
                var username, description, game, genre, language, app string
                var age int

                if err := users.Scan(&username, &age, &description, &game, &language, &app, &genre); err != nil {
                    log.Fatal(err)
                }
                users_v2 = append(users_v2, models.UserListShow{
                    Username: username,
                    Age: age,
                    Description: description,
                    MostLikeGame: game,
                    MostLikeGenre: genre,
                    Language: language,
                })
            }
            log.Println("users_v2", users_v2)
            data := map[string]interface{}{
                "MyUsername": user.Username,
                "User": users_v2,
            }
            tmpl, _ := template.ParseFiles("../frontend/main_search.html")
            tmpl.Execute(w, data)
        }
    }
}

func EmptyUserCheck(w http.ResponseWriter, r *http.Request) {
    if user.Username == "" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
    }
}

func handleGetProfile(w http.ResponseWriter, r *http.Request) {
    // Загружаем шаблон профиля
    tmpl, err := template.ParseFiles("../frontend/profile_look.html")
    if err != nil {
        log.Printf("Ошибка загрузки шаблона профиля: %v", err)
        http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
        return
    }

    profileData := getDataToShow(r, "GetProfile")

    // Рендерим шаблон
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if err := tmpl.Execute(w, profileData); err != nil {
        log.Printf("Ошибка рендеринга шаблона профиля: %v", err)
    }
}

func getDataToShow(r *http.Request, choise string) (map[string]interface{}){
    var data map[string]interface{}
    switch choise {
        case "main": {
            userCount, _ := aus.GetUserCount()
            data = map[string]interface{}{
                "Username": user.Username,
                "UserCount": userCount,
            }
        }   
        case "search": {
            languages, _ := aus.GetLanguages(r.Context())
            languages = append([]models.Language{{ID: -1, Lang: "Любой"}}, languages...)
            games, _ := aus.GetGames(r.Context())
            games = append([]models.Games{{ID: -1, Game: "Любая"}}, games...)
            genres, _ := aus.GetGenres(r.Context())
            genres = append([]models.Genres{{ID: -1, Genre: "Любой"}}, genres...)

            data = map[string]interface{}{
                "MyUsername": user.Username,
                "Languages": languages,
                "Games": games,
                "Genres": genres,
            }
        }
        case "GetProfile": {
            data = map[string]interface{}{
                "Username": user.Username,
                "Age": user.Age,
                "Description": user.Description,
                "SpeakingApp": user.App,
                "MLGame": user.MostLikeGame,
                "MLGenre": user.MostLikeGenre,
                "Language": user.Language,
            }
        } 
        case "UpdateProfile": {
            languages, _ := aus.GetLanguages(r.Context())
            games, _ := aus.GetGames(r.Context())
            genres, _ := aus.GetGenres(r.Context())
            data = map[string]interface{}{
                "Username": user.Username,
                "Age": user.Age,
                "Description": user.Description,
                "SpeakingApp": user.App,
                "Languages": languages,
                "Games": games,
                "Genres": genres,
            }
        }
    }
    return data
}


func handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        err := aus.UpdateUser(r, *user)
        if err == nil {
            log.Printf("Профиль пользователя %d обновлен", user.ID)
            UpdateUserValues(r)
        } else {
            log.Fatal("Не удалось обновить")
        }
        http.Redirect(w, r, "/profile/look", http.StatusSeeOther)
    }
    if r.Method == "GET" {
        tmpl, _ := template.ParseFiles("../frontend/profile_update.html")
        profileData := getDataToShow(r, "UpdateProfile")
    
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := tmpl.Execute(w, profileData); err != nil {
            log.Printf("Ошибка рендеринга шаблона профиля: %v", err)
        }
    }
}

func UpdateUserValues(r *http.Request) {
    user.Age, _ = strconv.Atoi(r.FormValue("age"))
    user.Description = r.FormValue("description")
    user.Language = r.FormValue("language")
    user.MostLikeGame = strings.SplitN(r.FormValue("game"), " ", 2)[1]
    user.MostLikeGenre = strings.SplitN(r.FormValue("genre"), " ", 2)[1]
    user.Language = strings.SplitN(r.FormValue("language"), " ", 2)[1]
    user.App = r.FormValue("app")
}
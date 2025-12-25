package ts_service_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"os"

	"html/template"
	"log"
	"strconv"
	"strings"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"github.com/DmitriySama/teammate_search/api/swagger"
	tsService "github.com/DmitriySama/teammate_search/internal/services/teammateSearchService"
	
	"github.com/DmitriySama/teammate_search/internal/models"
	"github.com/DmitriySama/teammate_search/internal/storage/pgstorage"
)

type API struct {
	service     *tsService.Service
	serviceName string
	once        sync.Once
	swaggerSpec []byte
    pg *pgstorage.PGstorage
    user *models.User
}

func New(service *tsService.Service, serviceName string, pg *pgstorage.PGstorage) *API {
	return &API{service: service, serviceName: serviceName, pg: pg}
}

func (a *API) Router() http.Handler {
	router := chi.NewRouter()
	http.DefaultServeMux.HandleFunc("/", a.MIMEProcessing)

	router.Get("/health", a.health)
	router.Get("/swagger", a.swaggerUI)
	router.Get("/swagger/web.swagger.json", a.swaggerSpecHandler)
	
	router.Get("/register", a.RegisterPage)
	router.Post("/register", a.RegisterHandler)

	router.Get("/login", a.LoginPage)
	router.Post("/login", a.LoginHandler)

	router.Get("/main/home", a.MainMainHandler)

	router.Get("/profile/look", a.HandleGetProfile)
	router.Get("/profile/update", a.HandleUpdateProfile)
	router.Post("/profile/update", a.HandleUpdateProfile)
	
	router.Get("/main/search", a.MainSearchHandler)
	router.Post("/main/search", a.MainSearchHandler)
	router.Post("/main/select-user", a.SelectUser)
	return router
}

type SelectUser struct {
    Username    string `json:"username"`
}

func (a *API) SelectUser(w http.ResponseWriter, r *http.Request) {    
    
    var req SelectUser
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
    }
    a.pg.SelectUser(req.Username)
}


func (a *API) RegisterPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/register.html")
}

func (a *API) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("Ошибка при разборе формы")
	} else {
		age, _ := strconv.Atoi(r.FormValue("age"))

		result, err := a.pg.Register(r.FormValue("username"), r.FormValue("password"), r.FormValue("description"), age)
		if err != nil {
			log.Fatal("Ошибка регистрации:", err)
		}
		if result.Success {
			log.Printf("Пользователь зарегистрирован: ID=%d", result.User.ID)
			a	.user = result.User
			http.Redirect(w, r, "/main/main", http.StatusSeeOther)
		} else {
			log.Printf("Ошибка регистрации: %s", result.Message)
		}
	}
}

func (a *API) LoginPage(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles(getFrontendPath()+"/login.html")).Execute(w, nil)
}

func getFrontendPath() string {
    // Директория исполняемого файла (X/cmd/app/)
    exePath, err := os.Executable()
    if err != nil {
        panic(err)
    }

    exeDir := filepath.Dir(exePath)
    
    // Поднимаемся на 2 уровня вверх: app/ -> cmd/ -> X/
    rootDir := filepath.Dir(filepath.Dir(exeDir))
    
    // Идём в internal/frontend/
    return filepath.Join(rootDir, "internal", "frontend")
}

func (a *API) LoginHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
		log.Println("Ошибка при разборе формы")
	} else {
		
		result, err := a.pg.Login(r.FormValue("username"), r.FormValue("password"))
		if err != nil {
			log.Fatal("Ошибка входа:", err)
		}
		if result.Success {
			log.Printf("Пользователь авторизован: ID=%d", result.User.ID)
			a.user = result.User
			http.Redirect(w, r, "/main/home", http.StatusSeeOther)
		} else {
			log.Printf("Ошибка авторизации: %s", result.Message)
		}
	}
}


func (a *API) MainMainHandler(w http.ResponseWriter, r *http.Request) {    
    a.EmptyUserCheck(w, r)

	profileData := a.GetDataToShow(r, "main")
	template.Must(template.ParseFiles(getFrontendPath()+"/main.html")).Execute(w, profileData)
}

func (a *API) MainSearchHandler(w http.ResponseWriter, r *http.Request) {    
    a.EmptyUserCheck(w, r)
    profileData := a.GetDataToShow(r, "search")
	
	if r.Method == "GET" {
		profileData["User"] = []models.User{}
		template.Must(template.ParseFiles(getFrontendPath()+"/main_search.html")).Execute(w, profileData)
	}

    if r.Method == "POST" {
        if err := r.ParseForm(); err != nil {
            log.Println("Ошибка при разборе формы")
        } else {
            // Отправление данных фильтров через KAFKA
            age, _ := strconv.Atoi(r.FormValue("age"))
            fd := models.FilterData{
                Age: age,
                Game: r.FormValue("game"),
                Genre: r.FormValue("genre"),
                Language: r.FormValue("language"),
                App: r.FormValue("app"),
            }
            a.pg.FilterData(fd)

            // Получение пользователей
            users, _ := a.pg.GetUsers(r)
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

            data := map[string]interface{}{
                "MyUsername": a.user.Username,
                "User": users_v2,
            }

            // Отрисовка пользователей
			template.Must(template.ParseFiles(getFrontendPath()+"/main_search.html")).Execute(w, data)
        }
    }
}

func (a *API) EmptyUserCheck(w http.ResponseWriter, r *http.Request) {
    if a.user.Username == "" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
    }
}

func (a *API) HandleGetProfile(w http.ResponseWriter, r *http.Request) {	
	profileData := a.GetDataToShow(r, "GetProfile")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	template.Must(template.ParseFiles(getFrontendPath()+"/profile_look.html")).Execute(w, profileData)
}

func (a *API) GetDataToShow(r *http.Request, choise string) (map[string]interface{}){
    var data map[string]interface{}
    switch choise {
        case "main": {
            userCount, _ := a.pg.GetUserCount()
            data = map[string]interface{}{
                "Username": a.user.Username,
                "UserCount": userCount,
            }
        }   
        case "search": {
            languages, _ := a.service.GetLanguages(r.Context())
            languages = append([]models.Language{{ID: -1, Lang: "Любой"}}, languages...)
            games, _ := a.service.GetGames(r.Context())
            games = append([]models.Games{{ID: -1, Game: "Любая"}}, games...)
            genres, _ := a.service.GetGenres(r.Context())
            genres = append([]models.Genres{{ID: -1, Genre: "Любой"}}, genres...)
            apps, _ := a.service.GetApps(r.Context())
            apps = append([]models.Apps{{ID: -1, App: "Любой"}}, apps...)

            data = map[string]interface{}{
                "MyUsername": a.user.Username,
                "Languages": languages,
                "Games": games,
                "Genres": genres,
                "Apps": apps,
            }
        }
        case "GetProfile": {
            data = map[string]interface{}{
                "Username": a.user.Username,
                "Age": a.user.Age,
                "Description": a.user.Description,
                "SpeakingApp": a.user.App,
                "MLGame": a.user.MostLikeGame,
                "MLGenre": a.user.MostLikeGenre,
                "Language": a.user.Language,
                "App": a.user.App,
            }
        } 
        case "UpdateProfile": {
            languages, _ := a.pg.GetLanguages(r.Context())
            games, _ := a.pg.GetGames(r.Context())
            genres, _ := a.pg.GetGenres(r.Context())
            apps, _ := a.pg.GetApps(r.Context())
            data = map[string]interface{}{
                "Username": a.user.Username,
                "Age": a.user.Age,
                "Description": a.user.Description,
                "SpeakingApp": a.user.App,
                "Languages": languages,
                "Games": games,
                "Apps": apps,
                "Genres": genres,
            }
        }
    }
    return data
}


func (a *API) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        err := a.pg.UpdateUser(r, *a.user)
        if err == nil {
            log.Printf("Профиль пользователя %d обновлен", a.user.ID)
            a.UpdateUserValues(r)
        } else {
            log.Fatal("Не удалось обновить")
        }
        http.Redirect(w, r, "/profile/look", http.StatusSeeOther)
    }
    if r.Method == "GET" {
        profileData := a.GetDataToShow(r, "UpdateProfile")
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
		template.Must(template.ParseFiles(getFrontendPath()+"/profile_update.html")).Execute(w, profileData)
    }
}

func (a *API) UpdateUserValues(r *http.Request) {
    a.user.Age, _ = strconv.Atoi(r.FormValue("age"))
    a.user.Description = r.FormValue("description")
    a.user.Language = r.FormValue("language")
    a.user.MostLikeGame = strings.SplitN(r.FormValue("game"), " ", 2)[1]
    a.user.MostLikeGenre = strings.SplitN(r.FormValue("genre"), " ", 2)[1]
    a.user.Language = strings.SplitN(r.FormValue("language"), " ", 2)[1]
    a.user.App = r.FormValue("app")
}

func (a *API) MIMEProcessing(w http.ResponseWriter, r *http.Request) {
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
}

func (a *API) health(w http.ResponseWriter, _ *http.Request) {
	body := map[string]string{
		"service": a.serviceName,
		"status":  "ok",
	}
	writeJSON(w, http.StatusOK, body)
}

func (a *API) swaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
		<html>
		<head>
		<title>Task Registry API</title>
		<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
		</head>
		<body>
		<div id="swagger-ui"></div>
		<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
		<script>
			window.onload = () => {
			SwaggerUIBundle({ url: '/swagger/web.swagger.json', dom_id: '#swagger-ui' });
			};
		</script>
		</body>
		</html>`)
}

func (a *API) swaggerSpecHandler(w http.ResponseWriter, _ *http.Request) {
	a.once.Do(func() {
		a.swaggerSpec = swagger.Teammate_search()
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(a.swaggerSpec)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

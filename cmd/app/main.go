package main

import (
    "log"
    "net/http"
    "encoding/json"
)


func main() {
    // Маршруты
    http.HandleFunc("/", mainHandler)
    http.HandleFunc("/register", registerPageHandler)
    //http.HandleFunc("/api/register", registerHandler)
    http.HandleFunc("/api/profile/details", profileDetailsHandler)
    http.HandleFunc("/dashboard", dashboardHandler)
    
    log.Println("Сервер запущен на http://127.0.0.1:3000")
    log.Fatal(http.ListenAndServe(":3000", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./templates/register.html")
}

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Static request: %s", r.URL.Path)

    http.ServeFile(w, r, "./static/css"+r.URL.Path)
}


// Обработчик сохранения дополнительной информации
func profileDetailsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Декодируем JSON данные
    var profileData map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&profileData); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // В реальном приложении здесь сохранение в БД
    // userID := getCurrentUserID(r)
    // saveProfileDetails(userID, profileData)
    
    log.Printf("Получены данные профиля: %+v", profileData)
    
    // Отправляем успешный ответ
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "message": "Данные профиля сохранены",
    })
}

// Обработчик дашборда (главной после регистрации)
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    // Временная заглушка
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(`
        <!DOCTYPE html>
        <html>
        <head>
            <title>TeamFind - Dashboard</title>
            <style>
                body { font-family: Arial, sans-serif; padding: 40px; }
                .container { max-width: 800px; margin: 0 auto; }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>Добро пожаловать в TeamFind!</h1>
                <p>Ваш профиль успешно создан. Скоро здесь будет ваша главная страница.</p>
                <a href="/">На главную</a>
            </div>
        </body>
        </html>
    `))
}


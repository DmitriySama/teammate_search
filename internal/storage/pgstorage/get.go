
package pgstorage

import (
	"log"
	"time"
	"context"
	"errors"
	"net/http"
	"database/sql"
	"github.com/DmitriySama/teammate_search/internal/models"
)


func (pg *PGstorage) FindUser(username, password string) (int, error) {
    log.Println(username, password)
    var id int
    err := pg.DB.QueryRow(`
        SELECT id
        FROM users 
        WHERE username = $1 and password = $2
    `, username, password).Scan(&id)
    
    return id, err
}


func (pg *PGstorage) GetLanguages(ctx context.Context) ([]models.Language, error) {
    query := `SELECT id_language, language FROM languages`
    
    rows, err := pg.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var languages []models.Language
    for rows.Next() {
        var lang models.Language
        if err := rows.Scan(&lang.ID, &lang.Lang); err != nil {
            return nil, err
        }
        languages = append(languages, lang)
    }
    
    return languages, nil
}

func (pg *PGstorage) GetGenres(ctx context.Context) ([]models.Genres, error) {
    query := `SELECT id_genre, genre FROM genres`
    
    rows, err := pg.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var genres []models.Genres
    for rows.Next() {
        var lang models.Genres
        if err := rows.Scan(&lang.ID, &lang.Genre); err != nil {
            return nil, err
        }
        genres = append(genres, lang)
    }
    
    return genres, nil
}

func (pg *PGstorage) GetGames(ctx context.Context) ([]models.Games, error) {
    query := `SELECT id_game, game FROM games`
    
    rows, err := pg.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var games []models.Games
    for rows.Next() {
        var game models.Games
        if err := rows.Scan(&game.ID, &game.Game); err != nil {
            return nil, err
        }
        games = append(games, game)
    }
    
    return games, nil
}

func (pg *PGstorage) GetApps(ctx context.Context) ([]models.Apps, error) {
    query := `SELECT id_app, app FROM apps`
    
    rows, err := pg.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var apps []models.Apps
    for rows.Next() {
        var app models.Apps
        if err := rows.Scan(&app.ID, &app.App); err != nil {
            return nil, err
        }
        apps = append(apps, app)
    }
    
    return apps, nil
}

func (pg *PGstorage) GetUserByID(userID int) (*models.User, error) {
    var username, password, f_game, f_genre, app, description, lang string
    var id, age int 
    var created_at time.Time 

    err := pg.DB.QueryRow(`
        SELECT 
            u.id, 
            u.username, 
            u.password,
            u.age, 
            u.description, 
            u.created_at,
            
            COALESCE(g1.game, '') AS f_game,
            COALESCE(g.genre, '') AS f_genre,
            COALESCE(a.app, '') AS app,
            COALESCE(l.language, '') AS lang
        FROM users u
        LEFT JOIN genres g ON u.most_like_genre = g.id_genre
        LEFT JOIN languages l ON u.language = l.id_language
        LEFT JOIN apps a ON u.speaking_app = a.id_app
        LEFT JOIN games g1 ON u.most_like_game = g1.id_game
        WHERE u.id = $1;
    `, userID).Scan(&id, &username, &password, &age, &description, &created_at, &f_game, &f_genre, &app, &lang)
    
    user := &models.User{
        ID:          id,
        Username:    username,
        Password:    password,
        Age:         age,
        Description: description,
        MostLikeGenre: f_genre,
        MostLikeGame: f_game,
        App: app,
        Language: lang,
        CreatedAt:   created_at,
    }
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("пользователь не найден")
        }
        return nil, err
    }
    
    return user, nil
}

func (pg *PGstorage) GetUserCount() (int, error) {
    var value int
    err := pg.DB.QueryRow(`SELECT count(*) from users`).Scan(&value)
    if err != nil {
        return 0, err
    }

    return value, err
}

func (pg *PGstorage) GetUsers(r *http.Request) (*sql.Rows, error) {
    
    query := `SELECT 
            u.username, 
            u.age, 
            u.description, 
            
            COALESCE(g1.game, '') AS f_game,
            COALESCE(g.genre, '') AS f_genre,
            COALESCE(a.app, '') AS app,
            COALESCE(l.language, '') AS lang
        FROM users u
        LEFT JOIN genres g ON u.most_like_genre = g.id_genre
        LEFT JOIN languages l ON u.language = l.id_language
        LEFT JOIN apps a ON u.speaking_app = a.id_app
        LEFT JOIN games g1 ON u.most_like_game = g1.id_game
        WHERE u.age between $1 and $2 `
    if r.FormValue("genre") != "-1" {
        query += ` and u.most_like_genre = ` + r.FormValue("genre")
    }
    if r.FormValue("game") != "-1" {
        query += ` and u.most_like_game = ` + r.FormValue("game") 
    }
    if r.FormValue("language") != "-1" {
        query += ` and u.language = ` + r.FormValue("language")
    }
    log.Println(query)
    rows, err := pg.DB.Query(query, r.FormValue("age0"), r.FormValue("age1"))
    if err != nil {
        log.Fatal("Ошибка при выполнении запроса: ", err.Error())
        return nil, err
    }

    return rows, err
}
package pgstorage

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"
    "github.com/DmitriySama/teammate_search/config"
    _ "github.com/lib/pq"

)

// PGstorage содержит бизнес-логику авторизации
type PGstorage struct {
    DB *sql.DB
}


func InitDB(connString string) (*PGstorage, error) {
    log.Println("Инициализация подключения к БД...")
    
    // 1. Загрузка конфигурации
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Printf("❌ Ошибка загрузки конфигурации: %v", err)
        return nil, fmt.Errorf("не удалось загрузить конфигурацию: %w", err)
    }
    
    if cfg == nil {
        log.Println("❌ Конфигурация загружена, но cfg == nil")
        return nil, fmt.Errorf("конфигурация не загружена (nil)")
    }
    
    log.Printf("✅ Конфиг загружен. База данных: %s", cfg.Database.DBName)
 
    
    // 4. Подключение к базе данных
    log.Println("Открытие соединения с PostgreSQL...")
    db, err := sql.Open("postgres", connString)
    if err != nil {
        log.Printf("❌ Ошибка открытия соединения: %v", err)
        return nil, fmt.Errorf("не удалось открыть соединение с БД: %w", err)
    }
    
    // 5. Настройка пула соединений
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(2 * time.Minute)
    
    // 6. Проверка подключения
    log.Println("Проверка подключения к БД (Ping)...")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        db.Close()
        log.Printf("❌ Ошибка подключения к БД: %v", err)
        
        // Детализация ошибки
        log.Printf("Параметры подключения: host=%s, port=%d, user=%s, dbname=%s", 
            cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.DBName)
        
        return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
    }

    storage := &PGstorage{
        DB: db,
    }
    
    log.Printf("✅ Успешно подключено к БД: %s@%s:%d/%s", 
        cfg.Database.Username, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
    
    return storage, nil
}

func (db PGstorage) Close() error {
    log.Println("Закрытие подключения к БД")
    return db.DB.Close()
}
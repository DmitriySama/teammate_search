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

type DB struct {
    *sql.DB
}


func InitDB() (*sql.DB, error) {
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
    
    // 2. Получение строки подключения
    connStr := cfg.Database.DatabaseURL()
    log.Printf("Получена строка подключения: %s", connStr)
    
    // 3. Проверка и преобразование формата строки подключения
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.Username,
        cfg.Database.Password,
        cfg.Database.DBName)
    
    // 4. Подключение к базе данных
    log.Println("Открытие соединения с PostgreSQL...")
    db, err := sql.Open("postgres", dsn)
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
    
    log.Printf("✅ Успешно подключено к БД: %s@%s:%d/%s", 
        cfg.Database.Username, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
    
    return db, nil
}



func (db *DB) Close() error {
    log.Println("Закрытие подключения к БД")
    return db.DB.Close()
}
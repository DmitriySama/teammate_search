package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ServiceName string         `yaml:"serviceName"`
	Port        int            `yaml:"port"`
	Database    DatabaseConfig `yaml:"database"`
	Kafka       KafkaConfig    `yaml:"kafka"`
	Topics      TopicsConfig   `yaml:"topics"`
	Redis       RedisConfig    `yaml:"redis"`
}

type DatabaseConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	ConnURL  string `yaml:"connURL"`
}


func LoadConfig() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("не удалось открыть файл конфигурации: %w", err)
    }
    defer file.Close()
    
    // Читаем содержимое
    var config Config
    decoder := yaml.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return nil, fmt.Errorf("не удалось декодировать YAML: %w", err)
    }
    
    return &config, nil
}


func (d *DatabaseConfig) DatabaseURL() string {
    if d == nil {
        return ""
    }
    
    // Правильный формат для lib/pq
    return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s ssl=disable",
        d.Host, d.Port, d.Username, d.Password, d.DBName)
}

func (c Config) RedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}


type KafkaConfig struct {
	Brokers         []string `yaml:"brokers"`
	Host            string   `yaml:"host"`
	Port            int      `yaml:"port"`
	GroupID         string   `yaml:"groupId"`
	MaxMessageBytes int      `yaml:"maxMessageBytes"`
}

type TopicsConfig struct {
	FilterData        string `yaml:"filterData"`
	UserPopylarity        string `yaml:"userPopylarity"`
	UpdateUserData        string `yaml:"UpdateUserData"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Username string    `yaml:"username"`
	DB   int    `yaml:"db"`
	TTL  int    `yaml:"ttlSeconds"`
}
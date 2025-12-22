package producer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/DmitriySama/teammate_search/config"
	"github.com/DmitriySama/teammate_search/internal/models"
)

type Manager struct {
	writer *kafka.Writer
	topic  string
}

func New(cfg *config.Config) *Manager {
	writerCfg := kafka.WriterConfig{
		Brokers:    cfg.Kafka.Brokers,
		Topic:      cfg.Topics.FilterData,
		Balancer:   &kafka.LeastBytes{},
		BatchBytes: cfg.Kafka.MaxMessageBytes,
	}

	return &Manager{
		writer: kafka.NewWriter(writerCfg),
		topic:  cfg.Topics.FilterData,
	}
}

func (m *Manager) SendFilterDataResponse(ctx context.Context, response models.FilterData) error {
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Kafka: ошибка сериализации данных фильтра")
		return err
	}

	log.Printf("Kafka: отправка данных поиска людей в топик %s", m.topic)
	if err := m.writer.WriteMessages(ctx, kafka.Message{Value: data}); err != nil {
		log.Printf("Kafka: ошибка данных поиска людей в топик %s", m.topic)
		return err
	}

	log.Printf("Kafka: успешно отправлены данные поиска людей")
	return nil
}

// func (m *Manager) SendNewUserDataResponse(ctx context.Context, response models.User) error {
// 	data, err := json.Marshal(response)
// 	if err != nil {
// 		log.Printf("Kafka: ошибка сериализации данных пользователя")
// 		return err
// 	}

// 	log.Printf("Kafka: отправка данных пользователя в топик %s", m.topic)
// 	if err := m.writer.WriteMessages(ctx, kafka.Message{Value: data}); err != nil {
// 		log.Printf("Kafka: ошибка данных поиска людей в топик %s", m.topic)
// 		return err
// 	}

// 	log.Printf("Kafka: успешно отправлены данные поиска людей")
// 	return nil
// }


// func (m *Manager) SendUpdateUserDataResponse(ctx context.Context, response models.UpdateUserData) error {
// 	data, err := json.Marshal(response)
// 	if err != nil {
// 		log.Printf("Kafka: ошибка сериализации данных пользователя")
// 		return err
// 	}

// 	log.Printf("Kafka: отправка данных пользователя в топик %s", m.topic)
// 	if err := m.writer.WriteMessages(ctx, kafka.Message{Value: data}); err != nil {
// 		log.Printf("Kafka: ошибка данных поиска людей в топик %s", m.topic)
// 		return err
// 	}

// 	log.Printf("Kafka: успешно отправлены данные поиска людей")
// 	return nil
// }


func (m *Manager) Close() error {
	if m.writer != nil {
		return m.writer.Close()
	}
	return nil
}

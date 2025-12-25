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
	filterTopic *kafka.Writer
	userPopularityTopic *kafka.Writer
	updateTopic *kafka.Writer
}

func NewWriter(cfg *config.Config, topic string) *kafka.Writer {
	log.Println("cfg != nil: ", cfg != nil)
	writerCfg := kafka.WriterConfig{
		Brokers:    cfg.Kafka.Brokers,
		Topic:      topic,
		Balancer:   &kafka.LeastBytes{},
		BatchBytes: cfg.Kafka.MaxMessageBytes,
	}

	return kafka.NewWriter(writerCfg)
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		filterTopic: NewWriter(cfg, "filter.data"),
		userPopularityTopic: NewWriter(cfg, "user.popularity"),
		updateTopic: NewWriter(cfg, "update.user.data"),
	}
}

func (m *Manager) SendFilterData(ctx context.Context, response models.FilterData) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Kafka: ошибка сериализации данных фильтрации %v", err)
	}

	log.Printf("Kafka: отправка ответа в топик %s", m.filterTopic.Topic,)
	if err := m.filterTopic.WriteMessages(ctx, kafka.Message{Value: data}); err != nil {
		log.Printf("Kafka: ошибка отправки сообщения в топик %s: %v", m.filterTopic.Topic, err)
	}

	log.Printf("Kafka: успешно отправлен ответ для приказа")
}

func (m *Manager) SendUserUpdateData(ctx context.Context, response models.UpdateUserData) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Kafka: ошибка сериализации новых данных пользователя %v", err)
	}

	log.Printf("Kafka: отправка ответа в топик %s", m.updateTopic.Topic,)
	if err := m.updateTopic.WriteMessages(ctx, kafka.Message{Value: data}); err != nil {
		log.Printf("Kafka: ошибка отправки сообщения в топик %s: %v", m.updateTopic.Topic, err)
	}

	log.Printf("Kafka: успешно отправлен ответ для приказа")
}

func (m *Manager) SendUserPopularityData(ctx context.Context, username string) {

	log.Printf("Kafka: отправка ответа в топик %s", m.userPopularityTopic.Topic,)
	if err := m.userPopularityTopic.WriteMessages(ctx, kafka.Message{Value: []byte(username)}); err != nil {
		log.Printf("Kafka: ошибка отправки сообщения в топик %s: %v", m.userPopularityTopic.Topic, err)
	}

	log.Printf("Kafka: успешно отправлен ответ для приказа")
}

package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/DmitriySama/teammate_search/config"
	"github.com/DmitriySama/teammate_search/internal/domain"
	"github.com/DmitriySama/teammate_search/internal/models"
	"github.com/DmitriySama/teammate_search/internal/producer"
	tsService "github.com/DmitriySama/teammate_search/internal/services/teammateSearchService"
)

type Manager struct {
	ordersReader   *kafka.Reader
	requestsReader *kafka.Reader
	updatesReader  *kafka.Reader
	producer       *producer.Manager
	service        *tsService.Service
	metrics        *metrics.Collector
}

func New(cfg *config.Config, service *tsService.Service, collector *metrics.Collector, prod *producer.Manager) *Manager {
	return &Manager{
		ordersReader:   newReader(cfg, cfg.Topics.OrdersTasks, cfg.Kafka.GroupID+"-orders"),
		requestsReader: newReader(cfg, cfg.Topics.GetOrderTasks, cfg.Kafka.GroupID+"-requests"),
		updatesReader:  newReader(cfg, cfg.Topics.UpdateTasksStatus, cfg.Kafka.GroupID+"-updates"),
		producer:       prod,
		service:        service,
		metrics:        collector,
	}
}

func newReader(cfg *config.Config, topic, group string) *kafka.Reader {
	readerCfg := kafka.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		GroupID:        group,
		Topic:          topic,
		MinBytes:       1,
		MaxBytes:       cfg.Kafka.MaxMessageBytes,
		CommitInterval: 0,
	}
	return kafka.NewReader(readerCfg)
}

func (m *Manager) Run(ctx context.Context) error {
	log.Printf("Kafka: запуск консьюмеров")
	errCh := make(chan error, 3)

	go func() { errCh <- m.consumeOrdersTasks(ctx) }()
	go func() { errCh <- m.consumeOrderTaskRequests(ctx) }()
	go func() { errCh <- m.consumeUpdateTasks(ctx) }()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("Kafka: ошибка консьюмера: %v", err)
			return err
		}
		return nil
	}
}

func (m *Manager) consumeOrdersTasks(ctx context.Context) error {
	log.Printf("Kafka: запуск консьюмера для топика orders.tasks")
	for {
		msg, err := m.ordersReader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("Kafka: остановка консьюмера orders.tasks")
				return nil
			}
			log.Printf("Kafka: ошибка получения сообщения из топика orders.tasks: %v", err)
			return err
		}

		m.metrics.RecordReceived(ctx)

		log.Printf("Kafka: получено сообщение из топика orders.tasks")
		var event models.OrdersTasksEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Kafka: невалидный payload в топике orders.tasks: %v", err)
			m.metrics.RecordError(ctx, metrics.CategoryInvalidInput)
			_ = m.ordersReader.CommitMessages(ctx, msg)
			continue
		}

		log.Printf("Kafka: обработка события orders.tasks для приказа %s", event.OrderID)
		err = m.service.HandleOrdersTasks(ctx, event)
		if err != nil {
			log.Printf("Kafka: ошибка сохранения задач для приказа %s: %v", event.OrderID, err)
			m.metrics.RecordError(ctx, classifyError(err))
			continue
		}

		log.Printf("Kafka: успешно обработано событие orders.tasks для приказа %s", event.OrderID)
		m.metrics.RecordProcessed(ctx)

		if err := m.ordersReader.CommitMessages(ctx, msg); err != nil {
			m.metrics.RecordError(ctx, metrics.CategoryUnknown)
			return err
		}
	}
}

func (m *Manager) consumeOrderTaskRequests(ctx context.Context) error {
	log.Printf("Kafka: запуск консьюмера для топика get.order.tasks")
	for {
		msg, err := m.requestsReader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("Kafka: остановка консьюмера get.order.tasks")
				return nil
			}
			log.Printf("Kafka: ошибка получения сообщения из топика get.order.tasks: %v", err)
			return err
		}

		m.metrics.RecordReceived(ctx)

		log.Printf("Kafka: получено сообщение из топика get.order.tasks")
		var request models.OrderTasksRequest
		if err := json.Unmarshal(msg.Value, &request); err != nil {
			log.Printf("Kafka: невалидный payload в топике get.order.tasks: %v", err)
			m.metrics.RecordError(ctx, metrics.CategoryInvalidInput)
			_ = m.requestsReader.CommitMessages(ctx, msg)
			continue
		}

		log.Printf("Kafka: обработка запроса get.order.tasks для приказа %s", request.OrderID)
		response, err := m.service.BuildOrderTasksResponse(ctx, request)
		if err != nil {
			log.Printf("Kafka: ошибка построения ответа для приказа %s: %v", request.OrderID, err)
			m.metrics.RecordError(ctx, classifyError(err))
			_ = m.requestsReader.CommitMessages(ctx, msg)
			continue
		}

		if err := m.producer.SendOrderTasksResponse(ctx, response); err != nil {
			m.metrics.RecordError(ctx, metrics.CategoryUnknown)
			return err
		}

		if err := m.requestsReader.CommitMessages(ctx, msg); err != nil {
			m.metrics.RecordError(ctx, metrics.CategoryUnknown)
			return err
		}

		m.metrics.RecordProcessed(ctx)
	}
}

func (m *Manager) consumeUpdateTasks(ctx context.Context) error {
	log.Printf("Kafka: запуск консьюмера для топика update.tasks.status")
	for {
		msg, err := m.updatesReader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("Kafka: остановка консьюмера update.tasks.status")
				return nil
			}
			log.Printf("Kafka: ошибка получения сообщения из топика update.tasks.status: %v", err)
			return err
		}

		m.metrics.RecordReceived(ctx)

		log.Printf("Kafka: получено сообщение из топика update.tasks.status")
		var event models.UpdateTasksEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Kafka: невалидный payload в топике update.tasks.status: %v", err)
			m.metrics.RecordError(ctx, metrics.CategoryInvalidInput)
			_ = m.updatesReader.CommitMessages(ctx, msg)
			continue
		}

		log.Printf("Kafka: обработка обновления статусов задач для приказа %s", event.OrderID)
		err = m.service.ApplyTaskUpdates(ctx, event)
		if err != nil {
			log.Printf("Kafka: ошибка обновления задач для приказа %s: %v", event.OrderID, err)
			m.metrics.RecordError(ctx, classifyError(err))
		} else {
			log.Printf("Kafka: успешно обновлены задачи для приказа %s", event.OrderID)
			m.metrics.RecordProcessed(ctx)
		}

		if err := m.updatesReader.CommitMessages(ctx, msg); err != nil {
			m.metrics.RecordError(ctx, metrics.CategoryUnknown)
			return err
		}
	}
}

func classifyError(err error) string {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return metrics.CategoryInvalidInput
	case errors.Is(err, domain.ErrTaskNotFound):
		return metrics.CategoryTaskNotFound
	case errors.Is(err, domain.ErrInvalidStatus):
		return metrics.CategoryInvalidStatus
	case errors.Is(err, domain.ErrInvalidDeadline):
		return metrics.CategoryInvalidDeadline
	default:
		return metrics.CategoryUnknown
	}
}

func (m *Manager) Close() {
	if m.ordersReader != nil {
		_ = m.ordersReader.Close()
	}
	if m.requestsReader != nil {
		_ = m.requestsReader.Close()
	}
	if m.updatesReader != nil {
		_ = m.updatesReader.Close()
	}
	if m.producer != nil {
		_ = m.producer.Close()
	}
}

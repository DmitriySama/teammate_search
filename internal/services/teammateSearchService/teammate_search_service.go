package registryservice

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/DmitriySama/teammate_search/internal/domain"
	"github.com/DmitriySama/teammate_search/internal/models"
)

type UsersStorage interface {
	ListGroups(ctx context.Context, orderID uuid.UUID) ([]models.Group, error)
	ListGroupTasks(ctx context.Context, orderID, groupID uuid.UUID) ([]models.Task, string, error)
	UpdateTaskStatus(ctx context.Context, orderID uuid.UUID, taskID int64, status string) (bool, error)
	GetTaskGroupID(ctx context.Context, orderID uuid.UUID, taskID int64) (uuid.UUID, error)
	OrderTasksForResponse(ctx context.Context, orderID uuid.UUID) ([]models.OrderTaskRecord, error)
	UpsertOrderTasks(ctx context.Context, evt models.OrdersTasksEvent) ([]models.GroupTasks, error)
	ApplyTaskUpdates(ctx context.Context, evt models.UpdateTasksEvent) error
}

type UsersCache interface {
	Get(ctx context.Context, orderID, groupID string) ([]models.Task, string, bool)
	Set(ctx context.Context, orderID, groupID, groupName string, tasks []models.Task)
	UpdateTaskInCache(ctx context.Context, orderID, groupID string, taskID int64, status string)
	InvalidateGroup(ctx context.Context, orderID, groupID string)
	InvalidateOrder(ctx context.Context, orderID string)
}

type Service struct {
	storage UsersStorage
	cache   UsersCache
}

func New(storage UsersStorage, cache UsersCache) *Service {
	return &Service{storage: storage, cache: cache}
}

func (s *Service) ListGroups(ctx context.Context, orderIDStr string) ([]models.Group, error) {
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return nil, fmt.Errorf("Неверный order_id: %w", err)
	}

	return s.storage.ListGroups(ctx, orderID)
}

func (s *Service) ListGroupTasks(ctx context.Context, orderIDStr, groupIDStr string) (string, []models.Task, error) {
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return "", nil, fmt.Errorf("Неверный order_id: %w", err)
	}

	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		return "", nil, fmt.Errorf("Неверный group_id: %w", err)
	}

	cachedTasks, cachedGroupName, ok := s.cache.Get(ctx, orderIDStr, groupIDStr)
	if ok {
		return cachedGroupName, cachedTasks, nil
	}

	tasks, groupName, err := s.storage.ListGroupTasks(ctx, orderID, groupID)
	if err != nil {
		return "", nil, err
	}

	s.cache.Set(ctx, orderIDStr, groupIDStr, groupName, tasks)
	return groupName, tasks, nil
}

func (s *Service) UpdateTaskStatus(ctx context.Context, orderIDStr, taskIDStr, status string) error {
	if !validateStatus(status) {
		return domain.ErrInvalidStatus
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return fmt.Errorf("Неверный order_id: %w", err)
	}

	taskID, err := strconv.ParseInt(strings.TrimSpace(taskIDStr), 10, 64)
	if err != nil {
		return fmt.Errorf("Неверный task_id: %w", err)
	}

	updated, err := s.storage.UpdateTaskStatus(ctx, orderID, taskID, status)
	if err != nil {
		return err
	}
	if !updated {
		return domain.ErrTaskNotFound
	}

	// Обновляем задачу в кэше без инвалидации всего кэша
	groupID, err := s.storage.GetTaskGroupID(ctx, orderID, taskID)
	if err == nil {
		s.cache.UpdateTaskInCache(ctx, orderIDStr, groupID.String(), taskID, status)
	}

	return nil
}

// Сохраняет входные задачи и обновляет кеш Redis для обработанных задач по группам
func (s *Service) HandleOrdersTasks(ctx context.Context, evt models.OrdersTasksEvent) error {
	for _, groupPayload := range evt.Groups {
		for _, taskPayload := range groupPayload.Tasks {
			if !validateDeadline(taskPayload.Deadline) {
				return domain.ErrInvalidDeadline
			}
		}
	}

	groupTasks, err := s.storage.UpsertOrderTasks(ctx, evt)
	if err != nil {
		return err
	}

	for _, group := range groupTasks {
		if len(group.Tasks) == 0 {
			s.cache.InvalidateGroup(ctx, evt.OrderID, group.GroupID)
			continue
		}
		s.cache.Set(ctx, evt.OrderID, group.GroupID, group.GroupName, group.Tasks)
	}

	return nil
}

func (s *Service) BuildOrderTasksResponse(ctx context.Context, req models.OrderTasksRequest) (models.OrderTasksResponse, error) {
	response := models.OrderTasksResponse{
		RequestID: req.RequestID,
		OrderID:   req.OrderID,
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return response, fmt.Errorf("Неверный order_id: %w", err)
	}

	records, err := s.storage.OrderTasksForResponse(ctx, orderID)
	if err != nil {
		return response, err
	}

	response.Tasks = records
	return response, nil
}

func (s *Service) ApplyTaskUpdates(ctx context.Context, evt models.UpdateTasksEvent) error {
	if err := s.storage.ApplyTaskUpdates(ctx, evt); err != nil {
		return err
	}

	s.cache.InvalidateOrder(ctx, evt.OrderID)
	return nil
}

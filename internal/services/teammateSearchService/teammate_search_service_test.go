package registryservice

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/DmitriySama/teammate_search/internal/models"
	"github.com/DmitriySama/teammate_search/internal/services/teammateSearchService/mocks"
)

type RegistryServiceSuite struct {
	suite.Suite
	ctx     context.Context
	cache   *mocks.MockTasksCache
	storage *mocks.MockStorage
	svc     *Service
}

func (s *RegistryServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.cache = mocks.NewMockTasksCache(s.T())
	s.storage = mocks.NewMockStorage(s.T())
	s.svc = New(s.storage, s.cache)
}

func (s *RegistryServiceSuite) TestListGroupTasks_CacheHit() {
	orderID := uuid.NewString()
	groupID := uuid.NewString()

	s.cache.EXPECT().Get(s.ctx, orderID, groupID).Return([]models.Task{{ID: 1}}, "group1", true)

	gotGroupName, gotTasks, err := s.svc.ListGroupTasks(s.ctx, orderID, groupID)
	s.Require().NoError(err)

	wantGroupName := "group1"
	assert.Equal(s.T(), wantGroupName, gotGroupName, "названия группы должны совпадать")

	wantTasksLen := 1
	gotTasksLen := len(gotTasks)
	assert.Equal(s.T(), wantTasksLen, gotTasksLen, "длина списка задач должна совпадать")

	wantTaskID := int64(1)
	gotTaskID := gotTasks[0].ID
	assert.Equal(s.T(), wantTaskID, gotTaskID, "ID задачи должен совпадать")

	s.cache.AssertNotCalled(s.T(), "Set")
}

func (s *RegistryServiceSuite) TestListGroupTasks_CacheMissStoresResult() {
	orderID := uuid.NewString()
	groupID := uuid.NewString()
	orderIDUUID, _ := uuid.Parse(orderID)
	groupIDUUID, _ := uuid.Parse(groupID)

	s.cache.EXPECT().Get(s.ctx, orderID, groupID).Return(nil, "", false)
	s.storage.EXPECT().ListGroupTasks(s.ctx, orderIDUUID, groupIDUUID).Return([]models.Task{{ID: 2}}, "group", nil)
	s.cache.EXPECT().Set(s.ctx, orderID, groupID, "group", []models.Task{{ID: 2}}).Return()

	_, gotTasks, err := s.svc.ListGroupTasks(s.ctx, orderID, groupID)
	s.Require().NoError(err)

	wantTasksLen := 1
	gotTasksLen := len(gotTasks)
	assert.Equal(s.T(), wantTasksLen, gotTasksLen, "длина списка задач должна совпадать")

	wantTaskID := int64(2)
	gotTaskID := gotTasks[0].ID
	assert.Equal(s.T(), wantTaskID, gotTaskID, "ID задачи должен совпадать")

	s.cache.AssertCalled(s.T(), "Set", s.ctx, orderID, groupID, "group", []models.Task{{ID: 2}})
}

func (s *RegistryServiceSuite) TestHandleOrdersTasks_InvalidatesGroups() {
	evt := models.OrdersTasksEvent{
		OrderID: uuid.NewString(),
		Groups: []models.GroupPayload{
			{GroupID: "group1"},
			{GroupID: "group2"},
		},
	}

	groupTasks := []models.GroupTasks{
		{
			GroupID:   "group1",
			GroupName: "Группа 1",
			Tasks: []models.Task{
				{ID: 1},
			},
		},
		{
			GroupID:   "group2",
			GroupName: "Группа 2",
		},
	}

	s.storage.EXPECT().UpsertOrderTasks(s.ctx, evt).Return(groupTasks, nil)
	s.cache.EXPECT().Set(s.ctx, evt.OrderID, "group1", "Группа 1", []models.Task{{ID: 1}}).Return()
	s.cache.EXPECT().InvalidateGroup(s.ctx, evt.OrderID, "group2").Return()

	err := s.svc.HandleOrdersTasks(s.ctx, evt)
	s.Require().NoError(err)
	s.cache.AssertCalled(s.T(), "Set", s.ctx, evt.OrderID, "group1", "Группа 1", []models.Task{{ID: 1}})
	s.cache.AssertCalled(s.T(), "InvalidateGroup", s.ctx, evt.OrderID, "group2")
}

func (s *RegistryServiceSuite) TestApplyTaskUpdates_InvalidatesOrder() {
	evt := models.UpdateTasksEvent{OrderID: uuid.NewString()}

	s.storage.EXPECT().ApplyTaskUpdates(s.ctx, evt).Return(nil)
	s.cache.EXPECT().InvalidateOrder(s.ctx, evt.OrderID).Return()

	err := s.svc.ApplyTaskUpdates(s.ctx, evt)
	s.Require().NoError(err)
	s.cache.AssertCalled(s.T(), "InvalidateOrder", s.ctx, evt.OrderID)
}

func (s *RegistryServiceSuite) TestUpdateTaskStatus_InvalidStatus() {
	gotErr := s.svc.UpdateTaskStatus(s.ctx, uuid.NewString(), "1", "test")
	wantErr := domain.ErrInvalidStatus
	assert.ErrorIs(s.T(), gotErr, wantErr, "должна возвращаться ошибка неверного статуса")
}

func (s *RegistryServiceSuite) TestUpdateTaskStatus_NotFound() {
	orderID := uuid.NewString()
	orderIDUUID, _ := uuid.Parse(orderID)
	taskID := int64(1)

	s.storage.EXPECT().UpdateTaskStatus(s.ctx, orderIDUUID, taskID, "Выполнено").Return(false, nil)

	gotErr := s.svc.UpdateTaskStatus(s.ctx, orderID, "1", "Выполнено")
	wantErr := domain.ErrTaskNotFound
	assert.ErrorIs(s.T(), gotErr, wantErr, "должна возвращаться ошибка не найденной записи задачи")
}

func (s *RegistryServiceSuite) TestListGroups_InvalidUUID() {
	_, gotErr := s.svc.ListGroups(s.ctx, "not-a-uuid")
	assert.Error(s.T(), gotErr, "должна возвращаться ошибка неверного UUID")
}

func TestRegistryServiceSuite(t *testing.T) {
	suite.Run(t, new(RegistryServiceSuite))
}

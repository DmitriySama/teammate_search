package teammateSearchService

import (
	"context"
	"testing"
	"errors"
	"net/http"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DmitriySama/teammate_search/internal/services/teammateSearchService/mocks"
	"github.com/DmitriySama/teammate_search/internal/storage/pgstorage"
	"github.com/DmitriySama/teammate_search/internal/models"
)

type TeammateSearchServiceSuite struct {
	suite.Suite
	ctx     context.Context
	cache   *mocks.MockUsersCache
	storage *mocks.MockUsersStorage
	svc     *Service
	req         *http.Request
}

func (s *TeammateSearchServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.cache = mocks.NewMockUsersCache(s.T())
	s.storage = mocks.NewMockUsersStorage(s.T())
	s.svc = New(s.storage, s.cache)
}


func TestTeammateSearchServiceSuite(t *testing.T) {
	suite.Run(t, new(TeammateSearchServiceSuite))
}



func (s *TeammateSearchServiceSuite) TestRegister_Success() {
    username, password, description:= "testuser", "pass123", "Test user"
    age := 25
    
    expectedResult := &pgstorage.AuthResult{User: &models.User{
		ID: 0,
		Username: username,
		Password: password,
		Description: description,
		Age: age,
		MostLikeGame: "",
		MostLikeGenre: "",
		Language: "",
		App: "",
		CreatedAt: time.Now(),
	}, Success: true}
    s.storage.On("Register", username, password, description, age).
        Return(expectedResult, nil)
    
    result, err := s.svc.storage.Register(username, password, description, age)
    
    s.NoError(err)
    s.Equal(expectedResult, result)
    s.storage.AssertExpectations(s.T())
}

func (s *TeammateSearchServiceSuite) TestRegister_UserExists() {
    username := "QQQ"
    s.storage.On("UserExists", username).Return(true, nil)
    
    _, err := s.svc.storage.Register(username, "pass", "desc", 20)
    
    s.Error(err)
    s.Contains(err.Error(), "такой ник существует")
}

func (s *TeammateSearchServiceSuite) TestUserExists_NotFound() {
    username := "newuser"
    s.storage.On("UserExists", username).Return(false, nil)
    
    exists, err := s.svc.storage.UserExists(username)
    
    s.NoError(err)
    s.False(exists)
}

func (s *TeammateSearchServiceSuite) TestLogin_Success() {
    username, password := "user", "pass"
    expected := &pgstorage.AuthResult{User: &models.User{}, Success: true}
    
    s.storage.On("Login", username, password).Return(expected, nil)
    
    result, err := s.svc.storage.Login(username, password)
    
    s.NoError(err)
    s.Equal(expected, result)
}

func (s *TeammateSearchServiceSuite) TestLogin_InvalidCredentials() {
    username, password := "user", "wrongpass"
    
    s.storage.On("Login", username, password).Return((*pgstorage.AuthResult)(nil), errors.New("invalid username or password"))
    
    result, err := s.svc.storage.Login(username, password)
    
    s.Error(err)
    s.Nil(result)
    s.Contains(err.Error(), "invalid")
    s.storage.AssertExpectations(s.T())
}

func (s *TeammateSearchServiceSuite) TestUpdateUser_Success() {
    user := models.User{ID: 1, Username: "Updated"}
    s.storage.On("UpdateUser", s.req, user).Return(nil)
    
    err := s.svc.storage.UpdateUser(s.req, user)
    
    s.NoError(err)
}

func (s *TeammateSearchServiceSuite) TestFindUser_Success() {
    s.storage.On("FindUser", "user", "pass").Return(1, nil)
    
    id, err := s.svc.storage.FindUser("user", "pass")
    
    s.NoError(err)
    s.Equal(1, id)
}


func (s *TeammateSearchServiceSuite) TestGetLanguages_Success() {
    expected := []models.Language{{ID: 1, Lang: "Russian"}, {ID: 2, Lang: "English"}}
    s.storage.On("GetLanguages", s.ctx).Return(expected, nil)
    
    langs, err := s.svc.GetLanguages(s.ctx)
    
    s.NoError(err)
    s.Len(langs, 2)
    s.Equal("Russian", langs[0].Lang)
}

func (s *TeammateSearchServiceSuite) TestGetGenres_Success() {
    expected := []models.Genres{{ID: 1, Genre: "Action"}}
    s.storage.On("GetGenres", s.ctx).Return(expected, nil)
    
    genres, err := s.svc.GetGenres(s.ctx)
    
    s.NoError(err)
    s.Equal(expected, genres)
}

func (s *TeammateSearchServiceSuite) TestGetGames_Success() {
    expected := []models.Games{{ID: 1, Game: "Game1"}}
    s.storage.On("GetGames", s.ctx).Return(expected, nil)
    
    games, err := s.svc.GetGames(s.ctx)
    
    s.NoError(err)
    s.Len(games, 1)
}

func (s *TeammateSearchServiceSuite) TestGetUserByID_Success() {
    expected := &models.User{ID: 1, Username: "John"}
    s.storage.On("GetUserByID", 1).Return(expected, nil)
    
    user, err := s.svc.storage.GetUserByID(1)
    
    s.NoError(err)
    s.Equal(expected, user)
}

func (s *TeammateSearchServiceSuite) TestGetUserCount_Success() {
    s.storage.On("GetUserCount").Return(100, nil)
    
    count, err := s.svc.storage.GetUserCount()
    
    s.NoError(err)
    s.Equal(100, count)
}


func (s *TeammateSearchServiceSuite) TestGetApps_Success() {
    expected := []models.Apps{{ID: 1, App: "App1"}}
    s.storage.On("GetApps", s.ctx).Return(expected, nil)
    
    apps, err := s.svc.GetApps(s.ctx)
    
    s.NoError(err)
    s.Len(apps, 1)
}

// MARK: Error Cases

func (s *TeammateSearchServiceSuite) TestGetLanguages_Error() {
    s.storage.On("GetLanguages", s.ctx).Return(nil, errors.New("db error"))
    
    _, err := s.svc.GetLanguages(s.ctx)
    
    s.Error(err)
}
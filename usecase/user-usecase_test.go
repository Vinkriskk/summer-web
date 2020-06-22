package usecase

import (
	"summer-web/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserMockRepository struct {
	mock.Mock
}

func (mock *UserMockRepository) GetUserByID(id uint, user *models.User) error {
	args := mock.Called()

	user.ID = id
	user.Username = "joko123"
	user.FollowerCount = 2
	user.FollowingCount = 3
	user.Email = "joko@joko.com"
	user.Name = "joko"
	user.Password = "123"

	return args.Error(0)
}

func (mock *UserMockRepository) AddUser(user *models.User) error {
	args := mock.Called()

	return args.Error(0)
}

func (mock *UserMockRepository) GetUserByUsername(username string, user *models.User) error {
	args := mock.Called()

	user.Username = username
	user.ID = 1
	user.FollowerCount = 2
	user.FollowingCount = 3
	user.Email = "joko@joko.com"
	user.Name = "joko"
	user.Password = "123"

	return args.Error(0)
}

func (mock *UserMockRepository) UpdateUser(updatedData models.User) error {
	args := mock.Called()

	return args.Error(0)
}

func TestAddingEmptyUsername(t *testing.T) {
	assert := assert.New(t)

	testUsecase := NewUserUsecase(nil)

	user := models.User{Email: "abcdefg@gmail.com", Name: "joko", FollowerCount: 1, FollowingCount: 2, Password: "ABcd"}

	err := testUsecase.AddUser(&user)

	assert.NotNil(err)

	assert.Equal("pg: can't be null \"users_username_key\"", err.Error())
}

func TestAddingEmptyName(t *testing.T) {
	assert := assert.New(t)

	testUsecase := NewUserUsecase(nil)

	user := models.User{Email: "abcdefg@gmail.com", Username: "joko", FollowerCount: 1, FollowingCount: 2, Password: "ABcd"}

	err := testUsecase.AddUser(&user)

	assert.NotNil(err)

	assert.Equal("pg: can't be null \"users_name_key\"", err.Error())
}

func TestAddingEmptyEmail(t *testing.T) {
	assert := assert.New(t)

	testUsecase := NewUserUsecase(nil)

	user := models.User{Name: "joko too", Username: "joko", FollowerCount: 1, FollowingCount: 2, Password: "ABcd"}

	err := testUsecase.AddUser(&user)

	assert.NotNil(err)

	assert.Equal("pg: can't be null \"users_email_key\"", err.Error())
}

func TestAddingInvalidEmail(t *testing.T) {
	assert := assert.New(t)

	testUsecase := NewUserUsecase(nil)

	user := models.User{Email: "asdasdasd", Name: "joko too", Username: "joko", FollowerCount: 1, FollowingCount: 2, Password: "ABcd"}

	err := testUsecase.AddUser(&user)

	assert.NotNil(err)

	assert.Equal("error: invalid \"users_email_key\"", err.Error())
}

func TestGetUserByID(t *testing.T) {
	id := uint(1)
	mockRepo := new(UserMockRepository)

	mockRepo.On("GetUserByID").Return(nil)

	testUsecase := NewUserUsecase(mockRepo)

	user := models.User{}

	err := testUsecase.GetUserByID(id, &user)

	mockRepo.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, id, user.ID)
}

func TestAddUser(t *testing.T) {
	mockRepo := new(UserMockRepository)

	user := models.User{Email: "asdasd@asd", Name: "joko too", Username: "joko", FollowerCount: 1, FollowingCount: 2, Password: "ABcd"}

	mockRepo.On("AddUser").Return(nil)

	testUsecase := NewUserUsecase(mockRepo)

	err := testUsecase.AddUser(&user)

	mockRepo.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(UserMockRepository)

	updatedData := models.User{ID: 1, Email: "asdasd@asd.com", Name: "joko too", Username: "joko", FollowerCount: 1, FollowingCount: 2, Password: "ABcd"}

	mockRepo.On("UpdateUser").Return(nil)

	testUsecase := NewUserUsecase(mockRepo)

	err := testUsecase.UpdateUser(updatedData)

	mockRepo.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestLogin(t *testing.T) {
	mockRepo := new(UserMockRepository)
	testUsecase := NewUserUsecase(mockRepo)

	mockRepo.On("GetUserByUsername").Return(nil)

	loginData := models.User{Username: "joko", Password: "123"}

	token, err := testUsecase.Login(loginData)

	mockRepo.AssertExpectations(t)
	assert.Nil(t, err)
	assert.NotNil(t, token)
}

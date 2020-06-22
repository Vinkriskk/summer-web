package usecase

import (
	"summer-web/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type PostMockRepository struct {
	mock.Mock
}

func (mock *PostMockRepository) GetPosts() ([]models.Post, error) {
	args := mock.Called()

	result := args.Get(0)

	return result.([]models.Post), args.Error(1)
}

func (mock *PostMockRepository) AddPost(post *models.Post) error {
	args := mock.Called()
	return args.Error(0)
}

func TestAddingEmptyCaption(t *testing.T) {
	assert := assert.New(t)

	testUsecase := NewPostUsecase(nil)

	post := models.Post{UserID: 1}

	err := testUsecase.AddPost(&post)

	assert.NotNil(err)
	assert.Equal("pg: can't be null \"posts_caption_key\"", err.Error())
}

func TestAddingEmptyUserID(t *testing.T) {
	assert := assert.New(t)

	testUsecase := NewPostUsecase(nil)

	post := models.Post{Caption: "ASDASDASD"}

	err := testUsecase.AddPost(&post)

	assert.NotNil(err)

	assert.Equal(("pg: can't be null \"posts_user_id_key\""), err.Error())
}

func TestFindAll(t *testing.T) {
	mockRepo := new(PostMockRepository)

	post := models.Post{Caption: "ADD", UserID: 123}
	// SETUP EXPECTATIONS
	mockRepo.On("GetPosts").Return([]models.Post{post}, nil)

	testUsecase := NewPostUsecase(mockRepo)

	result, err := testUsecase.GetPosts()

	// MOCK ASSERTIONS: BEHAVIOUR
	mockRepo.AssertExpectations(t)

	// DATA ASSERTION
	assert.Nil(t, err)

	assert.Equal(t, post.Caption, result[0].Caption)
	assert.Equal(t, post.UserID, result[0].UserID)
}

func TestCreate(t *testing.T) {
	mockRepo := new(PostMockRepository)

	post := models.Post{Caption: "ASDASD", UserID: 1}

	// SETUP EXPECTATIONS
	mockRepo.On("AddPost").Return(nil)

	testUsecase := NewPostUsecase(mockRepo)

	err := testUsecase.AddPost(&post)

	mockRepo.AssertExpectations(t)

	assert.Nil(t, err)
}

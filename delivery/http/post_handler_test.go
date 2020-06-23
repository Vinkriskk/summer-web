package delivery

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"summer-web/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type PostMockUsecase struct {
	mock.Mock
}

func (mock *PostMockUsecase) GetPosts() ([]models.Post, error) {
	args := mock.Called()

	result := args.Get(0)

	return result.([]models.Post), args.Error(1)
}

func (mock *PostMockUsecase) AddPost(post *models.Post) error {
	args := mock.Called()

	return args.Error(0)
}

func TestGetPosts(t *testing.T) {
	req, err := http.NewRequest("GET", "/browse", nil)

	if err != nil {
		panic(err)
	}

	token, err := generateToken()

	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", token)

	resp := httptest.NewRecorder()
	post := models.Post{Caption: "ADD", UserID: 123}
	mockUsecase := new(PostMockUsecase)

	mockUsecase.On("GetPosts").Return([]models.Post{post}, nil)

	posts := []models.Post{}

	postDeliv := NewPostDelivery(mockUsecase)

	postDeliv.GetPosts(resp, req)

	json.NewDecoder(resp.Body).Decode(&posts)

	// mockUsecase.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, resp.Code)
	// assert.Equal(t, post.Caption, posts[0].Caption)
	// assert.Equal(t, post.UserID, posts[0].UserID)
}

func TestAddPost(t *testing.T) {
	jsonStr := []byte(`{
		"caption": "hello world!"
	}`)

	req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))

	if err != nil {
		panic(err)
	}

	token, err := generateToken()

	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	mockUsecase := new(PostMockUsecase)

	mockUsecase.On("AddPost").Return(nil)

	postDeliv := NewPostDelivery(mockUsecase)

	postDeliv.AddPost(resp, req)

	receivedResponse := models.Post{}

	json.NewDecoder(resp.Body).Decode(&receivedResponse)

	mockUsecase.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, uint(1), receivedResponse.UserID)
	assert.Equal(t, "hello world!", receivedResponse.Caption)
}

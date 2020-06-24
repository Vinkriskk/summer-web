package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"summer-web/models"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserMockUsecase struct {
	mock.Mock
}

func (mock *UserMockUsecase) GetUserByID(id uint, user *models.User) error {
	args := mock.Called()
	user.ID = id
	return args.Error(0)
}

func (mock *UserMockUsecase) AddUser(user *models.User) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *UserMockUsecase) Login(loginData models.User) (string, error) {
	args := mock.Called()
	result := args.Get(0)

	return result.(string), args.Error(1)
}

func (mock *UserMockUsecase) UpdateUser(updatedData models.User) error {
	args := mock.Called()
	return args.Error(0)
}

func TestLoginSuccess(t *testing.T) {
	buf := new(bytes.Buffer)

	w := multipart.NewWriter(buf)

	err := w.WriteField("username", "us1")
	if err != nil {
		panic(err)
	}

	err = w.WriteField("password", "pw1")
	if err != nil {
		panic(err)
	}

	w.Close()

	req, err := http.NewRequest("POST", "/login", buf)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp := httptest.NewRecorder()
	mockUsecase := new(UserMockUsecase)

	mockUsecase.On("Login").Return("valid token", nil)

	userDeliv := NewUserDelivery(mockUsecase)

	userDeliv.Login(resp, req)

	receivedResponse := struct {
		AuthToken string `json:"auth_token"`
		Error     string `json:"error"`
	}{}

	mockUsecase.AssertExpectations(t)

	json.NewDecoder(resp.Body).Decode(&receivedResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "valid token", receivedResponse.AuthToken)
	assert.Equal(t, "", receivedResponse.Error)
}

func TestLoginFailed(t *testing.T) {
	buf := new(bytes.Buffer)

	w := multipart.NewWriter(buf)

	err := w.WriteField("username", "us1")
	if err != nil {
		panic(err)
	}

	err = w.WriteField("password", "pw1")
	if err != nil {
		panic(err)
	}

	w.Close()

	req, err := http.NewRequest("POST", "/login", buf)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp := httptest.NewRecorder()
	mockUsecase := new(UserMockUsecase)

	mockUsecase.On("Login").Return("", fmt.Errorf("please provide a correct credentials"))

	userDeliv := NewUserDelivery(mockUsecase)

	userDeliv.Login(resp, req)

	receivedResponse := struct {
		AuthToken string `json:"auth_token"`
		Error     string `json:"error"`
	}{}

	json.NewDecoder(resp.Body).Decode(&receivedResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, "", receivedResponse.AuthToken)
	assert.Equal(t, "please provide a correct credentials", receivedResponse.Error)
}

func TestSignUp(t *testing.T) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	err := w.WriteField("username", "us1")
	if err != nil {
		panic(err)
	}

	err = w.WriteField("password", "123")
	if err != nil {
		panic(err)
	}

	err = w.WriteField("email", "asdaweee@qsd.com")
	if err != nil {
		panic(err)
	}

	err = w.WriteField("name", "test")
	if err != nil {
		panic(err)
	}

	err = w.WriteField("follower_count", "1")
	if err != nil {
		panic(err)
	}

	err = w.WriteField("following_count", "2")
	if err != nil {
		panic(err)
	}

	w.Close()

	req, err := http.NewRequest("POST", "/sign_up", buf)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp := httptest.NewRecorder()
	mockUsecase := new(UserMockUsecase)

	mockUsecase.On("AddUser").Return(nil)

	userDeliv := NewUserDelivery(mockUsecase)
	userDeliv.AddUser(resp, req)

	receivedResponse := models.User{}

	json.NewDecoder(resp.Body).Decode(&receivedResponse)

	mockUsecase.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "us1", receivedResponse.Username)
	assert.Equal(t, "", receivedResponse.Password)
}

func TestGetUserByID(t *testing.T) {
	searchID := 1
	req, err := http.NewRequest("GET", "/users", nil)

	req = mux.SetURLVars(req, map[string]string{
		"id": strconv.Itoa(searchID),
	})
	if err != nil {
		panic(err)
	}

	token, err := generateToken()

	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", token)

	resp := httptest.NewRecorder()
	mockUsecase := new(UserMockUsecase)

	mockUsecase.On("GetUserByID").Return(nil)

	userDeliv := NewUserDelivery(mockUsecase)

	userDeliv.GetUserByID(resp, req)

	receivedResponse := models.User{}

	json.NewDecoder(resp.Body).Decode(&receivedResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, uint(searchID), receivedResponse.ID)
}

func TestUpdateUser(t *testing.T) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	err := w.WriteField("username", "joko")

	if err != nil {
		panic(err)
	}

	w.Close()

	req, err := http.NewRequest("PATCH", "/users/update", buf)

	if err != nil {
		panic(err)
	}

	token, err := generateToken()

	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp := httptest.NewRecorder()
	mockUsecase := new(UserMockUsecase)
	mockUsecase.On("UpdateUser").Return(nil)
	mockUsecase.On("GetUserByID").Return(nil)
	userDeliv := NewUserDelivery(mockUsecase)

	userDeliv.UpdateUser(resp, req)

	receivedResponse := models.User{}

	json.NewDecoder(resp.Body).Decode(&receivedResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, uint(1), receivedResponse.ID)
}

func generateToken() (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = 1
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("SECRET_JWT_KEY")))
	if err != nil {
		return "", err
	}
	return token, err
}

package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"summer-web/models"
	"summer-web/usecase"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// UserDelivery interface acts as Post Controller
type UserDelivery interface {
	GetUserByID(resp http.ResponseWriter, req *http.Request)
	AddUser(resp http.ResponseWriter, req *http.Request)
	Login(resp http.ResponseWriter, req *http.Request)
	UpdateUser(resp http.ResponseWriter, req *http.Request)
}

type userDelivery struct{}

// NewUserDelivery returns new userDelivery struct that implements UserDelivery
func NewUserDelivery() UserDelivery {
	return &userDelivery{}
}

var (
	userUsecase usecase.UserUsecase = usecase.NewUserUsecase()
)

func (*userDelivery) GetUserByID(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	uid := uint(id)

	var user models.User

	err = userUsecase.GetUserByID(uid, &user)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	sanitizePassword(&user)

	json.NewEncoder(resp).Encode(user)
}

func (*userDelivery) AddUser(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	var newUser models.User

	err := json.NewDecoder(req.Body).Decode(&newUser)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	err = userUsecase.AddUser(&newUser)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	sanitizePassword(&newUser)

	json.NewEncoder(resp).Encode(newUser)
}

func (*userDelivery) UpdateUser(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	token, err := jwt.Parse(req.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("SECRET_JWT_KEY")), nil
	})

	claims := token.Claims.(jwt.MapClaims)

	var user models.User

	userID := uint(claims["user_id"].(float64))

	err = userUsecase.GetUserByID(userID, &user)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	json.NewDecoder(req.Body).Decode(&user)
	err = userUsecase.UpdateUser(user)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	sanitizePassword(&user)

	json.NewEncoder(resp).Encode(user)
}

func (*userDelivery) Login(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	var loginData models.User

	json.NewDecoder(req.Body).Decode(&loginData)

	token, err := userUsecase.Login(loginData)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(token))
}

func sanitizePassword(user *models.User) {
	user.Password = ""
}

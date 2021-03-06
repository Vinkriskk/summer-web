package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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

var (
	userUsecase usecase.UserUsecase
)

// NewUserDelivery returns new userDelivery struct that implements UserDelivery
func NewUserDelivery(usecaseUser ...usecase.UserUsecase) UserDelivery {
	if len(usecaseUser) > 0 {
		userUsecase = usecaseUser[0]
	} else {
		userUsecase = usecase.NewUserUsecase()
	}
	return &userDelivery{}
}

func (*userDelivery) GetUserByID(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		key, value := trimError(err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{` + key + `:` + value + `}`))
		fmt.Println(err)
		return
	}

	uid := uint(id)

	var user models.User

	err = userUsecase.GetUserByID(uid, &user)

	if err != nil {
		key, value := trimError(err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{` + key + `:` + value + `}`))
		return
	}

	sanitizePassword(&user)

	json.NewEncoder(resp).Encode(user)
}

func (*userDelivery) AddUser(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	var newUser models.User

	err := addDataToUser(&newUser, req)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	err = userUsecase.AddUser(&newUser)

	if err != nil {
		key, value := trimError(err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{` + key + `:` + value + `}`))
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
		key, value := trimError(err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{` + key + `:` + value + `}`))
		return
	}

	err = addDataToUser(&user, req)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	err = userUsecase.UpdateUser(user)

	if err != nil {
		key, value := trimError(err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{` + key + `:` + value + `}`))
		return
	}

	sanitizePassword(&user)

	json.NewEncoder(resp).Encode(user)
}

func (*userDelivery) Login(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	var loginData models.User

	loginData.Username = req.FormValue("username")
	loginData.Password = req.FormValue("password")

	token, err := userUsecase.Login(loginData)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(`{"auth_token": "` + token + `"}`))
}

func sanitizePassword(user *models.User) {
	user.Password = ""
}

func trimError(err error) (string, string) {
	initial := err.Error()

	key := strings.Split(initial, "\"")[1]
	key = "\"" + key + "\""

	initial = strings.Replace(initial, key, "", 1)

	value := strings.Split(initial, ": ")[1]

	value = "\"" + value[0:len(value)-1] + "\""

	return key, value
}

func addDataToUser(user *models.User, data *http.Request) error {
	var err error

	if data.FormValue("follower_count") != "" {
		user.FollowerCount, err = strconv.Atoi(data.FormValue("follower_count"))
		if err != nil {
			return err
		}
	}

	if data.FormValue("following_count") != "" {
		user.FollowingCount, err = strconv.Atoi(data.FormValue("following_count"))
		if err != nil {
			return err
		}
	}

	if data.FormValue("password") != "" {
		user.Password = data.FormValue("password")
	}

	if data.FormValue("username") != "" {
		user.Username = data.FormValue("username")
	}

	if data.FormValue("name") != "" {
		user.Name = data.FormValue("name")
	}

	if data.FormValue("email") != "" {
		user.Email = data.FormValue("email")
	}

	return nil
}

package usecase

import (
	"fmt"
	"os"
	"summer-web/models"
	"summer-web/user/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// UserUsecase interface defines the methods that are going to be used in usecase
type UserUsecase interface {
	GetUserByID(id uint, user *models.User) error
	AddUser(user *models.User) error
	Login(loginData models.User) (string, error)
	UpdateUser(updatedData models.User) error
}

var (
	userRepo repository.UserRepository = repository.NewUserRepository()
)

type userUsecase struct{}

// NewUserUsecase creates a new usecase to fiddle around with repository
func NewUserUsecase() UserUsecase {
	return &userUsecase{}
}

func (*userUsecase) GetUserByID(id uint, user *models.User) error {
	return userRepo.GetUserByID(id, user)
}

func (*userUsecase) AddUser(user *models.User) error {
	return userRepo.AddUser(user)
}

func (*userUsecase) UpdateUser(updatedData models.User) error {
	return userRepo.UpdateUser(updatedData)
}

func (*userUsecase) Login(loginData models.User) (string, error) {
	var attemptedUser models.User

	err := userRepo.GetUserByUsername(loginData.Username, &attemptedUser)

	if err != nil || loginData.Password != attemptedUser.Password {
		return "", fmt.Errorf("please provide a correct credentials")
	}

	return createToken(attemptedUser.ID)
}

func createToken(id uint) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = id
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("SECRET_JWT_KEY")))
	if err != nil {
		return "", err
	}
	return token, err
}

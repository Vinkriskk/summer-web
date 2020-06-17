package repository

import (
	"fmt"
	"os"

	"summer-web/models"

	"github.com/jinzhu/gorm"
	// import postgres dialect from gorm lib
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// UserRepository is the repository interface for post
type UserRepository interface {
	GetUserByID(id uint, user *models.User) error
	AddUser(user *models.User) error
	GetUserByUsername(username string, user *models.User) error
	UpdateUser(updatedUser models.User) error
}

func init() {
	db, err := gorm.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}

	defer db.Close()

	db.AutoMigrate(&models.User{})
}

type repo struct{}

// NewUserRepository create a new post repository to fiddle around with database
func NewUserRepository() UserRepository {
	return &repo{}
}

// GetUserByID returns an error if there is any, modifies the user parameter with the found record
func (*repo) GetUserByID(id uint, user *models.User) error {
	db, err := gorm.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to database")
	}

	defer db.Close()

	return db.Where("id = ?", id).Find(&user).Error
}

// AddUser returns an error if there is any, otherwise creates a new user record into database
func (*repo) AddUser(user *models.User) error {
	db, err := gorm.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to database")
	}

	defer db.Close()

	return db.Create(&user).Error
}

// GetUserByUsername returns an error if there is any, otherwise modifies the user parameter with the found record
func (*repo) GetUserByUsername(username string, user *models.User) error {
	db, err := gorm.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to database")
	}

	defer db.Close()

	return db.Where("username = ?", username).Find(&user).Error
}

func (*repo) UpdateUser(updatedData models.User) error {
	db, err := gorm.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to database")
	}

	defer db.Close()

	return db.Model(&updatedData).Updates(updatedData).Error
}

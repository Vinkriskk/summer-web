package repository

import (
	"fmt"
	"os"
	"summer-web/models"

	"github.com/jinzhu/gorm"
	// import postgres dialect from gorm lib
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// PostRepository is the repository interface for post
type PostRepository interface {
	GetPosts() ([]models.Post, error)
	AddPost(post *models.Post) error
}

func init() {
	db, err := gorm.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}

	defer db.Close()

	db.AutoMigrate(&models.Post{})
}

type repo struct {
	db *gorm.DB
}

// NewPostRepository create a new post repository to fiddle around with database
func NewPostRepository(db *gorm.DB) PostRepository {
	if db == nil {
		gdb, err := gorm.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
		if err != nil {
			fmt.Println(err.Error())
			panic("Could not connect to database")
		}

		return &repo{db: gdb}
	}
	return &repo{db: db}
}

// GetPosts returns all post in database or an error if there is an error
func (r *repo) GetPosts() ([]models.Post, error) {
	// db, err := gorm.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	panic("Could not connect to database")
	// }

	// defer db.Close()

	var posts []models.Post

	err := r.db.Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// AddPost adds post into database, returns an error instead if there is an error
func (r *repo) AddPost(post *models.Post) error {
	// db, err := gorm.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	panic("Could not connect to database")
	// }

	// defer db.Close()

	return r.db.Create(&post).Error
}

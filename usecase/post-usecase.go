package usecase

import (
	"summer-web/models"
	"summer-web/post/repository"
)

// PostUsecase interface defines the methods that are going to be used in usecase
type PostUsecase interface {
	GetPosts() ([]models.Post, error)
	AddPost(post *models.Post) error
}

var (
	postRepo repository.PostRepository = repository.NewPostRepository()
)

type postUsecase struct{}

// NewPostUsecase creates a new usecase to fiddle around with repository
func NewPostUsecase() PostUsecase {
	return &postUsecase{}
}

// GetPosts accesses repo to get all the post records in database
func (*postUsecase) GetPosts() ([]models.Post, error) {
	return postRepo.GetPosts()
}

// AddPost accesses repo to add a post record to database
func (*postUsecase) AddPost(post *models.Post) error {
	return postRepo.AddPost(post)
}

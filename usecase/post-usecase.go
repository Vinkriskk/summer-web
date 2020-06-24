package usecase

import (
	"fmt"
	"summer-web/models"
	"summer-web/post/repository"
)

// PostUsecase interface defines the methods that are going to be used in usecase
type PostUsecase interface {
	GetPosts() ([]models.Post, error)
	AddPost(post *models.Post) error
}

var (
	postRepo repository.PostRepository
)

type postUsecase struct{}

// NewPostUsecase creates a new usecase to fiddle around with repository
func NewPostUsecase(repo ...repository.PostRepository) PostUsecase {
	if len(repo) > 0 {
		postRepo = repo[0]
	} else {
		postRepo = repository.NewPostRepository(nil)
	}
	return &postUsecase{}
}

// GetPosts accesses repo to get all the post records in database
func (*postUsecase) GetPosts() ([]models.Post, error) {
	return postRepo.GetPosts()
}

// AddPost accesses repo to add a post record to database
func (*postUsecase) AddPost(post *models.Post) error {
	if err := validatePost(post); err != nil {
		return err
	}
	return postRepo.AddPost(post)
}

func validatePost(post *models.Post) error {
	if post.Caption == "" {
		return fmt.Errorf("pg: can't be null \"posts_caption_key\"")
	}
	if post.UserID == 0 {
		return fmt.Errorf("pg: can't be null \"posts_user_id_key\"")
	}
	return nil
}

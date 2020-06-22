package delivery

import (
	"encoding/json"
	"net/http"
	"summer-web/models"
	"summer-web/usecase"
)

// PostDelivery interface acts as Post Controller
type PostDelivery interface {
	GetPosts(resp http.ResponseWriter, req *http.Request)
	AddPost(resp http.ResponseWriter, req *http.Request)
}

type postDelivery struct{}

var (
	postUsecase usecase.PostUsecase
)

// NewPostDelivery returns new postDelivery struct that implements PostDelivery
func NewPostDelivery() PostDelivery {
	postUsecase = usecase.NewPostUsecase(nil)
	return &postDelivery{}
}

func (*postDelivery) GetPosts(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	posts, err := postUsecase.GetPosts()

	if err != nil {
		key, value := trimError(err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{` + key + `:` + value + `}`))
		return
	}

	json.NewEncoder(resp).Encode(posts)
}

func (*postDelivery) AddPost(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	var newPost models.Post

	err := json.NewDecoder(req.Body).Decode(&newPost)

	if err != nil {
		key, value := trimError(err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{` + key + `:` + value + `}`))
		return
	}

	err = postUsecase.AddPost(&newPost)

	if err != nil {
		key, value := trimError(err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{` + key + `:` + value + `}`))
		return
	}

	json.NewEncoder(resp).Encode(newPost)
}

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

// NewPostDelivery returns new postDelivery struct that implements PostDelivery
func NewPostDelivery() PostDelivery {
	return &postDelivery{}
}

var (
	postUsecase usecase.PostUsecase = usecase.NewPostUsecase()
)

func (*postDelivery) GetPosts(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	posts, err := postUsecase.GetPosts()

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(resp).Encode(posts)
}

func (*postDelivery) AddPost(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	var newPost models.Post

	err := json.NewDecoder(req.Body).Decode(&newPost)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	err = postUsecase.AddPost(&newPost)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(resp).Encode(newPost)
}

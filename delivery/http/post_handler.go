package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"summer-web/models"
	"summer-web/usecase"

	"github.com/dgrijalva/jwt-go"
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
func NewPostDelivery(usecasePost usecase.PostUsecase) PostDelivery {
	if postUsecase != nil {
		postUsecase = usecasePost
	} else {
		postUsecase = usecase.NewPostUsecase(nil)
	}
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
		panic(err)
	}

	token, err := jwt.Parse(req.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("SECRET_JWT_KEY")), nil
	})

	claims := token.Claims.(jwt.MapClaims)

	userID := uint(claims["user_id"].(float64))

	newPost.UserID = userID

	err = postUsecase.AddPost(&newPost)

	if err != nil {
		key, value := trimError(err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{` + key + `:` + value + `}`))
		return
	}

	json.NewEncoder(resp).Encode(newPost)
}

package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// Middleware interfaces for authorizing, etc (if there is any)
type Middleware interface {
	IsAuthorized(endpoint func(resp http.ResponseWriter, req *http.Request)) http.Handler
}

type middleware struct{}

// NewMiddleware returns middleware struct that implements Middleware interface
func NewMiddleware() Middleware {
	return &middleware{}
}

func (*middleware) IsAuthorized(endpoint func(resp http.ResponseWriter, req *http.Request)) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")
		if req.Header["Authorization"] != nil {
			token, err := jwt.Parse(req.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method")
				}
				return []byte(os.Getenv("SECRET_JWT_KEY")), nil
			})

			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
			}

			if token.Valid {
				endpoint(resp, req)
			}
		} else {
			resp.WriteHeader(http.StatusUnauthorized)
			resp.Write([]byte(`{"error": "Not authorized"}`))
		}
	})
}

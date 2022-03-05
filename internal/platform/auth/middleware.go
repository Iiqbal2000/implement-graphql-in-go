package platform

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Login(auth Auth) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(rw, "the method was not allowed", http.StatusBadRequest)
			return
		}

		username, password, ok := r.BasicAuth()
		if ok {
			rw.Header().Set("Content-Type", "application/json")

			token, err := auth.authenticate(username, password)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusUnauthorized)
				return
			}

			err = json.NewEncoder(rw).Encode(map[string]interface{}{
				"access_token": token,
			})

			if err != nil {
				log.Println("failure when encoding to json: ", err.Error())
				http.Error(rw, "internal server error", http.StatusInternalServerError)
				return
			}

			return
		}

		http.Error(rw, "username and password are needed", http.StatusUnauthorized)
	}
}

func Authorize(next http.Handler, auth Auth) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		var (
			user_id      = "0"
			isAuthorized = false
		)

		// extract token from header
		tokenIn := strings.TrimPrefix(authorizationHeader, "Bearer ")
		userId, err := auth.authorize(tokenIn)

		// if there is no error it indicates that it is a true user
		if err == nil {
			isAuthorized = true
			user_id = userId
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", user_id)
		ctx = context.WithValue(ctx, "is_authorized", isAuthorized)
		next.ServeHTTP(rw, r.WithContext(ctx))
	}
}

package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MatthewFrisby/thesis-pieces/ent"
	"github.com/MatthewFrisby/thesis-pieces/ent/user"
)

type AuthMiddleware struct {
	db *ent.Client
}

func NewAuthMiddleware(db *ent.Client) *AuthMiddleware {
	return &AuthMiddleware{
		db: db,
	}
}

func (am *AuthMiddleware) Private(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the access token from the request
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		uuid, err := ParseAccessToken(accessToken)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		ctx := context.Background()
		// get the user from the database
		entUser, err := am.db.User.Query().Where(user.UUID(*uuid)).Only(ctx)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		// set the user in the context
		ctx = context.WithValue(r.Context(), "user", entUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the access token from the request
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		uuid, err := ParseAccessToken(accessToken)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		ctx := context.Background()
		// get the user from the database
		entUser, err := am.db.User.Query().Where(user.UUID(*uuid)).Only(ctx)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		if !entUser.IsAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// set the user in the context
		ctx = context.WithValue(r.Context(), "user", entUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

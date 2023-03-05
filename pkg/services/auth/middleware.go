package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/frisbm/graduateplace/pkg/store"
)

type AuthMiddleware struct {
	db   store.Querier
	auth *AuthService
}

func NewAuthMiddleware(db store.Querier, auth *AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		db:   db,
		auth: auth,
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

		uuid, err := am.auth.ParseAccessToken(accessToken)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		ctx := r.Context()
		// get the user from the database
		user, err := am.db.GetUserFromUUID(ctx, *uuid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		// set the user in the context
		ctx = context.WithValue(ctx, "user", user)
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

		uuid, err := am.auth.ParseAccessToken(accessToken)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		ctx := r.Context()
		// get the user from the database
		user, err := am.db.GetUserFromUUID(ctx, *uuid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		if !user.IsAdmin {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// set the user in the context
		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

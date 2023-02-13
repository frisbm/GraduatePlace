package routes

import (
	"context"
	"net/http"
)

type RouteWithContext func(context.Context, http.ResponseWriter, *http.Request)

func WithContext(withContext RouteWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		withContext(ctx, w, r)
	}
}

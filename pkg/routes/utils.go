package routes

import (
	"context"
	"encoding/json"
	"net/http"
)

type RouteWithContext func(context.Context, http.ResponseWriter, *http.Request)

func WithContext(withContext RouteWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		withContext(ctx, w, r)
	}
}

func responseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func Response(w http.ResponseWriter, statusCode int, response any) {
	responseHeaders(w)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

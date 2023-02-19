package routes

import (
	"context"
	"encoding/json"
	"io"
	"mime"
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

type Body interface {
	SetFile(file []byte)
	SetFileName(filename string)
	SetFileType(filetype string)
}

func ParseMultiPartFormWithFileAndBody[T Body](req *http.Request, body T) error {
	var (
		file     []byte
		filename string
		filetype string
	)
	mr, err := req.MultipartReader()
	if err != nil {
		return err
	}

	filetype, _, err = mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if part.FormName() == "file" {
			file, err = io.ReadAll(part)
			if err != nil {
				return err
			}
			filename = part.FileName()
		}

		if part.FormName() == "body" {
			err = json.NewDecoder(part).Decode(&body)
			if err != nil {
				return err
			}
		}
	}
	body.SetFile(file)
	body.SetFileName(filename)
	body.SetFileType(filetype)
	return nil
}

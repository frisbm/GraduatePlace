package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/MatthewFrisby/thesis-pieces/pkg/constants"

	"github.com/pkg/errors"
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

func getFiletypeFromFilename(filename string) (string, error) {
	parts := strings.Split(filename, ".")
	if len(parts) <= 1 {
		return "", errors.New("improper file name, file must have extension")
	}

	extension := parts[len(parts)-1]
	switch strings.ToUpper(extension) {
	case "PDF":
		return "PDF", nil
	case "DOC":
		return "DOC", nil
	case "DOCX":
		return "DOCX", nil
	case "TXT":
		return "TXT", nil
	default:
		return "", errors.New("wrong filetype")
	}
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
			filetype, err = getFiletypeFromFilename(filename)
			if err != nil {
				return err
			}
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

func SanitizePagination(limit, offset int32) (int32, int32) {
	if limit < constants.MIN_LIMIT {
		limit = constants.MIN_LIMIT
	} else if limit > constants.MAX_LIMIT {
		limit = constants.MIN_LIMIT
	}

	if offset < constants.MIN_OFFSET {
		offset = constants.MIN_OFFSET
	} else if offset > constants.MAX_OFFSET {
		offset = constants.MAX_OFFSET
	}
	return limit, offset
}

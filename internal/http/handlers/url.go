package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kiryu-dev/shorty/internal/http/validator"
)

type URLShortener interface {
	MakeShort(context.Context, string) (string, error)
}

type createRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type createResponse struct {
	Alias string `json:"alias"`
}

func CreateShortURL(v *validator.RequestValidator, s URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(createRequest)
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		if err := v.Validate(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		resp := new(createResponse)
		resp.Alias, err = s.MakeShort(r.Context(), req.URL)
		w.Header().Add("Content-Type", "encoding/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}

package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kiryu-dev/shorty/internal/http/validator"
	"github.com/kiryu-dev/shorty/internal/model"
)

type URLShortener interface {
	MakeShort(context.Context, string) (string, error)
	GetURL(context.Context, string) (string, error)
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
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Add("Content-Type", "encoding/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}

func Redirect(s URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid request"))
			return
		}
		url, err := s.GetURL(r.Context(), alias)
		if err == model.ErrURLNotFound {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
}

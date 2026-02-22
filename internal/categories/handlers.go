package categories

import (
	"encoding/json"
	"net/http"

	"github.com/kylerjohnsondev/quiz-app-api/internal/httperror"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) FetchCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.FetchCategories(r.Context())
	if err != nil {
		if httpErr, ok := err.(*httperror.HTTPError); ok {
			if httpErr.ContentType != "" {
				w.Header().Set("Content-Type", httpErr.ContentType)
			}
			w.WriteHeader(httpErr.StatusCode)
			w.Write(httpErr.Body)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categories)
}

package questions

import (
	"encoding/json"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) FetchQuestions(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	difficulty := r.URL.Query().Get("difficulty")
	limit := r.URL.Query().Get("limit")
	questions, err := h.service.FetchQuestions(r.Context(), category, difficulty, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(questions)
}

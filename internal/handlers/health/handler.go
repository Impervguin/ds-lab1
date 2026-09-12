package health

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Impervguin/ds-lab1/internal/handlers/common"
	"github.com/Impervguin/ds-lab1/internal/handlers/health/dto"
	"github.com/Impervguin/ds-lab1/internal/service/health"
)

type Handler struct {
	service *health.Service
}

func NewHandler(service *health.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/health", h.Health)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if h.service.IsReady(r.Context()) {
		common.WriteJSON(w, http.StatusServiceUnavailable, dto.HealthResponse{Status: "DOWN"})
		return
	}

	common.WriteJSON(w, http.StatusOK, dto.HealthResponse{Status: "UP"})
}

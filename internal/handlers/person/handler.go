package person

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/Impervguin/ds-lab1/internal/domain"
	"github.com/Impervguin/ds-lab1/internal/handlers/common"
	"github.com/Impervguin/ds-lab1/internal/handlers/person/dto"
)

type Handler struct {
	repo     domain.PersonRepository
	validate *validator.Validate
}

func NewHandler(repo domain.PersonRepository) *Handler {
	return &Handler{
		repo:     repo,
		validate: validator.New(),
	}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/api/v1/persons", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.Get)
			r.Patch("/", h.Update)
			r.Delete("/", h.Delete)
		})
	})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	persons, err := h.repo.List(r.Context())
	if err != nil {
		log.Printf("list persons: %v", err)
		common.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]*dto.PersonResponse, 0, len(persons))
	for _, p := range persons {
		responses = append(responses, dto.PersonResponseFromDomain(p))
	}

	common.WriteJSON(w, http.StatusOK, responses)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := dto.DeserializePersonRequest(r)
	if err != nil {
		common.WriteValidationError(w, "invalid request body", nil)
		return
	}

	created, err := h.repo.Create(r.Context(), req.ToDomain())
	if err != nil {
		log.Printf("create person: %v", err)
		common.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/api/v1/persons/%d", created.ID))
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := common.ParseID(r)
	if err != nil {
		common.WriteError(w, http.StatusNotFound, "person not found")
		return
	}

	p, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPersonNotFound) {
			common.WriteError(w, http.StatusNotFound, "person not found")
			return
		}
		log.Printf("get person: %v", err)
		common.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	common.WriteJSON(w, http.StatusOK, dto.PersonResponseFromDomain(p))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := common.ParseID(r)
	if err != nil {
		common.WriteError(w, http.StatusNotFound, "person not found")
		return
	}

	req, err := dto.DeserializePersonPatchRequest(r)
	if err != nil {
		common.WriteValidationError(w, "invalid request body", nil)
		return
	}

	existing, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPersonNotFound) {
			common.WriteError(w, http.StatusNotFound, "person not found")
			return
		}
		log.Printf("get person: %v", err)
		common.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	req.ApplyTo(existing)

	updated, err := h.repo.Update(r.Context(), id, existing)
	if err != nil {
		if errors.Is(err, domain.ErrPersonNotFound) {
			common.WriteError(w, http.StatusNotFound, "person not found")
			return
		}
		log.Printf("update person: %v", err)
		common.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	common.WriteJSON(w, http.StatusOK, dto.PersonResponseFromDomain(updated))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := common.ParseID(r)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil && !errors.Is(err, domain.ErrPersonNotFound) {
		log.Printf("delete person: %v", err)
		common.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

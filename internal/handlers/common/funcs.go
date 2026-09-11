package common

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Impervguin/ds-lab1/internal/handlers/common/dto"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func ParseID(r *http.Request) (int32, error) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, dto.ErrorResponse{Message: message})
}

func WriteValidationError(w http.ResponseWriter, message string, errs map[string]string) {
	WriteJSON(w, http.StatusBadRequest, dto.ValidationErrorResponse{
		Message: message,
		Errors:  errs,
	})
}

func ValidationErrors(err error) map[string]string {
	errs := make(map[string]string)
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, fe := range validationErrs {
			errs[fe.Field()] = fe.Tag()
		}
	}
	return errs
}

package dto

import (
	"encoding/json"
	"net/http"

	"github.com/Impervguin/ds-lab1/internal/domain"
	"github.com/go-playground/validator/v10"
)

type PersonRequest struct {
	Name    string  `json:"name" validate:"required"`
	Age     int32   `json:"age" validate:"omitempty,gte=0"`
	Address *string `json:"address"`
	Work    *string `json:"work"`
}

func (r *PersonRequest) ToDomain() *domain.Person {
	address := ""
	if r.Address != nil {
		address = *r.Address
	}
	work := ""
	if r.Work != nil {
		work = *r.Work
	}
	return &domain.Person{
		Name:    r.Name,
		Age:     r.Age,
		Address: address,
		Work:    work,
	}
}

func DeserializePersonRequest(r *http.Request) (*PersonRequest, error) {
	var req PersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

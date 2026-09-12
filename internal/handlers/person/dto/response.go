package dto

import "github.com/Impervguin/ds-lab1/internal/domain"

type PersonResponse struct {
	ID      int32  `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Age     int32  `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

func PersonResponseFromDomain(p *domain.Person) *PersonResponse {
	return &PersonResponse{
		ID:      p.ID,
		Name:    p.Name,
		Age:     p.Age,
		Address: p.Address,
		Work:    p.Work,
	}
}

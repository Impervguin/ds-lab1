package dto

import "github.com/Impervguin/ds-lab1/internal/domain"

type PersonRequest struct {
	Name    string `json:"name" validate:"required"`
	Age     int32  `json:"age" validate:"omitempty,gte=0"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

func (r *PersonRequest) ToDomain() *domain.Person {
	return &domain.Person{
		Name:    r.Name,
		Age:     r.Age,
		Address: r.Address,
		Work:    r.Work,
	}
}

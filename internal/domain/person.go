package domain

import (
	"context"
	"errors"
)

type Person struct {
	ID      int32
	Name    string
	Age     int32
	Address string
	Work    string
}

var ErrPersonNotFound = errors.New("person not found")

type PersonRepository interface {
	List(ctx context.Context) ([]*Person, error)
	Create(ctx context.Context, p *Person) (*Person, error)
	GetByID(ctx context.Context, id int32) (*Person, error)
	Update(ctx context.Context, id int32, p *Person) (*Person, error)
	Delete(ctx context.Context, id int32) error
}

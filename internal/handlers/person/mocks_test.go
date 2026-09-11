package person_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/Impervguin/ds-lab1/internal/domain"
)

type mockPersonRepository struct {
	mock.Mock
}

func (m *mockPersonRepository) List(ctx context.Context) ([]*domain.Person, error) {
	args := m.Called(ctx)
	var persons []*domain.Person
	if v := args.Get(0); v != nil {
		persons = v.([]*domain.Person)
	}
	return persons, args.Error(1)
}

func (m *mockPersonRepository) Create(ctx context.Context, p *domain.Person) (*domain.Person, error) {
	args := m.Called(ctx, p)
	var created *domain.Person
	if v := args.Get(0); v != nil {
		created = v.(*domain.Person)
	}
	return created, args.Error(1)
}

func (m *mockPersonRepository) GetByID(ctx context.Context, id int32) (*domain.Person, error) {
	args := m.Called(ctx, id)
	var p *domain.Person
	if v := args.Get(0); v != nil {
		p = v.(*domain.Person)
	}
	return p, args.Error(1)
}

func (m *mockPersonRepository) Update(ctx context.Context, id int32, p *domain.Person) (*domain.Person, error) {
	args := m.Called(ctx, id, p)
	var updated *domain.Person
	if v := args.Get(0); v != nil {
		updated = v.(*domain.Person)
	}
	return updated, args.Error(1)
}

func (m *mockPersonRepository) Delete(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

var _ domain.PersonRepository = (*mockPersonRepository)(nil)

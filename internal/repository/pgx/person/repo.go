package person

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Impervguin/ds-lab1/internal/domain"
)

type PgxPersonRepository struct {
	pool *pgxpool.Pool
}

func NewPgxPersonRepository(pool *pgxpool.Pool) *PgxPersonRepository {
	return &PgxPersonRepository{pool: pool}
}

var _ domain.PersonRepository = (*PgxPersonRepository)(nil)

func (r *PgxPersonRepository) List(ctx context.Context) ([]*domain.Person, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, COALESCE(age, 0), COALESCE(address, ''), COALESCE(work, '') FROM persons ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query persons: %w", err)
	}
	defer rows.Close()

	persons := make([]*domain.Person, 0)
	for rows.Next() {
		p, err := scanPerson(rows)
		if err != nil {
			return nil, fmt.Errorf("scan person: %w", err)
		}
		persons = append(persons, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate persons: %w", err)
	}

	return persons, nil
}

func (r *PgxPersonRepository) Create(ctx context.Context, p *domain.Person) (*domain.Person, error) {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO persons (name, age, address, work) VALUES ($1, $2, $3, $4)
		 RETURNING id, name, COALESCE(age, 0), COALESCE(address, ''), COALESCE(work, '')`,
		p.Name, p.Age, p.Address, p.Work,
	)

	created, err := scanPerson(row)
	if err != nil {
		return nil, fmt.Errorf("insert person: %w", err)
	}

	return created, nil
}

func (r *PgxPersonRepository) GetByID(ctx context.Context, id int32) (*domain.Person, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, COALESCE(age, 0), COALESCE(address, ''), COALESCE(work, '') FROM persons WHERE id = $1`, id,
	)

	p, err := scanPerson(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPersonNotFound
		}
		return nil, fmt.Errorf("get person: %w", err)
	}

	return p, nil
}

func (r *PgxPersonRepository) Update(ctx context.Context, id int32, p *domain.Person) (*domain.Person, error) {
	row := r.pool.QueryRow(ctx,
		`UPDATE persons SET name = $1, age = $2, address = $3, work = $4 WHERE id = $5
		 RETURNING id, name, COALESCE(age, 0), COALESCE(address, ''), COALESCE(work, '')`,
		p.Name, p.Age, p.Address, p.Work, id,
	)

	updated, err := scanPerson(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPersonNotFound
		}
		return nil, fmt.Errorf("update person: %w", err)
	}

	return updated, nil
}

func (r *PgxPersonRepository) Delete(ctx context.Context, id int32) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM persons WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete person: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrPersonNotFound
	}

	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanPerson(row rowScanner) (*domain.Person, error) {
	var p domain.Person
	if err := row.Scan(&p.ID, &p.Name, &p.Age, &p.Address, &p.Work); err != nil {
		return nil, err
	}
	return &p, nil
}

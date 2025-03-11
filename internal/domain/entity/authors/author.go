package authors

import (
	"context"
	"fmt"
)

type Author struct {
	ID      int
	Name    string
	Country string
}

type AuthorRepo interface {
	Create(ctx context.Context, author *Author) error
	GetById(ctx context.Context, id int) (*Author, error)
	Update(ctx context.Context, author *Author) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]Author, error)
}

func NewAuthor(id int, name string, country string) (*Author, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	return &Author{
		ID:      id,
		Name:    name,
		Country: country,
	}, nil
}

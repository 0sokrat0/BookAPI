package authors

import (
	"context"
	"fmt"
)

type Author struct {
	ID      string
	Name    string
	Country string
}

type AuthorRepo interface {
	Create(ctx context.Context, author *Author) error
	GetById(ctx context.Context, id string) (*Author, error)
	Update(ctx context.Context, author *Author) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]Author, error)
}

func NewAuthor(id string, name string, country string) (*Author, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	return &Author{
		ID:      id,
		Name:    name,
		Country: country,
	}, nil
}

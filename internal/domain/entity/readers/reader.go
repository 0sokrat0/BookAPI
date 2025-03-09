package readers

import (
	"context"
	"fmt"
)

type Reader struct {
	ID    string
	Name  string
	Phone string
	Email string
}

type ReaderRepo interface {
	Create(ctx context.Context, reader *Reader) error
	GetById(ctx context.Context, id string) (*Reader, error)
	Update(ctx context.Context, reader *Reader) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]Reader, error)
}

func NewReader(id string, name string, phone string, email string) (*Reader, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	return &Reader{
		ID:    id,
		Name:  name,
		Phone: phone,
		Email: email,
	}, nil
}

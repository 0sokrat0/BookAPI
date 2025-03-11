package readers

import (
	"context"
	"fmt"
)

type Reader struct {
	ID       int
	Name     string
	Phone    string
	Email    string
	Password string
	Admin    bool
}

type ReaderRepo interface {
	Create(ctx context.Context, reader *Reader) error
	GetById(ctx context.Context, id int) (*Reader, error)
	Update(ctx context.Context, reader *Reader) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]Reader, error)
}

func NewReader(id int, name string, phone string, email string, password string, admin bool) (*Reader, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if password == "" {

	}
	return &Reader{
		ID:       id,
		Name:     name,
		Phone:    phone,
		Email:    email,
		Password: password,
		Admin:    admin,
	}, nil
}

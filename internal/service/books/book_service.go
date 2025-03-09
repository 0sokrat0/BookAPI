package books

import (
	"api/internal/application/commands"
	"api/internal/domain/aggregate/books"
	"context"
)

type BookService interface {
	CreateBook(ctx context.Context, req commands.CreateBookRequest) (*books.Book, error)
}

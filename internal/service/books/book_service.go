package books

import (
	"context"

	"github.com/0sokrat0/BookAPI/internal/application/commands"
	"github.com/0sokrat0/BookAPI/internal/domain/aggregate/books"
)

type BookService interface {
	CreateBook(ctx context.Context, req commands.CreateBookRequest) (*books.Book, error)
}

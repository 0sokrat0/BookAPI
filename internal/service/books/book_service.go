package books

import (
	"context"

	"github.com/0sokrat0/BookAPI/internal/application/commands"
	"github.com/0sokrat0/BookAPI/internal/domain/aggregate/books"
)

type BookService interface {
	CreateBook(ctx context.Context, req commands.CreateBookRequest) (*books.Book, error)
	GetBook(ctx context.Context, id int) (*books.Book, error)
	UpdateBook(ctx context.Context, id int, req commands.UpdateBookRequest) (*books.Book, error)
	DeleteBook(ctx context.Context, id int) error
	ListBooks(ctx context.Context) ([]books.Book, error)
	ListBooksByAuthor(ctx context.Context, authorID int) ([]books.Book, error)
}

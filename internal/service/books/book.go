package books

import (
	"context"
	"fmt"

	"github.com/0sokrat0/BookAPI/internal/application/commands"

	"github.com/0sokrat0/BookAPI/internal/domain/aggregate/books"
)

type bookService struct {
	bookRepo books.BookRepo
}

func NewBookService(repo books.BookRepo) BookService {
	return &bookService{bookRepo: repo}
}

func (s *bookService) CreateBook(ctx context.Context, req commands.CreateBookRequest) (*books.Book, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	// id = 0, если база сама генерирует его
	newBook, err := books.NewBook(0, req.Title, req.Year, req.ISBN, req.Genre, req.AuthorIDs)
	if err != nil {
		return nil, err
	}
	err = s.bookRepo.Create(ctx, newBook)
	if err != nil {
		return nil, err
	}
	return newBook, nil
}

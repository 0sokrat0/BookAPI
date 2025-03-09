package books

import (
	"api/internal/application/commands"
	"api/internal/domain/aggregate/books"
	"context"
	"fmt"

	"github.com/google/uuid"
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
	bookID := uuid.New().String()
	newBook, err := books.NewBook(bookID, req.Title, req.Year, req.ISBN, req.Genre, req.AuthorIDs)
	if err != nil {
		return nil, err
	}
	err = s.bookRepo.Create(ctx, newBook)
	if err != nil {
		return nil, err
	}
	return newBook, nil
}

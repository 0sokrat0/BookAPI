package books

import (
	"context"
	"fmt"

	"github.com/0sokrat0/BookAPI/internal/application/commands"
	"github.com/0sokrat0/BookAPI/internal/domain/aggregate/books"
	genid "github.com/0sokrat0/BookAPI/pkg/GenID"
)

type bookService struct {
	bookRepo  books.BookRepo
	idCounter *genid.IDcounter
}

func NewBookService(repo books.BookRepo, counter *genid.IDcounter) BookService {
	return &bookService{
		bookRepo:  repo,
		idCounter: counter,
	}
}

func (s *bookService) CreateBook(ctx context.Context, req commands.CreateBookRequest) (*books.Book, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	newID := s.idCounter.GenerateID()
	newBook, err := books.NewBook(newID, req.Title, req.Year, req.ISBN, req.Genre, req.AuthorIDs)
	if err != nil {
		return nil, err
	}
	err = s.bookRepo.Create(ctx, newBook)
	if err != nil {
		return nil, err
	}
	return newBook, nil
}

func (s *bookService) GetBook(ctx context.Context, id int) (*books.Book, error) {
	return s.bookRepo.GetByID(ctx, id)
}

func (s *bookService) UpdateBook(ctx context.Context, id int, req commands.UpdateBookRequest) (*books.Book, error) {
	// Получаем существующую книгу, чтобы обновить её
	existingBook, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Обновляем поля книги на основе запроса. Например:
	if req.Title != "" {
		existingBook.Title = req.Title
	}
	existingBook.Year = req.Year
	existingBook.ISBN = req.ISBN
	existingBook.Genre = req.Genre
	existingBook.SetAuthorIDs(req.AuthorIDs)
	// Вызываем репозиторий для сохранения изменений
	if err := s.bookRepo.Update(ctx, existingBook); err != nil {
		return nil, err
	}
	return existingBook, nil
}

func (s *bookService) DeleteBook(ctx context.Context, id int) error {
	return s.bookRepo.Delete(ctx, id)
}

func (s *bookService) ListBooks(ctx context.Context) ([]books.Book, error) {
	return s.bookRepo.List(ctx)
}

func (s *bookService) ListBooksByAuthor(ctx context.Context, authorID int) ([]books.Book, error) {
	return s.bookRepo.ListBooksByAuthor(ctx, authorID)
}

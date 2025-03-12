package booksRepo

import (
	"context"
	"fmt"

	"github.com/0sokrat0/BookAPI/internal/domain/aggregate/books"
	"github.com/0sokrat0/BookAPI/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type bookRepo struct {
	db *pgxpool.Pool
}

func NewBookRepo(db *pgxpool.Pool) books.BookRepo {
	return &bookRepo{db: db}
}

func (r *bookRepo) Create(ctx context.Context, book *books.Book) error {
	lg := logger.FromContext(ctx)
	query := `
        INSERT INTO books (id, title, year, isbn, genre)
        VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, book.ID, book.Title, book.Year, book.ISBN, book.Genre)
	if err != nil {
		lg.Error("failed to create book", zap.Error(err))
		return err
	}
	// Вставляем связи в таблицу book_authors.
	if err := r.insertBookAuthors(ctx, book.ID, book.AuthorIDs()); err != nil {
		lg.Error("failed to insert book authors", zap.Error(err))
		return err
	}
	return nil
}

func (r *bookRepo) insertBookAuthors(ctx context.Context, bookID int, authorIDs []int) error {
	query := `INSERT INTO book_authors (book_id, author_id) VALUES ($1, $2)`
	for _, authorID := range authorIDs {
		if _, err := r.db.Exec(ctx, query, bookID, authorID); err != nil {
			return fmt.Errorf("failed to insert book author (book_id=%d, author_id=%d): %w", bookID, authorID, err)
		}
	}
	return nil
}

func (r *bookRepo) loadBookAuthors(ctx context.Context, bookID int) ([]int, error) {
	query := `SELECT author_id FROM book_authors WHERE book_id = $1`
	rows, err := r.db.Query(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authorIDs []int
	for rows.Next() {
		var authorID int
		if err := rows.Scan(&authorID); err != nil {
			return nil, err
		}
		authorIDs = append(authorIDs, authorID)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return authorIDs, nil
}

func (r *bookRepo) updateBookAuthors(ctx context.Context, bookID int, authorIDs []int) error {
	delQuery := `DELETE FROM book_authors WHERE book_id = $1`
	if _, err := r.db.Exec(ctx, delQuery, bookID); err != nil {
		return fmt.Errorf("failed to delete old book authors: %w", err)
	}
	return r.insertBookAuthors(ctx, bookID, authorIDs)
}

func (r *bookRepo) deleteBookAuthors(ctx context.Context, bookID int) error {
	query := `DELETE FROM book_authors WHERE book_id = $1`
	_, err := r.db.Exec(ctx, query, bookID)
	return err
}

func (r *bookRepo) GetByID(ctx context.Context, id int) (*books.Book, error) {
	lg := logger.FromContext(ctx)
	query := `
		SELECT id, title, year, isbn, genre
		FROM books
		WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var book books.Book
	err := row.Scan(&book.ID, &book.Title, &book.Year, &book.ISBN, &book.Genre)
	if err != nil {
		lg.Error("failed to get book by id", zap.Error(err))
		return nil, err
	}
	authorIDs, err := r.loadBookAuthors(ctx, book.ID)
	if err != nil {
		lg.Error("failed to load book authors", zap.Error(err))
		return nil, err
	}
	book.SetAuthorIDs(authorIDs)
	return &book, nil
}

func (r *bookRepo) Update(ctx context.Context, book *books.Book) error {
	lg := logger.FromContext(ctx)
	query := `
		UPDATE books
		SET title = $1, year = $2, isbn = $3, genre = $4
		WHERE id = $5`
	_, err := r.db.Exec(ctx, query, book.Title, book.Year, book.ISBN, book.Genre, book.ID)
	if err != nil {
		lg.Error("failed to update book", zap.Error(err))
		return fmt.Errorf("failed to update book: %w", err)
	}
	if err := r.updateBookAuthors(ctx, book.ID, book.AuthorIDs()); err != nil {
		lg.Error("failed to update book authors", zap.Error(err))
		return err
	}
	return nil
}

func (r *bookRepo) Delete(ctx context.Context, id int) error {
	lg := logger.FromContext(ctx)
	// Удаляем связи из таблицы book_authors.
	if err := r.deleteBookAuthors(ctx, id); err != nil {
		lg.Error("failed to delete book authors", zap.Error(err))
		return fmt.Errorf("failed to delete book authors: %w", err)
	}
	query := `
		DELETE FROM books
		WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		lg.Error("failed to delete book", zap.Error(err))
		return fmt.Errorf("failed to delete book: %w", err)
	}
	return nil
}

func (r *bookRepo) List(ctx context.Context) ([]books.Book, error) {
	lg := logger.FromContext(ctx)
	query := `
		SELECT id, title, year, isbn, genre
		FROM books`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		lg.Error("failed to list books", zap.Error(err))
		return nil, fmt.Errorf("failed to list books: %w", err)
	}
	defer rows.Close()

	var booksList []books.Book
	for rows.Next() {
		var book books.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Year, &book.ISBN, &book.Genre)
		if err != nil {
			lg.Error("failed to scan book", zap.Error(err))
			return nil, fmt.Errorf("failed to scan book: %w", err)
		}
		authorIDs, err := r.loadBookAuthors(ctx, book.ID)
		if err != nil {
			lg.Error("failed to load book authors", zap.Error(err))
			return nil, err
		}
		book.SetAuthorIDs(authorIDs)
		booksList = append(booksList, book)
	}
	if err = rows.Err(); err != nil {
		lg.Error("rows error", zap.Error(err))
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return booksList, nil
}

func (r *bookRepo) ListBooksByAuthor(ctx context.Context, authorID int) ([]books.Book, error) {
	query := `
		SELECT b.id, b.title, b.year, b.isbn, b.genre
		FROM books b
		JOIN book_authors ba ON b.id = ba.book_id
		WHERE ba.author_id = $1`
	rows, err := r.db.Query(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to list books by author: %w", err)
	}
	defer rows.Close()

	var booksList []books.Book
	for rows.Next() {
		var book books.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Year, &book.ISBN, &book.Genre)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book: %w", err)
		}
		booksList = append(booksList, book)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return booksList, nil
}

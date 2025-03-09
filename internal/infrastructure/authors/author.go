package authors

import (
	"api/internal/domain/entity/authors"
	"api/pkg/logger"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type authorRepo struct {
	db *pgxpool.Pool
}

func NewAuthorRepo(db *pgxpool.Pool) authors.AuthorRepo {
	return &authorRepo{db: db}
}

func (r *authorRepo) Create(ctx context.Context, author *authors.Author) error {
	lg := logger.FromContext(ctx)
	query := `
	    INSERT INTO authors (id, name, country)
		VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, author.ID, author.Name, author.Country)
	if err != nil {
		lg.Error("failed to create author", zap.Error(err))
		return err
	}
	return nil
}

func (r *authorRepo) GetById(ctx context.Context, id string) (*authors.Author, error) {
	lg := logger.FromContext(ctx)
	query := `
	    SELECT id, name, country
		FROM authors
		WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var author authors.Author
	err := row.Scan(&author.ID, &author.Name, &author.Country)
	if err != nil {
		lg.Error("failed to get author by id", zap.Error(err))
		return nil, err
	}
	return &author, nil
}

func (r *authorRepo) Delete(ctx context.Context, id string) error {
	lg := logger.FromContext(ctx)
	query := `
	    DELETE FROM authors
		WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		lg.Error("failed to delete author by id", zap.Error(err))
		return err
	}
	return nil
}

func (r *authorRepo) Update(ctx context.Context, author *authors.Author) error {
	lg := logger.FromContext(ctx)
	query := `
        UPDATE authors
        SET name = $2, country = $3
        WHERE id = $1`
	_, err := r.db.Exec(ctx, query, author.ID, author.Name, author.Country)
	if err != nil {
		lg.Error("failed to update author by id", zap.Error(err))
		return err
	}
	return nil
}

func (r *authorRepo) List(ctx context.Context) ([]authors.Author, error) {
	lg := logger.FromContext(ctx)
	query := `
        SELECT id, name, country
        FROM authors`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		lg.Error("failed to list authors", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var authorsList []authors.Author
	for rows.Next() {
		var author authors.Author
		err := rows.Scan(&author.ID, &author.Name, &author.Country)
		if err != nil {
			lg.Error("failed to scan author", zap.Error(err))
			return nil, fmt.Errorf("failed to scan author: %w", err)
		}
		authorsList = append(authorsList, author)
	}
	if err = rows.Err(); err != nil {
		lg.Error("rows error", zap.Error(err))
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return authorsList, nil
}

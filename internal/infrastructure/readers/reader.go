package readers

import (
	domainReaders "api/internal/domain/entity/readers"
	"api/pkg/logger"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type readerRepo struct {
	db *pgxpool.Pool
}

// NewReaderRepo создает новый репозиторий для читателей.
func NewReaderRepo(db *pgxpool.Pool) domainReaders.ReaderRepo {
	return &readerRepo{db: db}
}

func (r *readerRepo) Create(ctx context.Context, reader *domainReaders.Reader) error {
	lg := logger.FromContext(ctx)
	query := `
	    INSERT INTO readers (id, name, phone, email)
		VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, reader.ID, reader.Name, reader.Phone, reader.Email)
	if err != nil {
		lg.Error("failed to create reader", zap.Error(err))
		return err
	}
	return nil
}

func (r *readerRepo) GetById(ctx context.Context, id int) (*domainReaders.Reader, error) {
	lg := logger.FromContext(ctx)
	query := `
        SELECT id, name, phone, email
        FROM readers
        WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var reader domainReaders.Reader
	err := row.Scan(&reader.ID, &reader.Name, &reader.Phone, &reader.Email)
	if err != nil {
		lg.Error("failed to get reader by id", zap.Error(err))
		return nil, err
	}
	return &reader, nil
}

func (r *readerRepo) Delete(ctx context.Context, id int) error {
	lg := logger.FromContext(ctx)
	query := `
        DELETE FROM readers
        WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		lg.Error("failed to delete reader by id", zap.Error(err))
		return err
	}
	return nil
}

func (r *readerRepo) Update(ctx context.Context, reader *domainReaders.Reader) error {
	lg := logger.FromContext(ctx)
	query := `
        UPDATE readers
        SET name = $2, phone = $3, email = $4
        WHERE id = $1`
	_, err := r.db.Exec(ctx, query, reader.ID, reader.Name, reader.Phone, reader.Email)
	if err != nil {
		lg.Error("failed to update reader by id", zap.Error(err))
		return err
	}
	return nil
}

func (r *readerRepo) List(ctx context.Context) ([]domainReaders.Reader, error) {
	lg := logger.FromContext(ctx)
	query := `
        SELECT id, name, phone, email
        FROM readers`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		lg.Error("failed to list readers", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var readersList []domainReaders.Reader
	for rows.Next() {
		var reader domainReaders.Reader
		err := rows.Scan(&reader.ID, &reader.Name, &reader.Phone, &reader.Email)
		if err != nil {
			lg.Error("failed to scan reader", zap.Error(err))
			return nil, fmt.Errorf("failed to scan reader: %w", err)
		}
		readersList = append(readersList, reader)
	}
	if err = rows.Err(); err != nil {
		lg.Error("rows error", zap.Error(err))
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return readersList, nil
}

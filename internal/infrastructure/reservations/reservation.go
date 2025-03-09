package reservations

import (
	"context"
	"fmt"
	"time"

	"api/internal/domain/aggregate/books"
	"api/internal/domain/aggregate/reservations"
	"api/internal/domain/entity/readers"

	"github.com/jackc/pgx/v5/pgxpool"
)

type reservationRepo struct {
	db *pgxpool.Pool
}

func NewReservationRepo(db *pgxpool.Pool) reservations.ReservationRepo {
	return &reservationRepo{db: db}
}

// Create сохраняет новое бронирование в базу данных.
func (r *reservationRepo) Create(ctx context.Context, id int, book books.Book, reader readers.Reader, startDate, endDate time.Time) (*reservations.Reservation, error) {
	// Создаем объект бронирования через доменную фабрику.
	res, err := reservations.NewReservation(id, book, reader, startDate, endDate)
	if err != nil {
		return nil, err
	}
	query := `
		INSERT INTO reservations (id, book_id, reader_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = r.db.Exec(ctx, query, res.ID, res.Book.ID, res.Reader.ID, res.StartDate, res.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}
	return res, nil
}

func (r *reservationRepo) GetById(ctx context.Context, id int) (*reservations.Reservation, error) {
	query := `
		SELECT id, book_id, reader_id, start_date, end_date
		FROM reservations
		WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var resID, bookID, readerID int
	var startDate, endDate time.Time
	if err := row.Scan(&resID, &bookID, &readerID, &startDate, &endDate); err != nil {
		return nil, fmt.Errorf("failed to get reservation by id: %w", err)
	}
	book := books.Book{ID: bookID}
	reader := readers.Reader{ID: readerID}
	res, err := reservations.NewReservation(resID, book, reader, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *reservationRepo) Update(ctx context.Context, id int, book books.Book, reader readers.Reader, startDate, endDate time.Time) error {
	if endDate.Before(startDate) {
		return fmt.Errorf("end date cannot be before start date")
	}
	query := `
		UPDATE reservations
		SET book_id = $1, reader_id = $2, start_date = $3, end_date = $4
		WHERE id = $5`
	_, err := r.db.Exec(ctx, query, book.ID, reader.ID, startDate, endDate, id)
	if err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}
	return nil
}

func (r *reservationRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM reservations WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}
	return nil
}

func (r *reservationRepo) List(ctx context.Context, startDate, endDate time.Time) ([]reservations.Reservation, error) {
	query := `
		SELECT id, book_id, reader_id, start_date, end_date
		FROM reservations
		WHERE start_date >= $1 AND end_date <= $2`
	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to list reservations: %w", err)
	}
	defer rows.Close()

	var resList []reservations.Reservation
	for rows.Next() {
		var id, bookID, readerID int
		var sDate, eDate time.Time
		err := rows.Scan(&id, &bookID, &readerID, &sDate, &eDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservation: %w", err)
		}
		book := books.Book{ID: bookID}
		reader := readers.Reader{ID: readerID}
		res, err := reservations.NewReservation(id, book, reader, sDate, eDate)
		if err != nil {
			return nil, err
		}
		resList = append(resList, *res)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return resList, nil
}

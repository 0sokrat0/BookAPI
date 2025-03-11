package reservations

import (
	"context"
	"fmt"
	"time"

	"github.com/0sokrat0/BookAPI/internal/domain/aggregate/books"
	"github.com/0sokrat0/BookAPI/internal/domain/entity/readers"
)

type Reservation struct {
	ID        int
	Book      books.Book
	Reader    readers.Reader
	StartDate time.Time
	EndDate   time.Time
}

type ReservationRepo interface {
	Create(ctx context.Context, id int, book books.Book, reader readers.Reader, startDate, endDate time.Time) (*Reservation, error)
	GetById(ctx context.Context, id int) (*Reservation, error)
	Update(ctx context.Context, id int, book books.Book, reader readers.Reader, startDate, endDate time.Time) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, startDate, endDate time.Time) ([]Reservation, error)
}

func NewReservation(id int, book books.Book, reader readers.Reader, startDate, endDate time.Time) (*Reservation, error) {
	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}
	return &Reservation{
		ID:        id,
		Book:      book,
		Reader:    reader,
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

package reservations

import (
	"api/internal/domain/aggregate/books"
	"api/internal/domain/entity/readers"
	"context"
	"fmt"
	"time"
)

type Reservation struct {
	ID        string
	Book      books.Book
	Reader    readers.Reader
	StartDate time.Time
	EndDate   time.Time
}

type ReservationRepo interface {
	Create(ctx context.Context, id string, book books.Book, reader readers.Reader, startDate, endDate time.Time) (*Reservation, error)
	GetById(ctx context.Context, id string) (*Reservation, error)
	Update(ctx context.Context, id string, book books.Book, reader readers.Reader, startDate, endDate time.Time) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, startDate, endDate time.Time) ([]Reservation, error)
}

func NewReservation(id string, book books.Book, reader readers.Reader, startDate, endDate time.Time) (*Reservation, error) {
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

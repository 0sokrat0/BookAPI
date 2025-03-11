package reservations

import (
	"context"
	"fmt"
	"time"

	"github.com/0sokrat0/BookAPI/internal/domain/aggregate/books"
	"github.com/0sokrat0/BookAPI/internal/domain/aggregate/reservations"
	"github.com/0sokrat0/BookAPI/internal/domain/entity/readers"
)

// ReservationService определяет интерфейс сервиса бронирований.
type ReservationService interface {
	CreateReservation(ctx context.Context, req CreateReservationRequest) (*reservations.Reservation, error)
	GetReservationByID(ctx context.Context, id int) (*reservations.Reservation, error)
	UpdateReservation(ctx context.Context, req UpdateReservationRequest) error
	DeleteReservation(ctx context.Context, id int) error
	ListReservations(ctx context.Context, startDate, endDate time.Time) ([]reservations.Reservation, error)
}

// reservationService — реализация сервиса бронирований.
type reservationService struct {
	repo reservations.ReservationRepo
}

// NewReservationService создаёт новый сервис бронирований.
func NewReservationService(repo reservations.ReservationRepo) ReservationService {
	return &reservationService{repo: repo}
}

// CreateReservationRequest содержит данные для создания бронирования.
type CreateReservationRequest struct {
	ID        int
	Book      books.Book
	Reader    readers.Reader
	StartDate time.Time
	EndDate   time.Time
}

// UpdateReservationRequest содержит данные для обновления бронирования.
type UpdateReservationRequest struct {
	ID        int
	Book      books.Book
	Reader    readers.Reader
	StartDate time.Time
	EndDate   time.Time
}

func (s *reservationService) CreateReservation(ctx context.Context, req CreateReservationRequest) (*reservations.Reservation, error) {
	// Проверка бизнес-правил может быть добавлена здесь.
	if req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	// Создаём агрегат бронирования через доменную фабрику.
	res, err := reservations.NewReservation(req.ID, req.Book, req.Reader, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	// Сохраняем бронирование через репозиторий.
	return s.repo.Create(ctx, res.ID, req.Book, req.Reader, req.StartDate, req.EndDate)
}

func (s *reservationService) GetReservationByID(ctx context.Context, id int) (*reservations.Reservation, error) {
	return s.repo.GetById(ctx, id)
}

func (s *reservationService) UpdateReservation(ctx context.Context, req UpdateReservationRequest) error {
	if req.EndDate.Before(req.StartDate) {
		return fmt.Errorf("end date cannot be before start date")
	}
	return s.repo.Update(ctx, req.ID, req.Book, req.Reader, req.StartDate, req.EndDate)
}

func (s *reservationService) DeleteReservation(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *reservationService) ListReservations(ctx context.Context, startDate, endDate time.Time) ([]reservations.Reservation, error) {
	return s.repo.List(ctx, startDate, endDate)
}

package commands

import "time"

// CreateReservationRequestDTO содержит данные для создания бронирования.
type CreateReservationRequestDTO struct {
	ID        int       `json:"id" example:"0"`
	BookID    int       `json:"book_id" example:"1"`
	ReaderID  int       `json:"reader_id" example:"2"`
	StartDate time.Time `json:"start_date" example:"2025-03-15"`
	EndDate   time.Time `json:"end_date" example:"2025-03-20"`
}

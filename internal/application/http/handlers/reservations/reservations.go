package reservations

import (
	"time"

	"github.com/0sokrat0/BookAPI/internal/domain/aggregate/books"
	"github.com/0sokrat0/BookAPI/internal/domain/entity/readers"
	"github.com/0sokrat0/BookAPI/internal/service/reservations"
	"github.com/0sokrat0/BookAPI/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type CreateReservationRequestDTO struct {
	ID        int       `json:"id"`         // Если ID генерируется базой, можно опустить
	BookID    int       `json:"book_id"`    // Идентификатор книги
	ReaderID  int       `json:"reader_id"`  // Идентификатор читателя
	StartDate time.Time `json:"start_date"` // Начало бронирования
	EndDate   time.Time `json:"end_date"`   // Окончание бронирования
}

type UpdateReservationRequestDTO struct {
	ID        int       `json:"id"`
	BookID    int       `json:"book_id"`
	ReaderID  int       `json:"reader_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type Handler struct {
	reservationService reservations.ReservationService
}

func NewHandler(s reservations.ReservationService) *Handler {
	return &Handler{reservationService: s}
}

// CreateReservationHandler godoc
// @Summary      Create a new reservation
// @Description  Создаёт новое бронирование в системе.
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        request  body      CreateReservationRequestDTO  true  "Reservation creation request"
// @Success      200      {object}  response.BaseResponse "Бронирование создано успешно"
// @Failure      400      {object}  response.ErrorResponse  "Invalid request"
// @Failure      500      {object}  response.ErrorResponse  "Internal server error"
// @Router       /reservation [post]
func (h *Handler) CreateReservationHandler(c *fiber.Ctx) error {
	var req CreateReservationRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request",
		})
	}
	book := books.Book{ID: req.BookID}
	reader := readers.Reader{ID: req.ReaderID}
	serviceReq := reservations.CreateReservationRequest{
		ID:        req.ID,
		Book:      book,
		Reader:    reader,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
	reservation, err := h.reservationService.CreateReservation(c.UserContext(), serviceReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Reservation created successfully",
		Data:    reservation,
	})
}

// GetReservationHandler godoc
// @Summary      Get reservation by ID
// @Description  Возвращает бронирование по его уникальному идентификатору.
// @Tags         reservations
// @Produce      json
// @Param        id   path      int  true  "Reservation ID"
// @Success      200  {object}  response.BaseResponse "Данные бронирования"
// @Failure      400  {object}  response.ErrorResponse  "Invalid reservation ID"
// @Failure      404  {object}  response.ErrorResponse  "Not found"
// @Router       /reservation/{id} [get]
func (h *Handler) GetReservationHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid reservation ID",
		})
	}
	reservation, err := h.reservationService.GetReservationByID(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
			Code:    fiber.StatusNotFound,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Reservation retrieved successfully",
		Data:    reservation,
	})
}

// UpdateReservationHandler godoc
// @Summary      Update reservation
// @Description  Обновляет данные существующего бронирования.
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        request  body      UpdateReservationRequestDTO  true  "Reservation update request"
// @Success      200      {object}  response.BaseResponse "Бронирование обновлено успешно"
// @Failure      400      {object}  response.ErrorResponse  "Invalid request"
// @Failure      500      {object}  response.ErrorResponse  "Internal server error"
// @Router       /reservation [put]
func (h *Handler) UpdateReservationHandler(c *fiber.Ctx) error {
	var req UpdateReservationRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request",
		})
	}
	book := books.Book{ID: req.BookID}
	reader := readers.Reader{ID: req.ReaderID}
	serviceReq := reservations.UpdateReservationRequest{
		ID:        req.ID,
		Book:      book,
		Reader:    reader,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
	if err := h.reservationService.UpdateReservation(c.UserContext(), serviceReq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Reservation updated successfully",
	})
}

// DeleteReservationHandler godoc
// @Summary      Delete reservation
// @Description  Удаляет бронирование по его идентификатору.
// @Tags         reservations
// @Produce      json
// @Param        id   path      int  true  "Reservation ID"
// @Success      200  {object}  response.BaseResponse "Бронирование удалено успешно"
// @Failure      400  {object}  response.ErrorResponse  "Invalid ID"
// @Failure      500  {object}  response.ErrorResponse  "Internal server error"
// @Router       /reservation/{id} [delete]
func (h *Handler) DeleteReservationHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid reservation ID",
		})
	}
	if err := h.reservationService.DeleteReservation(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Reservation deleted successfully",
	})
}

// ListReservationsHandler godoc
// @Summary      List reservations
// @Description  Возвращает список бронирований в указанном диапазоне дат.
// @Tags         reservations
// @Produce      json
// @Param        startDate  query     string  true  "Start date (YYYY-MM-DD)"
// @Param        endDate    query     string  true  "End date (YYYY-MM-DD)"
// @Success      200      {object}  response.BaseResponse "Список бронирований"
// @Failure      400      {object}  response.ErrorResponse  "Invalid request"
// @Failure      500      {object}  response.ErrorResponse  "Internal server error"
// @Router       /reservations [get]
func (h *Handler) ListReservationsHandler(c *fiber.Ctx) error {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid startDate format",
		})
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid endDate format",
		})
	}
	resList, err := h.reservationService.ListReservations(c.UserContext(), startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Reservations list retrieved successfully",
		Data:    resList,
	})
}

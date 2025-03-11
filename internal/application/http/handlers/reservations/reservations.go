package reservations

import (
	"time"

	"github.com/0sokrat0/BookAPI/internal/domain/aggregate/books"
	"github.com/0sokrat0/BookAPI/internal/domain/entity/readers"
	service "github.com/0sokrat0/BookAPI/internal/service/reservations"
	"github.com/gofiber/fiber/v2"
)

// Handler представляет обработчик для операций с бронированиями.
type Handler struct {
	reservationService service.ReservationService
}

// NewHandler создаёт новый обработчик бронирований.
func NewHandler(s service.ReservationService) *Handler {
	return &Handler{reservationService: s}
}

// CreateReservationRequestDTO — структура для приема запроса через HTTP.
type CreateReservationRequestDTO struct {
	ID        int       `json:"id"`         // Если ID генерируется базой, можно опустить
	BookID    int       `json:"book_id"`    // Идентификатор книги
	ReaderID  int       `json:"reader_id"`  // Идентификатор читателя
	StartDate time.Time `json:"start_date"` // Начало бронирования
	EndDate   time.Time `json:"end_date"`   // Окончание бронирования
}

// CreateReservationHandler godoc
// @Summary      Create a new reservation
// @Description  Создаёт новое бронирование в системе.
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        request  body      CreateReservationRequestDTO  true  "Reservation creation request"
// @Failure      400      {object}  map[string]string  "Invalid request"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /reservation [post]
func (h *Handler) CreateReservationHandler(c *fiber.Ctx) error {
	var req CreateReservationRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	// Создаем минимальные объекты Book и Reader с заданными ID.
	book := books.Book{ID: req.BookID}
	reader := readers.Reader{ID: req.ReaderID}
	// Формируем запрос для сервиса.
	serviceReq := service.CreateReservationRequest{
		ID:        req.ID,
		Book:      book,
		Reader:    reader,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
	reservation, err := h.reservationService.CreateReservation(c.UserContext(), serviceReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(reservation)
}

// GetReservationHandler godoc
// @Summary      Get reservation by ID
// @Description  Возвращает бронирование по его уникальному идентификатору.
// @Tags         reservations
// @Produce      json
// @Param        id   path      int  true  "Reservation ID"
// @Failure      400  {object}  map[string]string  "Invalid reservation ID"
// @Failure      404  {object}  map[string]string  "Not found"
// @Router       /reservation/{id} [get]
func (h *Handler) GetReservationHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid reservation ID"})
	}
	reservation, err := h.reservationService.GetReservationByID(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(reservation)
}

// UpdateReservationRequestDTO — структура для обновления бронирования.
type UpdateReservationRequestDTO struct {
	ID        int       `json:"id"`
	BookID    int       `json:"book_id"`
	ReaderID  int       `json:"reader_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// UpdateReservationHandler godoc
// @Summary      Update reservation
// @Description  Обновляет данные существующего бронирования.
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        request  body      UpdateReservationRequestDTO  true  "Reservation update request"
// @Success      200      {object}  map[string]string  "Updated successfully"
// @Failure      400      {object}  map[string]string  "Invalid request"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /reservation [put]
func (h *Handler) UpdateReservationHandler(c *fiber.Ctx) error {
	var req UpdateReservationRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	book := books.Book{ID: req.BookID}
	reader := readers.Reader{ID: req.ReaderID}
	serviceReq := service.UpdateReservationRequest{
		ID:        req.ID,
		Book:      book,
		Reader:    reader,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
	err := h.reservationService.UpdateReservation(c.UserContext(), serviceReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Updated successfully"})
}

// DeleteReservationHandler godoc
// @Summary      Delete reservation
// @Description  Удаляет бронирование по его идентификатору.
// @Tags         reservations
// @Produce      json
// @Param        id   path      int  true  "Reservation ID"
// @Success      200  {object}  map[string]string  "Deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid ID"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /reservation/{id} [delete]
func (h *Handler) DeleteReservationHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid reservation ID"})
	}
	if err := h.reservationService.DeleteReservation(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Deleted successfully"})
}

// ListReservationsHandler godoc
// @Summary      List reservations
// @Description  Возвращает список бронирований в указанном диапазоне дат.
// @Tags         reservations
// @Produce      json
// @Param        startDate  query     string  true  "Start date (YYYY-MM-DD)"
// @Param        endDate    query     string  true  "End date (YYYY-MM-DD)"
// @Failure      400        {object}  map[string]string  "Invalid request"
// @Failure      500        {object}  map[string]string  "Internal server error"
// @Router       /reservations [get]
func (h *Handler) ListReservationsHandler(c *fiber.Ctx) error {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid startDate format"})
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid endDate format"})
	}
	resList, err := h.reservationService.ListReservations(c.UserContext(), startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resList)
}

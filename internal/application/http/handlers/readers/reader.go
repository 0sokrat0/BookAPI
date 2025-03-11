package readerhandlers

import (
	"strconv"

	"github.com/0sokrat0/BookAPI/internal/application/commands"
	"github.com/0sokrat0/BookAPI/internal/service/readers"
	"github.com/gofiber/fiber/v2"
)

// CreateReaderRequest содержит данные для создания читателя.
// swagger:model CreateReaderRequest
type CreateReaderRequest struct {
	Name     string `json:"name" example:"Ivan Ivanov"`
	Phone    string `json:"phone" example:"+79111234567"`
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"password123"`
	Admin    bool   `json:"admin" example:"false"`
}

// UpdateReaderRequest содержит данные для обновления читателя.
// swagger:model UpdateReaderRequest
type UpdateReaderRequest struct {
	Name     string `json:"name" example:"Ivan Ivanov"`
	Phone    string `json:"phone" example:"+79111234567"`
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"newpassword"`
	Admin    bool   `json:"admin" example:"false"`
}

// Handler представляет обработчик для операций с читателями.
type Handler struct {
	readerService readers.ReaderService
}

// NewHandler создаёт новый обработчик для читателей.
func NewHandler(service readers.ReaderService) *Handler {
	return &Handler{readerService: service}
}

// CreateReaderHandler godoc
// @Summary      Create a new reader
// @Description  Создаёт нового читателя с предоставленными данными.
// @Tags         readers
// @Accept       json
// @Produce      json
// @Param        reader  body      readerhandlers.CreateReaderRequest  true  "Параметры для создания читателя. Пример: {\"name\":\"Ivan Ivanov\", \"phone\":\"+79111234567\", \"email\":\"ivan@example.com\", \"password\":\"password123\", \"admin\":false}"
// @Success      200     {object}  map[string]interface{}  "Созданный читатель с уникальным ID"
// @Failure      400     {object}  map[string]string       "Неверный запрос"
// @Failure      500     {object}  map[string]string       "Ошибка сервера"
// @Router       /reader [post]
func (h *Handler) CreateReaderHandler(c *fiber.Ctx) error {
	var req CreateReaderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request: " + err.Error()})
	}

	// Преобразуем локальный тип в тип из пакета commands
	cmdReq := commands.CreateReaderRequest{
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
		Admin:    req.Admin,
	}

	reader, err := h.readerService.CreateReader(c.UserContext(), cmdReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(reader)
}

// GetReaderHandler godoc
// @Summary      Get a reader by ID
// @Description  Возвращает данные читателя по его уникальному идентификатору.
// @Tags         readers
// @Produce      json
// @Param        id   path      int  true  "Уникальный ID читателя"
// @Success      200  {object}  map[string]interface{}  "Данные читателя"
// @Failure      400  {object}  map[string]string       "Неверный ID"
// @Failure      404  {object}  map[string]string       "Читатель не найден"
// @Router       /reader/{id} [get]
func (h *Handler) GetReaderHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid reader ID"})
	}
	reader, err := h.readerService.GetReader(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Reader not found"})
	}
	return c.Status(fiber.StatusOK).JSON(reader)
}

// UpdateReaderHandler godoc
// @Summary      Update a reader
// @Description  Обновляет данные существующего читателя.
// @Tags         readers
// @Accept       json
// @Produce      json
// @Param        id      path      int  true  "Уникальный ID читателя"
// @Param        reader  body      readerhandlers.UpdateReaderRequest  true  "Новые данные читателя. Пример: {\"name\":\"Ivan Ivanov\", \"phone\":\"+79111234567\", \"email\":\"ivan@example.com\", \"password\":\"newpassword\", \"admin\":false}"
// @Success      200     {object}  map[string]interface{}  "Обновлённые данные читателя"
// @Failure      400     {object}  map[string]string       "Неверный запрос"
// @Failure      500     {object}  map[string]string       "Ошибка сервера"
// @Router       /reader/{id} [put]
func (h *Handler) UpdateReaderHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid reader ID"})
	}
	var req UpdateReaderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request: " + err.Error()})
	}

	// Преобразуем локальный тип в тип из пакета commands
	cmdReq := commands.UpdateReaderRequest{
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
		Admin:    req.Admin,
	}

	updatedReader, err := h.readerService.UpdateReader(c.UserContext(), id, cmdReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(updatedReader)
}

// DeleteReaderHandler godoc
// @Summary      Delete a reader
// @Description  Удаляет читателя по его уникальному идентификатору.
// @Tags         readers
// @Produce      json
// @Param        id   path      int  true  "Уникальный ID читателя"
// @Success      200  {object}  map[string]string  "Читатель успешно удалён"
// @Failure      400  {object}  map[string]string  "Неверный ID"
// @Failure      500  {object}  map[string]string  "Ошибка сервера"
// @Router       /reader/{id} [delete]
func (h *Handler) DeleteReaderHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid reader ID"})
	}
	if err := h.readerService.DeleteReader(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Reader deleted successfully"})
}

// ListReadersHandler godoc
// @Summary      List all readers
// @Description  Возвращает список всех читателей.
// @Tags         readers
// @Produce      json
// @Success      200  {array}   map[string]interface{}  "Список читателей"
// @Failure      500  {object}  map[string]string       "Ошибка сервера"
// @Router       /readers [get]
func (h *Handler) ListReadersHandler(c *fiber.Ctx) error {
	readersList, err := h.readerService.ListReaders(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(readersList)
}

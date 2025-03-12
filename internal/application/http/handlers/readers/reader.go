package readerhandlers

import (
	"strconv"

	"github.com/0sokrat0/BookAPI/internal/application/commands"
	"github.com/0sokrat0/BookAPI/internal/service/readers"
	"github.com/0sokrat0/BookAPI/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type CreateReaderRequest struct {
	Name     string `json:"name" example:"Ivan Ivanov"`
	Phone    string `json:"phone" example:"+79111234567"`
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"password123"`
	Admin    bool   `json:"admin" example:"false"`
}

type UpdateReaderRequest struct {
	Name     string `json:"name" example:"Ivan Ivanov"`
	Phone    string `json:"phone" example:"+79111234567"`
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"newpassword"`
	Admin    bool   `json:"admin" example:"false"`
}

// LoginRequest содержит данные для аутентификации пользователя.
// swagger:model LoginRequest
type LoginRequest struct {
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"password123"`
}

type Handler struct {
	readerService readers.ReaderService
}

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
// @Success      200     {object}  response.BaseResponse "Созданный читатель с уникальным ID"
// @Failure      400     {object}  response.ErrorResponse  "Неверный запрос"
// @Failure      500     {object}  response.ErrorResponse  "Ошибка сервера"
// @Router       /reader [post]
func (h *Handler) CreateReaderHandler(c *fiber.Ctx) error {
	var req CreateReaderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
	}

	cmdReq := commands.CreateReaderRequest{
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
		Admin:    req.Admin,
	}
	reader, err := h.readerService.CreateReader(c.UserContext(), cmdReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Reader created successfully",
		Data:    reader,
	})
}

// GetReaderHandler godoc
// @Summary      Get a reader by ID
// @Description  Возвращает данные читателя по его уникальному идентификатору.
// @Tags         readers
// @Produce      json
// @Param        id   path      int  true  "Уникальный ID читателя"
// @Success      200  {object}  response.BaseResponse "Данные читателя"
// @Failure      400  {object}  response.ErrorResponse  "Неверный ID"
// @Failure      404  {object}  response.ErrorResponse  "Читатель не найден"
// @Router       /reader/{id} [get]
func (h *Handler) GetReaderHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid reader ID",
		})
	}
	reader, err := h.readerService.GetReader(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
			Code:    fiber.StatusNotFound,
			Message: "Reader not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Reader retrieved successfully",
		Data:    reader,
	})
}

// UpdateReaderHandler godoc
// @Summary      Update a reader
// @Description  Обновляет данные существующего читателя.
// @Tags         readers
// @Accept       json
// @Produce      json
// @Param        id      path      int  true  "Уникальный ID читателя"
// @Param        reader  body      readerhandlers.UpdateReaderRequest  true  "Новые данные читателя. Пример: {\"name\":\"Ivan Ivanov\", \"phone\":\"+79111234567\", \"email\":\"ivan@example.com\", \"password\":\"newpassword\", \"admin\":false}"
// @Success      200     {object}  response.BaseResponse "Обновлённые данные читателя"
// @Failure      400     {object}  response.ErrorResponse  "Неверный запрос"
// @Failure      500     {object}  response.ErrorResponse  "Ошибка сервера"
// @Router       /reader/{id} [put]
func (h *Handler) UpdateReaderHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid reader ID",
		})
	}
	var req UpdateReaderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
	}
	cmdReq := commands.UpdateReaderRequest{
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
		Admin:    req.Admin,
	}
	updatedReader, err := h.readerService.UpdateReader(c.UserContext(), id, cmdReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Reader updated successfully",
		Data:    updatedReader,
	})
}

// DeleteReaderHandler godoc
// @Summary      Delete a reader
// @Description  Удаляет читателя по его уникальному идентификатору.
// @Tags         readers
// @Produce      json
// @Param        id   path      int  true  "Уникальный ID читателя"
// @Success      200  {object}  response.BaseResponse "Читатель успешно удалён"
// @Failure      400  {object}  response.ErrorResponse  "Неверный ID"
// @Failure      500  {object}  response.ErrorResponse  "Ошибка сервера"
// @Router       /reader/{id} [delete]
func (h *Handler) DeleteReaderHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid reader ID",
		})
	}
	if err := h.readerService.DeleteReader(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Reader deleted successfully",
	})
}

// ListReadersHandler godoc
// @Summary      List all readers
// @Description  Возвращает список всех читателей.
// @Tags         readers
// @Produce      json
// @Success      200  {object}  response.BaseResponse "Список читателей"
// @Failure      500  {object}  response.ErrorResponse  "Ошибка сервера"
// @Router       /readers [get]
func (h *Handler) ListReadersHandler(c *fiber.Ctx) error {
	readersList, err := h.readerService.ListReaders(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Readers list retrieved successfully",
		Data:    readersList,
	})
}

func (h *Handler) GetReadersByEmailHandler(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Email is required",
		})
	}
	reader, err := h.readerService.GetReaderByEmail(c.UserContext(), email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Readers retrieved successfully",
		Data:    reader,
	})
}

// AuthenticateReaderHandler godoc
// @Summary      Authenticate reader
// @Description  Аутентифицирует пользователя по email и паролю. При неверном пароле возвращает ошибку Unauthorized.
// @Tags         readers
// @Accept       json
// @Produce      json
// @Param        credentials  body      readerhandlers.LoginRequest  true  "Данные для аутентификации. Пример: {\"email\":\"ivan@example.com\", \"password\":\"password123\"}"
// @Success      200          {object}  response.BaseResponse  "Успешная аутентификация: данные пользователя"
// @Failure      400          {object}  response.ErrorResponse "Неверный формат запроса"
// @Failure      401          {object}  response.ErrorResponse "Неверный пароль или пользователь не найден"
// @Failure      500          {object}  response.ErrorResponse "Ошибка сервера"
// @Router       /login [post]
func (h *Handler) AuthenticateReaderHandler(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
	}

	reader, err := h.readerService.Authenticate(c.UserContext(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Authentication successful",
		Data: map[string]interface{}{
			"id":    reader.ID,
			"name":  reader.Name,
			"email": reader.Email,
			"admin": reader.Admin,
		},
	})
}

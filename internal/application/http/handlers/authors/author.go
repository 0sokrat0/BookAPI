package authors

import (
	"strconv"

	"github.com/0sokrat0/BookAPI/internal/application/commands"
	"github.com/0sokrat0/BookAPI/internal/service/authors"
	"github.com/0sokrat0/BookAPI/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// Локальные DTO, определённые в данном файле для Swagger

// CreateAuthorRequest содержит данные для создания автора.
// swagger:model CreateAuthorRequest
type CreateAuthorRequest struct {
	Name    string `json:"name" example:"Leo Tolstoy"`
	Country string `json:"country" example:"Russia"`
}

// UpdateAuthorRequest содержит данные для обновления автора.
// swagger:model UpdateAuthorRequest
type UpdateAuthorRequest struct {
	Name    string `json:"name" example:"Anton Chekhov"`
	Country string `json:"country" example:"Russia"`
}

// Handler представляет обработчик для операций с авторами.
type Handler struct {
	authorService authors.AuthorService
}

// NewHandler создаёт новый обработчик для авторов.
func NewHandler(service authors.AuthorService) *Handler {
	return &Handler{authorService: service}
}

// CreateAuthorHandler godoc
// @Summary      Create a new author
// @Description  Создаёт нового автора с указанными данными. Принимает JSON-представление автора и возвращает созданную запись.
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        author  body      authors.CreateAuthorRequest  true  "Параметры для создания автора. Пример: {\"name\":\"Leo Tolstoy\", \"country\":\"Russia\"}"
// @Success      200     {object}  map[string]interface{}  "Новый автор с уникальным ID"
// @Failure      400     {object}  map[string]string       "Неверный запрос"
// @Failure      500     {object}  map[string]string       "Ошибка сервера"
// @Router       /author [post]
func (h *Handler) CreateAuthorHandler(c *fiber.Ctx) error {
	var req CreateAuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
	}
	cmdReq := commands.CreateAuthorRequest{
		Name:    req.Name,
		Country: req.Country,
	}
	author, err := h.authorService.CreateAuthor(c.UserContext(), cmdReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Author created successfully",
		Data:    author,
	})
}

// GetAuthorHandler godoc
// @Summary      Get an author by ID
// @Description  Возвращает данные автора по его уникальному идентификатору.
// @Tags         authors
// @Produce      json
// @Param        id   path      int  true  "Уникальный ID автора"
// @Success      200  {object}  map[string]interface{}  "Информация об авторе"
// @Failure      400  {object}  map[string]string       "Неверный ID"
// @Failure      404  {object}  map[string]string       "Автор не найден"
// @Router       /author/{id} [get]
func (h *Handler) GetAuthorHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid author ID",
		})
	}
	author, err := h.authorService.GetAuthor(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
			Code:    fiber.StatusNotFound,
			Message: "Author not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Author retrieved successfully",
		Data:    author,
	})
}

// UpdateAuthorHandler godoc
// @Summary      Update an author
// @Description  Обновляет данные существующего автора по его уникальному идентификатору.
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        id      path      int  true  "Уникальный ID автора"
// @Param        author  body      authors.UpdateAuthorRequest  true  "Новые данные автора. Пример: {\"name\":\"Anton Chekhov\", \"country\":\"Russia\"}"
// @Success      200     {object}  map[string]interface{}  "Обновлённые данные автора"
// @Failure      400     {object}  map[string]string       "Неверный запрос или ID"
// @Failure      500     {object}  map[string]string       "Ошибка сервера"
// @Router       /author/{id} [put]
func (h *Handler) UpdateAuthorHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid author ID",
		})
	}
	var req UpdateAuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
	}
	cmdReq := commands.UpdateAuthorRequest{
		Name:    req.Name,
		Country: req.Country,
	}
	updatedAuthor, err := h.authorService.UpdateAuthor(c.UserContext(), id, cmdReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Author updated successfully",
		Data:    updatedAuthor,
	})
}

// DeleteAuthorHandler godoc
// @Summary      Delete an author
// @Description  Удаляет автора по его уникальному идентификатору.
// @Tags         authors
// @Produce      json
// @Param        id   path      int  true  "Уникальный ID автора"
// @Success      200  {object}  map[string]string  "Автор успешно удалён"
// @Failure      400  {object}  map[string]string  "Неверный ID"
// @Failure      500  {object}  map[string]string  "Ошибка сервера"
// @Router       /author/{id} [delete]
func (h *Handler) DeleteAuthorHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid author ID",
		})
	}
	if err := h.authorService.DeleteAuthor(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Author deleted successfully",
	})
}

// ListAuthorsHandler godoc
// @Summary      List all authors
// @Description  Возвращает список всех авторов, зарегистрированных в системе.
// @Tags         authors
// @Produce      json
// @Success      200  {array}   map[string]interface{}  "Массив объектов авторов"
// @Failure      500  {object}  map[string]string       "Ошибка сервера"
// @Router       /authors [get]
func (h *Handler) ListAuthorsHandler(c *fiber.Ctx) error {
	authorsList, err := h.authorService.ListAuthors(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Authors list retrieved successfully",
		Data:    authorsList,
	})
}

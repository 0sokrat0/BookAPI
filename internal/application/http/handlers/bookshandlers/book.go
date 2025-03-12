package bookshandlers

import (
	"strconv"

	"github.com/0sokrat0/BookAPI/internal/application/commands"
	"github.com/0sokrat0/BookAPI/pkg/response"

	"github.com/0sokrat0/BookAPI/internal/service/books"
	"github.com/gofiber/fiber/v2"
)

// swagger:model CreateBookRequest
type CreateBookRequest struct {
	Title     string `json:"title" example:"Go Programming"`
	Year      int    `json:"year" example:"2025"`
	ISBN      string `json:"isbn" example:"1234567890"`
	Genre     string `json:"genre" example:"Programming"`
	AuthorIDs []int  `json:"author_ids"`
}

// swagger:model UpdateBookRequest
type UpdateBookRequest struct {
	Title     string `json:"title" example:"Advanced Go"`
	Year      int    `json:"year" example:"2025"`
	ISBN      string `json:"isbn" example:"0987654321"`
	Genre     string `json:"genre" example:"Programming"`
	AuthorIDs []int  `json:"author_ids"`
}

type Handler struct {
	bookService books.BookService
}

func NewHandler(bookService books.BookService) *Handler {
	return &Handler{bookService: bookService}
}

// CreateBookHandler godoc
// @Summary      Create a new book
// @Description  Создаёт новую книгу в системе. Принимает данные книги в формате JSON и возвращает созданную запись.
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body       bookshandlers.CreateBookRequest  true  "Параметры для создания книги. Пример: {\"title\":\"Go Programming\",\"year\":2025,\"isbn\":\"1234567890\",\"genre\":\"Programming\",\"author_ids\":[1,2]}"
// @Success      200   {object}   response.BaseResponse "Созданная книга с её уникальным ID"
// @Failure      400   {object}   response.ErrorResponse "Неверный формат запроса или отсутствуют обязательные поля"
// @Failure      500   {object}   response.ErrorResponse "Ошибка сервера"
// @Router       /book [post]
func (h *Handler) CreateBookHandler(c *fiber.Ctx) error {
	var req commands.CreateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
	}
	book, err := h.bookService.CreateBook(c.UserContext(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Book created successfully",
		Data:    book,
	})
}

// GetBookHandler godoc
// @Summary      Retrieve a book by ID
// @Description  Возвращает данные книги по её уникальному идентификатору.
// @Tags         books
// @Produce      json
// @Param        id   path      int  true  "Уникальный ID книги"
// @Success      200  {object}  response.BaseResponse "Данные книги"
// @Failure      400  {object}  response.ErrorResponse "Неверный ID"
// @Failure      404  {object}  response.ErrorResponse "Книга не найдена"
// @Router       /book/{id} [get]
func (h *Handler) GetBookHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid book ID",
		})
	}
	book, err := h.bookService.GetBook(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
			Code:    fiber.StatusNotFound,
			Message: "Book not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Book retrieved successfully",
		Data:    book,
	})
}

// UpdateBookHandler godoc
// @Summary      Update a book
// @Description  Обновляет данные книги по её уникальному идентификатору. Принимает новые данные книги в формате JSON.
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "Уникальный ID книги"
// @Param        book  body       bookshandlers.UpdateBookRequest  true  "Данные для обновления книги. Пример: {\"title\":\"Advanced Go\",\"year\":2025,\"isbn\":\"0987654321\",\"genre\":\"Programming\",\"author_ids\":[3,4]}"
// @Success      200   {object}  response.BaseResponse "Обновлённые данные книги"
// @Failure      400   {object}  response.ErrorResponse "Неверный запрос или ID"
// @Failure      500   {object}  response.ErrorResponse "Ошибка сервера"
// @Router       /book/{id} [put]
func (h *Handler) UpdateBookHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid book ID",
		})
	}
	var req commands.UpdateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
	}
	updatedBook, err := h.bookService.UpdateBook(c.UserContext(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Book updated successfully",
		Data:    updatedBook,
	})
}

// DeleteBookHandler godoc
// @Summary      Delete a book
// @Description  Удаляет книгу из системы по её уникальному идентификатору.
// @Tags         books
// @Produce      json
// @Param        id   path      int  true  "Уникальный ID книги"
// @Success      200  {object}  response.BaseResponse "Сообщение об успешном удалении"
// @Failure      400  {object}  response.ErrorResponse "Неверный ID"
// @Failure      500  {object}  response.ErrorResponse "Ошибка сервера"
// @Router       /book/{id} [delete]
func (h *Handler) DeleteBookHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid book ID",
		})
	}
	if err := h.bookService.DeleteBook(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Book deleted successfully",
	})
}

// ListBooksHandler godoc
// @Summary      List all books
// @Description  Возвращает список всех книг, хранящихся в системе. Если указан параметр "author", возвращаются книги только этого автора. Дополнительно можно задать параметры сортировки: "sort" (поле сортировки) и "order" (asc или desc).
// @Tags         books
// @Produce      json
// @Param        author  query     int     false  "ID автора для фильтрации (например, 5)"
// @Param        sort    query     string  false  "Поле для сортировки (например, 'title', 'year')"
// @Param        order   query     string  false  "Порядок сортировки: 'asc' или 'desc' (по умолчанию: asc)"
// @Success      200     {object}  response.BaseResponse "Массив книг"
// @Failure      500     {object}  response.ErrorResponse  "Ошибка сервера"
// @Router       /books [get]
func (h *Handler) ListBooksHandler(c *fiber.Ctx) error {
	authorParam := c.Query("author")
	if authorParam != "" {
		authorID, err := strconv.Atoi(authorParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid author parameter"})
		}
		booksList, err := h.bookService.ListBooksByAuthor(c.UserContext(), authorID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
			Code:    fiber.StatusOK,
			Message: "Books list retrieved successfully",
			Data:    booksList,
		})
	}

	booksList, err := h.bookService.ListBooks(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Books list retrieved successfully",
		Data:    booksList,
	})
}

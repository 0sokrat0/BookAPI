package bookshandlers

import (
	"strconv"

	"github.com/0sokrat0/BookAPI/internal/application/commands"
	"github.com/0sokrat0/BookAPI/internal/service/books"
	"github.com/gofiber/fiber/v2"
)

// CreateBookRequest содержит данные для создания книги.
// swagger:model CreateBookRequest
type CreateBookRequest struct {
	Title     string `json:"title" example:"Go Programming"`
	Year      int    `json:"year" example:"2025"`
	ISBN      string `json:"isbn" example:"1234567890"`
	Genre     string `json:"genre" example:"Programming"`
	AuthorIDs []int  `json:"author_ids"` // Пример для массива убран, чтобы избежать ошибок преобразования
}

// UpdateBookRequest содержит данные для обновления книги.
// swagger:model UpdateBookRequest
type UpdateBookRequest struct {
	Title     string `json:"title" example:"Advanced Go"`
	Year      int    `json:"year" example:"2025"`
	ISBN      string `json:"isbn" example:"0987654321"`
	Genre     string `json:"genre" example:"Programming"`
	AuthorIDs []int  `json:"author_ids"`
}

// Handler представляет обработчик для операций с книгами.
type Handler struct {
	bookService books.BookService
}

// NewHandler создаёт новый обработчик для книг.
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
// @Success      200   {object}  map[string]interface{}  "Созданная книга с её уникальным ID"
// @Failure      400   {object}  map[string]string       "Неверный формат запроса или отсутствуют обязательные поля"
// @Failure      500   {object}  map[string]string       "Ошибка сервера"
// @Router       /book [post]
func (h *Handler) CreateBookHandler(c *fiber.Ctx) error {
	var req commands.CreateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request: " + err.Error()})
	}
	book, err := h.bookService.CreateBook(c.UserContext(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(book)
}

// GetBookHandler godoc
// @Summary      Retrieve a book by ID
// @Description  Возвращает данные книги по её уникальному идентификатору.
// @Tags         books
// @Produce      json
// @Param        id   path      int  true  "Уникальный ID книги"
// @Success      200  {object}  map[string]interface{}  "Данные книги"
// @Failure      400  {object}  map[string]string       "Неверный ID"
// @Failure      404  {object}  map[string]string       "Книга не найдена"
// @Router       /book/{id} [get]
func (h *Handler) GetBookHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}
	book, err := h.bookService.GetBook(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}
	return c.Status(fiber.StatusOK).JSON(book)
}

// UpdateBookHandler godoc
// @Summary      Update a book
// @Description  Обновляет данные книги по её уникальному идентификатору. Принимает новые данные книги в формате JSON.
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "Уникальный ID книги"
// @Param        book  body       bookshandlers.UpdateBookRequest  true  "Данные для обновления книги. Пример: {\"title\":\"Advanced Go\",\"year\":2025,\"isbn\":\"0987654321\",\"genre\":\"Programming\",\"author_ids\":[3,4]}"
// @Success      200   {object}  map[string]interface{}  "Обновлённые данные книги"
// @Failure      400   {object}  map[string]string       "Неверный запрос или ID"
// @Failure      500   {object}  map[string]string       "Ошибка сервера"
// @Router       /book/{id} [put]
func (h *Handler) UpdateBookHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}
	var req commands.UpdateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request: " + err.Error()})
	}
	updatedBook, err := h.bookService.UpdateBook(c.UserContext(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(updatedBook)
}

// DeleteBookHandler godoc
// @Summary      Delete a book
// @Description  Удаляет книгу из системы по её уникальному идентификатору.
// @Tags         books
// @Produce      json
// @Param        id   path      int  true  "Уникальный ID книги"
// @Success      200  {object}  map[string]string  "Сообщение об успешном удалении"
// @Failure      400  {object}  map[string]string  "Неверный ID"
// @Failure      500  {object}  map[string]string  "Ошибка сервера"
// @Router       /book/{id} [delete]
func (h *Handler) DeleteBookHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}
	if err := h.bookService.DeleteBook(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Book deleted successfully"})
}

// ListBooksHandler godoc
// @Summary      List all books
// @Description  Возвращает список всех книг, хранящихся в системе.
// @Tags         books
// @Produce      json
// @Success      200  {array}   map[string]interface{}  "Массив книг"
// @Failure      500  {object}  map[string]string       "Ошибка сервера"
// @Router       /books [get]
func (h *Handler) ListBooksHandler(c *fiber.Ctx) error {
	booksList, err := h.bookService.ListBooks(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(booksList)
}

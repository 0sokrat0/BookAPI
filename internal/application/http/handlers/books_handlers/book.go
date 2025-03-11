package books_handlers

import (
	"github.com/0sokrat0/BookAPI/internal/application/commands"
	"github.com/0sokrat0/BookAPI/internal/service/books"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	bookService books.BookService
}

func NewHandler(bookService books.BookService) *Handler {
	return &Handler{bookService: bookService}
}

// CreateBookHandler godoc
// @Summary      Create a new book
// @Description  Creates a new book record in the system.
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        request  body      github.com/0sokrat0/BookAPI/internal/application/commands.CreateBookRequest  true  "Book creation request"
// @Success      200      {object}  github.com/0sokrat0/BookAPI/internal/domain/aggregate/books.Book
// @Failure      400      {object}  map[string]string  "Invalid request"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /books [post]
func (h *Handler) CreateBookHandler(c *fiber.Ctx) error {
	var req commands.CreateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	book, err := h.bookService.CreateBook(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(book)
}

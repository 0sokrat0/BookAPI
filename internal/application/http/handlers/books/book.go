package books

import (
	"api/internal/application/commands"
	"api/internal/service/books"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	bookService books.BookService
}

func NewHandler(bookService books.BookService) *Handler {
	return &Handler{bookService: bookService}
}

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

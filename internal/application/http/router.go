package http

import (
	_ "github.com/0sokrat0/BookAPI/docs"
	authorhandlers "github.com/0sokrat0/BookAPI/internal/application/http/handlers/authors"
	"github.com/0sokrat0/BookAPI/internal/application/http/handlers/bookshandlers"
	readerhandlers "github.com/0sokrat0/BookAPI/internal/application/http/handlers/readers"
	reservationshandlers "github.com/0sokrat0/BookAPI/internal/application/http/handlers/reservations"

	"github.com/gofiber/contrib/swagger"
)

func (s *Server) registerRouter() {
	cfg := swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Swagger API Docs",
	}
	s.App.Use(swagger.New(cfg))

	handlerBooks := bookshandlers.NewHandler(s.bookService)
	handlerReader := readerhandlers.NewHandler(s.readerService)
	handlerAuthor := authorhandlers.NewHandler(s.authorService)
	handlerReservation := reservationshandlers.NewHandler(s.reservService)

	s.App.Post("/book", handlerBooks.CreateBookHandler)
	s.App.Get("/book/:id", handlerBooks.GetBookHandler)
	s.App.Put("/book/:id", handlerBooks.UpdateBookHandler)
	s.App.Delete("/book/:id", handlerBooks.DeleteBookHandler)
	s.App.Get("/books", handlerBooks.ListBooksHandler)

	s.App.Post("/reader", handlerReader.CreateReaderHandler)
	s.App.Get("/reader/:id", handlerReader.GetReaderHandler)
	s.App.Put("/reader/:id", handlerReader.UpdateReaderHandler)
	s.App.Delete("/reader/:id", handlerReader.DeleteReaderHandler)
	s.App.Get("/readers", handlerReader.ListReadersHandler)
	s.App.Post("/login", handlerReader.AuthenticateReaderHandler)

	s.App.Post("/author", handlerAuthor.CreateAuthorHandler)
	s.App.Get("/author/:id", handlerAuthor.GetAuthorHandler)
	s.App.Put("/author/:id", handlerAuthor.UpdateAuthorHandler)
	s.App.Delete("/author/:id", handlerAuthor.DeleteAuthorHandler)
	s.App.Get("/authors", handlerAuthor.ListAuthorsHandler)

	s.App.Post("/reservation", handlerReservation.CreateReservationHandler)
	s.App.Get("/reservation/:id", handlerReservation.GetReservationHandler)
	s.App.Put("/reservation/:id", handlerReservation.UpdateReservationHandler)
	s.App.Delete("/reservation/:id", handlerReservation.DeleteReservationHandler)
	s.App.Get("/reservations", handlerReservation.ListReservationsHandler)
}

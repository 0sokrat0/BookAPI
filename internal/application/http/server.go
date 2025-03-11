package http

import (
	"context"
	"fmt"

	"github.com/0sokrat0/BookAPI/internal/config"

	authorsrepo "github.com/0sokrat0/BookAPI/internal/infrastructure/authorsRepo"
	"github.com/0sokrat0/BookAPI/internal/infrastructure/booksRepo"
	readersrepo "github.com/0sokrat0/BookAPI/internal/infrastructure/readersRepo"
	reservrepo "github.com/0sokrat0/BookAPI/internal/infrastructure/reservations"
	"github.com/0sokrat0/BookAPI/internal/service/authors"
	"github.com/0sokrat0/BookAPI/internal/service/books"
	"github.com/0sokrat0/BookAPI/internal/service/readers"
	"github.com/0sokrat0/BookAPI/internal/service/reservations"
	genid "github.com/0sokrat0/BookAPI/pkg/GenID"
	"github.com/0sokrat0/BookAPI/pkg/db/postgres"
	"github.com/0sokrat0/BookAPI/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App    *fiber.App
	Config *config.Config

	bookService   books.BookService
	authorService authors.AuthorService
	readerService readers.ReaderService
	reservService reservations.ReservationService
}

func NewServer(ctx context.Context, cfg *config.Config, pool *postgres.Postgres, idCounter *genid.IDcounter) *Server {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       cfg.App.Name,
	})

	lg := logger.FromContext(ctx)
	app.Use(func(c *fiber.Ctx) error {
		c.SetUserContext(logger.WithLogger(c.UserContext(), lg))
		return c.Next()
	})

	bookRepos := booksRepo.NewBookRepo(pool.DB)
	// Передаём idCounter в конструктор BookService
	bookService := books.NewBookService(bookRepos, idCounter)

	authorRepos := authorsrepo.NewAuthorRepo(pool.DB)
	authorService := authors.NewAuthorService(authorRepos, idCounter)

	readerRepos := readersrepo.NewReaderRepo(pool.DB)
	readerService := readers.NewReaderService(readerRepos, idCounter)

	reservationsRepos := reservrepo.NewReservationRepo(pool.DB)
	reservationService := reservations.NewReservationService(reservationsRepos)

	srv := &Server{
		App:           app,
		Config:        cfg,
		bookService:   bookService,
		authorService: authorService,
		readerService: readerService,
		reservService: reservationService,
	}

	srv.registerRouter()
	return srv
}

func (s *Server) Start() error {
	address := fmt.Sprintf(":%d", s.Config.App.Port)
	return s.App.Listen(address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.App.ShutdownWithContext(ctx)
}

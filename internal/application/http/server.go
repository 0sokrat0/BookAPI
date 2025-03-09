package http

import (
	"api/internal/config"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App    *fiber.App
	Config *config.Config
}

func NewServer(ctx context.Context, cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       cfg.App.Name,
	})

	srv := &Server{
		App:    app,
		Config: cfg,
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

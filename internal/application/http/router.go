package http

import (
	_ "github.com/0sokrat0/BookAPI/docs"

	"github.com/gofiber/contrib/swagger"
)

func (s *Server) registerRouter() {
	cfg := swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Swagger API Docs",
		CacheAge: 3600,
	}
	s.App.Use(swagger.New(cfg))

}

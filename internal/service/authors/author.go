package authors

import (
	"context"

	"github.com/0sokrat0/BookAPI/internal/application/commands"
	"github.com/0sokrat0/BookAPI/internal/domain/entity/authors"
	genid "github.com/0sokrat0/BookAPI/pkg/GenID"
)

// AuthorService описывает бизнес-логику для авторов.
type AuthorService interface {
	CreateAuthor(ctx context.Context, req commands.CreateAuthorRequest) (*authors.Author, error)
	GetAuthor(ctx context.Context, id int) (*authors.Author, error)
	UpdateAuthor(ctx context.Context, id int, req commands.UpdateAuthorRequest) (*authors.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
	ListAuthors(ctx context.Context) ([]authors.Author, error)
}

type authorService struct {
	authorRepo authors.AuthorRepo
	idCounter  *genid.IDcounter
}

// NewAuthorService возвращает реализацию AuthorService.
func NewAuthorService(repo authors.AuthorRepo, counter *genid.IDcounter) AuthorService {
	return &authorService{
		authorRepo: repo,
		idCounter:  counter,
	}
}

func (s *authorService) CreateAuthor(ctx context.Context, req commands.CreateAuthorRequest) (*authors.Author, error) {
	newID := s.idCounter.GenerateID()
	newAuthor, err := authors.NewAuthor(newID, req.Name, req.Country)
	if err != nil {
		return nil, err
	}
	if err := s.authorRepo.Create(ctx, newAuthor); err != nil {
		return nil, err
	}
	return newAuthor, nil
}

func (s *authorService) GetAuthor(ctx context.Context, id int) (*authors.Author, error) {
	return s.authorRepo.GetById(ctx, id)
}

func (s *authorService) UpdateAuthor(ctx context.Context, id int, req commands.UpdateAuthorRequest) (*authors.Author, error) {
	existingAuthor, err := s.authorRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	// Обновляем поля
	if req.Name != "" {
		existingAuthor.Name = req.Name
	}
	existingAuthor.Country = req.Country
	if err := s.authorRepo.Update(ctx, existingAuthor); err != nil {
		return nil, err
	}
	return existingAuthor, nil
}

func (s *authorService) DeleteAuthor(ctx context.Context, id int) error {
	return s.authorRepo.Delete(ctx, id)
}

func (s *authorService) ListAuthors(ctx context.Context) ([]authors.Author, error) {
	return s.authorRepo.List(ctx)
}

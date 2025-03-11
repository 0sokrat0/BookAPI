package readers

import (
	"context"
	"fmt"

	"github.com/0sokrat0/BookAPI/internal/application/commands"
	domainReaders "github.com/0sokrat0/BookAPI/internal/domain/entity/readers"
	genid "github.com/0sokrat0/BookAPI/pkg/GenID"
)

type ReaderService interface {
	CreateReader(ctx context.Context, req commands.CreateReaderRequest) (*domainReaders.Reader, error)
	GetReader(ctx context.Context, id int) (*domainReaders.Reader, error)
	UpdateReader(ctx context.Context, id int, req commands.UpdateReaderRequest) (*domainReaders.Reader, error)
	DeleteReader(ctx context.Context, id int) error
	ListReaders(ctx context.Context) ([]domainReaders.Reader, error)
}

type readerService struct {
	readerRepo domainReaders.ReaderRepo
	idCounter  *genid.IDcounter
}

func NewReaderService(repo domainReaders.ReaderRepo, counter *genid.IDcounter) ReaderService {
	return &readerService{
		readerRepo: repo,
		idCounter:  counter,
	}
}

func (s *readerService) CreateReader(ctx context.Context, req commands.CreateReaderRequest) (*domainReaders.Reader, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	newID := s.idCounter.GenerateID()
	newReader, err := domainReaders.NewReader(newID, req.Name, req.Phone, req.Email, req.Password, req.Admin)
	if err != nil {
		return nil, err
	}
	if err := s.readerRepo.Create(ctx, newReader); err != nil {
		return nil, err
	}
	return newReader, nil
}

func (s *readerService) GetReader(ctx context.Context, id int) (*domainReaders.Reader, error) {
	return s.readerRepo.GetById(ctx, id)
}

func (s *readerService) UpdateReader(ctx context.Context, id int, req commands.UpdateReaderRequest) (*domainReaders.Reader, error) {
	existingReader, err := s.readerRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		existingReader.Name = req.Name
	}
	existingReader.Phone = req.Phone
	existingReader.Email = req.Email
	existingReader.Password = req.Password
	existingReader.Admin = req.Admin
	if err := s.readerRepo.Update(ctx, existingReader); err != nil {
		return nil, err
	}
	return existingReader, nil
}

func (s *readerService) DeleteReader(ctx context.Context, id int) error {
	return s.readerRepo.Delete(ctx, id)
}

func (s *readerService) ListReaders(ctx context.Context) ([]domainReaders.Reader, error) {
	return s.readerRepo.List(ctx)
}

package books

import (
	"context"
	"fmt"
)

type Book struct {
	ID        int
	Title     string
	Year      int
	ISBN      string
	Genre     string
	authorIDs []int
}

type BookRepo interface {
	Create(ctx context.Context, book *Book) error
	GetByID(ctx context.Context, id int) (*Book, error)
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]Book, error)
	ListBooksByAuthor(ctx context.Context, authorID int) ([]Book, error)
}

func NewBook(id int, title string, year int, isbn string, genre string, authorIDs []int) (*Book, error) {
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	return &Book{
		ID:        id,
		Title:     title,
		Year:      year,
		ISBN:      isbn,
		Genre:     genre,
		authorIDs: authorIDs,
	}, nil
}

// AuthorIDs возвращает копию списка идентификаторов авторов.
func (b *Book) AuthorIDs() []int {
	ids := make([]int, len(b.authorIDs))
	copy(ids, b.authorIDs)
	return ids
}

// SetAuthorIDs устанавливает список идентификаторов авторов.
func (b *Book) SetAuthorIDs(ids []int) {
	b.authorIDs = ids
}

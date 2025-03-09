package commands

type CreateBookRequest struct {
	Title     string
	Year      int
	ISBN      string
	Genre     string
	AuthorIDs []int
}

package commands

type CreateBookRequest struct {
	Title     string `json:"title"`
	Year      int    `json:"year"`
	ISBN      string `json:"isbn"`
	Genre     string `json:"genre"`
	AuthorIDs []int  `json:"author_ids"`
}

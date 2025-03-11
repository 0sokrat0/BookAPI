package commands

// CreateBookRequest содержит данные для создания книги.
// swagger:parameters CreateBookRequest
type CreateBookRequest struct {
	Title     string `json:"title" example:"Go Programming"`
	Year      int    `json:"year" example:"2025"`
	ISBN      string `json:"isbn" example:"1234567890"`
	Genre     string `json:"genre" example:"Programming"`
	AuthorIDs []int  `json:"author_ids" `
}

// UpdateBookRequest содержит данные для обновления книги.
type UpdateBookRequest struct {
	Title     string `json:"title" example:"Advanced Go"`
	Year      int    `json:"year" example:"2025"`
	ISBN      string `json:"isbn" example:"0987654321"`
	Genre     string `json:"genre" example:"Programming"`
	AuthorIDs []int  `json:"author_ids" `
}

package commands

// CreateAuthorRequest содержит данные для создания автора.
type CreateAuthorRequest struct {
	Name    string `json:"name" example:"Leo Tolstoy"`
	Country string `json:"country" example:"Russia"`
}

// UpdateAuthorRequest содержит данные для обновления автора.
type UpdateAuthorRequest struct {
	Name    string `json:"name" example:"Leo Tolstoy"`
	Country string `json:"country" example:"Russia"`
}

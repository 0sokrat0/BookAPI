package commands

// CreateReaderRequest содержит данные для создания читателя.
type CreateReaderRequest struct {
	Name     string `json:"name" example:"Ivan Ivanov"`
	Phone    string `json:"phone" example:"+79111234567"`
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"password123"`
	Admin    bool   `json:"admin" example:"false"`
}

// UpdateReaderRequest содержит данные для обновления читателя.
type UpdateReaderRequest struct {
	Name     string `json:"name" example:"Ivan Ivanov"`
	Phone    string `json:"phone" example:"+79111234567"`
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"newpassword"` // добавляем поле, если требуется обновление пароля
	Admin    bool   `json:"admin" example:"false"`          // добавляем поле, если требуется обновление прав администратора
}

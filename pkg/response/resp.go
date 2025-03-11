package response

// BaseResponse — стандартный формат ответа API
type BaseResponse struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse — формат ответа в случае ошибки
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}

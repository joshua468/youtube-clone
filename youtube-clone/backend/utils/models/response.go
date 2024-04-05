package models

import "github.com/google/uuid"

type VideoResponse struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"userID"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}

type ErrorResponse struct {
	ID            string `json:"id"`
	Handler       string `json:"handler"`
	PublicMessage string `json:"publicMessage"`
}

type GenericResponse struct {
	Code    int           `json:"code"`
	Data    interface{}   `json:"data,omitempty"`
	Message *string       `json:"message,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

func NewErrorResponse(id, handler, publicMessage string) *ErrorResponse {
	return &ErrorResponse{
		ID:            id,
		Handler:       handler,
		PublicMessage: publicMessage,
	}
}

func NewGenericResponse(code int, data interface{}, message *string, err *ErrorResponse) *GenericResponse {
	return &GenericResponse{
		Code:    code,
		Data:    data,
		Message: message,
		Error:   err,
	}
}

package models

import "github.com/google/uuid"

type CreateVideoRequest struct {
    UserID  uuid.UUID `json:"userID" validate:"required,uuid"`
    Title   string    `json:"title" validate:"required"`
    Content string    `json:"content"`
}

type UpdateVideoRequest struct {
    Title   string    `json:"title" validate:"required"`
    Content string    `json:"content" validate:"required"`
    ID      uuid.UUID `json:"id" validate:"required"`
}

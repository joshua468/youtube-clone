// user.go

package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joshua468/youtube-clone/backend/repository"
	"github.com/joshua468/youtube-clone/backend/utils/models"
)

// GetUserByID retrieves a user by ID
func (a *App) GetUserByID(ctx *gin.Context, userID uuid.UUID) (*models.User, error) {
	user, err := a.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to get user by ID")
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user
func (a *App) CreateUser(ctx *gin.Context, userRequest models.CreateUserRequest) (*models.User, error) {
	user, err := a.userRepository.CreateUser(ctx, userRequest)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}
	return user, nil
}

// Login performs user login
func (a *App) Login(ctx *gin.Context, loginReq models.LoginRequest) (*models.User, error) {
	user, err := a.userRepository.Login(ctx, loginReq)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to login")
		return nil, err
	}
	return user, nil
}

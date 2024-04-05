package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/joshua468/youtube-clone/backend/repository"
	"github.com/joshua468/youtube-clone/backend/utils/helpers"
	"github.com/joshua468/youtube-clone/backend/utils/models"
)

// App represents the core application struct
type App struct {
	env             models.Env
	logger          zerolog.Logger
	userRepository  repository.UserRepository
	videoRepository repository.VideoRepository
}

// Operations defines the operations supported by the App
type Operations interface {
	GetUserByID(ctx *gin.Context, userID uuid.UUID) (*models.User, error)
	CreateUser(ctx *gin.Context, userRequest models.CreateUserRequest) (*models.User, error)
	Login(ctx *gin.Context, loginReq models.LoginRequest) (*models.User, error)
	CreateVideo(ctx context.Context, video models.Video) (*models.Video, error)
	GetVideos(ctx context.Context, page helpers.Page) ([]*models.Video, helpers.PageInfo, error)
	GetUserVideos(ctx context.Context, userID uuid.UUID, page helpers.Page) ([]*models.Video, helpers.PageInfo, error)
}

// New creates a new instance of App
func New(env models.Env, store repository.Store, logger zerolog.Logger) *App {
	appLogger := logger.With().Str("package", "app").Logger()

	userRepo := repository.NewUserRepository(store)
	videoRepo := repository.NewVideoRepository(store)

	return &App{
		env:             env,
		logger:          appLogger,
		userRepository:  userRepo,
		videoRepository: videoRepo,
	}
}

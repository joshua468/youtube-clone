// video.go

package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joshua468/youtube-clone/backend/repository"
	"github.com/joshua468/youtube-clone/backend/utils/helpers"
	"github.com/joshua468/youtube-clone/backend/utils/models"
)

// CreateVideo creates a new video
func (a *App) CreateVideo(ctx context.Context, video models.Video) (*models.Video, error) {
	newVideo, err := a.videoRepository.CreateVideo(ctx, video)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to create video")
		return nil, err
	}
	return newVideo, nil
}

// GetVideos retrieves a list of videos
func (a *App) GetVideos(ctx context.Context, page helpers.Page) ([]*models.Video, helpers.PageInfo, error) {
	videos, pageInfo, err := a.videoRepository.GetVideos(ctx, page)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to get videos")
		return nil, helpers.PageInfo{}, err
	}
	return videos, pageInfo, nil
}

// GetUserVideos retrieves videos for a specific user
func (a *App) GetUserVideos(ctx context.Context, userID uuid.UUID, page helpers.Page) ([]*models.Video, helpers.PageInfo, error) {
	videos, pageInfo, err := a.videoRepository.GetUserVideos(ctx, userID, page)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to get user videos")
		return nil, helpers.PageInfo{}, err
	}
	return videos, pageInfo, nil
}

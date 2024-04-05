package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/joshua468/youtube-clone/backend/utils/helpers"
	"github.com/joshua468/youtube-clone/backend/utils/models"
	"github.com/joshua468/youtube-clone/backend/utils/helpers"
	"github.com/joshua468/youtube-clone/backend/repository/models" 
)

type VideoRepo interface {
	Create(ctx context.Context, v models.Video) (*models.Video, error)
	GetAllVideos(ctx context.Context, query models.Video, p helpers.Page) ([]*models.Video, helpers.PageInfo, error)
	SoftDeleteByID(ctx context.Context, ID uuid.UUID) error
	GetByID(ctx context.Context, ID uuid.UUID) (*models.Video, error)
	UpdateVideo(ctx context.Context, video models.Video) (*models.Video, error)
	CountVideos(ctx context.Context) (int64, error)
}

type Video struct {
	logger  zerolog.Logger
	storage *Store
}

// NewVideo creates a new reference to the Video storage entity
func NewVideo(s *Store) *Video {
	l := s.logger.With().Str("LEVEL_NAME", "video").Logger()
	video := &Video{
		logger:  l,
		storage: s,
	}
	videoDatabase := VideoRepo(video)
	return &videoDatabase
}

func (v *Video) CountVideos(ctx context.Context) (int64, error) {
	log := v.logger.With().Str(helpers.LogStrRequestIDLevel, v.storage.getRequestID(ctx)).
		Str(helpers.LogStrKeyMethod, "repository.video.CountVideos").Logger()

	var count int64
	db := v.storage.DB.WithContext(ctx).Model(&models.Video{}).Count(&count)
	if db.Error != nil {
		log.Err(db.Error).Msg("count not possible")
		return count, helpers.ErrRecordNotFound
	}
	return count, nil
}

func (v *Video) UpdateVideo(ctx context.Context, video models.Video) (*models.Video, error) {
	log := v.logger.With().Str(helpers.LogStrRequestIDLevel, v.storage.getRequestID(ctx)).
		Str(helpers.LogStrKeyMethod, "repository.video.UpdateVideo").Logger()

	db := v.storage.DB.WithContext(ctx).Model(models.Video{
		ID: video.ID,
	}).Updates(models.Video{
		Title:       video.Title,
		Description: video.Description,
		UpdatedAt:   video.UpdatedAt,
	})
	if db.Error != nil {
		log.Err(db.Error).Msg("unable to update video")
		return nil, helpers.ErrRecordUpdateFail
	}
	return &video, nil
}

func (v *Video) GetByID(ctx context.Context, ID uuid.UUID) (*models.Video, error) {
	log := v.logger.With().Str(helpers.LogStrRequestIDLevel, v.storage.getRequestID(ctx)).
		Str(helpers.LogStrKeyMethod, "repository.video.GetByID").Logger()

	var video models.Video
	db := v.storage.DB.WithContext(ctx).Where("id = ?", ID.String()).Find(&video)
	if db.Error != nil || strings.EqualFold(video.ID.String(), helpers.ZeroUUID) {
		log.Err(db.Error).Msg("record not found")
		return nil, helpers.ErrRecordNotFound
	}
	return &video, nil
}

func (v *Video) GetAllVideos(ctx context.Context, query models.Video, page helpers.Page) ([]*models.Video, helpers.PageInfo, error) {
	log := v.logger.With().Str(helpers.LogStrRequestIDLevel, v.storage.getRequestID(ctx)).
		Str(helpers.LogStrKeyMethod, "repository.video.GetAllVideos").Logger()

	var videos []*models.Video
	offset := 0
	// load defaults
	if page.Number == nil {
		tmpPageNumber := helpers.PageDefaultNumber
		page.Number = &tmpPageNumber
	}
	if page.Size == nil {
		tmpPageSize := helpers.PageDefaultSize
		page.Size = &tmpPageSize
	}
	if page.SortBy == nil {
		tmpPageSortBy := helpers.PageDefaultSortBy
		page.SortBy = &tmpPageSortBy
	}
	if page.SortDirectionDesc == nil {
		tmpPageSortDirectionDesc := helpers.PageDefaultSortDirectionDesc
		page.SortDirectionDesc = &tmpPageSortDirectionDesc
	}

	if *page.Number > 1 {
		offset = *page.Size * (*page.Number - 1)
	}
	sortDirection := helpers.PageSortDirectionDescending
	if !*page.SortDirectionDesc {
		sortDirection = helpers.PageSortDirectionAscending
	}

	queryDraft := v.storage.DB.WithContext(ctx).Model(models.Video{}).Where(query)

	// then do counting
	var count int64
	queryDraft.Count(&count)

	db := queryDraft.Offset(offset).Limit(*page.Size).
		Order(fmt.Sprintf("%s %s", *page.SortBy, sortDirection)).
		Find(&videos)

	if db.Error != nil {
		log.Err(db.Error).Msg("could not fetch list of videos")
		return nil, helpers.PageInfo{}, helpers.ErrEmptyResult
	}

	return videos, helpers.PageInfo{
		Page:            *page.Number,
		Size:            *page.Size,
		HasNextPage:     int64(offset+*page.Size) < count,
		HasPreviousPage: *page.Number > 1,
		TotalCount:      count,
	}, nil
}

func (v *Video) Create(ctx context.Context, video models.Video) (*models.Video, error) {
	log := v.logger.With().Str(helpers.LogStrRequestIDLevel, v.storage.getRequestID(ctx)).
		Str(helpers.LogStrKeyMethod, "repository.video.Create").Logger()

	db := v.storage.DB.WithContext(ctx).Model(&models.Video{}).Create(&video)
	if db.Error != nil {
		log.Err(db.Error).Msg("unable to insert new row")
		if strings.Contains(db.Error.Error(), "duplicate key value") {
			return nil, errors.New("duplicate record error")
		}
		return nil, helpers.ErrRecordCreationFailed
	}

	return &video, nil
}

func (v *Video) SoftDeleteByID(ctx context.Context, id uuid.UUID) error {
	log := v.logger.With().Str(helpers.LogStrRequestIDLevel, v.storage.getRequestID(ctx)).
		Str(helpers.LogStrKeyMethod, "repository.video.SoftDeleteByID").Logger()

	db := v.storage.DB.WithContext(ctx).Model(models.Video{}).Where("id = ?", id).UpdateColumns(models.Video{
		DeletedAt: &gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if db.Error != nil {
		log.Err(db.Error).Msg("soft delete failed")
		return helpers.ErrDeleteFailed
	}

	return nil
}

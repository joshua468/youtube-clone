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
)

type UserRepo interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetAllUsers(ctx context.Context, query models.User, page helpers.Page) ([]*models.User, helpers.PageInfo, error)
	CountUsers(ctx context.Context) (int64, error)
}

type User struct {
	logger  zerolog.Logger
	storage *Store
}

// NewUser creates a new reference to the User storage entity
func NewUser(s *Store) *User {
	l := s.logger.With().Str("LEVEL_NAME", "user").Logger()
	user := &User{
		logger:  l,
		storage: s,
	}
	userDatabase := UserRepo(user)
	return &userDatabase
}

func (u *User) CountUsers(ctx context.Context) (int64, error) {
	log := u.logger.With().Str(helpers.LogStrRequestIDLevel, u.storage.getRequestID(ctx)).
		Str(helpers.LogStrKeyMethod, "repository.user.CountUsers").Logger()

	var count int64
	db := u.storage.DB.WithContext(ctx).Model(&models.User{}).Count(&count)
	if db.Error != nil {
		log.Err(db.Error).Msg("count not possible")
		return count, helpers.ErrRecordNotFound
	}
	return count, nil
}

func (u *User) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	log := u.logger.With().Str(helpers.LogStrRequestIDLevel, u.storage.getRequestID(ctx)).
		Str(helpers.LogStrKeyMethod, "repository.user.GetUserByUsername").Logger()

	var user models.User
	db := u.storage.DB.WithContext(ctx).Where("username = ?", username).First(&user)
	if db.Error != nil || strings.EqualFold(user.ID.String(), helpers.ZeroUUID) {
		log.Err(db.Error).Msg("user not found")
		return nil, helpers.ErrRecordNotFound
	}
	return &user, nil
}

func (u *User) GetAllUsers(ctx context.Context, query models.User, page helpers.Page) ([]*models.User, helpers.PageInfo, error) {
	log := u.logger.With().Str(helpers.LogStrRequestIDLevel, u.storage.getRequestID(ctx)).
		Str(helpers.LogStrKeyMethod, "repository.user.GetAllUser").Logger()

	var users []*models.User
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

	queryDraft := u.storage.DB.WithContext(ctx).Model(models.User{}).Where(query)

	// then do counting
	var count int64
	queryDraft.Count(&count)

	db := queryDraft.Offset(offset).Limit(*page.Size).
		Order(fmt.Sprintf("%s %s", *page.SortBy, sortDirection)).
		Find(&users)

	if db.Error != nil {
		log.Err(db.Error).Msg("could not fetch list of users")
		return nil, helpers.PageInfo{}, helpers.ErrEmptyResult
	}

	return users, helpers.PageInfo{
		Page:            *page.Number,
		Size:            *page.Size,
		HasNextPage:     int64(offset+*page.Size) < count,
		HasPreviousPage: *page.Number > 1,
		TotalCount:      count,
	}, nil
}

func (u *User) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	log := u.logger.With().Str(helpers.LogStrRequestIDLevel, u.storage.getRequestID(ctx)).
		Str(helpers.LogStrKeyMethod, "repository.user.Create").Logger()

	db := u.storage.DB.WithContext(ctx).Model(&models.User{}).Create(&user)
	if db.Error != nil {
		log.Err(db.Error).Msg("unable to insert new row")
		if strings.Contains(db.Error.Error(), "duplicate key value") {
			return nil, errors.New("duplicate record error")
		}
		return nil, helpers.ErrRecordCreationFailed
	}

	return &user, nil
}

func (u *User) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	log := u.logger.With().Str(helpers.LogStrRequestIDLevel, u.storage.getRequestID(ctx)).
		Str(helpers.LogStrKeyMethod, "repository.user.GetUserByID").Logger()

	var user models.User
	db := u.storage.DB.WithContext(ctx).Where("id = ?", userID.String()).First(&user)
	if db.Error != nil || strings.EqualFold(user.ID.String(), helpers.ZeroUUID) {
		log.Err(db.Error).Msg("user not found")
		return nil, helpers.ErrRecordNotFound
	}
	return &user, nil
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	log := u.logger.With().Str(helpers.LogStrRequestIDLevel, u.storage.getRequestID(ctx)).
		Str(helpers.LogStrKeyMethod, "repository.user.GetUserByEmail").Logger()

	var user models.User
	db := u.storage.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if db.Error != nil || strings.EqualFold(user.ID.String(), helpers.ZeroUUID) {
		log.Err(db.Error).Msg("user not found")
		return nil, helpers.ErrRecordNotFound
	}
	return &user, nil
}

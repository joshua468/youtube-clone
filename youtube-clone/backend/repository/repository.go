package repository

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/joshua468/youtube-clone/backend/utils/helpers"
	"github.com/joshua468/youtube-clone/backend/utils/models"
)

const (
	packageName = "backend.repository"
)

// Store object
type Store struct {
	logger *zerolog.Logger
	DB     *gorm.DB
}

// New creates new instance of a Store
func New(z zerolog.Logger, env models.Env) *Store {
	log := z.With().Str("PACKAGE", packageName).Logger()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.DBUsername,
		env.DBPassword,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		z.Fatal().Err(err).Msgf("could not connect to the DB %+v", err.Error())
		panic(err)
	}

	z.Debug().Msg("connected to the database")

	err = db.AutoMigrate(&models.User{}, &models.Video{}) // Adjust the models as per your requirements
	if err != nil {
		z.Fatal().Err(err).Msg("unable to auto migrate models")
		panic(err)
	}

	return &Store{
		logger: &log,
		DB:     db,
	}
}

func (s *Store) Close() {
	sqlDB, _ := s.DB.DB()
	_ = sqlDB.Close()
}

func (s *Store) getRequestID(ctx context.Context) string {
	rID := ctx.Value("RequestIDContextKey")
	if rID != nil {
		return rID.(string)
	}

	return helpers.ZeroUUID
}

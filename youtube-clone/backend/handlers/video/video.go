package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joshua468/youtube-clone/backend/app"
	"github.com/joshua468/youtube-clone/backend/utils/middlewares"
	"github.com/joshua468/youtube-clone/backend/utils/models"
	"github.com/rs/zerolog"
)

const handlerNameVideo = "video"

type videoHandler struct {
	logger     *zerolog.Logger
	app        *app.App
	env        *models.Env
	middleware middlewares.Middleware
}

func NewVideoHandler(r *gin.RouterGroup, l *zerolog.Logger, a *app.App, e *models.Env, m middlewares.Middleware) {
	video := videoHandler{
		logger:     l,
		app:        a,
		env:        e,
		middleware: m,
	}

	videoGroup := r.Group("/video")

	videoGroup.POST("", m.AuthMiddleware(false), video.create())
	videoGroup.PUT("/update/:id", m.AuthMiddleware(false), video.update()) // Added :id param
	videoGroup.GET("/mine", m.AuthMiddleware(false), video.getMyVideos())
	videoGroup.GET("/:id/user", m.AuthMiddleware(true), video.getUserVideos())
	videoGroup.GET("/:id", m.AuthMiddleware(false), video.getVideoByID())
	videoGroup.GET("/all", m.AuthMiddleware(true), video.getAllVideos())
}

func (v *videoHandler) create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateVideoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			// Handle validation errors
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call app.CreateVideo method passing req
		video, err := v.app.CreateVideo(c, req)
		if err != nil {
			// Handle error from app.CreateVideo
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create video"})
			return
		}

		// Return appropriate response
		c.JSON(http.StatusCreated, gin.H{"message": "Video created successfully", "video": video})
	}
}

func (v *videoHandler) update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.UpdateVideoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			// Handle validation errors
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Extract videoID from the URL param
		videoID := c.Param("id")
		videoUUID, err := uuid.Parse(videoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
			return
		}

		// Call app.UpdateVideo method passing req
		updatedVideo, err := v.app.UpdateVideo(c, videoUUID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update video"})
			return
		}

		// Return appropriate response
		c.JSON(http.StatusOK, gin.H{"message": "Video updated successfully", "video": updatedVideo})
	}
}

func (v *videoHandler) getMyVideos() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user ID from the context
		userID, exists := c.Get(middlewares.UserIDInContext)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			return
		}

		// Convert the user ID to UUID
		userUUID, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Call the app method to get videos by user ID
		videos, err := v.app.GetVideosByUserID(c, userUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
			return
		}

		// Return the videos in the response
		c.JSON(http.StatusOK, gin.H{"videos": videos})
	}
}

func (v *videoHandler) getUserVideos() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user ID from the request parameters
		userID := c.Param("id")

		// Convert the user ID to UUID format
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Call the app method to get videos by user ID
		videos, err := v.app.GetVideosByUserID(c, userUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user's videos"})
			return
		}

		// Return the videos in the response
		c.JSON(http.StatusOK, gin.H{"videos": videos})
	}
}

func (v *videoHandler) getVideoByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract video ID from the URL parameters
		videoID := c.Param("id")

		// Convert the videoID string to a UUID
		uuid, err := uuid.Parse(videoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
			return
		}

		// Call the app method to get the video by ID
		video, err := v.app.GetVideoByID(c, uuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video"})
			return
		}

		// Return the video in the response
		c.JSON(http.StatusOK, gin.H{"video": video})
	}
}

func (v *videoHandler) getAllVideos() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call the app method to get all videos
		videos, err := v.app.GetAllVideos(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
			return
		}

		// Return the videos in the response
		c.JSON(http.StatusOK, gin.H{"videos": videos})
	}
}

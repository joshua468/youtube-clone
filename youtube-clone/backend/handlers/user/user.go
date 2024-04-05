package user

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joshua468/youtube-clone/app"
	"github.com/joshua468/youtube-clone/logger"
	"github.com/joshua468/youtube-clone/models"
	"github.com/joshua468/youtube-clone/utils/helpers" // Import helpers package
	"github.com/joshua468/youtube-clone/utils/middlewares"
)

const handlerNameUser = "user"

type userHandler struct {
	logger     *logger.Logger
	app        *app.App
	env        *models.Env
	middleware middlewares.Middleware
}

func NewUserHandler(r *gin.RouterGroup, l *logger.Logger, a *app.App, e *models.Env, m middlewares.Middleware) {
	user := userHandler{
		app:        a,
		env:        e,
		middleware: m,
		logger:     l,
	}

	userGroup := r.Group("/user")

	userGroup.POST("/login", user.login())
	userGroup.POST("/signup", user.signup())

	userGroup.GET("/me", m.AuthMiddleware(false), user.me())
	userGroup.GET("/all", m.AuthMiddleware(true), user.getUsers())
	userGroup.GET("/search", m.AuthMiddleware(true), user.searchUsers())
	userGroup.GET("/:id", m.AuthMiddleware(true), user.getUserByID())
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *userHandler) signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SignupRequest
		var err error

		// Retrieve request body and bind it to SignupRequest struct
		if err := c.ShouldBindJSON(&req); err != nil {
			u.logger.Err(err).Msg("error binding request body")
			models.ErrorResponse(c, http.StatusBadRequest, models.ErrorData{
				PublicMessage: "Invalid request body",
			})
			return
		}

		// Validate request parameters
		if err := helpers.ValidateRequest(req); err != nil {
			u.logger.Err(err).Msg("error validating request")
			models.ErrorResponse(c, http.StatusBadRequest, models.ErrorData{
				PublicMessage: err.Error(),
			})
			return
		}

		// Check if the email is already registered
		if exists, err := u.app.UserExistsByEmail(req.Email); err != nil {
			u.logger.Err(err).Msg("error checking if user exists by email")
			models.ErrorResponse(c, http.StatusInternalServerError, models.ErrorData{
				PublicMessage: "Failed to check email availability",
			})
			return
		} else if exists {
			models.ErrorResponse(c, http.StatusConflict, models.ErrorData{
				PublicMessage: "Email already registered",
			})
			return
		}

		// Create user
		user, err := u.app.CreateUser(c, req)
		if err != nil {
			u.logger.Err(err).Msg("error creating user")
			models.ErrorResponse(c, http.StatusInternalServerError, models.ErrorData{
				PublicMessage: "Failed to create user",
			})
			return
		}

		// Hide sensitive information
		user.Password = "********"

		// Generate JWT token for the newly created user
		token, err := u.middleware.CreateToken(c, u.env, user.ID.String(), user.IsAdmin)
		if err != nil {
			u.logger.Err(err).Msg("error generating JWT token")
			models.ErrorResponse(c, http.StatusInternalServerError, models.ErrorData{
				PublicMessage: "Failed to generate authentication token",
			})
			return
		}

		// Send success response with user details and JWT token
		models.OkResponse(c, http.StatusCreated, "User created successfully", struct {
			User  models.User        `json:"user"`
			Token middlewares.Tokens `json:"token"`
		}{
			User:  *user,
			Token: *token,
		})
	}
}

func (u *userHandler) login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest
		requestID := requestid.Get(c)

		if err := c.ShouldBind(&req); err != nil {
			u.logger.Err(err).Msg("bad request")
			models.ErrorResponse(c, http.StatusBadRequest, models.ErrorData{
				ID:            requestID,
				Handler:       handlerNameUser,
				PublicMessage: "Bad request",
			})
			return
		}

		if err := helpers.ValidateRequest(req); err != nil {
			u.logger.Err(err).Msg("login validation failed")
			models.ErrorResponse(c, http.StatusBadRequest, models.ErrorData{
				ID:            requestID,
				Handler:       handlerNameUser,
				PublicMessage: "Invalid login data",
			})
			return
		}

		user, err := u.app.Login(c, req)
		if err != nil {
			u.logger.Err(err).Msg("login error")
			models.ErrorResponse(c, http.StatusBadRequest, models.ErrorData{
				ID:            requestID,
				Handler:       handlerNameUser,
				PublicMessage: "Failed to log in",
			})
			return
		}

		user.Password = helpers.StarPassword

		token, err := u.middleware.CreateToken(c, u.env, user.ID.String(), user.IsAdmin)
		if err != nil {
			u.logger.Err(err).Msg("token generation error")
			models.ErrorResponse(c, http.StatusInternalServerError, models.ErrorData{
				ID:            requestID,
				Handler:       handlerNameUser,
				PublicMessage: "Failed to generate token",
			})
			return
		}

		models.OkResponse(c, http.StatusCreated, "User logged in successfully", struct {
			User  models.User        `json:"user"`
			Token middlewares.Tokens `json:"token"`
		}{
			User:  *user,
			Token: *token,
		})
	}
}

func (u *userHandler) getUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := requestid.Get(c)

		// Retrieve users from the application layer
		users, err := u.app.GetUsers(c)
		if err != nil {
			u.logger.Err(err).Msg("error getting users")
			models.ErrorResponse(c, http.StatusInternalServerError, models.ErrorData{
				ID:            requestID,
				Handler:       handlerNameUser,
				PublicMessage: "Failed to fetch users",
			})
			return
			}

		models.OkResponse(c, http.StatusOK, "Users fetched successfully", users)
	}
}

func (u *userHandler) me() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := requestid.Get(c)
		userID := c.GetString("userId")

		// Retrieve user details from the application layer
		user, err := u.app.GetUserByID(c, userID)
		if err != nil {
			u.logger.Err(err).Msg("error getting user")
			models.ErrorResponse(c, http.StatusInternalServerError, models.ErrorData{
				ID:            requestID,
				Handler:       handlerNameUser,
				PublicMessage: "Failed to fetch user",
			})
			return
		}

		models.OkResponse(c, http.StatusOK, "User fetched successfully", user)
	}
}

func (u *userHandler) searchUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := requestid.Get(c)

		query := c.Query("q")

		// Validate query parameter
		if query == "" {
			models.ErrorResponse(c, http.StatusBadRequest, models.ErrorData{
				ID:            requestID,
				Handler:       handlerNameUser,
				PublicMessage: "Missing search query parameter",
			})
			return
		}

		// Search users with the provided query
		users, err := u.app.SearchUsers(c, query)
		if err != nil {
			u.logger.Err(err).Msg("error searching users")
			models.ErrorResponse(c, http.StatusInternalServerError, models.ErrorData{
				ID:            requestID,
				Handler:       handlerNameUser,
				PublicMessage: "Failed to search users",
			})
			return
		}

		models.OkResponse(c, http.StatusOK, "User search successful", users)
	}
}

func (u *userHandler) getUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := requestid.Get(c)

		userID := c.Param("id")
		if userID == "" {
			models.ErrorResponse(c, http.StatusBadRequest, models.ErrorData{
				ID:            requestID,
				Handler:       handlerNameUser,
				PublicMessage: "User ID is required",
			})
			return
		}

		user, err := u.app.GetUserByID(c, userID)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				models.ErrorResponse(c, http.StatusNotFound, models.ErrorData{
					ID:            requestID,
					Handler:       handlerNameUser,
					PublicMessage: "User not found",
				})
				return
			}
			u.logger.Err(err).Msg("error getting user by ID")
			models.ErrorResponse(c, http.StatusInternalServerError, models.ErrorData{
				ID:            requestID,
				Handler:       handlerNameUser,
				PublicMessage: "Failed to fetch user",
			})
			return
		}

		models.OkResponse(c, http.StatusOK, "User fetched successfully", user)
	}
}

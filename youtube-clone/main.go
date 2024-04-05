package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joshua468/youtube-clone/backend/app"
	"github.com/joshua468/youtube-clone/backend/handlers"
	"github.com/joshua468/youtube-clone/backend/utils/middlewares"
	"github.com/joshua468/youtube-clone/backendlogger"
	"github.com/joshua468/youtube-clone/backend/models"
	"github.com/joshua468/youtube-clone/backend/repository"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()

	// Load environment variables
	env := models.NewEnv()

	// Initialize Gin router
	router := gin.New()

	// Set Gin middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize repository
	store := repository.New(log, env)

	// Initialize application
	application := app.NewApp(store)

	// Initialize middleware
	middleware := middlewares.NewMiddleware()

	// Initialize user handler
	handlers.NewUserHandler(router.Group("/api"), log, application, env, middleware)

	// Initialize video handler
	handlers.NewVideoHandler(router.Group("/api"), log, application, env, middleware)

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	addr := fmt.Sprintf(":%s", port)
	log.Info().Str("address", addr).Msg("Starting server")
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}

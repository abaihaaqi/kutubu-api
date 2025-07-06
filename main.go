package main

import (
	"os"

	"github.com/abaihaaqi/kutubu-api/config"
	"github.com/abaihaaqi/kutubu-api/handler"
	"github.com/abaihaaqi/kutubu-api/middleware"
	"github.com/abaihaaqi/kutubu-api/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.ConnectDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}))

	repo := repository.NewBookRepository(config.DB)
	h := handler.NewBookHandler(repo)

	// Auth
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	// Books
	protected := r.Group("/books")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("", h.GetBooks)
		protected.GET("/:id", h.GetBookByID)
		protected.POST("", h.CreateBook)
		protected.PUT("/:id", h.UpdateBook)
		protected.DELETE("/:id", h.DeleteBook)
		protected.GET("/search", h.SearchBooks)

	}

	// Static file server
	r.Static("/uploads", "./uploads")

	port := os.Getenv("PORT")
	r.Run(":" + port)
}

package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ConnectDatabase()

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/books", GetBooks)
	r.GET("/books/:id", GetBookByID)
	r.POST("/books", CreateBook)
	r.PUT("/books/:id", UpdateBook)
	r.DELETE("/books/:id", DeleteBook)

	port := os.Getenv("PORT")
	r.Run(":" + port)
}

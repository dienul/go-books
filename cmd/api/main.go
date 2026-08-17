package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/dienul/go-books/config"
	"github.com/dienul/go-books/internal/handler"
	"github.com/dienul/go-books/internal/middleware"
	"github.com/dienul/go-books/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("gagal terhubung ke database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("gagal menutup koneksi database: %v", err)
		}
	}()

	if err := config.RunMigrations(db); err != nil {
		log.Fatalf("gagal menyiapkan database: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET wajib diatur dan minimal 32 karakter")
	}

	userRepository := repository.NewUserRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	bookRepository := repository.NewBookRepository(db)

	if err := config.SeedAdmin(context.Background(), userRepository); err != nil {
		log.Fatalf("gagal menyiapkan akun admin: %v", err)
	}

	authHandler := handler.NewAuthHandler(userRepository, jwtSecret)
	categoryHandler := handler.NewCategoryHandler(categoryRepository, bookRepository)
	bookHandler := handler.NewBookHandler(bookRepository)

	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("gagal mengatur trusted proxy: %v", err)
	}
	router.GET("/health", func(c *gin.Context) {
		if err := db.PingContext(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "database tidak tersedia",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "API dan database terhubung",
		})
	})

	api := router.Group("/api")
	api.POST("/users/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(middleware.JWTAuth(jwtSecret))
	{
		protected.GET("/categories", categoryHandler.List)
		protected.POST("/categories", categoryHandler.Create)
		protected.GET("/categories/:id", categoryHandler.Get)
		protected.PUT("/categories/:id", categoryHandler.Update)
		protected.DELETE("/categories/:id", categoryHandler.Delete)
		protected.GET("/categories/:id/books", categoryHandler.Books)

		protected.GET("/books", bookHandler.List)
		protected.POST("/books", bookHandler.Create)
		protected.GET("/books/:id", bookHandler.Get)
		protected.PUT("/books/:id", bookHandler.Update)
		protected.DELETE("/books/:id", bookHandler.Delete)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server berjalan pada port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}

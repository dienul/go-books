package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dienul/go-books/config"
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

	router := gin.Default()
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server berjalan pada port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}

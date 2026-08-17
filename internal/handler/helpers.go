package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dienul/go-books/internal/repository"
	"github.com/gin-gonic/gin"
)

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id harus berupa bilangan bulat positif"})
		return 0, false
	}
	return id, true
}

func actor(c *gin.Context) string {
	username, _ := c.Get("username")
	value, _ := username.(string)
	return value
}

func databaseError(c *gin.Context, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"message": "terjadi kesalahan pada server"})
}

package handler

import (
	"errors"
	"net/http"

	"github.com/dienul/go-books/internal/model"
	"github.com/dienul/go-books/internal/repository"
	"github.com/gin-gonic/gin"
)

type BookHandler struct{ books *repository.BookRepository }

func NewBookHandler(books *repository.BookRepository) *BookHandler { return &BookHandler{books: books} }

func (h *BookHandler) List(c *gin.Context) {
	books, err := h.books.List(c.Request.Context())
	if err != nil {
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}

func (h *BookHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	book, err := h.books.Get(c.Request.Context(), id)
	if err != nil {
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (h *BookHandler) Create(c *gin.Context) {
	var input model.BookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "data buku tidak valid; release_year harus 1980-2024 dan total_page harus lebih dari 0"})
		return
	}
	book, err := h.books.Create(c.Request.Context(), input, actor(c))
	if err != nil {
		if repository.IsForeignKeyViolation(err) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "category_id tidak tersedia"})
			return
		}
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "buku berhasil ditambahkan", "data": book})
}

func (h *BookHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var input model.BookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "data buku tidak valid; release_year harus 1980-2024 dan total_page harus lebih dari 0"})
		return
	}
	book, err := h.books.Update(c.Request.Context(), id, input, actor(c))
	if err != nil {
		if repository.IsForeignKeyViolation(err) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "category_id tidak tersedia"})
			return
		}
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "buku berhasil diperbarui", "data": book})
}

func (h *BookHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	err := h.books.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "buku tidak ditemukan"})
			return
		}
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "buku berhasil dihapus"})
}

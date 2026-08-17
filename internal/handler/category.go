package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dienul/go-books/internal/repository"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categories *repository.CategoryRepository
	books      *repository.BookRepository
}

func NewCategoryHandler(categories *repository.CategoryRepository, books *repository.BookRepository) *CategoryHandler {
	return &CategoryHandler{categories: categories, books: books}
}

type categoryInput struct {
	Name string `json:"name" binding:"required,max=100"`
}

func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categories.List(c.Request.Context())
	if err != nil {
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *CategoryHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	category, err := h.categories.Get(c.Request.Context(), id)
	if err != nil {
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": category})
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var input categoryInput
	if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "nama kategori wajib diisi, maksimal 100 karakter"})
		return
	}
	category, err := h.categories.Create(c.Request.Context(), strings.TrimSpace(input.Name), actor(c))
	if err != nil {
		if repository.IsUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"message": "nama kategori sudah tersedia"})
			return
		}
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "kategori berhasil ditambahkan", "data": category})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var input categoryInput
	if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "nama kategori wajib diisi, maksimal 100 karakter"})
		return
	}
	category, err := h.categories.Update(c.Request.Context(), id, strings.TrimSpace(input.Name), actor(c))
	if err != nil {
		if repository.IsUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"message": "nama kategori sudah tersedia"})
			return
		}
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "kategori berhasil diperbarui", "data": category})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	err := h.categories.Delete(c.Request.Context(), id)
	if err != nil {
		if repository.IsForeignKeyViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"message": "kategori masih memiliki buku dan tidak dapat dihapus"})
			return
		}
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "kategori tidak ditemukan"})
			return
		}
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "kategori berhasil dihapus"})
}

func (h *CategoryHandler) Books(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	category, err := h.categories.Get(c.Request.Context(), id)
	if err != nil {
		databaseError(c, err)
		return
	}
	books, err := h.books.ListByCategory(c.Request.Context(), id)
	if err != nil {
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"category": category, "data": books})
}

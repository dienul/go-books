package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/dienul/go-books/internal/middleware"
	"github.com/dienul/go-books/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	users     *repository.UserRepository
	jwtSecret string
}

func NewAuthHandler(users *repository.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{users: users, jwtSecret: jwtSecret}
}

type loginInput struct {
	Username string `json:"username" binding:"required,max=100"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input loginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username dan password wajib diisi"})
		return
	}

	user, err := h.users.FindByUsername(c.Request.Context(), input.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "username atau password salah"})
			return
		}
		databaseError(c, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "username atau password salah"})
		return
	}

	now := time.Now()
	claims := middleware.Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login berhasil", "token": token, "expires_in": 86400})
}

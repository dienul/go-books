package config

import (
	"context"
	"fmt"
	"os"

	"github.com/dienul/go-books/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(ctx context.Context, users *repository.UserRepository) error {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	if username == "" || password == "" {
		hasUsers, err := users.HasUsers(ctx)
		if err != nil {
			return fmt.Errorf("memeriksa akun user: %w", err)
		}
		if hasUsers {
			return nil
		}
		return fmt.Errorf("ADMIN_USERNAME dan ADMIN_PASSWORD wajib diatur")
	}
	if len(password) < 8 {
		return fmt.Errorf("ADMIN_PASSWORD minimal 8 karakter")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("membuat hash password admin: %w", err)
	}
	if err := users.CreateAdmin(ctx, username, string(hash)); err != nil {
		return fmt.Errorf("membuat akun admin: %w", err)
	}
	return nil
}

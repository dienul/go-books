package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/dienul/go-books/internal/model"
	"github.com/lib/pq"
)

var ErrNotFound = errors.New("data tidak ditemukan")

func IsUniqueViolation(err error) bool {
	var state interface{ SQLState() string }
	if errors.As(err, &state) && state.SQLState() == "23505" {
		return true
	}
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

func IsForeignKeyViolation(err error) bool {
	var state interface{ SQLState() string }
	if errors.As(err, &state) && state.SQLState() == "23503" {
		return true
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23503" {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "foreign key constraint")
}

type UserRepository struct{ db *sql.DB }

func NewUserRepository(db *sql.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) HasUsers(ctx context.Context) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users)`).Scan(&exists)
	return exists, err
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, password, created_at, created_by, modified_at, modified_by
		FROM users WHERE username = $1`, username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.CreatedBy, &user.ModifiedAt, &user.ModifiedBy)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrNotFound
	}
	return user, err
}

func (r *UserRepository) CreateAdmin(ctx context.Context, username, passwordHash string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (username, password, created_by, modified_by)
		VALUES ($1, $2, $1, $1)
		ON CONFLICT (username) DO NOTHING`, username, passwordHash)
	return err
}

type CategoryRepository struct{ db *sql.DB }

func NewCategoryRepository(db *sql.DB) *CategoryRepository { return &CategoryRepository{db: db} }

func (r *CategoryRepository) List(ctx context.Context) ([]model.Category, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, created_at, created_by, modified_at, modified_by
		FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]model.Category, 0)
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.CreatedBy, &category.ModifiedAt, &category.ModifiedBy); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *CategoryRepository) Get(ctx context.Context, id int64) (model.Category, error) {
	var category model.Category
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, created_at, created_by, modified_at, modified_by
		FROM categories WHERE id = $1`, id,
	).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.CreatedBy, &category.ModifiedAt, &category.ModifiedBy)
	if errors.Is(err, sql.ErrNoRows) {
		return category, ErrNotFound
	}
	return category, err
}

func (r *CategoryRepository) Create(ctx context.Context, name, actor string) (model.Category, error) {
	var category model.Category
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO categories (name, created_by, modified_by)
		VALUES ($1, $2, $2)
		RETURNING id, name, created_at, created_by, modified_at, modified_by`, name, actor,
	).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.CreatedBy, &category.ModifiedAt, &category.ModifiedBy)
	return category, err
}

func (r *CategoryRepository) Update(ctx context.Context, id int64, name, actor string) (model.Category, error) {
	var category model.Category
	err := r.db.QueryRowContext(ctx, `
		UPDATE categories SET name = $1, modified_at = NOW(), modified_by = $2
		WHERE id = $3
		RETURNING id, name, created_at, created_by, modified_at, modified_by`, name, actor, id,
	).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.CreatedBy, &category.ModifiedAt, &category.ModifiedBy)
	if errors.Is(err, sql.ErrNoRows) {
		return category, ErrNotFound
	}
	return category, err
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

type BookRepository struct{ db *sql.DB }

func NewBookRepository(db *sql.DB) *BookRepository { return &BookRepository{db: db} }

const bookColumns = `id, title, description, image_url, release_year, price, total_page,
	thickness, category_id, created_at, created_by, modified_at, modified_by`

func scanBook(scanner interface{ Scan(...any) error }) (model.Book, error) {
	var book model.Book
	err := scanner.Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear,
		&book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt,
		&book.CreatedBy, &book.ModifiedAt, &book.ModifiedBy)
	return book, err
}

func (r *BookRepository) List(ctx context.Context) ([]model.Book, error) {
	return r.listQuery(ctx, "SELECT "+bookColumns+" FROM books ORDER BY id")
}

func (r *BookRepository) ListByCategory(ctx context.Context, categoryID int64) ([]model.Book, error) {
	return r.listQuery(ctx, "SELECT "+bookColumns+" FROM books WHERE category_id = $1 ORDER BY id", categoryID)
}

func (r *BookRepository) listQuery(ctx context.Context, query string, args ...any) ([]model.Book, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]model.Book, 0)
	for rows.Next() {
		book, err := scanBook(rows)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, rows.Err()
}

func (r *BookRepository) Get(ctx context.Context, id int64) (model.Book, error) {
	book, err := scanBook(r.db.QueryRowContext(ctx, "SELECT "+bookColumns+" FROM books WHERE id = $1", id))
	if errors.Is(err, sql.ErrNoRows) {
		return book, ErrNotFound
	}
	return book, err
}

func (r *BookRepository) Create(ctx context.Context, input model.BookInput, actor string) (model.Book, error) {
	query := `INSERT INTO books
		(title, description, image_url, release_year, price, total_page, thickness, category_id, created_by, modified_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$9) RETURNING ` + bookColumns
	return scanBook(r.db.QueryRowContext(ctx, query, input.Title, input.Description, input.ImageURL,
		input.ReleaseYear, input.Price, input.TotalPage, model.ThicknessFromPages(input.TotalPage), input.CategoryID, actor))
}

func (r *BookRepository) Update(ctx context.Context, id int64, input model.BookInput, actor string) (model.Book, error) {
	query := `UPDATE books SET title=$1, description=$2, image_url=$3, release_year=$4,
		price=$5, total_page=$6, thickness=$7, category_id=$8, modified_at=NOW(), modified_by=$9
		WHERE id=$10 RETURNING ` + bookColumns
	book, err := scanBook(r.db.QueryRowContext(ctx, query, input.Title, input.Description, input.ImageURL,
		input.ReleaseYear, input.Price, input.TotalPage, model.ThicknessFromPages(input.TotalPage), input.CategoryID, actor, id))
	if errors.Is(err, sql.ErrNoRows) {
		return book, ErrNotFound
	}
	return book, err
}

func (r *BookRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("membaca hasil delete: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

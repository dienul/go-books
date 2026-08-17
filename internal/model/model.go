package model

import "time"

type User struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedBy string    `json:"modified_by"`
}

type Category struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedBy string    `json:"modified_by"`
}

type Book struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ReleaseYear int       `json:"release_year"`
	Price       int       `json:"price"`
	TotalPage   int       `json:"total_page"`
	Thickness   string    `json:"thickness"`
	CategoryID  int64     `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	ModifiedAt  time.Time `json:"modified_at"`
	ModifiedBy  string    `json:"modified_by"`
}

type BookInput struct {
	Title       string `json:"title" binding:"required,max=255"`
	Description string `json:"description" binding:"required,max=1000"`
	ImageURL    string `json:"image_url" binding:"required,url,max=1000"`
	ReleaseYear int    `json:"release_year" binding:"required,gte=1980,lte=2024"`
	Price       int    `json:"price" binding:"gte=0"`
	TotalPage   int    `json:"total_page" binding:"required,gt=0"`
	CategoryID  int64  `json:"category_id" binding:"required,gt=0"`
}

func ThicknessFromPages(totalPage int) string {
	if totalPage > 100 {
		return "tebal"
	}
	return "tipis"
}

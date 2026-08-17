-- +migrate Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100) NOT NULL,
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_by VARCHAR(100) NOT NULL
);

CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100) NOT NULL,
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_by VARCHAR(100) NOT NULL
);

CREATE TABLE books (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    image_url VARCHAR(1000) NOT NULL,
    release_year INTEGER NOT NULL CHECK (release_year BETWEEN 1980 AND 2024),
    price INTEGER NOT NULL CHECK (price >= 0),
    total_page INTEGER NOT NULL CHECK (total_page > 0),
    thickness VARCHAR(10) NOT NULL CHECK (thickness IN ('tipis', 'tebal')),
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100) NOT NULL,
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_by VARCHAR(100) NOT NULL
);

CREATE INDEX idx_books_category_id ON books(category_id);

-- +migrate Down
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

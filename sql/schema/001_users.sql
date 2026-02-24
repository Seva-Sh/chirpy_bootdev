-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;

-- psql "postgres://postgres:2991@localhost:5432/chirpy?sslmode=disable"
-- goose postgres "postgres://postgres:2991@localhost:5432/chirpy?sslmode=disable" up
-- goose postgres "postgres://postgres:2991@localhost:5432/chirpy?sslmode=disable" down
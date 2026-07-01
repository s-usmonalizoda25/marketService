package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigration(ctx context.Context, pool *pgxpool.Pool) error {
	const query = `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(64) NOT NULL,
		email VARCHAR(128) UNIQUE NOT NULL,
		phone VARCHAR(12) NOT NULL,
		password_hash VARCHAR(256) NOT NULL,
		role VARCHAR(20) NOT NULL DEFAULT 'user',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
	);

	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id SERIAL PRIMARY KEY,
		token_hash VARCHAR(255) UNIQUE NOT NULL,
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		product VARCHAR(256) NOT NULL,
		price DECIMAL(10, 2) NOT NULL,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
		status VARCHAR(20) NOT NULL DEFAULT 'new',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS login_history (
    	id SERIAL PRIMARY KEY,
    	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    	ip VARCHAR(45) NOT NULL,
    	user_agent TEXT NOT NULL,
    	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	
	`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to run table migration: %w", err)
	}
	return nil
}

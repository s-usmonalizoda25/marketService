package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s-usmonalizoda25/marketService/internal/models"
	"github.com/s-usmonalizoda25/marketService/pkg/errs"
)

type UserRepo interface {
	CreateUser(ctx context.Context, req *models.RegisterRequest, passwordHash string) (uint, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserById(ctx context.Context, id uint) (*models.User, error)
	SaveRefreshTokens(ctx context.Context, userID uint, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	UpdateUser(ctx context.Context, id uint, name, phone string) error
	SoftDeleteUser(ctx context.Context, id uint) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	UpdateUserRole(ctx context.Context, id uint, role models.UserRole) error
}

type PostgresUserRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepo(pool *pgxpool.Pool) *PostgresUserRepo {
	return &PostgresUserRepo{pool: pool}
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, req *models.RegisterRequest, passwordHash string) (uint, error) {
	const query = `
	INSERT INTO users (name, email, phone, password_hash, role)
	VALUES($1, $2, $3, $4, 'user')
	RETURNING id;
	`
	var id uint
	err := r.pool.QueryRow(ctx, query, req.Name, req.Email, req.Phone, passwordHash).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const query = `
	SELECT id, name, email, phone, password_hash, role, created_at, deleted_at
	FROM users
	WHERE email = $1 and deleted_at IS NULL;
	`
	var user models.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.PassworHash, &user.Role, &user.CreatedAt, &user.DeletedAT,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepo) GetUserById(ctx context.Context, id uint) (*models.User, error) {
	const query = `
	SELECT id, name , email, phone, password_hash, role, created_at, deleted_at
	FROM users
	WHERE id = $1 AND deleted_at IS NULL;
	`
	var user models.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone,
		&user.PassworHash, &user.Role, &user.CreatedAt, &user.DeletedAT,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepo) SaveRefreshTokens(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	const query = `
        INSERT INTO refresh_tokens (token, expires_at, is_revoked, user_id)
        VALUES ($1, $2, FALSE, $3);
    `
	_, err := r.pool.Exec(ctx, query, token, expiresAt, userID)
	return err
}

func (r *PostgresUserRepo) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	const query = `
        SELECT id, token, expires_at, is_revoked, user_id, created_at
        FROM refresh_tokens
        WHERE token = $1;
    `
	var rt models.RefreshToken
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&rt.ID, &rt.Token, &rt.ExpiresAt, &rt.IsRevoked, &rt.UserID, &rt.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrTokenInvalid
		}
		return nil, err
	}
	return &rt, nil
}

func (r *PostgresUserRepo) RevokeRefreshToken(ctx context.Context, token string) error {
	const query = `
        UPDATE refresh_tokens
        SET is_revoked = TRUE
        WHERE token = $1;
    `
	_, err := r.pool.Exec(ctx, query, token)
	return err
}

func (r *PostgresUserRepo) UpdateUser(ctx context.Context, id uint, name, phone string) error {
	const query = `UPDATE users SET name = $1, phone = $2 WHERE id = $3 AND deleted_at IS NULL;`
	res, err := r.pool.Exec(ctx, query, name, phone, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}

func (r *PostgresUserRepo) SoftDeleteUser(ctx context.Context, id uint) error {
	const query = `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL;`
	res, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}

func (r *PostgresUserRepo) GetAllUsers(ctx context.Context) ([]models.User, error) {
	const query = `SELECT id, name, email, phone, role, created_at FROM users WHERE deleted_at IS NULL ORDER BY id ASC;`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *PostgresUserRepo) UpdateUserRole(ctx context.Context, id uint, role models.UserRole) error {
	const query = `UPDATE users SET role = $1 WHERE id = $2 AND deleted_at IS NULL;`
	res, err := r.pool.Exec(ctx, query, role, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}

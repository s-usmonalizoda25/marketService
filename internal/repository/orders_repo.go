package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s-usmonalizoda25/marketService/internal/models"
	"github.com/s-usmonalizoda25/marketService/pkg/errs"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, userID uint, req *models.CreateOrderRequest) (uint, error)
	GetOrderById(ctx context.Context, id uint) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, id uint, status models.OrderStatus) error
	DeleteOrder(ctx context.Context, id uint) error
	GetAllOrders(ctx context.Context) ([]models.Order, error)
}

type PostgresOrderRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresOrderRepo(pool *pgxpool.Pool) *PostgresOrderRepo {
	return &PostgresOrderRepo{pool: pool}
}

func (r *PostgresOrderRepo) CreateOrder(ctx context.Context, userID uint, req *models.CreateOrderRequest) (uint, error) {
	const query = `
        INSERT INTO orders (user_id, product, price, status)
        VALUES ($1, $2, $3, 'new')
        RETURNING id;
    `
	var id uint
	err := r.pool.QueryRow(ctx, query, userID, req.Product, req.Price).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresOrderRepo) GetOrderById(ctx context.Context, id uint) (*models.Order, error) {
	const query = `
        SELECT id, user_id, product, price, status, created_at
        FROM orders
        WHERE id = $1;
    `
	var order models.Order
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&order.ID, &order.UserID, &order.Product, &order.Price, &order.Status, &order.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *PostgresOrderRepo) GetOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	const query = `
        SELECT id, user_id, product, price, status, created_at
        FROM orders
        WHERE user_id = $1
        ORDER BY created_at DESC;
    `
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.Product, &order.Price, &order.Status, &order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *PostgresOrderRepo) UpdateOrderStatus(ctx context.Context, id uint, status models.OrderStatus) error {
	const query = `
        UPDATE orders
        SET status = $1
        WHERE id = $2;
    `
	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}
	return nil
}

func (r *PostgresOrderRepo) DeleteOrder(ctx context.Context, id uint) error {
	const query = `DELETE FROM orders WHERE id = $1;`
	res, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}
	return nil
}

func (r *PostgresOrderRepo) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	const query = `SELECT id, user_id, product, price, status, created_at FROM orders ORDER BY created_at DESC;`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Product, &o.Price, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

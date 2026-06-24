package models

import "time"

type OrderStatus string

const (
	StatusNew        OrderStatus = "new"
	StatusInProgress OrderStatus = "in_progress"
	StatusDone       OrderStatus = "done"
	StatusCanceled   OrderStatus = "canceled"
)

type Order struct {
	ID        uint        `json:"id"`
	Product   string      `json:"product"`
	Price     float64     `json:"price"`
	UserID    uint        `json:"user_id"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}

type CreateOrderRequest struct {
	Product string  `json:"product"`
	Price   float64 `json:"price"`
}

type UpdateOrderRequest struct {
	Status OrderStatus `json:"status"`
}

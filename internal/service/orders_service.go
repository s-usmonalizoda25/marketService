package service

import (
	"context"

	"github.com/s-usmonalizoda25/marketService/internal/models"
	"github.com/s-usmonalizoda25/marketService/internal/repository"
	"github.com/s-usmonalizoda25/marketService/pkg/errs"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID uint, req *models.CreateOrderRequest) (uint, error)
	GetOrderById(ctx context.Context, id uint) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, id uint, status models.OrderStatus) error

	// Новые методы под таблицу
	DeleteOrder(ctx context.Context, id uint) error
	GetAllOrders(ctx context.Context) ([]models.Order, error)
}

type MyOrderService struct {
	repo repository.OrderRepo
}

func NewMyOrderService(repo repository.OrderRepo) *MyOrderService {
	return &MyOrderService{repo: repo}
}

func (s *MyOrderService) CreateOrder(ctx context.Context, userID uint, req *models.CreateOrderRequest) (uint, error) {
	if req.Product == "" {
		return 0, errs.ErrEmptyProductStatus
	}

	if req.Price <= 0 {
		return 0, errs.ErrInvalidOrderPrice
	}

	return s.repo.CreateOrder(ctx, userID, req)
}

func (s *MyOrderService) GetOrderById(ctx context.Context, id uint) (*models.Order, error) {
	return s.repo.GetOrderById(ctx, id)
}

func (s *MyOrderService) GetOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	return s.repo.GetOrdersByUserID(ctx, userID)
}

func (s *MyOrderService) UpdateOrderStatus(ctx context.Context, id uint, status models.OrderStatus) error {
	currentOrder, err := s.repo.GetOrderById(ctx, id)
	if err != nil {
		return err
	}

	if currentOrder.Status == models.StatusDone && status == models.StatusCanceled {
		return errs.ErrCannotCancelOrder
	}

	return s.repo.UpdateOrderStatus(ctx, id, status)
}

func (s *MyOrderService) DeleteOrder(ctx context.Context, id uint) error {
	currentOrder, err := s.repo.GetOrderById(ctx, id)
	if err != nil {
		return err
	}

	if currentOrder.Status == models.StatusInProgress || currentOrder.Status == models.StatusDone {
		return errs.ErrCannotCancelOrder
	}

	return s.repo.DeleteOrder(ctx, id)
}

func (s *MyOrderService) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	return s.repo.GetAllOrders(ctx)
}

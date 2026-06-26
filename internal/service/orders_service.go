package service

import (
	"context"

	"github.com/s-usmonalizoda25/marketService/internal/models"
	"github.com/s-usmonalizoda25/marketService/internal/repository"
	"github.com/s-usmonalizoda25/marketService/pkg/errs"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID uint, req *models.CreateOrderRequest) (uint, error)
	GetOrderById(ctx context.Context, userID uint, orderID uint) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, userID uint, userRole string, id uint, status models.OrderStatus) error

	DeleteOrder(ctx context.Context, userId uint, id uint) error
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

func (s *MyOrderService) GetOrderById(ctx context.Context, userID uint, orderID uint) (*models.Order, error) {
	order, err := s.repo.GetOrderById(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errs.ErrAccessDenied
	}

	return order, nil
}

func (s *MyOrderService) CheckOrderOwnership(ctx context.Context, userID uint, orderID uint) error {
	order, err := s.repo.GetOrderById(ctx, orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errs.ErrAccessDenied
	}

	return nil
}

func (s *MyOrderService) GetOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	return s.repo.GetOrdersByUserID(ctx, userID)
}

func (s *MyOrderService) UpdateOrderStatus(ctx context.Context, userID uint, userRole string, id uint, status models.OrderStatus) error {
	order, err := s.repo.GetOrderById(ctx, id)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errs.ErrAccessDenied
	}
	if userRole != "admin" {
		if status != models.StatusCanceled {
			return errs.ErrAccessDenied
		}

		if order.Status != models.StatusNew {
			return errs.ErrCannotCancelOrder
		}
	}
	switch status {
	case models.StatusNew, models.StatusInProgress, models.StatusDone, models.StatusCanceled:
	default:
		return errs.ErrInvalidStatus
	}

	return s.repo.UpdateOrderStatus(ctx, id, status)
}

func (s *MyOrderService) DeleteOrder(ctx context.Context, userID uint, id uint) error {
	if err := s.CheckOrderOwnership(ctx, userID, id); err != nil {
		return err
	}
	return s.repo.DeleteOrder(ctx, id)
}

func (s *MyOrderService) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	return s.repo.GetAllOrders(ctx)
}

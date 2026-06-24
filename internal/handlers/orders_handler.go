package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/s-usmonalizoda25/marketService/internal/models"
	"github.com/s-usmonalizoda25/marketService/internal/service"
	"github.com/s-usmonalizoda25/marketService/pkg/logger"
)

type OrderHandler struct {
	service service.OrderService
	log     *logger.Logger
}

func NewOrderHandler(s service.OrderService, log *logger.Logger) *OrderHandler {
	return &OrderHandler{service: s, log: log}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(UserIDKey).(uint)
	if !ok {
		http.Error(w, "unauthorized: user id not found in context", http.StatusUnauthorized)
		return
	}

	id, err := h.service.CreateOrder(r.Context(), userID, &req)
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uint{"order_id": id})
}

func (h *OrderHandler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid order ID format", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrderById(r.Context(), uint(id))
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uint)
	if !ok {
		http.Error(w, "unauthorized: user id not found in context", http.StatusUnauthorized)
		return
	}

	orders, err := h.service.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid order ID format", http.StatusBadRequest)
		return
	}

	var req models.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		http.Error(w, "status parameter is required", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateOrderStatus(r.Context(), uint(id), req.Status)
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "order status updated successfully"}`))
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid order ID format", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteOrder(r.Context(), uint(id))
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "order deleted successfully"}`))
}


func (h *OrderHandler) AdminGetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders(r.Context())
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
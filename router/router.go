package router

import (
	"net/http"

	"github.com/s-usmonalizoda25/marketService/internal/handlers"
	"github.com/s-usmonalizoda25/marketService/internal/infrastructure/security"
)

func NewRouter(jwtManager *security.JWTManager, userHandler *handlers.UserHandler, orderHandler *handlers.OrderHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", userHandler.Register)
	mux.HandleFunc("POST /auth/login", userHandler.Login)
	mux.HandleFunc("POST /auth/refresh", userHandler.Refresh)

	//профиль юзера
	mux.HandleFunc("GET /users/me", handlers.AuthMiddleware(jwtManager, userHandler.GetProfile))
	mux.HandleFunc("PUT /users/me", handlers.AuthMiddleware(jwtManager, userHandler.UpdateProfile))
	mux.HandleFunc("DELETE /users/me", handlers.AuthMiddleware(jwtManager, userHandler.DeleteMe))

	//заказы польщователя
	mux.HandleFunc("POST /orders", handlers.AuthMiddleware(jwtManager, orderHandler.CreateOrder))
	mux.HandleFunc("GET /orders", handlers.AuthMiddleware(jwtManager, orderHandler.GetMyOrders))
	mux.HandleFunc("GET /orders/{id}", handlers.AuthMiddleware(jwtManager, orderHandler.GetOrderById))
	mux.HandleFunc("PUT /orders/{id}", handlers.AuthMiddleware(jwtManager, orderHandler.UpdateStatus))
	mux.HandleFunc("DELETE /orders/{id}", handlers.AuthMiddleware(jwtManager, orderHandler.DeleteOrder))

	//админ
	mux.HandleFunc("GET /admin/users", handlers.AuthMiddleware(jwtManager, handlers.AdminMiddleware(userHandler.AdminGetAllUsers)))
	mux.HandleFunc("PUT /admin/users/{id}/role", handlers.AuthMiddleware(jwtManager, handlers.AdminMiddleware(userHandler.AdminChangeRole)))
	mux.HandleFunc("GET /admin/orders", handlers.AuthMiddleware(jwtManager, handlers.AdminMiddleware(orderHandler.AdminGetAllOrders)))

	return mux
}


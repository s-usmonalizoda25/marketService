package handlers

import (
	"errors"
	"net/http"

	"github.com/s-usmonalizoda25/marketService/pkg/errs"
	"github.com/s-usmonalizoda25/marketService/pkg/logger"
	"go.uber.org/zap"
)

func HandleError(w http.ResponseWriter, log *logger.Logger, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, errs.ErrInvalidCredentials):
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid email or password"}`))

	case errors.Is(err, errs.WrongOldPassword):
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))

	case errors.Is(err, errs.ErrTokenInvalid) || errors.Is(err, errs.ErrTokenExpired):
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized or invalid token"}`))

	case errors.Is(err, errs.ErrUserNotFound):
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "user not found"}`))

	case errors.Is(err, errs.ErrOrderNotFound):
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "order not found"}`))

	case errors.Is(err, errs.ErrAccessDenied):
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "access denied"}`))

	case errors.Is(err, errs.ErrEmailAlreadyExists):
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error": "email already exists"}`))

	case errors.Is(err, errs.ErrInvalidOrderPrice) || errors.Is(err, errs.ErrEmptyProductStatus) || errors.Is(err, errs.ErrCannotCancelOrder):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))

	default:
		log.Error("Internal server error occurred", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
	}
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/s-usmonalizoda25/marketService/internal/models"
	"github.com/s-usmonalizoda25/marketService/internal/service"
	"github.com/s-usmonalizoda25/marketService/pkg/logger"
	"go.uber.org/zap"
)

type UserHandler struct {
	service service.UserService
	log     *logger.Logger
}

func NewUserHandler(s service.UserService, log *logger.Logger) *UserHandler {
	return &UserHandler{service: s, log: log}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Email) == "" || !strings.Contains(req.Email, "@") {
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "password must be at least 6 characters long", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "name cannot be empty", http.StatusBadRequest)
		return
	}

	err := h.service.RegisterRequest(req)
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "OTP code sent to email"}`))
}

func (h *UserHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.VerifyEmail(r.Context(), req)
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uint{"user_id": id})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.Login(r.Context(), &req)
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "refresh token is required", http.StatusBadRequest)
		return
	}

	res, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uint)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uint)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "profile updated successfully"}`))
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uint)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.service.DeleteMe(r.Context(), userID)
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "account deleted successfully"}`))
}

func (h *UserHandler) AdminGetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) AdminChangeRole(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user ID format", http.StatusBadRequest)
		return
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.ChangeRole(r.Context(), uint(id), req.Role)
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "user role updated successfully"}`))
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(UserIDKey).(uint)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("failed to decode change password request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.NewPassword) < 6 {
		http.Error(w, "new password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	err := h.service.ChangePassword(r.Context(), userId, req.OldPassword, req.NewPassword)
	if err != nil {

		h.log.Error("failed to change password",
			zap.Uint("user_id", userId),
			zap.Error(err),
		)

		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "password changed successfully"}`))
}

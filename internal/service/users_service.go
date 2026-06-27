package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/s-usmonalizoda25/marketService/internal/infrastructure/security"
	"github.com/s-usmonalizoda25/marketService/internal/models"
	"github.com/s-usmonalizoda25/marketService/internal/repository"
	"github.com/s-usmonalizoda25/marketService/pkg/cache"
	"github.com/s-usmonalizoda25/marketService/pkg/errs"
)

type UserService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (uint, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.TokenResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*models.TokenResponse, error)
	GetProfile(ctx context.Context, id uint) (*models.User, error)
	UpdateProfile(ctx context.Context, id uint, req *models.UpdateProfileRequest) error
	DeleteMe(ctx context.Context, id uint) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	ChangeRole(ctx context.Context, id uint, role string) error
	RegisterRequest(req models.RegisterReq) error
	VerifyEmail(ctx context.Context, req models.VerifyReq) (uint, error)
	ChangePassword(ctx context.Context, id uint, oldPass, newPass string) error
}

type MyUserService struct {
	repo       repository.UserRepo
	hasher     *security.BcryptHasher
	jwtManager *security.JWTManager
	cache      cache.MemoryCache
}

func NewMyUserService(repo repository.UserRepo, hasher *security.BcryptHasher, jwtManager *security.JWTManager, cache cache.MemoryCache) *MyUserService {
	return &MyUserService{
		repo:       repo,
		hasher:     hasher,
		jwtManager: jwtManager,
		cache:      cache,
	}
}

func (s *MyUserService) RegisterRequest(req models.RegisterReq) error {

	otp := "252525"
	hashedPassword, err := s.hasher.Hash(req.Password)
	if err != nil {
		return err
	}

	waitingUser := models.WaitingUser{
		User: models.User{
			Name:        req.Name,
			Email:       req.Email,
			Phone:       req.Phone,
			PassworHash: hashedPassword,
		},
		Otp:       otp,
		CreatedAt: time.Now(),
	}

	s.cache.Set(req.Email, waitingUser, 5*time.Minute)

	return nil
}

func (s *MyUserService) VerifyEmail(ctx context.Context, req models.VerifyReq) (uint, error) {
	val, found := s.cache.Get(req.Email)
	if !found {
		return 0, errs.ErrInvalidOtp
	}
	inwaiting := val.(models.WaitingUser)
	if inwaiting.Otp != req.Code {
		return 0, errs.ErrInvalidOtp
	}

	id, err := s.repo.CreateUser(ctx, &models.RegisterRequest{
		Name:     inwaiting.User.Name,
		Email:    inwaiting.User.Email,
		Phone:    inwaiting.User.Phone,
		Password: "",
	}, inwaiting.User.PassworHash)

	if err != nil {
		return 0, err
	}
	s.cache.Remove(req.Email)
	return id, nil
}

func (s *MyUserService) Register(ctx context.Context, req *models.RegisterRequest) (uint, error) {
	passwordHash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.CreateUser(ctx, req, passwordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, errs.ErrEmailAlreadyExists
		}
		return 0, err
	}

	return id, nil
}

func (s *MyUserService) Login(ctx context.Context, req *models.LoginRequest) (*models.TokenResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = s.hasher.Compare(user.PassworHash, req.Password)
	if err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshExpiresAt := time.Now().Add(30 * 24 * time.Hour)

	err = s.repo.SaveRefreshTokens(ctx, user.ID, refreshToken, refreshExpiresAt)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *MyUserService) Refresh(ctx context.Context, refreshToken string) (*models.TokenResponse, error) {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if rt.IsRevoked || time.Now().After(rt.ExpiresAt) {
		return nil, errs.ErrTokenInvalid
	}

	user, err := s.repo.GetUserById(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	err = s.repo.RevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	newAccessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshExpiresAt := time.Now().Add(30 * 24 * time.Hour)

	err = s.repo.SaveRefreshTokens(ctx, user.ID, newRefreshToken, refreshExpiresAt)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *MyUserService) GetProfile(ctx context.Context, id uint) (*models.User, error) {
	return s.repo.GetUserById(ctx, id)
}

func (s *MyUserService) UpdateProfile(ctx context.Context, id uint, req *models.UpdateProfileRequest) error {
	if req.Name == "" {
		return errs.ErrEmptyName
	}
	return s.repo.UpdateUser(ctx, id, req.Name, req.Phone)
}

func (s *MyUserService) DeleteMe(ctx context.Context, id uint) error {
	return s.repo.SoftDeleteUser(ctx, id)
}

func (s *MyUserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *MyUserService) ChangeRole(ctx context.Context, id uint, role string) error {
	var userRole models.UserRole
	switch role {
	case "admin":
		userRole = models.RoleAdmin
	case "user":
		userRole = models.RoleUser
	default:
		return errs.ErrInvalidRole
	}
	return s.repo.UpdateUserRole(ctx, id, userRole)
}

func (s *MyUserService) ChangePassword(ctx context.Context, id uint, oldPass, newPass string) error {
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return err
	}
	err = s.hasher.Compare(user.PassworHash, oldPass)
	if err != nil {
		return errs.WrongOldPassword
	}

	newHash, err := s.hasher.Hash(newPass)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, id, newHash)
}

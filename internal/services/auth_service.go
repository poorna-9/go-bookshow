package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/poorna-9/goshow/internal/models"
	"github.com/poorna-9/goshow/internal/repositories"
)

type AuthService struct {
	Repo            *repositories.UserRepository
	JWTSecret       string
	AdminSignupCode string
}

func NewAuthService(repo *repositories.UserRepository, jwtSecret, adminSignupCode string) *AuthService {
	return &AuthService{Repo: repo, JWTSecret: jwtSecret, AdminSignupCode: adminSignupCode}
}

type SignupInput struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	AdminCode string `json:"admin_code"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResult struct {
	Token string   `json:"token"`
	User  UserView `json:"user"`
}

type UserView struct {
	ID    uuid.UUID   `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email"`
	Role  models.Role `json:"role"`
}

func (s *AuthService) Signup(input SignupInput) (*AuthResult, error) {
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return nil, errors.New("name, email and password are required")
	}
	if len(input.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	if _, err := s.Repo.FindByEmail(input.Email); err == nil {
		return nil, errors.New("an account with this email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	role := models.RoleUser
	if input.AdminCode != "" && s.AdminSignupCode != "" && input.AdminCode == s.AdminSignupCode {
		role = models.RoleAdmin
	}

	user := &models.User{
		ID:           uuid.New(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
		Phone:        input.Phone,
		Role:         role,
	}
	if err := s.Repo.Create(user); err != nil {
		return nil, err
	}

	return s.buildAuthResult(user)
}

func (s *AuthService) Login(input LoginInput) (*AuthResult, error) {
	user, err := s.Repo.FindByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return s.buildAuthResult(user)
}

func (s *AuthService) buildAuthResult(user *models.User) (*AuthResult, error) {
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token: token,
		User: UserView{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    string(user.Role),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTSecret))
}

func (s *AuthService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.Repo.FindByID(id)
}

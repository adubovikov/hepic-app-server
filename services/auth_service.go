package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"hepic-app-server/v2/database"
	"hepic-app-server/v2/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	clickhouse *database.ClickHouseDB
	jwtSecret  string
	jwtExpire  int
}

// NewAuthService creates a new authentication service
func NewAuthService(clickhouse *database.ClickHouseDB, jwtSecret string, jwtExpire int) *AuthService {
	return &AuthService{
		clickhouse: clickhouse,
		jwtSecret:  jwtSecret,
		jwtExpire:  jwtExpire,
	}
}

// Register creates a new user
func (s *AuthService) Register(ctx context.Context, req *models.UserCreateRequest) (*models.User, error) {
	slog.Info("Registering new user", "username", req.Username, "email", req.Email)

	// Check if user already exists
	existingUser, err := s.GetUserByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with username %s already exists", req.Username)
	}

	// Check if email already exists
	existingEmail, err := s.GetUserByEmail(ctx, req.Email)
	if err == nil && existingEmail != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "user"
	}

	// Create user
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user to database
	userID, err := s.clickhouse.InsertUser(ctx, user)
	if err != nil {
		slog.Error("Failed to create user", "error", err, "username", req.Username)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = userID
	user.Password = "" // Don't return password

	slog.Info("User registered successfully", "user_id", userID, "username", req.Username)
	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	slog.Info("User login attempt", "username", req.Username)

	// Get user by username
	user, err := s.GetUserByUsername(ctx, req.Username)
	if err != nil {
		slog.Error("User not found", "username", req.Username, "error", err)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		slog.Error("Inactive user login attempt", "username", req.Username)
		return nil, fmt.Errorf("account is disabled")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		slog.Error("Invalid password", "username", req.Username)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, expiresAt, err := s.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		slog.Error("Failed to generate JWT", "error", err, "username", req.Username)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login
	now := time.Now()
	err = s.clickhouse.UpdateUserLastLogin(ctx, user.ID, now)
	if err != nil {
		slog.Warn("Failed to update last login", "error", err, "user_id", user.ID)
	}

	user.LastLogin = &now
	user.Password = "" // Don't return password

	slog.Info("User logged in successfully", "user_id", user.ID, "username", req.Username)

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *user,
	}, nil
}

// GenerateJWT generates a JWT token for a user
func (s *AuthService) GenerateJWT(userID int64, username, role string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(s.jwtExpire) * time.Hour)

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      expiresAt.Unix(),
		"iat":      now.Unix(),
		"jti":      s.generateJTI(), // JWT ID for token tracking
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateJWT validates a JWT token and returns the payload
func (s *AuthService) ValidateJWT(tokenString string) (*models.JWTPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Extract user information
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid role in token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid exp in token")
	}

	iat, ok := claims["iat"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid iat in token")
	}

	return &models.JWTPayload{
		UserID:   int64(userID),
		Username: username,
		Role:     role,
		Exp:      int64(exp),
		Iat:      int64(iat),
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	return s.clickhouse.GetUserByID(ctx, userID)
}

// GetUserByUsername retrieves a user by username
func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.clickhouse.GetUserByUsername(ctx, username)
}

// GetUserByEmail retrieves a user by email
func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.clickhouse.GetUserByEmail(ctx, email)
}

// UpdateUser updates a user
func (s *AuthService) UpdateUser(ctx context.Context, userID int64, req *models.UserUpdateRequest) (*models.User, error) {
	slog.Info("Updating user", "user_id", userID)

	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update fields if provided
	if req.Username != "" {
		// Check if new username is already taken
		existingUser, err := s.GetUserByUsername(ctx, req.Username)
		if err == nil && existingUser != nil && existingUser.ID != userID {
			return nil, fmt.Errorf("username %s is already taken", req.Username)
		}
		user.Username = req.Username
	}

	if req.Email != "" {
		// Check if new email is already taken
		existingUser, err := s.GetUserByEmail(ctx, req.Email)
		if err == nil && existingUser != nil && existingUser.ID != userID {
			return nil, fmt.Errorf("email %s is already taken", req.Email)
		}
		user.Email = req.Email
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	user.UpdatedAt = time.Now()

	err = s.clickhouse.UpdateUser(ctx, user)
	if err != nil {
		slog.Error("Failed to update user", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	user.Password = "" // Don't return password
	slog.Info("User updated successfully", "user_id", userID)

	return user, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID int64, req *models.UserChangePasswordRequest) error {
	slog.Info("Changing password", "user_id", userID)

	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	err = s.clickhouse.UpdateUserPassword(ctx, userID, string(hashedPassword))
	if err != nil {
		slog.Error("Failed to update password", "error", err, "user_id", userID)
		return fmt.Errorf("failed to update password: %w", err)
	}

	slog.Info("Password changed successfully", "user_id", userID)
	return nil
}

// GetUsers retrieves a list of users with pagination
func (s *AuthService) GetUsers(ctx context.Context, page, perPage int, role string) (*models.UserListResponse, error) {
	return s.clickhouse.GetUsers(ctx, page, perPage, role)
}

// GetUserStats retrieves user statistics
func (s *AuthService) GetUserStats(ctx context.Context) (*models.UserStats, error) {
	return s.clickhouse.GetUserStats(ctx)
}

// DeleteUser deletes a user
func (s *AuthService) DeleteUser(ctx context.Context, userID int64) error {
	slog.Info("Deleting user", "user_id", userID)

	err := s.clickhouse.DeleteUser(ctx, userID)
	if err != nil {
		slog.Error("Failed to delete user", "error", err, "user_id", userID)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	slog.Info("User deleted successfully", "user_id", userID)
	return nil
}

// generateJTI generates a unique JWT ID
func (s *AuthService) generateJTI() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
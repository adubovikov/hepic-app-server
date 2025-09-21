package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"hepic-app-server/v2/models"
	"hepic-app-server/v2/services"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserCreateRequest true "User registration data"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	slog.Info("User registration request",
		"method", c.Request().Method,
		"path", c.Request().URL.Path,
		"remote_addr", c.Request().RemoteAddr,
	)

	var req models.UserCreateRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		slog.Error("Validation failed", "error", err)
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	user, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		slog.Error("Registration failed", "error", err, "username", req.Username)
		return c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	slog.Info("User registered successfully", "user_id", user.ID, "username", user.Username)

	return c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    user,
		Message: "User registered successfully",
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	slog.Info("User login request",
		"method", c.Request().Method,
		"path", c.Request().URL.Path,
		"remote_addr", c.Request().RemoteAddr,
	)

	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		slog.Error("Validation failed", "error", err)
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	response, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		slog.Error("Login failed", "error", err, "username", req.Username)
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	slog.Info("User logged in successfully", "user_id", response.User.ID, "username", response.User.Username)

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
		Message: "Login successful",
	})
}

// Me godoc
// @Summary Get current user info
// @Description Get information about the currently authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c echo.Context) error {
	// Get user ID from JWT context (set by middleware)
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		slog.Error("User ID not found in context")
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	slog.Info("Get current user info", "user_id", userID)

	user, err := h.authService.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		slog.Error("Failed to get user", "error", err, "user_id", userID)
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get user information",
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    user,
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the current user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UserUpdateRequest true "User update data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	// Get user ID from JWT context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		slog.Error("User ID not found in context")
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	slog.Info("Update user profile", "user_id", userID)

	var req models.UserUpdateRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		slog.Error("Validation failed", "error", err)
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	user, err := h.authService.UpdateUser(c.Request().Context(), userID, &req)
	if err != nil {
		slog.Error("Failed to update user", "error", err, "user_id", userID)
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	slog.Info("User profile updated successfully", "user_id", userID)

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    user,
		Message: "Profile updated successfully",
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the current user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UserChangePasswordRequest true "Password change data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	// Get user ID from JWT context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		slog.Error("User ID not found in context")
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	slog.Info("Change user password", "user_id", userID)

	var req models.UserChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		slog.Error("Validation failed", "error", err)
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	err := h.authService.ChangePassword(c.Request().Context(), userID, &req)
	if err != nil {
		slog.Error("Failed to change password", "error", err, "user_id", userID)
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	slog.Info("Password changed successfully", "user_id", userID)

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// GetUsers godoc
// @Summary Get users list
// @Description Get a paginated list of users (admin only)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param role query string false "Filter by role"
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/users [get]
func (h *AuthHandler) GetUsers(c echo.Context) error {
	// Check if user is admin
	userRole, ok := c.Get("user_role").(string)
	if !ok || userRole != "admin" {
		slog.Error("Access denied - admin role required")
		return c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Access denied - admin role required",
		})
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	role := c.QueryParam("role")

	slog.Info("Get users list", "page", page, "per_page", perPage, "role", role)

	users, err := h.authService.GetUsers(c.Request().Context(), page, perPage, role)
	if err != nil {
		slog.Error("Failed to get users", "error", err)
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get users",
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    users,
	})
}

// GetUserStats godoc
// @Summary Get user statistics
// @Description Get user statistics (admin only)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/stats [get]
func (h *AuthHandler) GetUserStats(c echo.Context) error {
	// Check if user is admin
	userRole, ok := c.Get("user_role").(string)
	if !ok || userRole != "admin" {
		slog.Error("Access denied - admin role required")
		return c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Access denied - admin role required",
		})
	}

	slog.Info("Get user statistics")

	stats, err := h.authService.GetUserStats(c.Request().Context())
	if err != nil {
		slog.Error("Failed to get user stats", "error", err)
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get user statistics",
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    stats,
	})
}
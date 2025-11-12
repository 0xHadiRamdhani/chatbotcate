package handlers

import (
	"net/http"
	"time"

	"whatsapp-bot/internal/models"
	"whatsapp-bot/internal/services"
	"whatsapp-bot/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	serviceManager *services.ServiceManager
}

func NewAuthHandler(sm *services.ServiceManager) *AuthHandler {
	return &AuthHandler{serviceManager: sm}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	User         UserInfo `json:"user"`
}

type UserInfo struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	IsAdmin     bool      `json:"is_admin"`
	Points      int       `json:"points"`
	Level       int       `json:"level"`
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	PhoneNumber string `json:"phone_number"`
	DisplayName string `json:"display_name"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate user
	user, err := h.serviceManager.UserService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	accessToken, err := h.generateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := h.generateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Update last login
	h.serviceManager.UserService.UpdateLastActivity(user.ID)

	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400, // 24 hours
		User: UserInfo{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			IsAdmin:     user.IsAdmin,
			Points:      user.Points,
			Level:       user.Level,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	existingUser, _ := h.serviceManager.UserService.GetUserByUsername(req.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Check if email already exists
	existingUser, _ = h.serviceManager.UserService.GetUserByEmail(req.Email)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Create user
	user, err := h.serviceManager.UserService.CreateUser(req.Username, req.Email, req.Password, req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Update display name if provided
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
		h.serviceManager.UserService.UpdateUser(user.ID, map[string]interface{}{"display_name": req.DisplayName})
	}

	// Generate tokens
	accessToken, err := h.generateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := h.generateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400,
		User: UserInfo{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			IsAdmin:     user.IsAdmin,
			Points:      user.Points,
			Level:       user.Level,
		},
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.serviceManager.Config.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Extract user ID from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get user
	user, err := h.serviceManager.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Generate new access token
	accessToken, err := h.generateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken, // Return the same refresh token
		TokenType:    "Bearer",
		ExpiresIn:    86400,
		User: UserInfo{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			IsAdmin:     user.IsAdmin,
			Points:      user.Points,
			Level:       user.Level,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"username": user.Username,
		"email": user.Email,
		"is_admin": user.IsAdmin,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.serviceManager.Config.JWT.Secret))
}

func (h *AuthHandler) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.serviceManager.Config.JWT.Secret))
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	user, err := h.serviceManager.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	preferences, err := h.serviceManager.UserService.GetUserPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user preferences"})
		return
	}

	response := gin.H{
		"id":          user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"display_name": user.DisplayName,
		"phone_number": user.PhoneNumber,
		"avatar":      user.Avatar,
		"is_active":   user.IsActive,
		"is_admin":    user.IsAdmin,
		"points":      user.Points,
		"level":       user.Level,
		"last_login":  user.LastLoginAt,
		"created_at":  user.CreatedAt,
		"preferences": preferences,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate allowed fields
	allowedFields := []string{"display_name", "phone_number", "avatar"}
	for field := range updates {
		if !utils.Contains(allowedFields, field) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field: " + field})
			return
		}
	}

	err := h.serviceManager.UserService.UpdateUser(userID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (h *AuthHandler) UpdatePreferences(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var preferences map[string]interface{}
	if err := c.ShouldBindJSON(&preferences); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate allowed fields
	allowedFields := []string{"language", "timezone", "enable_games", "enable_business", "enable_utils", "enable_moderation"}
	for field := range preferences {
		if !utils.Contains(allowedFields, field) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field: " + field})
			return
		}
	}

	err := h.serviceManager.UserService.UpdateUserPreferences(userID, preferences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preferences updated successfully"})
}

func (h *AuthHandler) GetUserStats(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	stats, err := h.serviceManager.UserService.GetUserStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *AuthHandler) GetUserDashboard(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	dashboard, err := h.serviceManager.UserService.GetUserDashboard(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user dashboard"})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user
	user, err := h.serviceManager.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update password
	err = h.serviceManager.UserService.UpdateUser(userID, map[string]interface{}{"password": string(hashedPassword)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// In a real implementation, you might want to invalidate the token
	// For now, we'll just return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
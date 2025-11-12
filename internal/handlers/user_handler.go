package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/services"
	"kilocode.dev/whatsapp-bot/pkg/utils"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile gets user profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ResponseError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(userUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.ResponseSuccess(c, user)
}

// UpdateProfile updates user profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ResponseError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email" binding:"email"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
		Bio      string `json:"bio"`
		Language string `json:"language"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.UpdateProfile(userUUID, req.Name, req.Email, req.Phone, req.Avatar, req.Bio, req.Language)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, user)
}

// DeleteProfile deletes user profile
func (h *UserHandler) DeleteProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ResponseError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.userService.DeleteUser(userUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "User profile deleted successfully"})
}

// ChangePassword changes user password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ResponseError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.userService.ChangePassword(userUUID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Password changed successfully"})
}

// GetSettings gets user settings
func (h *UserHandler) GetSettings(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ResponseError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	settings, err := h.userService.GetUserSettings(userUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "Settings not found")
		return
	}

	utils.ResponseSuccess(c, settings)
}

// UpdateSettings updates user settings
func (h *UserHandler) UpdateSettings(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ResponseError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req struct {
		NotificationsEnabled bool   `json:"notifications_enabled"`
		AutoReplyEnabled     bool   `json:"auto_reply_enabled"`
		GameNotifications    bool   `json:"game_notifications"`
		BusinessAlerts       bool   `json:"business_alerts"`
		Language             string `json:"language"`
		Timezone             string `json:"timezone"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	settings, err := h.userService.UpdateUserSettings(userUUID, req.NotificationsEnabled, req.AutoReplyEnabled, req.GameNotifications, req.BusinessAlerts, req.Language, req.Timezone)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, settings)
}
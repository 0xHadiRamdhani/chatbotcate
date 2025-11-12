package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/services"
	"kilocode.dev/whatsapp-bot/pkg/utils"
)

type AutoReplyHandler struct {
	autoReplyService *services.AutoReplyService
}

func NewAutoReplyHandler(autoReplyService *services.AutoReplyService) *AutoReplyHandler {
	return &AutoReplyHandler{
		autoReplyService: autoReplyService,
	}
}

// GetAutoReplies gets all auto-replies
func (h *AutoReplyHandler) GetAutoReplies(c *gin.Context) {
	userID := c.Query("user_id")
	status := c.Query("status")
	page := 1
	limit := 10

	if userID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	autoReplies, total, err := h.autoReplyService.GetAutoReplies(userUUID, status, page, limit)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Paginate(total, page, limit)
	utils.ResponsePaginated(c, autoReplies, pagination)
}

// CreateAutoReply creates a new auto-reply
func (h *AutoReplyHandler) CreateAutoReply(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Trigger     string   `json:"trigger" binding:"required"`
		Response    string   `json:"response" binding:"required"`
		IsActive    bool     `json:"is_active"`
		MatchType   string   `json:"match_type" binding:"required,oneof=exact contains regex"`
		Keywords    []string `json:"keywords" binding:"required"`
		UserID      string   `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	autoReply, err := h.autoReplyService.CreateAutoReply(req.Name, req.Trigger, req.Response, req.IsActive, req.MatchType, req.Keywords, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, autoReply)
}

// GetAutoReply gets an auto-reply
func (h *AutoReplyHandler) GetAutoReply(c *gin.Context) {
	replyID := c.Param("reply_id")
	if replyID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Reply ID is required")
		return
	}

	replyUUID, err := uuid.Parse(replyID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid reply ID")
		return
	}

	autoReply, err := h.autoReplyService.GetAutoReply(replyUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "Auto-reply not found")
		return
	}

	utils.ResponseSuccess(c, autoReply)
}

// UpdateAutoReply updates an auto-reply
func (h *AutoReplyHandler) UpdateAutoReply(c *gin.Context) {
	replyID := c.Param("reply_id")
	if replyID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Reply ID is required")
		return
	}

	var req struct {
		Name      string   `json:"name"`
		Trigger   string   `json:"trigger"`
		Response  string   `json:"response"`
		IsActive  bool     `json:"is_active"`
		MatchType string   `json:"match_type" binding:"oneof=exact contains regex"`
		Keywords  []string `json:"keywords"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	replyUUID, err := uuid.Parse(replyID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid reply ID")
		return
	}

	autoReply, err := h.autoReplyService.UpdateAutoReply(replyUUID, req.Name, req.Trigger, req.Response, req.IsActive, req.MatchType, req.Keywords)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, autoReply)
}

// DeleteAutoReply deletes an auto-reply
func (h *AutoReplyHandler) DeleteAutoReply(c *gin.Context) {
	replyID := c.Param("reply_id")
	if replyID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Reply ID is required")
		return
	}

	replyUUID, err := uuid.Parse(replyID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid reply ID")
		return
	}

	err = h.autoReplyService.DeleteAutoReply(replyUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Auto-reply deleted successfully"})
}

// ToggleAutoReply toggles an auto-reply
func (h *AutoReplyHandler) ToggleAutoReply(c *gin.Context) {
	replyID := c.Param("reply_id")
	if replyID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Reply ID is required")
		return
	}

	replyUUID, err := uuid.Parse(replyID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid reply ID")
		return
	}

	autoReply, err := h.autoReplyService.ToggleAutoReply(replyUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, autoReply)
}
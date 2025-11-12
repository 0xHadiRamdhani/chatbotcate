package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/services"
	"kilocode.dev/whatsapp-bot/pkg/utils"
)

type ModerationHandler struct {
	moderationService *services.ModerationService
}

func NewModerationHandler(moderationService *services.ModerationService) *ModerationHandler {
	return &ModerationHandler{
		moderationService: moderationService,
	}
}

// GetBlockedUsers gets blocked users
func (h *ModerationHandler) GetBlockedUsers(c *gin.Context) {
	userID := c.Query("user_id")
	page := 1
	limit := 10

	var userUUID *uuid.UUID
	if userID != "" {
		uuid, err := uuid.Parse(userID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
			return
		}
		userUUID = &uuid
	}

	blockedUsers, total, err := h.moderationService.GetBlockedUsers(userUUID, page, limit)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Paginate(total, page, limit)
	utils.ResponsePaginated(c, blockedUsers, pagination)
}

// BlockUser blocks a user
func (h *ModerationHandler) BlockUser(c *gin.Context) {
	var req struct {
		UserID        string `json:"user_id" binding:"required"`
		BlockedUserID string `json:"blocked_user_id" binding:"required"`
		Reason        string `json:"reason"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	blockedUserID, err := uuid.Parse(req.BlockedUserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid blocked user ID")
		return
	}

	blockedUser, err := h.moderationService.BlockUser(userID, blockedUserID, req.Reason)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, blockedUser)
}

// UnblockUser unblocks a user
func (h *ModerationHandler) UnblockUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.moderationService.UnblockUser(userUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "User unblocked successfully"})
}

// GetReportedMessages gets reported messages
func (h *ModerationHandler) GetReportedMessages(c *gin.Context) {
	status := c.Query("status")
	page := 1
	limit := 10

	reportedMessages, total, err := h.moderationService.GetReportedMessages(status, page, limit)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Paginate(total, page, limit)
	utils.ResponsePaginated(c, reportedMessages, pagination)
}

// ReportMessage reports a message
func (h *ModerationHandler) ReportMessage(c *gin.Context) {
	var req struct {
		MessageID string `json:"message_id" binding:"required"`
		ReporterID string `json:"reporter_id" binding:"required"`
		Reason    string `json:"reason" binding:"required"`
		Details   string `json:"details"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	reporterID, err := uuid.Parse(req.ReporterID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid reporter ID")
		return
	}

	report, err := h.moderationService.ReportMessage(req.MessageID, reporterID, req.Reason, req.Details)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, report)
}

// UpdateReportStatus updates report status
func (h *ModerationHandler) UpdateReportStatus(c *gin.Context) {
	reportID := c.Param("report_id")
	if reportID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Report ID is required")
		return
	}

	var req struct {
		Status  string `json:"status" binding:"required,oneof=pending reviewed resolved dismissed"`
		Notes   string `json:"notes"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	reportUUID, err := uuid.Parse(reportID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid report ID")
		return
	}

	report, err := h.moderationService.UpdateReportStatus(reportUUID, req.Status, req.Notes)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, report)
}

// GetSpamDetections gets spam detections
func (h *ModerationHandler) GetSpamDetections(c *gin.Context) {
	userID := c.Query("user_id")
	page := 1
	limit := 10

	var userUUID *uuid.UUID
	if userID != "" {
		uuid, err := uuid.Parse(userID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
			return
		}
		userUUID = &uuid
	}

	spamDetections, total, err := h.moderationService.GetSpamDetections(userUUID, page, limit)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Paginate(total, page, limit)
	utils.ResponsePaginated(c, spamDetections, pagination)
}

// MarkAsSpam marks message as spam
func (h *ModerationHandler) MarkAsSpam(c *gin.Context) {
	var req struct {
		MessageID string `json:"message_id" binding:"required"`
		UserID    string `json:"user_id" binding:"required"`
		Reason    string `json:"reason"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	spamDetection, err := h.moderationService.MarkAsSpam(req.MessageID, userID, req.Reason)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, spamDetection)
}

// GetContentFilters gets content filters
func (h *ModerationHandler) GetContentFilters(c *gin.Context) {
	userID := c.Query("user_id")
	page := 1
	limit := 10

	var userUUID *uuid.UUID
	if userID != "" {
		uuid, err := uuid.Parse(userID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
			return
		}
		userUUID = &uuid
	}

	filters, total, err := h.moderationService.GetContentFilters(userUUID, page, limit)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Paginate(total, page, limit)
	utils.ResponsePaginated(c, filters, pagination)
}

// CreateContentFilter creates a content filter
func (h *ModerationHandler) CreateContentFilter(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required,oneof=keyword regex pattern"`
		Pattern     string `json:"pattern" binding:"required"`
		Action      string `json:"action" binding:"required,oneof=block replace flag"`
		Replacement string `json:"replacement"`
		IsActive    bool   `json:"is_active"`
		UserID      string `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	filter, err := h.moderationService.CreateContentFilter(req.Name, req.Type, req.Pattern, req.Action, req.Replacement, req.IsActive, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, filter)
}

// UpdateContentFilter updates a content filter
func (h *ModerationHandler) UpdateContentFilter(c *gin.Context) {
	filterID := c.Param("filter_id")
	if filterID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Filter ID is required")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Type        string `json:"type" binding:"oneof=keyword regex pattern"`
		Pattern     string `json:"pattern"`
		Action      string `json:"action" binding:"oneof=block replace flag"`
		Replacement string `json:"replacement"`
		IsActive    bool   `json:"is_active"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	filterUUID, err := uuid.Parse(filterID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid filter ID")
		return
	}

	filter, err := h.moderationService.UpdateContentFilter(filterUUID, req.Name, req.Type, req.Pattern, req.Action, req.Replacement, req.IsActive)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, filter)
}

// DeleteContentFilter deletes a content filter
func (h *ModerationHandler) DeleteContentFilter(c *gin.Context) {
	filterID := c.Param("filter_id")
	if filterID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Filter ID is required")
		return
	}

	filterUUID, err := uuid.Parse(filterID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid filter ID")
		return
	}

	err = h.moderationService.DeleteContentFilter(filterUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Content filter deleted successfully"})
}
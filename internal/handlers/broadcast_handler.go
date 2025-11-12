package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/services"
	"kilocode.dev/whatsapp-bot/pkg/utils"
)

type BroadcastHandler struct {
	broadcastService *services.BroadcastService
}

func NewBroadcastHandler(broadcastService *services.BroadcastService) *BroadcastHandler {
	return &BroadcastHandler{
		broadcastService: broadcastService,
	}
}

// GetBroadcasts gets all broadcasts
func (h *BroadcastHandler) GetBroadcasts(c *gin.Context) {
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

	broadcasts, total, err := h.broadcastService.GetBroadcasts(userUUID, status, page, limit)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Paginate(total, page, limit)
	utils.ResponsePaginated(c, broadcasts, pagination)
}

// CreateBroadcast creates a new broadcast
func (h *BroadcastHandler) CreateBroadcast(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Message     string   `json:"message" binding:"required"`
		Recipients  []string `json:"recipients" binding:"required"`
		ScheduleAt  string   `json:"schedule_at"`
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

	broadcast, err := h.broadcastService.CreateBroadcast(req.Name, req.Message, req.Recipients, req.ScheduleAt, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, broadcast)
}

// GetBroadcast gets a broadcast
func (h *BroadcastHandler) GetBroadcast(c *gin.Context) {
	broadcastID := c.Param("broadcast_id")
	if broadcastID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Broadcast ID is required")
		return
	}

	broadcastUUID, err := uuid.Parse(broadcastID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid broadcast ID")
		return
	}

	broadcast, err := h.broadcastService.GetBroadcast(broadcastUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "Broadcast not found")
		return
	}

	utils.ResponseSuccess(c, broadcast)
}

// UpdateBroadcast updates a broadcast
func (h *BroadcastHandler) UpdateBroadcast(c *gin.Context) {
	broadcastID := c.Param("broadcast_id")
	if broadcastID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Broadcast ID is required")
		return
	}

	var req struct {
		Name        string   `json:"name"`
		Message     string   `json:"message"`
		Recipients  []string `json:"recipients"`
		ScheduleAt  string   `json:"schedule_at"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	broadcastUUID, err := uuid.Parse(broadcastID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid broadcast ID")
		return
	}

	broadcast, err := h.broadcastService.UpdateBroadcast(broadcastUUID, req.Name, req.Message, req.Recipients, req.ScheduleAt)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, broadcast)
}

// DeleteBroadcast deletes a broadcast
func (h *BroadcastHandler) DeleteBroadcast(c *gin.Context) {
	broadcastID := c.Param("broadcast_id")
	if broadcastID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Broadcast ID is required")
		return
	}

	broadcastUUID, err := uuid.Parse(broadcastID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid broadcast ID")
		return
	}

	err = h.broadcastService.DeleteBroadcast(broadcastUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Broadcast deleted successfully"})
}

// SendBroadcast sends a broadcast
func (h *BroadcastHandler) SendBroadcast(c *gin.Context) {
	broadcastID := c.Param("broadcast_id")
	if broadcastID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Broadcast ID is required")
		return
	}

	broadcastUUID, err := uuid.Parse(broadcastID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid broadcast ID")
		return
	}

	err = h.broadcastService.SendBroadcast(broadcastUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Broadcast sent successfully"})
}

// GetBroadcastStats gets broadcast statistics
func (h *BroadcastHandler) GetBroadcastStats(c *gin.Context) {
	broadcastID := c.Param("broadcast_id")
	if broadcastID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Broadcast ID is required")
		return
	}

	broadcastUUID, err := uuid.Parse(broadcastID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid broadcast ID")
		return
	}

	stats, err := h.broadcastService.GetBroadcastStats(broadcastUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, stats)
}
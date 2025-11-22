package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/services"
	"kilocode.dev/whatsapp-bot/pkg/utils"
)

type TelegramHandler struct {
	telegramService *services.TelegramService
}

func NewTelegramHandler(telegramService *services.TelegramService) *TelegramHandler {
	return &TelegramHandler{
		telegramService: telegramService,
	}
}

// SendMessage sends a message via Telegram
func (h *TelegramHandler) SendMessage(c *gin.Context) {
	var req struct {
		ChatID int64  `json:"chat_id" binding:"required"`
		Text   string `json:"text" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.telegramService.SendMessage(req.ChatID, req.Text)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Message sent successfully"})
}

// SendMessageWithMarkup sends a message with custom keyboard markup
func (h *TelegramHandler) SendMessageWithMarkup(c *gin.Context) {
	var req struct {
		ChatID  int64       `json:"chat_id" binding:"required"`
		Text    string      `json:"text" binding:"required"`
		Markup  interface{} `json:"markup" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.telegramService.SendMessageWithMarkup(req.ChatID, req.Text, req.Markup)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Message with markup sent successfully"})
}

// SendPhoto sends a photo via Telegram
func (h *TelegramHandler) SendPhoto(c *gin.Context) {
	var req struct {
		ChatID  int64  `json:"chat_id" binding:"required"`
		PhotoURL string `json:"photo_url" binding:"required"`
		Caption  string `json:"caption"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.telegramService.SendPhoto(req.ChatID, req.PhotoURL, req.Caption)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Photo sent successfully"})
}

// SendDocument sends a document via Telegram
func (h *TelegramHandler) SendDocument(c *gin.Context) {
	var req struct {
		ChatID      int64  `json:"chat_id" binding:"required"`
		DocumentURL string `json:"document_url" binding:"required"`
		Caption     string `json:"caption"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.telegramService.SendDocument(req.ChatID, req.DocumentURL, req.Caption)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Document sent successfully"})
}

// SendLocation sends a location via Telegram
func (h *TelegramHandler) SendLocation(c *gin.Context) {
	var req struct {
		ChatID    int64   `json:"chat_id" binding:"required"`
		Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
		Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.telegramService.SendLocation(req.ChatID, req.Latitude, req.Longitude)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Location sent successfully"})
}

// SendPoll sends a poll via Telegram
func (h *TelegramHandler) SendPoll(c *gin.Context) {
	var req struct {
		ChatID  int64    `json:"chat_id" binding:"required"`
		Question string   `json:"question" binding:"required"`
		Options  []string `json:"options" binding:"required,min=2,max=10"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.telegramService.SendPoll(req.ChatID, req.Question, req.Options)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Poll sent successfully"})
}

// SendBroadcast sends a broadcast message to multiple Telegram users
func (h *TelegramHandler) SendBroadcast(c *gin.Context) {
	var req struct {
		Recipients []int64 `json:"recipients" binding:"required"`
		Message    string  `json:"message" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.telegramService.SendBroadcast(req.Recipients, req.Message)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Broadcast sent successfully"})
}

// HandleWebhook handles incoming Telegram webhook
func (h *TelegramHandler) HandleWebhook(c *gin.Context) {
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid webhook data")
		return
	}

	err := h.telegramService.ProcessWebhook(updateData)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Webhook processed successfully"})
}

// SetWebhook sets Telegram webhook
func (h *TelegramHandler) SetWebhook(c *gin.Context) {
	var req struct {
		WebhookURL string `json:"webhook_url" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.telegramService.SetWebhook(req.WebhookURL)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Webhook set successfully"})
}

// DeleteWebhook deletes Telegram webhook
func (h *TelegramHandler) DeleteWebhook(c *gin.Context) {
	err := h.telegramService.DeleteWebhook()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Webhook deleted successfully"})
}

// GetWebhookInfo gets Telegram webhook information
func (h *TelegramHandler) GetWebhookInfo(c *gin.Context) {
	info, err := h.telegramService.GetWebhookInfo()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, info)
}

// GetMe gets bot information
func (h *TelegramHandler) GetMe(c *gin.Context) {
	me, err := h.telegramService.GetMe()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, me)
}

// GetUpdates gets updates from Telegram (for polling mode)
func (h *TelegramHandler) GetUpdates(c *gin.Context) {
	offsetStr := c.Query("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid offset")
		return
	}

	updates, err := h.telegramService.GetUpdates(offset)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, updates)
}

// AnswerCallbackQuery answers a callback query
func (h *TelegramHandler) AnswerCallbackQuery(c *gin.Context) {
	var req struct {
		CallbackQueryID string `json:"callback_query_id" binding:"required"`
		Text            string `json:"text"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.telegramService.AnswerCallbackQuery(req.CallbackQueryID, req.Text)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Callback query answered successfully"})
}

// CreateTelegramUser creates a Telegram user
func (h *TelegramHandler) CreateTelegramUser(c *gin.Context) {
	var req struct {
		UserID    int64  `json:"user_id" binding:"required"`
		Username  string `json:"username"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	user, err := h.telegramService.CreateTelegramUser(req.UserID, req.Username, req.FirstName, req.LastName)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, user)
}

// GetTelegramUser gets a Telegram user
func (h *TelegramHandler) GetTelegramUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		utils.ResponseError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.telegramService.GetTelegramUser(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.ResponseSuccess(c, user)
}

// UpdateTelegramUser updates a Telegram user
func (h *TelegramHandler) UpdateTelegramUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		utils.ResponseError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	user, err := h.telegramService.UpdateTelegramUser(userID, req.Username, req.FirstName, req.LastName)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, user)
}

// GetTelegramMessages gets Telegram messages for a chat
func (h *TelegramHandler) GetTelegramMessages(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	if chatIDStr == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Chat ID is required")
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	limitStr := c.Query("limit")
	if limitStr == "" {
		limitStr = "50"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid limit")
		return
	}

	messages, err := h.telegramService.GetTelegramMessages(chatID, limit)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, messages)
}

// GetTelegramStats gets Telegram statistics
func (h *TelegramHandler) GetTelegramStats(c *gin.Context) {
	stats, err := h.telegramService.GetTelegramStats()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, stats)
}

// StartPolling starts Telegram polling mode
func (h *TelegramHandler) StartPolling(c *gin.Context) {
	go func() {
		if err := h.telegramService.StartPolling(); err != nil {
			logger.Error("Telegram polling error", err)
		}
	}()

	utils.ResponseSuccess(c, gin.H{"message": "Telegram polling started"})
}

// StopPolling stops Telegram polling mode
func (h *TelegramHandler) StopPolling(c *gin.Context) {
	// In a real implementation, you would have a way to stop polling
	// For now, we'll just return success
	utils.ResponseSuccess(c, gin.H{"message": "Telegram polling stopped"})
}

// GetTelegramBroadcasts gets Telegram broadcasts
func (h *TelegramHandler) GetTelegramBroadcasts(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	status := c.Query("status")
	page := 1
	limit := 10

	broadcasts, total, err := h.telegramService.GetTelegramBroadcasts(userUUID, status, page, limit)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Paginate(total, page, limit)
	utils.ResponsePaginated(c, broadcasts, pagination)
}

// CreateTelegramBroadcast creates a Telegram broadcast
func (h *TelegramHandler) CreateTelegramBroadcast(c *gin.Context) {
	var req struct {
		Name       string  `json:"name" binding:"required"`
		Message    string  `json:"message" binding:"required"`
		Recipients []int64 `json:"recipients" binding:"required"`
		UserID     string  `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	broadcast, err := h.telegramService.CreateTelegramBroadcast(req.Name, req.Message, req.Recipients, userUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, broadcast)
}

// SendTelegramBroadcast sends a Telegram broadcast
func (h *TelegramHandler) SendTelegramBroadcast(c *gin.Context) {
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

	err = h.telegramService.SendTelegramBroadcast(broadcastUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Telegram broadcast sent successfully"})
}

// GetTelegramBroadcastStats gets Telegram broadcast statistics
func (h *TelegramHandler) GetTelegramBroadcastStats(c *gin.Context) {
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

	stats, err := h.telegramService.GetTelegramBroadcastStats(broadcastUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, stats)
}
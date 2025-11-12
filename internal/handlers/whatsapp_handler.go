package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kilocode.dev/whatsapp-bot/internal/services"
	"kilocode.dev/whatsapp-bot/pkg/utils"
)

type WhatsAppHandler struct {
	whatsappService *services.WhatsAppService
}

func NewWhatsAppHandler(whatsappService *services.WhatsAppService) *WhatsAppHandler {
	return &WhatsAppHandler{
		whatsappService: whatsappService,
	}
}

// SendMessage sends a WhatsApp message
func (h *WhatsAppHandler) SendMessage(c *gin.Context) {
	var req struct {
		To      string `json:"to" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.whatsappService.SendMessage(req.To, req.Message)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Message sent successfully"})
}

// SendMedia sends a WhatsApp media message
func (h *WhatsAppHandler) SendMedia(c *gin.Context) {
	var req struct {
		To          string `json:"to" binding:"required"`
		MediaURL    string `json:"media_url" binding:"required"`
		MediaType   string `json:"media_type" binding:"required,oneof=image video audio document"`
		Caption     string `json:"caption"`
		Filename    string `json:"filename"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.whatsappService.SendMedia(req.To, req.MediaURL, req.MediaType, req.Caption, req.Filename)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Media sent successfully"})
}

// SendTemplate sends a WhatsApp template message
func (h *WhatsAppHandler) SendTemplate(c *gin.Context) {
	var req struct {
		To           string   `json:"to" binding:"required"`
		TemplateName string   `json:"template_name" binding:"required"`
		Language     string   `json:"language" binding:"required"`
		Parameters   []string `json:"parameters"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.whatsappService.SendTemplate(req.To, req.TemplateName, req.Language, req.Parameters)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Template sent successfully"})
}

// GetContacts gets WhatsApp contacts
func (h *WhatsAppHandler) GetContacts(c *gin.Context) {
	contacts, err := h.whatsappService.GetContacts()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, contacts)
}

// GetChats gets WhatsApp chats
func (h *WhatsAppHandler) GetChats(c *gin.Context) {
	limit := 50
	offset := 0

	chats, err := h.whatsappService.GetChats(limit, offset)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, chats)
}

// GetChatMessages gets messages from a chat
func (h *WhatsAppHandler) GetChatMessages(c *gin.Context) {
	chatID := c.Param("chat_id")
	if chatID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Chat ID is required")
		return
	}

	limit := 50
	offset := 0

	messages, err := h.whatsappService.GetChatMessages(chatID, limit, offset)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, messages)
}

// MarkAsRead marks messages as read
func (h *WhatsAppHandler) MarkAsRead(c *gin.Context) {
	var req struct {
		MessageIDs []string `json:"message_ids" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.whatsappService.MarkAsRead(req.MessageIDs)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Messages marked as read"})
}

// GetStatus gets WhatsApp status
func (h *WhatsAppHandler) GetStatus(c *gin.Context) {
	status, err := h.whatsappService.GetStatus()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, status)
}

// HandleWebhook handles WhatsApp webhook
func (h *WhatsAppHandler) HandleWebhook(c *gin.Context) {
	var webhookData map[string]interface{}
	if err := c.ShouldBindJSON(&webhookData); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid webhook data")
		return
	}

	err := h.whatsappService.ProcessWebhook(webhookData)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Webhook processed successfully"})
}

// VerifyWebhook verifies WhatsApp webhook
func (h *WhatsAppHandler) VerifyWebhook(c *gin.Context) {
	verifyToken := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if verifyToken == "" || challenge == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Missing verification parameters")
		return
	}

	isValid, err := h.whatsappService.VerifyWebhook(verifyToken)
	if err != nil || !isValid {
		utils.ResponseError(c, http.StatusForbidden, "Invalid verification token")
		return
	}

	c.String(http.StatusOK, challenge)
}
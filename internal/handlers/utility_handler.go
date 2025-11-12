package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/models"
	"kilocode.dev/whatsapp-bot/internal/services"
	"kilocode.dev/whatsapp-bot/pkg/utils"
)

type UtilityHandler struct {
	utilityService *services.UtilityService
}

func NewUtilityHandler(utilityService *services.UtilityService) *UtilityHandler {
	return &UtilityHandler{
		utilityService: utilityService,
	}
}

// CreateQRCode creates a QR code
func (h *UtilityHandler) CreateQRCode(c *gin.Context) {
	var req struct {
		Data    string `json:"data" binding:"required"`
		Size    int    `json:"size" binding:"min=1,max=1000"`
		Format  string `json:"format" binding:"oneof=png svg"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	qrCode, err := h.utilityService.CreateQRCode(req.Data, req.Size, req.Format)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"qr_code": qrCode,
		"format":  req.Format,
	})
}

// CreateShortLink creates a short link
func (h *UtilityHandler) CreateShortLink(c *gin.Context) {
	var req struct {
		OriginalURL string `json:"original_url" binding:"required,url"`
		CustomAlias string `json:"custom_alias"`
		ExpiryDays  int    `json:"expiry_days" binding:"min=1,max=365"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	shortLink, err := h.utilityService.CreateShortLink(req.OriginalURL, req.CustomAlias, req.ExpiryDays)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, shortLink)
}

// GetShortLink retrieves a short link
func (h *UtilityHandler) GetShortLink(c *gin.Context) {
	alias := c.Param("alias")
	if alias == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Alias is required")
		return
	}

	shortLink, err := h.utilityService.GetShortLink(alias)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "Short link not found")
		return
	}

	utils.ResponseSuccess(c, shortLink)
}

// ConvertCurrency converts currency
func (h *UtilityHandler) ConvertCurrency(c *gin.Context) {
	var req struct {
		Amount float64 `json:"amount" binding:"required,min=0"`
		From   string  `json:"from" binding:"required,len=3"`
		To     string  `json:"to" binding:"required,len=3"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	result, err := h.utilityService.ConvertCurrency(req.Amount, req.From, req.To)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"amount":   req.Amount,
		"from":     req.From,
		"to":       req.To,
		"result":   result,
		"rate":     result / req.Amount,
	})
}

// GetWeather gets weather information
func (h *UtilityHandler) GetWeather(c *gin.Context) {
	var req struct {
		City    string `json:"city" binding:"required"`
		Country string `json:"country"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	weather, err := h.utilityService.GetWeather(req.City, req.Country)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, weather)
}

// TranslateText translates text
func (h *UtilityHandler) TranslateText(c *gin.Context) {
	var req struct {
		Text       string `json:"text" binding:"required"`
		TargetLang string `json:"target_lang" binding:"required,len=2"`
		SourceLang string `json:"source_lang" binding:"len=2"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	translation, err := h.utilityService.TranslateText(req.Text, req.TargetLang, req.SourceLang)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"original_text": req.Text,
		"translated_text": translation,
		"source_lang":    req.SourceLang,
		"target_lang":    req.TargetLang,
	})
}

// GetLocationInfo gets location information
func (h *UtilityHandler) GetLocationInfo(c *gin.Context) {
	var req struct {
		Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
		Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	location, err := h.utilityService.GetLocationInfo(req.Latitude, req.Longitude)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, location)
}

// CreatePoll creates a poll
func (h *UtilityHandler) CreatePoll(c *gin.Context) {
	var req struct {
		Question string   `json:"question" binding:"required"`
		Options  []string `json:"options" binding:"required,min=2,max=10"`
		UserID   string   `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	poll, err := h.utilityService.CreatePoll(req.Question, req.Options, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, poll)
}

// VotePoll votes on a poll
func (h *UtilityHandler) VotePoll(c *gin.Context) {
	pollID := c.Param("poll_id")
	if pollID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Poll ID is required")
		return
	}

	var req struct {
		OptionID int    `json:"option_id" binding:"required"`
		UserID   string `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	pollUUID, err := uuid.Parse(pollID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.utilityService.VotePoll(pollUUID, req.OptionID, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Vote submitted successfully"})
}

// GetPollResults gets poll results
func (h *UtilityHandler) GetPollResults(c *gin.Context) {
	pollID := c.Param("poll_id")
	if pollID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Poll ID is required")
		return
	}

	pollUUID, err := uuid.Parse(pollID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	results, err := h.utilityService.GetPollResults(pollUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, results)
}

// CreateReminder creates a reminder
func (h *UtilityHandler) CreateReminder(c *gin.Context) {
	var req struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		RemindAt    time.Time `json:"remind_at" binding:"required"`
		UserID      string    `json:"user_id" binding:"required"`
		Repeat      string    `json:"repeat" binding:"oneof=none daily weekly monthly yearly"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	reminder, err := h.utilityService.CreateReminder(req.Title, req.Description, req.RemindAt, userID, req.Repeat)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, reminder)
}

// GetReminders gets user reminders
func (h *UtilityHandler) GetReminders(c *gin.Context) {
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

	reminders, err := h.utilityService.GetReminders(userUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, reminders)
}

// DeleteReminder deletes a reminder
func (h *UtilityHandler) DeleteReminder(c *gin.Context) {
	reminderID := c.Param("reminder_id")
	if reminderID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Reminder ID is required")
		return
	}

	reminderUUID, err := uuid.Parse(reminderID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid reminder ID")
		return
	}

	err = h.utilityService.DeleteReminder(reminderUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Reminder deleted successfully"})
}

// CreateNote creates a note
func (h *UtilityHandler) CreateNote(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		UserID   string `json:"user_id" binding:"required"`
		Category string `json:"category"`
		Tags     []string `json:"tags"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	note, err := h.utilityService.CreateNote(req.Title, req.Content, userID, req.Category, req.Tags)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, note)
}

// GetNotes gets user notes
func (h *UtilityHandler) GetNotes(c *gin.Context) {
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

	category := c.Query("category")
	tag := c.Query("tag")

	notes, err := h.utilityService.GetNotes(userUUID, category, tag)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, notes)
}

// UpdateNote updates a note
func (h *UtilityHandler) UpdateNote(c *gin.Context) {
	noteID := c.Param("note_id")
	if noteID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Note ID is required")
		return
	}

	var req struct {
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		Category string   `json:"category"`
		Tags     []string `json:"tags"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	noteUUID, err := uuid.Parse(noteID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid note ID")
		return
	}

	note, err := h.utilityService.UpdateNote(noteUUID, req.Title, req.Content, req.Category, req.Tags)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, note)
}

// DeleteNote deletes a note
func (h *UtilityHandler) DeleteNote(c *gin.Context) {
	noteID := c.Param("note_id")
	if noteID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Note ID is required")
		return
	}

	noteUUID, err := uuid.Parse(noteID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid note ID")
		return
	}

	err = h.utilityService.DeleteNote(noteUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Note deleted successfully"})
}

// SearchNotes searches notes
func (h *UtilityHandler) SearchNotes(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	query := c.Query("query")
	if query == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Search query is required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	notes, err := h.utilityService.SearchNotes(userUUID, query)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, notes)
}

// CreateTimer creates a timer
func (h *UtilityHandler) CreateTimer(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Duration int    `json:"duration" binding:"required,min=1,max=86400"` // max 24 hours
		UserID   string `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	timer, err := h.utilityService.CreateTimer(req.Name, req.Duration, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, timer)
}

// GetTimer gets timer status
func (h *UtilityHandler) GetTimer(c *gin.Context) {
	timerID := c.Param("timer_id")
	if timerID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Timer ID is required")
		return
	}

	timerUUID, err := uuid.Parse(timerID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid timer ID")
		return
	}

	timer, err := h.utilityService.GetTimer(timerUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "Timer not found")
		return
	}

	utils.ResponseSuccess(c, timer)
}

// StopTimer stops a timer
func (h *UtilityHandler) StopTimer(c *gin.Context) {
	timerID := c.Param("timer_id")
	if timerID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Timer ID is required")
		return
	}

	timerUUID, err := uuid.Parse(timerID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid timer ID")
		return
	}

	err = h.utilityService.StopTimer(timerUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Timer stopped successfully"})
}

// CreateFileUpload creates file upload
func (h *UtilityHandler) CreateFileUpload(c *gin.Context) {
	var req struct {
		Filename string `json:"filename" binding:"required"`
		FileSize int64  `json:"file_size" binding:"required,min=1,max=10485760"` // max 10MB
		FileType string `json:"file_type" binding:"required"`
		UserID   string `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	upload, err := h.utilityService.CreateFileUpload(req.Filename, req.FileSize, req.FileType, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, upload)
}

// GetFileUpload gets file upload info
func (h *UtilityHandler) GetFileUpload(c *gin.Context) {
	uploadID := c.Param("upload_id")
	if uploadID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Upload ID is required")
		return
	}

	uploadUUID, err := uuid.Parse(uploadID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid upload ID")
		return
	}

	upload, err := h.utilityService.GetFileUpload(uploadUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "File upload not found")
		return
	}

	utils.ResponseSuccess(c, upload)
}

// DeleteFileUpload deletes file upload
func (h *UtilityHandler) DeleteFileUpload(c *gin.Context) {
	uploadID := c.Param("upload_id")
	if uploadID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Upload ID is required")
		return
	}

	uploadUUID, err := uuid.Parse(uploadID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid upload ID")
		return
	}

	err = h.utilityService.DeleteFileUpload(uploadUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "File upload deleted successfully"})
}
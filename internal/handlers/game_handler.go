package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/services"
	"kilocode.dev/whatsapp-bot/pkg/utils"
)

type GameHandler struct {
	gameService *services.GameService
}

func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// GetGames gets available games
func (h *GameHandler) GetGames(c *gin.Context) {
	games, err := h.gameService.GetGames()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, games)
}

// StartGame starts a game
func (h *GameHandler) StartGame(c *gin.Context) {
	var req struct {
		GameType string `json:"game_type" binding:"required,oneof=trivia math word memory"`
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

	game, err := h.gameService.StartGame(req.GameType, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, game)
}

// PlayGame plays a game move
func (h *GameHandler) PlayGame(c *gin.Context) {
	gameID := c.Param("game_id")
	if gameID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Game ID is required")
		return
	}

	var req struct {
		UserID string `json:"user_id" binding:"required"`
		Answer string `json:"answer" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	gameUUID, err := uuid.Parse(gameID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid game ID")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	result, err := h.gameService.PlayGame(gameUUID, userID, req.Answer)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, result)
}

// GetLeaderboard gets game leaderboard
func (h *GameHandler) GetLeaderboard(c *gin.Context) {
	gameID := c.Param("game_id")
	if gameID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Game ID is required")
		return
	}

	gameUUID, err := uuid.Parse(gameID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid game ID")
		return
	}

	leaderboard, err := h.gameService.GetLeaderboard(gameUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, leaderboard)
}

// GetGameStats gets game statistics
func (h *GameHandler) GetGameStats(c *gin.Context) {
	gameID := c.Param("game_id")
	if gameID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Game ID is required")
		return
	}

	gameUUID, err := uuid.Parse(gameID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid game ID")
		return
	}

	stats, err := h.gameService.GetGameStats(gameUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, stats)
}

// StartTrivia starts a trivia game
func (h *GameHandler) StartTrivia(c *gin.Context) {
	var req struct {
		Category string `json:"category"`
		Difficulty string `json:"difficulty" binding:"oneof=easy medium hard"`
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

	trivia, err := h.gameService.StartTrivia(req.Category, req.Difficulty, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, trivia)
}

// AnswerTrivia answers a trivia question
func (h *GameHandler) AnswerTrivia(c *gin.Context) {
	var req struct {
		GameID   string `json:"game_id" binding:"required"`
		QuestionID string `json:"question_id" binding:"required"`
		Answer   string `json:"answer" binding:"required"`
		UserID   string `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	gameUUID, err := uuid.Parse(req.GameID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid game ID")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	result, err := h.gameService.AnswerTrivia(gameUUID, req.QuestionID, req.Answer, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, result)
}
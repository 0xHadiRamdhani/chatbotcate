package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"kilocode.dev/whatsapp-bot/internal/config"
	"kilocode.dev/whatsapp-bot/internal/database"
	"kilocode.dev/whatsapp-bot/internal/handlers"
	"kilocode.dev/whatsapp-bot/internal/services"
)

func setupTestServer() *gin.Engine {
	// Setup test configuration
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "postgres123",
		DBName:     "whatsapp_bot_test",
		RedisHost:  "localhost",
		RedisPort:  "6379",
		JWTSecret:  "test-secret-key",
	}

	// Setup test database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Auto migrate test database
	if err := db.AutoMigrate(); err != nil {
		panic("Failed to migrate test database")
	}

	// Create service manager
	serviceManager := services.NewServiceManager(db)

	// Setup routes
	router := gin.New()
	handlers.SetupRoutes(router, serviceManager)

	return router
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

func TestUserRegistration(t *testing.T) {
	router := setupTestServer()

	// Test user registration
	payload := map[string]interface{}{
		"name":     "Test User",
		"email":    "test@example.com",
		"phone":    "+1234567890",
		"password": "password123",
	}

	jsonPayload, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.NotEmpty(t, response["data"])
}

func TestUserLogin(t *testing.T) {
	router := setupTestServer()

	// First register a user
	registerPayload := map[string]interface{}{
		"name":     "Login Test User",
		"email":    "logintest@example.com",
		"phone":    "+1234567891",
		"password": "password123",
	}

	jsonRegister, _ := json.Marshal(registerPayload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonRegister))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Then test login
	loginPayload := map[string]interface{}{
		"email":    "logintest@example.com",
		"password": "password123",
	}

	jsonLogin, _ := json.Marshal(loginPayload)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.NotEmpty(t, response["token"])
	assert.NotEmpty(t, response["refresh_token"])
}

func TestProtectedEndpointWithoutToken(t *testing.T) {
	router := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/profile", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProtectedEndpointWithToken(t *testing.T) {
	router := setupTestServer()

	// First register and login to get token
	registerPayload := map[string]interface{}{
		"name":     "Protected Test User",
		"email":    "protectedtest@example.com",
		"phone":    "+1234567892",
		"password": "password123",
	}

	jsonRegister, _ := json.Marshal(registerPayload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonRegister))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	loginPayload := map[string]interface{}{
		"email":    "protectedtest@example.com",
		"password": "password123",
	}

	jsonLogin, _ := json.Marshal(loginPayload)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"].(string)

	// Test protected endpoint with token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.NotEmpty(t, response["data"])
}

func TestAutoReplyCreation(t *testing.T) {
	router := setupTestServer()

	// Register and login to get token
	registerPayload := map[string]interface{}{
		"name":     "AutoReply Test User",
		"email":    "autoreplytest@example.com",
		"phone":    "+1234567893",
		"password": "password123",
	}

	jsonRegister, _ := json.Marshal(registerPayload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonRegister))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	loginPayload := map[string]interface{}{
		"email":    "autoreplytest@example.com",
		"password": "password123",
	}

	jsonLogin, _ := json.Marshal(loginPayload)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"].(string)

	// Test auto-reply creation
	autoReplyPayload := map[string]interface{}{
		"name":       "Test AutoReply",
		"trigger":    "hello",
		"response":   "Hello! How can I help you?",
		"match_type": "exact",
		"keywords":   []string{"hello", "hi"},
		"user_id":    loginResponse["user_id"],
	}

	jsonAutoReply, _ := json.Marshal(autoReplyPayload)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auto-replies", bytes.NewBuffer(jsonAutoReply))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.NotEmpty(t, response["data"])
}

func TestBroadcastCreation(t *testing.T) {
	router := setupTestServer()

	// Register and login to get token
	registerPayload := map[string]interface{}{
		"name":     "Broadcast Test User",
		"email":    "broadcasttest@example.com",
		"phone":    "+1234567894",
		"password": "password123",
	}

	jsonRegister, _ := json.Marshal(registerPayload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonRegister))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	loginPayload := map[string]interface{}{
		"email":    "broadcasttest@example.com",
		"password": "password123",
	}

	jsonLogin, _ := json.Marshal(loginPayload)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"].(string)

	// Test broadcast creation
	broadcastPayload := map[string]interface{}{
		"name":       "Test Broadcast",
		"message":    "This is a test broadcast message",
		"recipients": []string{"+1234567890", "+1234567891"},
		"user_id":    loginResponse["user_id"],
	}

	jsonBroadcast, _ := json.Marshal(broadcastPayload)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/broadcasts", bytes.NewBuffer(jsonBroadcast))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.NotEmpty(t, response["data"])
}

func TestUtilityServices(t *testing.T) {
	router := setupTestServer()

	// Test QR code generation
	qrPayload := map[string]interface{}{
		"data":   "https://example.com",
		"size":   200,
		"format": "png",
	}

	jsonQR, _ := json.Marshal(qrPayload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/utils/qr-code", bytes.NewBuffer(jsonQR))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.NotEmpty(t, response["qr_code"])

	// Test short link creation
	shortLinkPayload := map[string]interface{}{
		"original_url": "https://example.com/very/long/url",
		"expiry_days":  7,
	}

	jsonShortLink, _ := json.Marshal(shortLinkPayload)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/utils/short-link", bytes.NewBuffer(jsonShortLink))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.NotEmpty(t, response["data"])
}

func TestRateLimiting(t *testing.T) {
	router := setupTestServer()

	// Test rate limiting on auth endpoint
	payload := map[string]interface{}{
		"email":    "ratelimit@example.com",
		"password": "password123",
	}

	jsonPayload, _ := json.Marshal(payload)

	// Make multiple requests quickly
	for i := 0; i < 15; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if i >= 10 {
			// Should be rate limited after 10 requests
			assert.Equal(t, http.StatusTooManyRequests, w.Code)
		}
	}
}

func TestValidation(t *testing.T) {
	router := setupTestServer()

	// Test invalid email format
	payload := map[string]interface{}{
		"name":     "Test User",
		"email":    "invalid-email",
		"phone":    "+1234567890",
		"password": "password123",
	}

	jsonPayload, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["error"], "email")
}

func TestErrorHandling(t *testing.T) {
	router := setupTestServer()

	// Test with invalid JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

func TestDatabaseConnection(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "postgres123",
		DBName:     "whatsapp_bot_test",
	}

	db, err := database.NewDatabase(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Test connection
	err = db.DB.Raw("SELECT 1").Error
	assert.NoError(t, err)
}

func TestServiceManager(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "postgres123",
		DBName:     "whatsapp_bot_test",
		RedisHost:  "localhost",
		RedisPort:  "6379",
		JWTSecret:  "test-secret-key",
	}

	db, err := database.NewDatabase(cfg)
	assert.NoError(t, err)

	serviceManager := services.NewServiceManager(db)
	assert.NotNil(t, serviceManager)
	assert.NotNil(t, serviceManager.UserService)
	assert.NotNil(t, serviceManager.WhatsAppService)
	assert.NotNil(t, serviceManager.AutoReplyService)
	assert.NotNil(t, serviceManager.BroadcastService)
	assert.NotNil(t, serviceManager.GameService)
	assert.NotNil(t, serviceManager.UtilityService)
	assert.NotNil(t, serviceManager.BusinessService)
	assert.NotNil(t, serviceManager.ModerationService)
	assert.NotNil(t, serviceManager.AnalyticsService)
	assert.NotNil(t, serviceManager.CleanupService)
	assert.NotNil(t, serviceManager.ReminderService)
}

func TestMain(t *testing.T) {
	// This test ensures the main function can be called without errors
	// In a real scenario, you might want to test the actual main function
	// or create a test version of it
	
	// For now, we'll just test that our setup works
	router := setupTestServer()
	assert.NotNil(t, router)
}
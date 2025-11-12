package handlers

import (
	"net/http"
	"strconv"

	"whatsapp-bot/internal/services"
	"whatsapp-bot/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UtilsHandler struct {
	serviceManager *services.ServiceManager
}

func NewUtilsHandler(sm *services.ServiceManager) *UtilsHandler {
	return &UtilsHandler{serviceManager: sm}
}

func (h *UtilsHandler) GetWeather(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	city := c.Param("city")

	weather, err := h.serviceManager.UtilityService.GetWeather(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weather data"})
		return
	}

	// Log analytics
	h.serviceManager.AnalyticsService.LogEvent(userID, "weather_checked", 1, map[string]interface{}{
		"city": city,
	})

	c.JSON(http.StatusOK, weather)
}

func (h *UtilsHandler) ConvertCurrency(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")

	if from == "" || to == "" || amountStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from, to, and amount parameters are required"})
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
		return
	}

	result, err := h.serviceManager.UtilityService.ConvertCurrency(from, to, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert currency"})
		return
	}

	// Log analytics
	h.serviceManager.AnalyticsService.LogEvent(userID, "currency_converted", 1, map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	})

	c.JSON(http.StatusOK, gin.H{
		"from":      from,
		"to":        to,
		"amount":    amount,
		"result":    result,
		"rate":      result / amount,
	})
}

func (h *UtilsHandler) Translate(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req struct {
		Text       string `json:"text" binding:"required"`
		TargetLang string `json:"target_lang" binding:"required"`
		SourceLang string `json:"source_lang"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Translate text (placeholder - implement actual translation)
	translated := map[string]interface{}{
		"original": req.Text,
		"translated": "Terjemahan dari " + req.Text,
		"source_lang": req.SourceLang,
		"target_lang": req.TargetLang,
	}

	// Log analytics
	h.serviceManager.AnalyticsService.LogEvent(userID, "text_translated", 1, map[string]interface{}{
		"source_lang": req.SourceLang,
		"target_lang": req.TargetLang,
		"text_length": len(req.Text),
	})

	c.JSON(http.StatusOK, translated)
}

func (h *UtilsHandler) GenerateQRCode(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	text := c.Query("text")

	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "text parameter is required"})
		return
	}

	qrCode, err := h.serviceManager.UtilityService.GenerateQRCode(text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Log analytics
	h.serviceManager.AnalyticsService.LogEvent(userID, "qr_code_generated", 1, map[string]interface{}{
		"text_length": len(text),
	})

	// Return QR code as base64 or file
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", "attachment; filename=qrcode.png")
	c.Data(http.StatusOK, "image/png", qrCode)
}

func (h *UtilsHandler) GetPrayerTimes(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	city := c.Query("city")

	if city == "" {
		city = "jakarta"
	}

	// Get prayer times (placeholder - implement actual prayer time API)
	prayerTimes := map[string]interface{}{
		"city": city,
		"date": time.Now().Format("YYYY-MM-DD"),
		"times": map[string]string{
			"fajr":    "04:30",
			"dhuhr":   "12:00",
			"asr":     "15:30",
			"maghrib": "18:00",
			"isha":    "19:30",
		},
		"source": "Kementerian Agama RI",
	}

	// Log analytics
	h.serviceManager.AnalyticsService.LogEvent(userID, "prayer_times_checked", 1, map[string]interface{}{
		"city": city,
	})

	c.JSON(http.StatusOK, prayerTimes)
}

func (h *UtilsHandler) GetNews(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	category := c.Query("category")
	limitStr := c.Query("limit")

	limit := 5
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	// Get news (placeholder - implement actual news API)
	news := []map[string]interface{}{
		{
			"id":       "news_1",
			"title":    "Teknologi AI Semakin Canggih",
			"summary":  "Kemajuan teknologi AI semakin pesat dalam berbagai bidang.",
			"category": "technology",
			"source":   "Tech News",
			"date":     time.Now().Format("YYYY-MM-DD"),
			"url":      "https://example.com/news/1",
		},
		{
			"id":       "news_2",
			"title":    "Ekonomi Indonesia Tumbuh Positif",
			"summary":  "Pertumbuhan ekonomi Indonesia mencapai 5% di kuartal ini.",
			"category": "economy",
			"source":   "Economic Times",
			"date":     time.Now().Format("YYYY-MM-DD"),
			"url":      "https://example.com/news/2",
		},
	}

	// Filter by category if specified
	if category != "" {
		filteredNews := []map[string]interface{}{}
		for _, item := range news {
			if item["category"] == category {
				filteredNews = append(filteredNews, item)
			}
		}
		news = filteredNews
	}

	// Limit results
	if len(news) > limit {
		news = news[:limit]
	}

	// Log analytics
	h.serviceManager.AnalyticsService.LogEvent(userID, "news_viewed", 1, map[string]interface{}{
		"category": category,
		"news_count": len(news),
	})

	c.JSON(http.StatusOK, gin.H{
		"news":  news,
		"count": len(news),
	})
}

func (h *UtilsHandler) GetStockPrice(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	symbol := c.Query("symbol")

	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol parameter is required"})
		return
	}

	// Get stock price (placeholder - implement actual stock API)
	stock := map[string]interface{}{
		"symbol":      symbol,
		"name":        "PT Bank Central Asia Tbk",
		"price":       10000.0,
		"change":      150.0,
		"change_percent": 1.5,
		"volume":      1000000,
		"market_cap":  2500000000000.0,
		"last_update": time.Now().Format("YYYY-MM-DD HH:mm:ss"),
	}

	// Log analytics
	h.serviceManager.AnalyticsService.LogEvent(userID, "stock_checked", 1, map[string]interface{}{
		"symbol": symbol,
	})

	c.JSON(http.StatusOK, stock)
}

func (h *UtilsHandler) GeneratePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	lengthStr := c.Query("length")

	length := 12
	if lengthStr != "" {
		length, _ = strconv.Atoi(lengthStr)
		if length < 8 || length > 32 {
			length = 12
		}
	}

	// Generate password (placeholder - implement actual password generation)
	password := map[string]interface{}{
		"password": "Abc123!@#XYZ",
		"length":   length,
		"strength": "strong",
		"contains": map[string]bool{
			"uppercase": true,
			"lowercase": true,
			"numbers":   true,
			"symbols":   true,
		},
	}

	// Log analytics
	h.serviceManager.AnalyticsService.LogEvent(userID, "password_generated", 1, map[string]interface{}{
		"length": length,
	})

	c.JSON(http.StatusOK, password)
}

func (h *UtilsHandler) ConvertUnits(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	from := c.Query("from")
	to := c.Query("to")
	valueStr := c.Query("value")

	if from == "" || to == "" || valueStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from, to, and value parameters are required"})
		return
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value format"})
		return
	}

	// Convert units (placeholder - implement actual unit conversion)
	result := map[string]interface{}{
		"from":     from,
		"to":       to,
		"value":    value,
		"result":   value * 1000, // Simple conversion placeholder
		"formula":  "1 " + from + " = 1000 " + to,
	}

	// Log analytics
	h.serviceManager.AnalyticsService.LogEvent(userID, "units_converted", 1, map[string]interface{}{
		"from":  from,
		"to":    to,
		"value": value,
	})

	c.JSON(http.StatusOK, result)
}

func (h *UtilsHandler) GetExchangeRates(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	base := c.Query("base")
	if base == "" {
		base = "IDR"
	}

	// Get exchange rates (placeholder - implement actual exchange rate API)
	rates := map[string]interface{}{
		"base": base,
		"rates": map[string]float64{
			"USD": 0.000067,
			"EUR": 0.000061,
			"GBP": 0.000052,
			"JPY": 0.0093,
			"SGD": 0.000089,
			"MYR": 0.00029,
		},
		"last_update": time.Now().Format("YYYY-MM-DD HH:mm:ss"),
	}

	// Log analytics
	h.serviceManager.AnalyticsService.LogEvent(userID, "exchange_rates_viewed", 1, map[string]interface{}{
		"base": base,
	})

	c.JSON(http.StatusOK, rates)
}

func (h *UtilsHandler) GetTimeZones(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	// Get time zones (placeholder - implement actual time zone data)
	timezones := []map[string]interface{}{
		{
			"timezone":   "Asia/Jakarta",
			"offset":     "+07:00",
			"current_time": time.Now().Format("YYYY-MM-DD HH:mm:ss"),
		},
		{
			"timezone":   "Asia/Singapore",
			"offset":     "+08:00",
			"current_time": time.Now().Add(time.Hour).Format("YYYY-MM-DD HH:mm:ss"),
		},
		{
			"timezone":   "UTC",
			"offset":     "+00:00",
			"current_time": time.Now().UTC().Format("YYYY-MM-DD HH:mm:ss"),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"timezones": timezones,
		"count":     len(timezones),
	})
}

func (h *UtilsHandler) GetCalendar(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	year := time.Now().Year()
	month := time.Now().Month()

	if yearStr != "" {
		year, _ = strconv.Atoi(yearStr)
	}
	if monthStr != "" {
		monthInt, _ := strconv.Atoi(monthStr)
		month = time.Month(monthInt)
	}

	// Get calendar data (placeholder - implement actual calendar data)
	calendar := map[string]interface{}{
		"year":  year,
		"month": month,
		"days":  30,
		"holidays": []map[string]interface{}{
			{
				"date":  "2024-01-01",
				"name":  "New Year",
				"type":  "national",
			},
			{
				"date":  "2024-12-25",
				"name":  "Christmas",
				"type":  "religious",
			},
		},
	}

	c.JSON(http.StatusOK, calendar)
}
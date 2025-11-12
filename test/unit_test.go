package test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"kilocode.dev/whatsapp-bot/internal/models"
	"kilocode.dev/whatsapp-bot/pkg/utils"
)

func TestUtils(t *testing.T) {
	t.Run("GenerateRandomString", func(t *testing.T) {
		str := utils.GenerateRandomString(10)
		assert.Equal(t, 10, len(str))
		assert.NotEmpty(t, str)
	})

	t.Run("GenerateRandomHex", func(t *testing.T) {
		hex := utils.GenerateRandomHex(16)
		assert.Equal(t, 16, len(hex))
		assert.NotEmpty(t, hex)
	})

	t.Run("ValidateEmail", func(t *testing.T) {
		assert.True(t, utils.ValidateEmail("test@example.com"))
		assert.True(t, utils.ValidateEmail("user.name@domain.co.id"))
		assert.False(t, utils.ValidateEmail("invalid-email"))
		assert.False(t, utils.ValidateEmail("@example.com"))
	})

	t.Run("ValidatePhone", func(t *testing.T) {
		assert.True(t, utils.ValidatePhone("+1234567890"))
		assert.True(t, utils.ValidatePhone("1234567890"))
		assert.False(t, utils.ValidatePhone("invalid"))
		assert.False(t, utils.ValidatePhone("123"))
	})

	t.Run("SanitizeString", func(t *testing.T) {
		input := "<script>alert('xss')</script>Hello World"
		sanitized := utils.SanitizeString(input)
		assert.NotContains(t, sanitized, "<script>")
		assert.Contains(t, sanitized, "Hello World")
	})

	t.Run("TruncateString", func(t *testing.T) {
		longString := "This is a very long string that should be truncated"
		truncated := utils.TruncateString(longString, 20)
		assert.Equal(t, 23, len(truncated)) // 20 + "..."
		assert.Contains(t, truncated, "...")
	})

	t.Run("FormatCurrency", func(t *testing.T) {
		formatted := utils.FormatCurrency(1234.56, "USD")
		assert.Equal(t, "USD 1234.56", formatted)
	})

	t.Run("FormatDateTime", func(t *testing.T) {
		now := time.Now()
		formatted := utils.FormatDateTime(now)
		assert.NotEmpty(t, formatted)
	})

	t.Run("IsWeekend", func(t *testing.T) {
		saturday := time.Date(2023, 1, 7, 0, 0, 0, 0, time.UTC) // Saturday
		assert.True(t, utils.IsWeekend(saturday))

		monday := time.Date(2023, 1, 9, 0, 0, 0, 0, time.UTC) // Monday
		assert.False(t, utils.IsWeekend(monday))
	})

	t.Run("CalculateAge", func(t *testing.T) {
		birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
		age := utils.CalculateAge(birthDate)
		assert.Greater(t, age, 30)
	})

	t.Run("GenerateOTP", func(t *testing.T) {
		otp := utils.GenerateOTP()
		assert.Equal(t, 6, len(otp))
		assert.NotEmpty(t, otp)
	})

	t.Run("Slugify", func(t *testing.T) {
		input := "Hello World! This is a Test."
		slug := utils.Slugify(input)
		assert.Equal(t, "hello-world-this-is-a-test", slug)
	})

	t.Run("ExtractDomain", func(t *testing.T) {
		domain := utils.ExtractDomain("user@example.com")
		assert.Equal(t, "example.com", domain)
	})

	t.Run("IsValidURL", func(t *testing.T) {
		assert.True(t, utils.IsValidURL("https://example.com"))
		assert.True(t, utils.IsValidURL("http://subdomain.example.com/path"))
		assert.False(t, utils.IsValidURL("not-a-url"))
	})

	t.Run("ExtractURLs", func(t *testing.T) {
		text := "Check out https://example.com and http://test.org"
		urls := utils.ExtractURLs(text)
		assert.Equal(t, 2, len(urls))
		assert.Contains(t, urls, "https://example.com")
		assert.Contains(t, urls, "http://test.org")
	})

	t.Run("RemoveDuplicates", func(t *testing.T) {
		input := []string{"apple", "banana", "apple", "orange", "banana"}
		result := utils.RemoveDuplicates(input)
		assert.Equal(t, 3, len(result))
		assert.Contains(t, result, "apple")
		assert.Contains(t, result, "banana")
		assert.Contains(t, result, "orange")
	})

	t.Run("Paginate", func(t *testing.T) {
		pagination := utils.Paginate(100, 2, 10)
		assert.Equal(t, 2, pagination["page"])
		assert.Equal(t, 10, pagination["limit"])
		assert.Equal(t, 10, pagination["offset"])
		assert.Equal(t, 100, pagination["total_items"])
		assert.Equal(t, 10, pagination["total_pages"])
		assert.True(t, pagination["has_next"].(bool))
		assert.True(t, pagination["has_prev"].(bool))
	})

	t.Run("MaskString", func(t *testing.T) {
		masked := utils.MaskString("sensitive_data", 3)
		assert.Contains(t, masked, "sen")
		assert.Contains(t, masked, "ata")
		assert.Contains(t, masked, "***")
	})

	t.Run("MaskEmail", func(t *testing.T) {
		masked := utils.MaskEmail("user@example.com")
		assert.Contains(t, masked, "us")
		assert.Contains(t, masked, "@example.com")
		assert.Contains(t, masked, "****")
	})

	t.Run("MaskPhone", func(t *testing.T) {
		masked := utils.MaskPhone("+1234567890")
		assert.Equal(t, "+12****90", masked)
	})

	t.Run("BytesToSize", func(t *testing.T) {
		assert.Equal(t, "100 B", utils.BytesToSize(100))
		assert.Equal(t, "1.0 KB", utils.BytesToSize(1024))
		assert.Equal(t, "1.0 MB", utils.BytesToSize(1024*1024))
		assert.Equal(t, "1.0 GB", utils.BytesToSize(1024*1024*1024))
	})

	t.Run("DurationToString", func(t *testing.T) {
		assert.Equal(t, "30 seconds", utils.DurationToString(30*time.Second))
		assert.Equal(t, "5 minutes", utils.DurationToString(5*time.Minute))
		assert.Equal(t, "2 hours", utils.DurationToString(2*time.Hour))
		assert.Equal(t, "1 days 0 hours", utils.DurationToString(24*time.Hour))
	})
}

func TestModels(t *testing.T) {
	t.Run("UserModel", func(t *testing.T) {
		user := &models.User{
			ID:        uuid.New(),
			Name:      "Test User",
			Email:     "test@example.com",
			Phone:     "+1234567890",
			Password:  "hashed_password",
			Role:      "user",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, user.ID)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "user", user.Role)
		assert.True(t, user.IsActive)
	})

	t.Run("AutoReplyModel", func(t *testing.T) {
		autoReply := &models.AutoReply{
			ID:        uuid.New(),
			Name:      "Test AutoReply",
			Trigger:   "hello",
			Response:  "Hello! How can I help?",
			IsActive:  true,
			MatchType: "exact",
			Keywords:  []string{"hello", "hi"},
			UserID:    uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, autoReply.ID)
		assert.Equal(t, "Test AutoReply", autoReply.Name)
		assert.Equal(t, "hello", autoReply.Trigger)
		assert.Equal(t, "Hello! How can I help?", autoReply.Response)
		assert.True(t, autoReply.IsActive)
		assert.Equal(t, "exact", autoReply.MatchType)
		assert.Equal(t, 2, len(autoReply.Keywords))
	})

	t.Run("BroadcastModel", func(t *testing.T) {
		broadcast := &models.Broadcast{
			ID:         uuid.New(),
			Name:       "Test Broadcast",
			Message:    "Test message",
			Recipients: []string{"+1234567890", "+1234567891"},
			Status:     "pending",
			UserID:     uuid.New(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		assert.NotNil(t, broadcast.ID)
		assert.Equal(t, "Test Broadcast", broadcast.Name)
		assert.Equal(t, "Test message", broadcast.Message)
		assert.Equal(t, 2, len(broadcast.Recipients))
		assert.Equal(t, "pending", broadcast.Status)
	})

	t.Run("ProductModel", func(t *testing.T) {
		product := &models.Product{
			ID:          uuid.New(),
			Name:        "Test Product",
			Description: "Test description",
			Price:       99.99,
			Category:    "Electronics",
			ImageURL:    "https://example.com/image.jpg",
			Stock:       100,
			UserID:      uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		assert.NotNil(t, product.ID)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, 99.99, product.Price)
		assert.Equal(t, "Electronics", product.Category)
		assert.Equal(t, 100, product.Stock)
	})

	t.Run("OrderModel", func(t *testing.T) {
		order := &models.Order{
			ID:            uuid.New(),
			CustomerName:  "John Doe",
			CustomerPhone: "+1234567890",
			TotalAmount:   199.98,
			Status:        "pending",
			UserID:        uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		assert.NotNil(t, order.ID)
		assert.Equal(t, "John Doe", order.CustomerName)
		assert.Equal(t, 199.98, order.TotalAmount)
		assert.Equal(t, "pending", order.Status)
	})

	t.Run("GameModel", func(t *testing.T) {
		game := &models.Game{
			ID:        uuid.New(),
			GameType:  "trivia",
			Status:    "active",
			UserID:    uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, game.ID)
		assert.Equal(t, "trivia", game.GameType)
		assert.Equal(t, "active", game.Status)
	})

	t.Run("PollModel", func(t *testing.T) {
		poll := &models.Poll{
			ID:        uuid.New(),
			Question:  "What is your favorite color?",
			Options:   []string{"Red", "Blue", "Green"},
			UserID:    uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, poll.ID)
		assert.Equal(t, "What is your favorite color?", poll.Question)
		assert.Equal(t, 3, len(poll.Options))
	})

	t.Run("ReminderModel", func(t *testing.T) {
		reminder := &models.Reminder{
			ID:          uuid.New(),
			Title:       "Test Reminder",
			Description: "Test description",
			RemindAt:    time.Now().Add(1 * time.Hour),
			Repeat:      "daily",
			IsActive:    true,
			UserID:      uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		assert.NotNil(t, reminder.ID)
		assert.Equal(t, "Test Reminder", reminder.Title)
		assert.Equal(t, "daily", reminder.Repeat)
		assert.True(t, reminder.IsActive)
	})

	t.Run("NoteModel", func(t *testing.T) {
		note := &models.Note{
			ID:        uuid.New(),
			Title:     "Test Note",
			Content:   "Test content",
			Category:  "Personal",
			Tags:      []string{"important", "todo"},
			UserID:    uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, note.ID)
		assert.Equal(t, "Test Note", note.Title)
		assert.Equal(t, "Personal", note.Category)
		assert.Equal(t, 2, len(note.Tags))
	})

	t.Run("TimerModel", func(t *testing.T) {
		timer := &models.Timer{
			ID:        uuid.New(),
			Name:      "Test Timer",
			Duration:  300,
			UserID:    uuid.New(),
			StartTime: time.Now(),
			EndTime:   time.Now().Add(5 * time.Minute),
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, timer.ID)
		assert.Equal(t, "Test Timer", timer.Name)
		assert.Equal(t, 300, timer.Duration)
		assert.True(t, timer.IsActive)
	})

	t.Run("CustomerModel", func(t *testing.T) {
		customer := &models.Customer{
			ID:        uuid.New(),
			Name:      "John Doe",
			Email:     "john@example.com",
			Phone:     "+1234567890",
			Address:   "123 Main St, City",
			UserID:    uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, customer.ID)
		assert.Equal(t, "John Doe", customer.Name)
		assert.Equal(t, "john@example.com", customer.Email)
		assert.Equal(t, "+1234567890", customer.Phone)
	})

	t.Run("InvoiceModel", func(t *testing.T) {
		invoice := &models.Invoice{
			ID:          uuid.New(),
			OrderID:     uuid.New(),
			Amount:      299.99,
			Description: "Test invoice",
			DueDate:     time.Now().Add(30 * 24 * time.Hour),
			Status:      "pending",
			UserID:      uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		assert.NotNil(t, invoice.ID)
		assert.Equal(t, 299.99, invoice.Amount)
		assert.Equal(t, "pending", invoice.Status)
	})

	t.Run("PaymentModel", func(t *testing.T) {
		payment := &models.Payment{
			ID:            uuid.New(),
			InvoiceID:     uuid.New(),
			Amount:        299.99,
			PaymentMethod: "credit_card",
			Reference:     "REF123",
			Status:        "completed",
			UserID:        uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		assert.NotNil(t, payment.ID)
		assert.Equal(t, 299.99, payment.Amount)
		assert.Equal(t, "credit_card", payment.PaymentMethod)
		assert.Equal(t, "completed", payment.Status)
	})

	t.Run("ContentFilterModel", func(t *testing.T) {
		filter := &models.ContentFilter{
			ID:          uuid.New(),
			Name:        "Profanity Filter",
			Type:        "keyword",
			Pattern:     "badword",
			Action:      "block",
			Replacement: "",
			IsActive:    true,
			UserID:      uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		assert.NotNil(t, filter.ID)
		assert.Equal(t, "Profanity Filter", filter.Name)
		assert.Equal(t, "keyword", filter.Type)
		assert.Equal(t, "block", filter.Action)
		assert.True(t, filter.IsActive)
	})

	t.Run("ShortLinkModel", func(t *testing.T) {
		shortLink := &models.ShortLink{
			ID:          uuid.New(),
			OriginalURL: "https://example.com/very/long/url",
			Alias:       "abc123",
			ClickCount:  42,
			ExpiryAt:    time.Now().Add(7 * 24 * time.Hour),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		assert.NotNil(t, shortLink.ID)
		assert.Equal(t, "https://example.com/very/long/url", shortLink.OriginalURL)
		assert.Equal(t, "abc123", shortLink.Alias)
		assert.Equal(t, 42, shortLink.ClickCount)
	})

	t.Run("FileUploadModel", func(t *testing.T) {
		fileUpload := &models.FileUpload{
			ID:        uuid.New(),
			Filename:  "document.pdf",
			FileSize:  1024 * 1024,
			FileType:  "application/pdf",
			UploadURL: "/uploads/user123/document.pdf",
			Status:    "completed",
			UserID:    uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, fileUpload.ID)
		assert.Equal(t, "document.pdf", fileUpload.Filename)
		assert.Equal(t, 1024*1024, fileUpload.FileSize)
		assert.Equal(t, "completed", fileUpload.Status)
	})
}

func TestValidation(t *testing.T) {
	t.Run("EmailValidation", func(t *testing.T) {
		validEmails := []string{
			"test@example.com",
			"user.name@domain.co.id",
			"test+tag@example.org",
		}

		invalidEmails := []string{
			"invalid-email",
			"@example.com",
			"user@",
			"user@.com",
			"user@domain",
		}

		for _, email := range validEmails {
			assert.True(t, utils.ValidateEmail(email), "Email should be valid: %s", email)
		}

		for _, email := range invalidEmails {
			assert.False(t, utils.ValidateEmail(email), "Email should be invalid: %s", email)
		}
	})

	t.Run("PhoneValidation", func(t *testing.T) {
		validPhones := []string{
			"+1234567890",
			"1234567890",
			"+62-812-3456-7890",
			"081234567890",
		}

		invalidPhones := []string{
			"invalid",
			"123",
			"abc123",
			"+123abc456",
		}

		for _, phone := range validPhones {
			assert.True(t, utils.ValidatePhone(phone), "Phone should be valid: %s", phone)
		}

		for _, phone := range invalidPhones {
			assert.False(t, utils.ValidatePhone(phone), "Phone should be invalid: %s", phone)
		}
	})

	t.Run("URLValidation", func(t *testing.T) {
		validURLs := []string{
			"https://example.com",
			"http://subdomain.example.com/path",
			"https://example.com/path/to/resource",
			"http://localhost:8080",
		}

		invalidURLs := []string{
			"not-a-url",
			"example.com",
			"http://",
			"https://",
		}

		for _, url := range validURLs {
			assert.True(t, utils.IsValidURL(url), "URL should be valid: %s", url)
		}

		for _, url := range invalidURLs {
			assert.False(t, utils.IsValidURL(url), "URL should be invalid: %s", url)
		}
	})
}

func TestSecurity(t *testing.T) {
	t.Run("SanitizeString", func(t *testing.T) {
		dangerousInputs := []string{
			"<script>alert('xss')</script>",
			"SELECT * FROM users",
			"DROP TABLE users",
			"<img src='x' onerror='alert(1)'>",
		}

		for _, input := range dangerousInputs {
			sanitized := utils.SanitizeString(input)
			assert.NotContains(t, sanitized, "<script>")
			assert.NotContains(t, sanitized, "SELECT")
			assert.NotContains(t, sanitized, "DROP")
			assert.NotContains(t, sanitized, "onerror")
		}
	})

	t.Run("MaskSensitiveData", func(t *testing.T) {
		email := "user@example.com"
		maskedEmail := utils.MaskEmail(email)
		assert.Contains(t, maskedEmail, "****")
		assert.Contains(t, maskedEmail, "@example.com")

		phone := "+1234567890"
		maskedPhone := utils.MaskPhone(phone)
		assert.Contains(t, maskedPhone, "****")

		sensitive := "sensitive_data_123"
		masked := utils.MaskString(sensitive, 3)
		assert.Contains(t, masked, "***")
	})
}

func TestPagination(t *testing.T) {
	t.Run("PaginationCalculation", func(t *testing.T) {
		// Test various pagination scenarios
		testCases := []struct {
			totalItems int
			page       int
			limit      int
			expected   map[string]interface{}
		}{
			{
				totalItems: 100,
				page:       1,
				limit:      10,
				expected: map[string]interface{}{
					"page":        1,
					"limit":       10,
					"offset":      0,
					"total_items": 100,
					"total_pages": 10,
					"has_next":    true,
					"has_prev":    false,
				},
			},
			{
				totalItems: 100,
				page:       5,
				limit:      10,
				expected: map[string]interface{}{
					"page":        5,
					"limit":       10,
					"offset":      40,
					"total_items": 100,
					"total_pages": 10,
					"has_next":    true,
					"has_prev":    true,
				},
			},
			{
				totalItems: 100,
				page:       10,
				limit:      10,
				expected: map[string]interface{}{
					"page":        10,
					"limit":       10,
					"offset":      90,
					"total_items": 100,
					"total_pages": 10,
					"has_next":    false,
					"has_prev":    true,
				},
			},
		}

		for _, tc := range testCases {
			result := utils.Paginate(tc.totalItems, tc.page, tc.limit)
			
			for key, expectedValue := range tc.expected {
				assert.Equal(t, expectedValue, result[key], "Pagination %s mismatch", key)
			}
		}
	})

	t.Run("EdgeCases", func(t *testing.T) {
		// Test edge cases
		result := utils.Paginate(0, 1, 10)
		assert.Equal(t, 1, result["total_pages"])
		assert.False(t, result["has_next"].(bool))
		assert.False(t, result["has_prev"].(bool))

		result = utils.Paginate(5, 1, 10)
		assert.Equal(t, 1, result["total_pages"])
		assert.False(t, result["has_next"].(bool))
		assert.False(t, result["has_prev"].(bool))

		result = utils.Paginate(100, 0, 10) // Invalid page
		assert.Equal(t, 1, result["page"])

		result = utils.Paginate(100, 20, 10) // Page beyond total
		assert.Equal(t, 10, result["page"])
	})
}

func TestTimeFunctions(t *testing.T) {
	t.Run("TimeFormatting", func(t *testing.T) {
		testTime := time.Date(2023, 12, 25, 14, 30, 45, 0, time.UTC)
		
		dateStr := utils.FormatDate(testTime)
		assert.NotEmpty(t, dateStr)
		
		dateTimeStr := utils.FormatDateTime(testTime)
		assert.NotEmpty(t, dateTimeStr)
	})

	t.Run("TimeBoundaries", func(t *testing.T) {
		now := time.Now()
		
		startOfDay := utils.GetStartOfDay(now)
		assert.Equal(t, 0, startOfDay.Hour())
		assert.Equal(t, 0, startOfDay.Minute())
		assert.Equal(t, 0, startOfDay.Second())
		
		endOfDay := utils.GetEndOfDay(now)
		assert.Equal(t, 23, endOfDay.Hour())
		assert.Equal(t, 59, endOfDay.Minute())
		assert.Equal(t, 59, endOfDay.Second())
	})

	t.Run("WeekBoundaries", func(t *testing.T) {
		// Test with a Monday
		monday := time.Date(2023, 12, 25, 12, 0, 0, 0, time.UTC) // Monday
		
		startOfWeek := utils.GetStartOfWeek(monday)
		assert.Equal(t, monday.Year(), startOfWeek.Year())
		assert.Equal(t, monday.Month(), startOfWeek.Month())
		assert.Equal(t, monday.Day(), startOfWeek.Day())
		
		endOfWeek := utils.GetEndOfWeek(monday)
		assert.Equal(t, monday.AddDate(0, 0, 6).Day(), endOfWeek.Day())
	})

	t.Run("MonthBoundaries", func(t *testing.T) {
		testDate := time.Date(2023, 12, 15, 12, 0, 0, 0, time.UTC)
		
		startOfMonth := utils.GetStartOfMonth(testDate)
		assert.Equal(t, 1, startOfMonth.Day())
		
		endOfMonth := utils.GetEndOfMonth(testDate)
		assert.Equal(t, 31, endOfMonth.Day())
	})
}

func TestContainsFunctions(t *testing.T) {
	t.Run("ContainsString", func(t *testing.T) {
		slice := []string{"apple", "banana", "orange"}
		
		assert.True(t, utils.Contains(slice, "banana"))
		assert.False(t, utils.Contains(slice, "grape"))
		assert.True(t, utils.Contains(slice, "apple"))
		assert.True(t, utils.Contains(slice, "orange"))
	})

	t.Run("ContainsInt", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		
		assert.True(t, utils.ContainsInt(slice, 3))
		assert.False(t, utils.ContainsInt(slice, 6))
		assert.True(t, utils.ContainsInt(slice, 1))
		assert.True(t, utils.ContainsInt(slice, 5))
	})
}
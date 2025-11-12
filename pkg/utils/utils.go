package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Contains checks if a string is in a slice
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsInt checks if an int is in a slice
func ContainsInt(slice []int, item int) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

// GenerateRandomString generates a random string of given length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return ""
		}
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

// GenerateRandomHex generates a random hex string of given length
func GenerateRandomHex(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) bool {
	// Remove spaces and dashes
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	
	// Check if it starts with + and has 10-15 digits
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	return phoneRegex.MatchString(phone)
}

// SanitizeString removes potentially harmful characters
func SanitizeString(input string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	input = re.ReplaceAllString(input, "")
	
	// Remove SQL injection patterns
	re = regexp.MustCompile(`(['";\\])`)
	input = re.ReplaceAllString(input, "")
	
	return strings.TrimSpace(input)
}

// TruncateString truncates string to specified length
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// FormatCurrency formats number as currency
func FormatCurrency(amount float64, currency string) string {
	return fmt.Sprintf("%s %.2f", currency, amount)
}

// FormatDate formats time as date string
func FormatDate(t time.Time) string {
	return t.Format("YYYY-MM-DD")
}

// FormatDateTime formats time as datetime string
func FormatDateTime(t time.Time) string {
	return t.Format("YYYY-MM-DD HH:mm:ss")
}

// ParseDate parses date string to time
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("YYYY-MM-DD", dateStr)
}

// ParseDateTime parses datetime string to time
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return time.Parse("YYYY-MM-DD HH:mm:ss", dateTimeStr)
}

// IsWeekend checks if the given time is weekend
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// GetStartOfDay returns start of the day for given time
func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay returns end of the day for given time
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// GetStartOfWeek returns start of the week (Monday) for given time
func GetStartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return t.AddDate(0, 0, -weekday+1)
}

// GetEndOfWeek returns end of the week (Sunday) for given time
func GetEndOfWeek(t time.Time) time.Time {
	startOfWeek := GetStartOfWeek(t)
	return startOfWeek.AddDate(0, 0, 6)
}

// GetStartOfMonth returns start of the month for given time
func GetStartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// GetEndOfMonth returns end of the month for given time
func GetEndOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 999999999, t.Location())
}

// CalculateAge calculates age from birth date
func CalculateAge(birthDate time.Time) int {
	now := time.Now()
	years := now.Year() - birthDate.Year()
	
	if now.YearDay() < birthDate.YearDay() {
		years--
	}
	
	return years
}

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() string {
	num, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "000000"
	}
	return fmt.Sprintf("%06d", num.Int64())
}

// GenerateToken generates a secure token
func GenerateToken(length int) string {
	return GenerateRandomHex(length)
}

// HashString generates a hash of the input string
func HashString(input string) string {
	// Simple hash function (not cryptographically secure)
	hash := 0
	for _, char := range input {
		hash = hash*31 + int(char)
	}
	return fmt.Sprintf("%x", hash)
}

// Slugify converts string to URL-friendly slug
func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	
	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	
	// Remove non-alphanumeric characters except hyphens
	re := regexp.MustCompile(`[^a-z0-9-]`)
	s = re.ReplaceAllString(s, "")
	
	// Remove multiple hyphens
	re = regexp.MustCompile(`-+`)
	s = re.ReplaceAllString(s, "-")
	
	// Trim hyphens
	s = strings.Trim(s, "-")
	
	return s
}

// ExtractDomain extracts domain from email
func ExtractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// IsValidURL checks if string is valid URL
func IsValidURL(url string) bool {
	urlRegex := regexp.MustCompile(`^(https?://)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(/.*)?$`)
	return urlRegex.MatchString(url)
}

// ExtractURLs extracts URLs from text
func ExtractURLs(text string) []string {
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	return urlRegex.FindAllString(text, -1)
}

// RemoveDuplicates removes duplicate strings from slice
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if _, exists := keys[item]; !exists {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// Paginate calculates pagination metadata
func Paginate(totalItems int, page int, limit int) map[string]interface{} {
	totalPages := (totalItems + limit - 1) / limit
	
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}
	
	offset := (page - 1) * limit
	
	return map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"offset":      offset,
		"total_items": totalItems,
		"total_pages": totalPages,
		"has_next":    page < totalPages,
		"has_prev":    page > 1,
	}
}

// ResponseSuccess returns successful JSON response
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// ResponseError returns error JSON response
func ResponseError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   message,
	})
}

// ResponsePaginated returns paginated JSON response
func ResponsePaginated(c *gin.Context, data interface{}, pagination map[string]interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       data,
		"pagination": pagination,
	})
}

// BindAndValidate binds request data and validates
func BindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		ResponseError(c, http.StatusBadRequest, err.Error())
		return err
	}
	return nil
}

// GetClientIP gets client IP address
func GetClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		// Take the first IP if multiple are present
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}
	
	// Check X-Real-IP header
	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}
	
	// Fall back to RemoteAddr
	ip = c.RemoteIP()
	return ip
}

// IsLocalhost checks if request is from localhost
func IsLocalhost(c *gin.Context) bool {
	ip := GetClientIP(c)
	return ip == "127.0.0.1" || ip == "::1" || ip == "localhost"
}

// MaskString masks part of a string for privacy
func MaskString(s string, visible int) string {
	if len(s) <= visible*2 {
		return s
	}
	
	start := s[:visible]
	end := s[len(s)-visible:]
	
	masked := ""
	for i := visible; i < len(s)-visible; i++ {
		masked += "*"
	}
	
	return start + masked + end
}

// MaskEmail masks email for privacy
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	
	local := parts[0]
	domain := parts[1]
	
	if len(local) <= 3 {
		return email
	}
	
	maskedLocal := local[:2] + "****" + local[len(local)-1:]
	return maskedLocal + "@" + domain
}

// MaskPhone masks phone number for privacy
func MaskPhone(phone string) string {
	if len(phone) < 8 {
		return phone
	}
	
	return phone[:3] + "****" + phone[len(phone)-2:]
}

// BytesToSize converts bytes to human-readable size
func BytesToSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// DurationToString converts duration to human-readable string
func DurationToString(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%d seconds", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%d minutes", int(d.Minutes()))
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%d hours", int(d.Hours()))
	} else {
		days := int(d.Hours() / 24)
		hours := int(d.Hours()) % 24
		return fmt.Sprintf("%d days %d hours", days, hours)
	}
}
package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"whatsapp-bot/internal/models"
	"whatsapp-bot/pkg/logger"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type GameService struct {
	sm *ServiceManager
}

// Khodam generator data
var khodamList = []string{
	"Khodam Macan Putih - Memberikan kekuatan dan keberanian",
	"Khodam Naga Emas - Melindungi dari kejahatan dan memberikan keberuntungan",
	"Khodam Garuda - Memberikan kebebasan dan visi yang tajam",
	"Khodam Gajah - Memberikan kekuatan fisik dan stabilitas",
	"Khodam Harimau - Memberikan keberanian dan kekuatan",
	"Khodam Elang - Memberikan ketajaman mata dan kecepatan",
	"Khodam Lumba-lumba - Memberikan kecerdasan dan kebahagiaan",
	"Khodam Kucing - Memberikan kelincahan dan intuisi",
	"Khodam Anjing - Memberikan kesetiaan dan perlindungan",
	"Khodam Burung Hantu - Memberikan kebijaksanaan dan penglihatan malam",
	"Khodam Ular - Memberikan kekuatan spiritual dan penyembuhan",
	"Khodam Kupu-kupu - Memberikan transformasi dan keindahan",
	"Khodam Singa - Memberikan kepemimpinan dan keberanian",
	"Khodam Serigala - Memberikan kecerdasan dan kerja sama tim",
	"Khodam Merak - Memberikan keindahan dan percaya diri",
}

var zodiacSigns = map[string]string{
	"aries":       "Aries (21 Maret - 19 April) - Berapi-api, penuh energi, dan suka memimpin",
	"taurus":      "Taurus (20 April - 20 Mei) - Stabil, setia, dan menyukai kenyamanan",
	"gemini":      "Gemini (21 Mei - 20 Juni) - Komunikatif, cerdas, dan serba bisa",
	"cancer":      "Cancer (21 Juni - 22 Juli) - Penuh perhatian, emosional, dan protektif",
	"leo":         "Leo (23 Juli - 22 Agustus) - Percaya diri, kreatif, dan suka perhatian",
	"virgo":       "Virgo (23 Agustus - 22 September) - Detail, praktis, dan perfeksionis",
	"libra":       "Libra (23 September - 22 Oktober) - Seimbang, diplomatis, dan suka keadilan",
	"scorpio":     "Scorpio (23 Oktober - 21 November) - Intens, misterius, dan penuh gairah",
	"sagittarius": "Sagittarius (22 November - 21 Desember) - Petualang, optimis, dan bebas",
	"capricorn":   "Capricorn (22 Desember - 19 Januari) - Ambisius, disiplin, dan bertanggung jawab",
	"aquarius":    "Aquarius (20 Januari - 18 Februari) - Inovatif, independen, dan humanis",
	"pisces":      "Pisces (19 Februari - 20 Maret) - Empati, imajinatif, dan intuitif",
}

func (s *GameService) ProcessGameCommand(contact *models.Contact, message *models.Message) error {
	content := strings.ToLower(strings.TrimSpace(message.Content))

	// Check khodam command
	if strings.Contains(content, "cek khodam") || strings.Contains(content, "khodam") {
		return s.handleKhodamCommand(contact)
	}

	// Check zodiac command
	if strings.Contains(content, "zodiak") || strings.Contains(content, "zodiac") {
		return s.handleZodiacCommand(contact, content)
	}

	// Check love calculator
	if strings.Contains(content, "love calculator") || strings.Contains(content, "kalkulator cinta") {
		return s.handleLoveCalculatorCommand(contact, content)
	}

	// Check quiz command
	if strings.Contains(content, "kuis") || strings.Contains(content, "quiz") {
		return s.handleQuizCommand(contact)
	}

	// Check tebak gambar
	if strings.Contains(content, "tebak gambar") {
		return s.handleTebakGambarCommand(contact)
	}

	// Check math challenge
	if strings.Contains(content, "math challenge") || strings.Contains(content, "tantangan matematika") {
		return s.handleMathChallengeCommand(contact)
	}

	// Check jokes
	if strings.Contains(content, "joke") || strings.Contains(content, "jokes") || strings.Contains(content, "lucu") {
		return s.handleJokesCommand(contact)
	}

	// Check story
	if strings.Contains(content, "cerita") || strings.Contains(content, "story") {
		return s.handleStoryCommand(contact)
	}

	return nil
}

func (s *GameService) handleKhodamCommand(contact *models.Contact) error {
	// Generate random khodam
	rand.Seed(time.Now().UnixNano())
	khodam := khodamList[rand.Intn(len(khodamList))]

	// Create personalized message
	message := fmt.Sprintf("✨ KHODAM ANDA ✨\n\nNama: %s\n\n%s\n\nKhodam ini akan melindungi dan membantu Anda dalam perjalanan hidup. Semangat! 🙏", khodam, getKhodamDescription(khodam))

	// Send message
	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
	if err != nil {
		return err
	}

	// Update game score
	s.updateGameScore(contact.UserID, "khodam", 10)

	// Log analytics
	s.sm.AnalyticsService.LogEvent(contact.UserID, "khodam_checked", 1, map[string]interface{}{
		"khodam": khodam,
	})

	return nil
}

func (s *GameService) handleZodiacCommand(contact *models.Contact, content string) error {
	// Extract zodiac sign from message
	words := strings.Fields(content)
	var zodiacSign string

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		if _, exists := zodiacSigns[lowerWord]; exists {
			zodiacSign = lowerWord
			break
		}
	}

	if zodiacSign == "" {
		// Send zodiac list
		message := "🌟 ZODIAK HARI INI 🌟\n\nPilih zodiak Anda:\n"
		for sign, description := range zodiacSigns {
			message += fmt.Sprintf("• %s\n", strings.Title(sign))
		}
		message += "\nContoh: \"zodiak aries\""

		_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
		return err
	}

	// Generate daily horoscope
	horoscope := s.generateDailyHoroscope(zodiacSign)
	zodiacInfo := zodiacSigns[zodiacSign]

	message := fmt.Sprintf("🔮 RAMALAN %s HARI INI 🔮\n\n%s\n\n%s\n\nSemoga harimu menyenangkan! ✨", strings.Title(zodiacSign), zodiacInfo, horoscope)

	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
	if err != nil {
		return err
	}

	// Update game score
	s.updateGameScore(contact.UserID, "zodiac", 5)

	return nil
}

func (s *GameService) handleLoveCalculatorCommand(contact *models.Contact, content string) error {
	// Extract names from message
	words := strings.Fields(content)
	var names []string

	for _, word := range words {
		if len(word) > 2 && word != "love" && word != "calculator" && word != "kalkulator" && word != "cinta" {
			names = append(names, word)
		}
	}

	if len(names) < 2 {
		message := "💕 KALKULATOR CINTA 💕\n\nGunakan format: \"love calculator nama1 nama2\"\n\nContoh: \"love calculator budi ani\""
		_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
		return err
	}

	// Calculate love percentage
	percentage := s.calculateLovePercentage(names[0], names[1])
	compatibility := s.getLoveCompatibility(percentage)

	message := fmt.Sprintf("💑 KECOCOKAN CINTA 💑\n\n%s ❤️ %s\n\nKecocokan: %d%%\n\n%s\n\nSemoga berbahagia! 💖", names[0], names[1], percentage, compatibility)

	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
	if err != nil {
		return err
	}

	// Update game score
	s.updateGameScore(contact.UserID, "love_calculator", 15)

	return nil
}

func (s *GameService) handleQuizCommand(contact *models.Contact) error {
	// Start new quiz session
	quizSession := &models.QuizSession{
		UserID:     contact.UserID,
		Status:     "active",
		StartedAt:  time.Now(),
	}

	// Get random quiz questions
	var questions []models.Quiz
	s.sm.DB.Order("RANDOM()").Limit(5).Find(&questions)

	if len(questions) == 0 {
		// Create sample questions if none exist
		questions = s.createSampleQuizQuestions()
	}

	quizSession.TotalQuestions = len(questions)
	s.sm.DB.Create(quizSession)

	// Send first question
	return s.sendQuizQuestion(contact, quizSession, questions[0])
}

func (s *GameService) handleTebakGambarCommand(contact *models.Contact) error {
	// Create tebak gambar game
	rand.Seed(time.Now().UnixNano())
	games := []struct {
		imageURL string
		answer   string
		hint     string
	}{
		{"https://example.com/cat.jpg", "kucing", "Hewan kesayangan yang suya tidur"},
		{"https://example.com/tree.jpg", "pohon", "Tumbuhan besar yang berdaun"},
		{"https://example.com/car.jpg", "mobil", "Kendaraan beroda empat"},
	}

	game := games[rand.Intn(len(games))]

	// Send image with question
	message := fmt.Sprintf("🖼️ TEBAK GAMBAR 🖼️\n\nApa nama benda/hewan ini?\nHint: %s", game.hint)
	_, err := s.sm.WhatsApp.SendImageMessage(contact.PhoneNumber, game.imageURL, message)
	if err != nil {
		return err
	}

	// Store game state in Redis
	gameData := map[string]interface{}{
		"type":   "tebak_gambar",
		"answer": game.answer,
		"hint":   game.hint,
	}
	gameJSON, _ := json.Marshal(gameData)
	s.sm.Redis.Set(s.sm.Redis.Context(), fmt.Sprintf("game:%s:tebak_gambar", contact.ID.String()), gameJSON, 5*time.Minute)

	return nil
}

func (s *GameService) handleMathChallengeCommand(contact *models.Contact) error {
	// Generate math problem
	rand.Seed(time.Now().UnixNano())
	operations := []string{"+", "-", "*"}
	operation := operations[rand.Intn(len(operations))]
	
	num1 := rand.Intn(100) + 1
	num2 := rand.Intn(100) + 1
	
	var question string
	var answer int
	
	switch operation {
	case "+":
		answer = num1 + num2
		question = fmt.Sprintf("Berapa %d + %d?", num1, num2)
	case "-":
		answer = num1 - num2
		question = fmt.Sprintf("Berapa %d - %d?", num1, num2)
	case "*":
		answer = num1 * num2
		question = fmt.Sprintf("Berapa %d × %d?", num1, num2)
	}

	message := fmt.Sprintf("🧮 TANTANGAN MATEMATIKA 🧮\n\n%s\n\nJawab dalam 60 detik!", question)
	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
	if err != nil {
		return err
	}

	// Store answer in Redis
	s.sm.Redis.Set(s.sm.Redis.Context(), fmt.Sprintf("game:%s:math_answer", contact.ID.String()), answer, 60*time.Second)

	return nil
}

func (s *GameService) handleJokesCommand(contact *models.Contact) error {
	jokes := []string{
		"Kenapa kucing tidak pernah kalah dalam pertandingan? Karena dia selalu punya semangat 'meong'!",
		"Apa bedanya matematika dan pacar? Kalau matematika ada jawabannya, kalau pacar... ya sudahlah 😅",
		"Kenapa komputer tidak bisa tidur? Karena dia selalu 'ter-log in'!",
		"Apa yang dilakukan kucing saat bosan? Dia 'meong'-kel!",
		"Kenapa buku tidak pernah bosan? Karena dia punya banyak 'halaman'!",
	}

	rand.Seed(time.Now().UnixNano())
	joke := jokes[rand.Intn(len(jokes))]

	message := fmt.Sprintf("😂 JOKE HARI INI 😂\n\n%s\n\nSemoga harimu lebih ceria! 🌟", joke)
	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
	
	if err != nil {
		return err
	}

	// Update game score
	s.updateGameScore(contact.UserID, "jokes", 3)

	return nil
}

func (s *GameService) handleStoryCommand(contact *models.Contact) error {
	stories := []string{
		"🏰 CERITA HARI INI 🏰\n\nAda seekor kucing kecil yang ingin menjadi singa. Setiap hari dia berlatih mengaum di depan cermin. Suatu hari, saat ada tikus mengganggu rumah, kucing itu mengaum dengan keras dan tikus itu lari ketakutan. Pemilik rumah pun berkata, 'Kamu memang singa kecilku!'\n\nMoral: Percayalah pada dirimu sendiri! 💪",
		
		"🌟 KISAH INSPIRATIF 🌟\n\nSeorang anak kecil menanam biji kacang. Setiap hari dia menyiramnya tapi tidak ada yang tumbuh. 30 hari berlalu, tetap saja tidak tumbuh. Tapi dia tidak berhenti menyiram. Di hari ke-31, tumbuhlah tanaman kecil yang indah.\n\nKesabaran dan ketekunan selalu membuahkan hasil! 🌱",
	}

	rand.Seed(time.Now().UnixNano())
	story := stories[rand.Intn(len(stories))]

	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, story, false)
	if err != nil {
		return err
	}

	// Update game score
	s.updateGameScore(contact.UserID, "story", 5)

	return nil
}

func (s *GameService) sendQuizQuestion(contact *models.Contact, session *models.QuizSession, question models.Quiz) error {
	// Parse options
	var options []string
	json.Unmarshal([]byte(question.Options), &options)

	message := fmt.Sprintf("🧠 KUIS NO. %d 🧠\n\n%s\n\n", session.CurrentQuestion+1, question.Question)
	
	for i, option := range options {
		message += fmt.Sprintf("%d. %s\n", i+1, option)
	}

	message += "\nBalas dengan angka jawaban Anda!"

	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
	return err
}

func (s *GameService) ProcessQuizAnswer(contact *models.Contact, answer int) error {
	// Get active quiz session
	var session models.QuizSession
	err := s.sm.DB.Where("user_id = ? AND status = ?", contact.UserID, "active").First(&session).Error
	if err != nil {
		return err
	}

	// Get current question
	var currentQuestion models.Quiz
	err = s.sm.DB.Where("id = ?", session.QuizID).First(&currentQuestion).Error
	if err != nil {
		return err
	}

	// Check answer
	if answer == currentQuestion.CorrectAnswer {
		session.Score += currentQuestion.Points
		
		// Send correct answer message
		message := "✅ BENAR! ✅\n\nSelamat! Kamu mendapatkan %d poin!\n\nSkor sementara: %d poin"
		s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, fmt.Sprintf(message, currentQuestion.Points, session.Score), false)
	} else {
		// Send wrong answer message
		message := "❌ SALAH ❌\n\nJawaban yang benar adalah: %d\n\nSkor sementara: %d poin"
		s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, fmt.Sprintf(message, currentQuestion.CorrectAnswer, session.Score), false)
	}

	// Move to next question
	session.CurrentQuestion++
	
	if session.CurrentQuestion >= session.TotalQuestions {
		// Quiz completed
		session.Status = "completed"
		session.CompletedAt = &time.Now()
		
		// Send final score
		finalMessage := fmt.Sprintf("🎉 KUIS SELESAI! 🎉\n\nSkor akhir: %d/%d poin\n\nTerima kasih sudah bermain!", session.Score, session.TotalQuestions*10)
		s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, finalMessage, false)
		
		// Update game score
		s.updateGameScore(contact.UserID, "quiz", session.Score)
	} else {
		// Send next question
		var nextQuestion models.Quiz
		s.sm.DB.Where("id != ?", currentQuestion.ID).Order("RANDOM()").First(&nextQuestion)
		session.QuizID = nextQuestion.ID
		s.sendQuizQuestion(contact, &session, nextQuestion)
	}

	return s.sm.DB.Save(&session).Error
}

func (s *GameService) GetLeaderboard(gameType string, limit int) ([]models.GameScore, error) {
	var leaderboard []models.GameScore
	
	query := s.sm.DB.Where("game_type = ?", gameType)
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Order("score DESC, best_score DESC").Find(&leaderboard).Error
	return leaderboard, err
}

func (s *GameService) ResetDailyLeaderboard() error {
	// Reset daily scores for all games
	return s.sm.DB.Model(&models.GameScore{}).
		Where("game_type IN ('quiz', 'tebak_gambar', 'math_challenge')").
		Update("score", 0).Error
}

func (s *GameService) updateGameScore(userID uuid.UUID, gameType string, points int) {
	var gameScore models.GameScore
	err := s.sm.DB.Where("user_id = ? AND game_type = ?", userID, gameType).First(&gameScore).Error
	
	if err != nil {
		// Create new game score
		gameScore = models.GameScore{
			UserID:    userID,
			GameType:  gameType,
			Score:     points,
			PlayCount: 1,
			BestScore: points,
		}
		s.sm.DB.Create(&gameScore)
	} else {
		// Update existing score
		gameScore.Score += points
		gameScore.PlayCount++
		if points > gameScore.BestScore {
			gameScore.BestScore = points
		}
		s.sm.DB.Save(&gameScore)
	}
}

func (s *GameService) generateDailyHoroscope(zodiacSign string) string {
	rand.Seed(time.Now().UnixNano())
	
	horoscopes := map[string][]string{
		"general": []string{
			"Hari ini adalah hari yang baik untuk memulai sesuatu yang baru.",
			"Kesabaran Anda akan membuahkan hasil yang manis.",
			"Keberuntungan berpihak pada Anda hari ini.",
			"Jangan ragu untuk meminta bantuan jika membutuhkannya.",
			"Hari ini cocok untuk bersantai dan merenung.",
		},
		"love": []string{
			"Asmara Anda akan berjalan dengan lancar.",
			"Saatnya untuk membuka hati kepada seseorang yang spesial.",
			"Jangan takut untuk mengungkapkan perasaan Anda.",
			"Cinta sejati akan datang pada waktunya.",
			"Hubungan Anda akan semakin kuat hari ini.",
		},
		"career": []string{
			"Kesuksesan profesional sedang menunggu Anda.",
			"Waktu yang tepat untuk mengajukan promosi.",
			"Kreativitas Anda akan menghasilkan ide brilian.",
			"Kerja keras Anda akan dihargai oleh atasan.",
			"Kesempatan baru akan segera datang.",
		},
	}
	
	// Select random horoscope from each category
	general := horoscopes["general"][rand.Intn(len(horoscopes["general"]))]
	love := horoscopes["love"][rand.Intn(len(horoscopes["love"]))]
	career := horoscopes["career"][rand.Intn(len(horoscopes["career"]))]
	
	return fmt.Sprintf("💫 Umum: %s\n❤️ Cinta: %s\n💼 Karier: %s", general, love, career)
}

func (s *GameService) calculateLovePercentage(name1, name2 string) int {
	// Simple love calculator algorithm
	combined := strings.ToLower(name1 + name2)
	sum := 0
	for _, char := range combined {
		sum += int(char)
	}
	
	return (sum % 100) + 1
}

func (s *GameService) getLoveCompatibility(percentage int) string {
	if percentage >= 90 {
		return "Cinta sejati! Kalian sangat cocok bersama! 💕"
	} else if percentage >= 70 {
		return "Hubungan yang baik! Banyak kesamaan dan pemahaman! 💝"
	} else if percentage >= 50 {
		return "Hubungan yang menjanjikan! Perlu usaha dari kedua belah pihak! 💗"
	} else if percentage >= 30 {
		return "Tantangan dalam hubungan, tapi bisa diatasi dengan komunikasi! 💓"
	} else {
		return "Mungkin lebih baik sebagai teman! Tetap semangat! 💔"
	}
}

func (s *GameService) createSampleQuizQuestions() []models.Quiz {
	questions := []models.Quiz{
		{
			Question:    "Apa ibu kota Indonesia?",
			Options:     `["Jakarta", "Surabaya", "Bandung", "Medan"]`,
			CorrectAnswer: 1,
			Category:    "Geography",
			Difficulty:  "easy",
			Points:      10,
		},
		{
			Question:    "Berapa jumlah provinsi di Indonesia?",
			Options:     `["32", "33", "34", "35"]`,
			CorrectAnswer: 3,
			Category:    "Geography",
			Difficulty:  "easy",
			Points:      10,
		},
		{
			Question:    "Siapa penemu lampu pijar?",
			Options:     `["Thomas Edison", "Albert Einstein", "Isaac Newton", "Galileo"]`,
			CorrectAnswer: 1,
			Category:    "Science",
			Difficulty:  "medium",
			Points:      15,
		},
	}

	for _, question := range questions {
		s.sm.DB.Create(&question)
	}

	return questions
}

func getKhodamDescription(khodam string) string {
	// Extract description from khodam string
	parts := strings.Split(khodam, " - ")
	if len(parts) > 1 {
		return parts[1]
	}
	return "Khodam ini akan membantu dan melindungi Anda"
}

func (s *GameService) createSampleQuizQuestions() []models.Quiz {
	questions := []models.Quiz{
		{
			Question:    "Apa ibu kota Indonesia?",
			Options:     `["Jakarta", "Surabaya", "Bandung", "Medan"]`,
			CorrectAnswer: 1,
			Category:    "Geography",
			Difficulty:  "easy",
			Points:      10,
		},
		{
			Question:    "Berapa jumlah provinsi di Indonesia?",
			Options:     `["32", "33", "34", "35"]`,
			CorrectAnswer: 3,
			Category:    "Geography",
			Difficulty:  "easy",
			Points:      10,
		},
		{
			Question:    "Siapa penemu lampu pijar?",
			Options:     `["Thomas Edison", "Albert Einstein", "Isaac Newton", "Galileo"]`,
			CorrectAnswer: 1,
			Category:    "Science",
			Difficulty:  "medium",
			Points:      15,
		},
	}

	for _, question := range questions {
		s.sm.DB.Create(&question)
	}

	return questions
}
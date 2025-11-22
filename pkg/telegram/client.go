package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"kilocode.dev/whatsapp-bot/pkg/logger"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type Message struct {
	ChatID      int64  `json:"chat_id"`
	Text        string `json:"text"`
	ParseMode   string `json:"parse_mode,omitempty"`
	ReplyMarkup interface{} `json:"reply_markup,omitempty"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  *Message `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
}

type CallbackQuery struct {
	ID   string `json:"id"`
	From User   `json:"from"`
	Data string `json:"data"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

type SendMessageResponse struct {
	Ok     bool   `json:"ok"`
	Result Message `json:"result"`
}

type GetUpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://api.telegram.org/bot",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) SendMessage(chatID int64, text string, parseMode string) error {
	message := Message{
		ChatID:    chatID,
		Text:      text,
		ParseMode: parseMode,
	}

	return c.sendMessage(message)
}

func (c *Client) SendMessageWithMarkup(chatID int64, text string, markup interface{}) error {
	message := Message{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: markup,
	}

	return c.sendMessage(message)
}

func (c *Client) sendMessage(message Message) error {
	url := fmt.Sprintf("%s%s/sendMessage", c.baseURL, c.apiKey)

	jsonData, err := json.Marshal(message)
	if err != nil {
		logger.Error("Failed to marshal message", err)
		return err
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Failed to send message", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("telegram API error: %s", string(body))
		logger.Error("Telegram API error", err)
		return err
	}

	var response SendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logger.Error("Failed to decode response", err)
		return err
	}

	if !response.Ok {
		err := fmt.Errorf("telegram API returned error")
		logger.Error("Telegram API error", err)
		return err
	}

	return nil
}

func (c *Client) GetUpdates(offset int) ([]Update, error) {
	url := fmt.Sprintf("%s%s/getUpdates?offset=%d&timeout=30", c.baseURL, c.apiKey, offset)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		logger.Error("Failed to get updates", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("telegram API error: %s", string(body))
		logger.Error("Telegram API error", err)
		return nil, err
	}

	var response GetUpdatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logger.Error("Failed to decode response", err)
		return nil, err
	}

	if !response.Ok {
		err := fmt.Errorf("telegram API returned error")
		logger.Error("Telegram API error", err)
		return nil, err
	}

	return response.Result, nil
}

func (c *Client) AnswerCallbackQuery(callbackQueryID string, text string) error {
	url := fmt.Sprintf("%s%s/answerCallbackQuery?callback_query_id=%s&text=%s", 
		c.baseURL, c.apiKey, callbackQueryID, url.QueryEscape(text))

	resp, err := c.httpClient.Get(url)
	if err != nil {
		logger.Error("Failed to answer callback query", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("telegram API error: %s", string(body))
		logger.Error("Telegram API error", err)
		return err
	}

	return nil
}

func (c *Client) SendPhoto(chatID int64, photoURL string, caption string) error {
	url := fmt.Sprintf("%s%s/sendPhoto?chat_id=%d&photo=%s&caption=%s", 
		c.baseURL, c.apiKey, chatID, url.QueryEscape(photoURL), url.QueryEscape(caption))

	resp, err := c.httpClient.Get(url)
	if err != nil {
		logger.Error("Failed to send photo", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("telegram API error: %s", string(body))
		logger.Error("Telegram API error", err)
		return err
	}

	return nil
}

func (c *Client) SendDocument(chatID int64, documentURL string, caption string) error {
	url := fmt.Sprintf("%s%s/sendDocument?chat_id=%d&document=%s&caption=%s", 
		c.baseURL, c.apiKey, chatID, url.QueryEscape(documentURL), url.QueryEscape(caption))

	resp, err := c.httpClient.Get(url)
	if err != nil {
		logger.Error("Failed to send document", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("telegram API error: %s", string(body))
		logger.Error("Telegram API error", err)
		return err
	}

	return nil
}

func (c *Client) SendLocation(chatID int64, latitude float64, longitude float64) error {
	url := fmt.Sprintf("%s%s/sendLocation?chat_id=%d&latitude=%f&longitude=%f", 
		c.baseURL, c.apiKey, chatID, latitude, longitude)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		logger.Error("Failed to send location", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("telegram API error: %s", string(body))
		logger.Error("Telegram API error", err)
		return err
	}

	return nil
}

func (c *Client) GetMe() (*User, error) {
	url := fmt.Sprintf("%s%s/getMe", c.baseURL, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		logger.Error("Failed to get bot info", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("telegram API error: %s", string(body))
		logger.Error("Telegram API error", err)
		return nil, err
	}

	var response struct {
		Ok     bool `json:"ok"`
		Result User `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logger.Error("Failed to decode response", err)
		return nil, err
	}

	if !response.Ok {
		err := fmt.Errorf("telegram API returned error")
		logger.Error("Telegram API error", err)
		return nil, err
	}

	return &response.Result, nil
}

func (c *Client) SetWebhook(webhookURL string) error {
	url := fmt.Sprintf("%s%s/setWebhook?url=%s", c.baseURL, c.apiKey, url.QueryEscape(webhookURL))

	resp, err := c.httpClient.Get(url)
	if err != nil {
		logger.Error("Failed to set webhook", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("telegram API error: %s", string(body))
		logger.Error("Telegram API error", err)
		return err
	}

	return nil
}

func (c *Client) DeleteWebhook() error {
	url := fmt.Sprintf("%s%s/deleteWebhook", c.baseURL, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		logger.Error("Failed to delete webhook", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("telegram API error: %s", string(body))
		logger.Error("Telegram API error", err)
		return err
	}

	return nil
}

func (c *Client) GetWebhookInfo() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s/getWebhookInfo", c.baseURL, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		logger.Error("Failed to get webhook info", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("telegram API error: %s", string(body))
		logger.Error("Telegram API error", err)
		return nil, err
	}

	var response struct {
		Ok     bool                   `json:"ok"`
		Result map[string]interface{} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logger.Error("Failed to decode response", err)
		return nil, err
	}

	if !response.Ok {
		err := fmt.Errorf("telegram API returned error")
		logger.Error("Telegram API error", err)
		return nil, err
	}

	return response.Result, nil
}

func (c *Client) ProcessUpdate(update Update) error {
	logger.Info("Processing Telegram update", map[string]interface{}{
		"update_id": update.UpdateID,
		"has_message": update.Message != nil,
		"has_callback": update.CallbackQuery != nil,
	})

	// Handle message updates
	if update.Message != nil {
		return c.handleMessage(update.Message)
	}

	// Handle callback query updates
	if update.CallbackQuery != nil {
		return c.handleCallbackQuery(update.CallbackQuery)
	}

	return nil
}

func (c *Client) handleMessage(message *Message) error {
	logger.Info("Handling Telegram message", map[string]interface{}{
		"chat_id": message.ChatID,
		"text":    message.Text,
	})

	// Echo the message back (for testing)
	responseText := fmt.Sprintf("You said: %s", message.Text)
	return c.SendMessage(message.ChatID, responseText, "")
}

func (c *Client) handleCallbackQuery(callbackQuery *CallbackQuery) error {
	logger.Info("Handling Telegram callback query", map[string]interface{}{
		"callback_id": callbackQuery.ID,
		"data":        callbackQuery.Data,
	})

	// Answer the callback query
	responseText := fmt.Sprintf("You clicked: %s", callbackQuery.Data)
	return c.AnswerCallbackQuery(callbackQuery.ID, responseText)
}

func (c *Client) StartPolling(handler func(Update) error) error {
	logger.Info("Starting Telegram polling")
	
	offset := 0
	for {
		updates, err := c.GetUpdates(offset)
		if err != nil {
			logger.Error("Failed to get updates", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range updates {
			if err := handler(update); err != nil {
				logger.Error("Failed to handle update", err)
			}
			offset = update.UpdateID + 1
		}

		if len(updates) == 0 {
			time.Sleep(1 * time.Second)
		}
	}
}

func (c *Client) GetStatus() map[string]interface{} {
	me, err := c.GetMe()
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}

	webhookInfo, err := c.GetWebhookInfo()
	if err != nil {
		webhookInfo = map[string]interface{}{"error": err.Error()}
	}

	return map[string]interface{}{
		"status":        "connected",
		"bot_name":      me.FirstName,
		"bot_username":  me.Username,
		"webhook_info":  webhookInfo,
	}
}
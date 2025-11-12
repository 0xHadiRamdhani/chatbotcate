package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"whatsapp-bot/internal/config"
	"whatsapp-bot/pkg/logger"

	"github.com/sirupsen/logrus"
)

type Client struct {
	config      config.WhatsAppConfig
	httpClient  *http.Client
	baseURL     string
	accessToken string
}

type MessageRequest struct {
	MessagingProduct string                `json:"messaging_product"`
	RecipientType    string                `json:"recipient_type"`
	To               string                `json:"to"`
	Type             string                `json:"type"`
	Text             *TextMessage          `json:"text,omitempty"`
	Image            *MediaMessage         `json:"image,omitempty"`
	Audio            *MediaMessage         `json:"audio,omitempty"`
	Video            *MediaMessage         `json:"video,omitempty"`
	Document         *MediaMessage         `json:"document,omitempty"`
	Template         *TemplateMessage      `json:"template,omitempty"`
	Interactive      *InteractiveMessage   `json:"interactive,omitempty"`
}

type TextMessage struct {
	PreviewURL bool   `json:"preview_url,omitempty"`
	Body       string `json:"body"`
}

type MediaMessage struct {
	ID   string `json:"id,omitempty"`
	Link string `json:"link,omitempty"`
	Caption string `json:"caption,omitempty"`
}

type TemplateMessage struct {
	Name       string                 `json:"name"`
	Language   Language               `json:"language"`
	Components []TemplateComponent    `json:"components,omitempty"`
}

type Language struct {
	Code string `json:"code"`
}

type TemplateComponent struct {
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters,omitempty"`
}

type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type InteractiveMessage struct {
	Type   string          `json:"type"`
	Action *InteractiveAction `json:"action"`
	Body   *InteractiveBody   `json:"body,omitempty"`
	Footer *InteractiveFooter `json:"footer,omitempty"`
}

type InteractiveAction struct {
	Button string           `json:"button,omitempty"`
	Buttons []ReplyButton   `json:"buttons,omitempty"`
	Sections []Section      `json:"sections,omitempty"`
}

type ReplyButton struct {
	Type  string `json:"type"`
	Reply Reply  `json:"reply"`
}

type Reply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type InteractiveBody struct {
	Text string `json:"text"`
}

type InteractiveFooter struct {
	Text string `json:"text"`
}

type Section struct {
	Title string          `json:"title"`
	Rows  []SectionRow    `json:"rows"`
}

type SectionRow struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type MessageResponse struct {
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
	Error *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type WebhookPayload struct {
	Object string `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Value Value `json:"value"`
	Field string `json:"field"`
}

type Value struct {
	MessagingProduct string     `json:"messaging_product"`
	Metadata         Metadata   `json:"metadata"`
	Contacts         []Contact  `json:"contacts"`
	Messages         []Message  `json:"messages"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

type Profile struct {
	Name string `json:"name"`
}

type Message struct {
	From      string    `json:"from"`
	ID        string    `json:"id"`
	Timestamp string    `json:"timestamp"`
	Text      *Text     `json:"text,omitempty"`
	Image     *Media    `json:"image,omitempty"`
	Audio     *Media    `json:"audio,omitempty"`
	Video     *Media    `json:"video,omitempty"`
	Document  *Media    `json:"document,omitempty"`
	Type      string    `json:"type"`
}

type Text struct {
	Body string `json:"body"`
}

type Media struct {
	ID       string `json:"id"`
	MimeType string `json:"mime_type"`
	Caption  string `json:"caption"`
}

func Initialize(cfg config.WhatsAppConfig) (*Client, error) {
	client := &Client{
		config:      cfg,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		baseURL:     fmt.Sprintf("%s/%s", cfg.BaseURL, cfg.APIVersion),
		accessToken: cfg.AccessToken,
	}

	// Test connection
	if err := client.testConnection(); err != nil {
		return nil, fmt.Errorf("failed to test WhatsApp connection: %v", err)
	}

	logger.Log.Info("WhatsApp client initialized successfully")
	return client, nil
}

func (c *Client) testConnection() error {
	url := fmt.Sprintf("%s/%s", c.baseURL, c.config.PhoneNumberID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("WhatsApp API test failed: %s", string(body))
	}

	return nil
}

func (c *Client) SendMessage(message MessageRequest) (*MessageResponse, error) {
	url := fmt.Sprintf("%s/%s/messages", c.baseURL, c.config.PhoneNumberID)

	jsonData, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("WhatsApp API error: %s", string(body))
	}

	var result MessageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("WhatsApp error %d: %s", result.Error.Code, result.Error.Message)
	}

	logger.Log.WithFields(logrus.Fields{
		"message_id": result.Messages[0].ID,
		"recipient":  message.To,
		"type":       message.Type,
	}).Info("Message sent successfully")

	return &result, nil
}

func (c *Client) SendTextMessage(to, text string, previewURL bool) (*MessageResponse, error) {
	message := MessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "text",
		Text: &TextMessage{
			PreviewURL: previewURL,
			Body:       text,
		},
	}

	return c.SendMessage(message)
}

func (c *Client) SendImageMessage(to, imageURL, caption string) (*MessageResponse, error) {
	message := MessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "image",
		Image: &MediaMessage{
			Link:    imageURL,
			Caption: caption,
		},
	}

	return c.SendMessage(message)
}

func (c *Client) SendInteractiveMessage(to, body string, buttons []ReplyButton) (*MessageResponse, error) {
	message := MessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "interactive",
		Interactive: &InteractiveMessage{
			Type: "button",
			Body: &InteractiveBody{
				Text: body,
			},
			Action: &InteractiveAction{
				Buttons: buttons,
			},
		},
	}

	return c.SendMessage(message)
}

func (c *Client) SendTemplateMessage(to, templateName string, components []TemplateComponent) (*MessageResponse, error) {
	message := MessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "template",
		Template: &TemplateMessage{
			Name:       templateName,
			Language:   Language{Code: "id"},
			Components: components,
		},
	}

	return c.SendMessage(message)
}

func (c *Client) Disconnect() error {
	logger.Log.Info("WhatsApp client disconnected")
	return nil
}

func (c *Client) GetPhoneNumberID() string {
	return c.config.PhoneNumberID
}

func (c *Client) VerifyWebhookSignature(payload []byte, signature string) bool {
	// Implement webhook signature verification
	// This is a placeholder - implement actual signature verification
	return true
}
package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Client struct {
	BaseURL	string
	Model	string
}

type chatRequest struct {
	Model		string		`json:"model"`
	Messages	[]Message	`json:"messages"`
	Stream		bool		`json:"stream"`
	Options		chatOptions	`json:"options"`
}

type chatOptions struct {
	NumPredict	int	`json:"num_predict"`
	Temperature	float64	`json:"temperature"`
	TopK		int	`json:"top_k"`
	TopP		float64	`json:"top_p"`
}

type Message struct {
	Role	string	`json:"role"`
	Content	string	`json:"content"`
}

type chatResponse struct {
	Message Message `json:"message"`
}

func NewClient(baseURL, model string) *Client {
	return &Client{BaseURL: baseURL, Model: model}
}

var thinkRegex = regexp.MustCompile(`(?s)<\|channel>.*?<channel\|>`)

func stripThinkingTokens(s string) string {
	s = thinkRegex.ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}

var httpClient = &http.Client{Timeout: 90 * time.Second}

func (c *Client) Chat(history []map[string]string, userMessage string) (string, error) {
	log.Printf("[SİSTEM] AI Modeline İstek Gönderiliyor: %s", c.Model)
	log.Printf("[SİSTEM] Kullanıcı Mesajı: %s", userMessage)

	var messages []Message

	execPath, _ := os.Executable()
	promptJsonPath := filepath.Join(filepath.Dir(execPath), "prompt.json")
	promptTxtPath := filepath.Join(filepath.Dir(execPath), "prompt.txt")

	var systemPrompt string
	
	// First try to read JSON prompt
	promptBytes, err := os.ReadFile(promptJsonPath)
	if err != nil {
		promptBytes, err = os.ReadFile("prompt.json")
	}

	if err == nil && len(promptBytes) > 0 {
		systemPrompt = string(promptBytes)
	} else {
		// Fallback to text prompt
		promptBytes, err = os.ReadFile(promptTxtPath)
		if err != nil {
			promptBytes, err = os.ReadFile("prompt.txt")
		}
		
		systemPrompt = "You are a professional and helpful WhatsApp assistant. Your system rules have not been defined by the administrator yet. Please provide general and short answers to users for now and wait for the administrator's configuration."
		if err == nil {
			systemPrompt = string(promptBytes)
		}
	}

	messages = append(messages, Message{
		Role:		"system",
		Content:	systemPrompt,
	})

	for _, m := range history {
		messages = append(messages, Message{
			Role:		m["role"],
			Content:	m["content"],
		})
	}

	messages = append(messages, Message{
		Role:		"user",
		Content:	userMessage,
	})

	reqBody, err := json.Marshal(chatRequest{
		Model:		c.Model,
		Messages:	messages,
		Stream:		false,
		Options: chatOptions{
			NumPredict:	1500,
			Temperature:	0.7,
			TopK:		64,
			TopP:		0.95,
		},
	})
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}

	log.Printf("[SİSTEM] %d adet mesaj geçmişi eklendi. Ollama'dan yanıt bekleniyor...", len(messages))
	
	baseURL := strings.TrimRight(c.BaseURL, "/")
	baseURL = strings.TrimSuffix(baseURL, "/v1")
	baseURL = strings.TrimRight(baseURL, "/")

	resp, err := httpClient.Post(
		baseURL+"/api/chat",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", fmt.Errorf("ollama isteği: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("yanıt okuma: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama Hatası (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		log.Printf("[AI DEBUG] Unmarshal hatası. Ham yanıt: %s", string(body))
		return "", fmt.Errorf("json unmarshal: %w", err)
	}

	if chatResp.Message.Content == "" {
		log.Printf("[AI DEBUG] Model boş cevap döndürdü. Ham yanıt: %s", string(body))
	} else {
		// Düşünme sürecini (varsa) logla
		log.Printf("[SİSTEM] AI Düşünme ve Ham Çıktı: %s", chatResp.Message.Content)
	}

	clean := stripThinkingTokens(chatResp.Message.Content)
	log.Printf("[SİSTEM] AI Temizlenmiş Yanıt: %s", clean)
	return clean, nil
}

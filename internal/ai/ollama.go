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

// appDataDir returns ~/Library/Application Support/Kuliss on macOS (same as config.AppDataDir).
func appDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	dir := filepath.Join(home, "Library", "Application Support", "Kuliss")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

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

var thinkRegex = regexp.MustCompile(`(?s)<\|channel\>.*?<channel\|>`)

func stripThinkingTokens(s string) string {
	s = thinkRegex.ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}

// promptJSON holds the parsed structure of a prompt.json file.
type promptJSON struct {
	Role               string            `json:"role"`
	Tone               string            `json:"tone"`
	Constraints        []string          `json:"constraints"`
	ProductsAndServices []struct {
		Service string `json:"service"`
		Price   string `json:"price"`
	} `json:"products_and_services"`
	FAQ []struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	} `json:"faq"`
	Examples []struct {
		User string `json:"user"`
		AI   string `json:"ai"`
	} `json:"examples"`
}

// buildSystemPromptFromJSON converts the structured JSON prompt into a clear
// natural-language system prompt that LLMs follow reliably.
func buildSystemPromptFromJSON(data []byte) (string, error) {
	var p promptJSON
	if err := json.Unmarshal(data, &p); err != nil {
		return "", err
	}

	var sb strings.Builder

	if p.Role != "" {
		sb.WriteString("## Your Role\n")
		sb.WriteString(p.Role)
		sb.WriteString("\n\n")
	}

	if p.Tone != "" {
		sb.WriteString("## Tone\n")
		sb.WriteString(p.Tone)
		sb.WriteString("\n\n")
	}

	if len(p.Constraints) > 0 {
		sb.WriteString("## Rules (follow strictly)\n")
		for i, c := range p.Constraints {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, c))
		}
		sb.WriteString("\n")
	}

	if len(p.ProductsAndServices) > 0 {
		sb.WriteString("## Products & Services\n")
		for _, ps := range p.ProductsAndServices {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", ps.Service, ps.Price))
		}
		sb.WriteString("\n")
	}

	if len(p.FAQ) > 0 {
		sb.WriteString("## Frequently Asked Questions\n")
		for _, faq := range p.FAQ {
			sb.WriteString(fmt.Sprintf("Q: %s\nA: %s\n\n", faq.Question, faq.Answer))
		}
	}

	if len(p.Examples) > 0 {
		sb.WriteString("## Few-Shot Examples (Use these as reference for tone and structure)\n")
		for _, ex := range p.Examples {
			sb.WriteString(fmt.Sprintf("User: %s\nAI: %s\n\n", ex.User, ex.AI))
		}
	}

	return strings.TrimSpace(sb.String()), nil
}

var httpClient = &http.Client{Timeout: 90 * time.Second}

// longHttpClient is used for tasks that may process very large prompts (e.g. JSON generation).
var longHttpClient = &http.Client{Timeout: 300 * time.Second}

func (c *Client) Chat(history []map[string]string, userMessage string) (string, error) {
	log.Printf("[SYSTEM] Sending Request to AI Model: %s", c.Model)
	log.Printf("[SYSTEM] User Message: %s", userMessage)

	var messages []Message

	execPath, _ := os.Executable()
	// Search paths: GUI app data dir (highest priority), executable dir, cwd
	appData := appDataDir()
	jsonSearchPaths := []string{
		filepath.Join(appData, "prompt.json"),
		filepath.Join(filepath.Dir(execPath), "prompt.json"),
		"prompt.json",
	}
	txtSearchPaths := []string{
		filepath.Join(appData, "prompt.txt"),
		filepath.Join(filepath.Dir(execPath), "prompt.txt"),
		"prompt.txt",
	}

	var systemPrompt string

	// First try to read JSON prompt from all search paths
	var promptBytes []byte
	var err error
	for _, p := range jsonSearchPaths {
		promptBytes, err = os.ReadFile(p)
		if err == nil && len(promptBytes) > 0 {
			log.Printf("[SYSTEM] JSON prompt loaded from: %s", p)
			break
		}
	}

	if err == nil && len(promptBytes) > 0 {
		// Parse JSON and convert to natural-language instructions
		if built, parseErr := buildSystemPromptFromJSON(promptBytes); parseErr == nil && built != "" {
			systemPrompt = built
			log.Printf("[SYSTEM] JSON prompt successfully parsed into system instructions.")
		} else {
			// Fallback: use raw JSON string if parsing fails
			log.Printf("[SYSTEM] JSON parse warning (%v), using raw JSON as prompt.", parseErr)
			systemPrompt = string(promptBytes)
		}
	} else {
		// Fallback to text prompt
		for _, p := range txtSearchPaths {
			promptBytes, err = os.ReadFile(p)
			if err == nil && len(promptBytes) > 0 {
				log.Printf("[SYSTEM] Text prompt loaded from: %s", p)
				break
			}
		}

		systemPrompt = "You are a professional and helpful WhatsApp assistant. Your system rules have not been defined by the administrator yet. Please provide general and short answers to users for now and wait for the administrator's configuration."
		if err == nil && len(promptBytes) > 0 {
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

	log.Printf("[SYSTEM] %d message(s) added to history. Waiting for response from Ollama...", len(messages))
	
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
		log.Printf("[AI DEBUG] Model returned empty response. Raw response: %s", string(body))
	} else {
		// Düşünme sürecini (varsa) logla
		log.Printf("[SYSTEM] AI Thinking and Raw Output: %s", chatResp.Message.Content)
	}

	clean := stripThinkingTokens(chatResp.Message.Content)
	log.Printf("[SYSTEM] AI Cleaned Response: %s", clean)
	return clean, nil
}

// ChatDirect sends a request to Ollama with an explicit system prompt,
// bypassing any prompt file on disk. Used for internal tasks like JSON generation.
func (c *Client) ChatDirect(systemPrompt string, userMessage string) (string, error) {
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userMessage},
	}

	reqBody, err := json.Marshal(chatRequest{
		Model:    c.Model,
		Messages: messages,
		Stream:   false,
		Options: chatOptions{
			NumPredict:  2000,
			Temperature: 0.3,
			TopK:        40,
			TopP:        0.9,
		},
	})
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}

	baseURL := strings.TrimRight(c.BaseURL, "/")
	baseURL = strings.TrimSuffix(baseURL, "/v1")
	baseURL = strings.TrimRight(baseURL, "/")

	log.Printf("[SYSTEM] ChatDirect → Ollama (%s)...", c.Model)
	resp, err := longHttpClient.Post(
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
		return "", fmt.Errorf("json unmarshal: %w", err)
	}

	result := stripThinkingTokens(chatResp.Message.Content)
	log.Printf("[SYSTEM] ChatDirect cleaned response: %s", result)
	return result, nil
}

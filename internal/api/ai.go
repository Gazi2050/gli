package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const encodedAIEndpoint = "aHR0cHM6Ly9kaW55LWNsaS52ZXJjZWwuYXBwL2FwaS92Mi9jb21taXQ="

type AIService struct {
	Client   *http.Client
	Endpoint string
}

type AIRequest struct {
	GitDiff  string   `json:"gitDiff"`
	Version  string   `json:"version"`
	Name     string   `json:"name"`
	RepoName string   `json:"repoName"`
	System   string   `json:"system"`
	Config   AIConfig `json:"config"`
}

type AIConfig struct {
	Theme  string         `json:"Theme"`
	Commit AICommitConfig `json:"Commit"`
}

type AICommitConfig struct {
	Conventional       bool              `json:"Conventional"`
	ConventionalFormat []string          `json:"ConventionalFormat"`
	Emoji              bool              `json:"Emoji"`
	EmojiMap           map[string]string `json:"EmojiMap"`
	Tone               string            `json:"Tone"`
	Length             string            `json:"Length"`
	CustomInstructions string            `json:"CustomInstructions"`
	HashAfterCommit    bool              `json:"HashAfterCommit"`
}

type AIResponse struct {
	Data AIResponseData `json:"data"`
}

type AIResponseData struct {
	CommitMessage string `json:"commitMessage"`
}

func NewAIService(client *http.Client) (*AIService, error) {
	endpoint, err := decodeEndpoint(encodedAIEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to decode AI endpoint: %w", err)
	}

	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	return &AIService{Client: client, Endpoint: endpoint}, nil
}

func (s *AIService) GenerateCommitMessage(diff, username, repoName, customInstructions string) (string, error) {
	if s.Client == nil {
		s.Client = &http.Client{Timeout: 30 * time.Second}
	}
	if strings.TrimSpace(s.Endpoint) == "" {
		endpoint, err := decodeEndpoint(encodedAIEndpoint)
		if err != nil {
			return "", fmt.Errorf("failed to decode AI endpoint: %w", err)
		}
		s.Endpoint = endpoint
	}

	payload := AIRequest{
		GitDiff:  diff,
		Version:  "v1.0.0",
		Name:     username,
		RepoName: repoName,
		System:   runtime.GOOS,
		Config: AIConfig{
			Theme: "catppuccin",
			Commit: AICommitConfig{
				Conventional:       true,
				ConventionalFormat: []string{"feat", "fix", "docs", "chore", "style", "refactor", "test", "perf"},
				Emoji:              false,
				EmojiMap: map[string]string{
					"feat":     "🚀",
					"fix":      "🐛",
					"docs":     "📚",
					"style":    "🎨",
					"refactor": "♻️",
					"test":     "✅",
					"chore":    "🔧",
					"perf":     "⚡",
				},
				Tone:               "casual",
				Length:             "short",
				CustomInstructions: customInstructions,
				HashAfterCommit:    false,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode AI request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.Endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create AI request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call AI service: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read AI response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, truncate(string(respBody), 300))
	}

	var parsed AIResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("failed to parse AI response: %w", err)
	}

	message := strings.TrimSpace(parsed.Data.CommitMessage)
	if message == "" {
		return "", fmt.Errorf("AI response missing commit message")
	}

	return message, nil
}

func decodeEndpoint(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

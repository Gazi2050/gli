package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const baseURL = "https://diny-cli.vercel.app"

var conventionalTypes = []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "build", "ci", "chore", "revert"}

var conventionalPattern = regexp.MustCompile(`^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(.+\))?:\s.+`)

type AIService struct {
	Client *http.Client
}

type dinyRequest struct {
	Type       string      `json:"type"`
	UserPrompt string      `json:"userPrompt"`
	Config     dinyConfig  `json:"config"`
	Version    string      `json:"version"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	RepoName   string      `json:"repoName"`
	System     string      `json:"system"`
}

type dinyConfig struct {
	Theme    string          `json:"Theme"`
	Commit   dinyCommitCfg   `json:"Request"`
	Prompts  dinyPromptsCfg  `json:"Prompts"`
}

type dinyCommitCfg struct {
	Conventional       bool   `json:"Conventional"`
	Emoji              bool   `json:"Emoji"`
	Tone               string `json:"Tone"`
	Length             string `json:"Length"`
	CustomInstructions string `json:"CustomInstructions"`
	HashAfterCommit    bool   `json:"HashAfterCommit"`
}

type dinyPromptsCfg struct {
	Enabled bool `json:"enabled"`
}

type dinyResponse struct {
	Error *string        `json:"error,omitempty"`
	Data  *dinyRespData  `json:"data,omitempty"`
}

type dinyRespData struct {
	Message string `json:"message"`
}

func NewAIService(client *http.Client) (*AIService, error) {
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second}
	}
	return &AIService{Client: client}, nil
}

func (s *AIService) GenerateCommitMessage(diff, username, repoName, customInstructions string) (string, error) {
	if s.Client == nil {
		s.Client = &http.Client{Timeout: 60 * time.Second}
	}

	instructions := "You MUST output ONLY a Conventional Commit message in the exact format: <type>[optional scope]: <description>. Allowed types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert. No explanation, no markdown, no quotes, no backticks. Just the commit message."
	if customInstructions != "" {
		instructions += " " + customInstructions
	}

	payload := dinyRequest{
		Type:       "commit",
		UserPrompt: diff,
		Version:    "v1.0.0",
		Name:       username,
		RepoName:   repoName,
		System:     runtime.GOOS,
		Config: dinyConfig{
			Theme: "catppuccin",
			Commit: dinyCommitCfg{
				Conventional:       true,
				Emoji:              false,
				Tone:               "casual",
				Length:             "short",
				CustomInstructions: instructions,
				HashAfterCommit:    false,
			},
			Prompts: dinyPromptsCfg{Enabled: false},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/requests", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call AI service: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, truncate(string(respBody), 300))
	}

	var parsed dinyResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if parsed.Error != nil {
		return "", fmt.Errorf("AI service error: %s", *parsed.Error)
	}

	if parsed.Data == nil {
		return "", fmt.Errorf("no data in AI response")
	}

	message := cleanMessage(parsed.Data.Message)
	if message == "" {
		return "", fmt.Errorf("AI response missing commit message")
	}

	if !conventionalPattern.MatchString(message) {
		message = enforceConventional(message)
	}

	return message, nil
}

func cleanMessage(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(s, "`")
	s = strings.TrimSuffix(s, "`")
	s = strings.Trim(s, "\"'")
	return strings.TrimSpace(s)
}

func enforceConventional(msg string) string {
	msg = strings.TrimSpace(msg)
	lines := strings.SplitN(msg, "\n", 2)
	first := strings.TrimSpace(lines[0])

	lower := strings.ToLower(first)
	for _, t := range conventionalTypes {
		if strings.HasPrefix(lower, t+":") || strings.HasPrefix(lower, t+"(") {
			return msg
		}
	}

	if lower != "" {
		return "chore: " + first
	}
	return msg
}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

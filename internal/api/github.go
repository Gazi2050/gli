package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultGitHubBaseURL = "https://api.github.com/"

type GitHubClient struct {
	Client  *http.Client
	BaseURL string
}

type User struct {
	Login           string `json:"login"`
	Name            string `json:"name"`
	Bio             string `json:"bio"`
	Location        string `json:"location"`
	PublicRepos     int    `json:"public_repos"`
	Followers       int    `json:"followers"`
	Following       int    `json:"following"`
	CreatedAt       string `json:"created_at"`
	TwitterUsername string `json:"twitter_username"`
	Blog            string `json:"blog"`
	AvatarURL       string `json:"avatar_url"`
	HTMLURL         string `json:"html_url"`
	Company         string `json:"company"`
}

type Repo struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	StargazersCount int    `json:"stargazers_count"`
	Language        string `json:"language"`
	UpdatedAt       string `json:"updated_at"`
	Fork            bool   `json:"fork"`
	HTMLURL         string `json:"html_url"`
}

func NewGitHubClient(client *http.Client, baseURL string) *GitHubClient {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	if strings.TrimSpace(baseURL) == "" {
		baseURL = defaultGitHubBaseURL
	}
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return &GitHubClient{
		Client:  client,
		BaseURL: baseURL,
	}
}

func (c *GitHubClient) GetUser(username string) (*User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	path := "users/" + url.PathEscape(username)
	endpoint, err := c.buildURL(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build user URL: %w", err)
	}

	respBody, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(respBody, &user); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub user response: %w", err)
	}

	return &user, nil
}

func (c *GitHubClient) GetUserRepos(username string, limit int) ([]Repo, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if limit <= 0 {
		limit = 5
	}

	query := url.Values{}
	query.Set("sort", "updated")
	query.Set("per_page", strconv.Itoa(limit))

	path := "users/" + url.PathEscape(username) + "/repos"
	endpoint, err := c.buildURL(path, query)
	if err != nil {
		return nil, fmt.Errorf("failed to build repos URL: %w", err)
	}

	respBody, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var repos []Repo
	if err := json.Unmarshal(respBody, &repos); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub repos response: %w", err)
	}

	return repos, nil
}

func (c *GitHubClient) buildURL(path string, query url.Values) (string, error) {
	base := c.BaseURL
	if strings.TrimSpace(base) == "" {
		base = defaultGitHubBaseURL
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	rel, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	full := parsed.ResolveReference(rel)
	if query != nil {
		full.RawQuery = query.Encode()
	}
	return full.String(), nil
}

func (c *GitHubClient) get(endpoint string) ([]byte, error) {
	if c.Client == nil {
		c.Client = &http.Client{Timeout: 30 * time.Second}
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gli-go-client")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GitHub API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GitHub response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, truncate(string(body), 300))
	}

	return body, nil
}

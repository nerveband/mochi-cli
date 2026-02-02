package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nerveband/mochi-cli/internal/models"
)

const (
	baseURL   = "https://app.mochi.cards/api"
	timeout   = 30 * time.Second
	userAgent = "mochi-cli/1.0"
)

// Client represents the Mochi API client
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// NewClientWithTimeout creates a new API client with custom timeout
func NewClientWithTimeout(apiKey string, timeout time.Duration) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// setAuth sets the Authorization header for basic auth
func (c *Client) setAuth(req *http.Request) {
	// Basic auth: username is API key, no password
	auth := base64.StdEncoding.EncodeToString([]byte(c.apiKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
}

// doRequest performs an HTTP request and returns the response
func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	return c.httpClient.Do(req)
}

// handleError processes error responses
func handleError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)

	var errResp models.ErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil {
		return fmt.Errorf("API error (%d): %v", resp.StatusCode, errResp.Errors)
	}

	return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
}

// === Card Operations ===

// ListCards lists cards with optional filtering
func (c *Client) ListCards(deckID string, limit int, bookmark string) (*models.PaginatedResponse, error) {
	params := url.Values{}
	if deckID != "" {
		params.Set("deck-id", deckID)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	if bookmark != "" {
		params.Set("bookmark", bookmark)
	}

	url := baseURL + "/cards"
	if len(params) > 0 {
		url = url + "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.setAuth(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleError(resp); err != nil {
		return nil, err
	}

	var result models.PaginatedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCard retrieves a specific card
func (c *Client) GetCard(cardID string) (*models.Card, error) {
	url := baseURL + "/cards/" + cardID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.setAuth(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleError(resp); err != nil {
		return nil, err
	}

	var card models.Card
	if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
		return nil, err
	}

	return &card, nil
}

// CreateCard creates a new card
func (c *Client) CreateCard(card *models.Card) (*models.Card, error) {
	url := baseURL + "/cards"

	body, err := json.Marshal(card)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.setAuth(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleError(resp); err != nil {
		return nil, err
	}

	var created models.Card
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}

	return &created, nil
}

// UpdateCard updates an existing card
func (c *Client) UpdateCard(cardID string, card *models.Card) (*models.Card, error) {
	url := baseURL + "/cards/" + cardID

	body, err := json.Marshal(card)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.setAuth(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleError(resp); err != nil {
		return nil, err
	}

	var updated models.Card
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteCard permanently deletes a card
func (c *Client) DeleteCard(cardID string) error {
	url := baseURL + "/cards/" + cardID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	c.setAuth(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return handleError(resp)
}

// AddAttachment adds an attachment to a card
func (c *Client) AddAttachment(cardID string, filename string, fileData []byte) error {
	url := fmt.Sprintf("%s/cards/%s/attachments/%s", baseURL, cardID, filename)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}

	if _, err := part.Write(fileData); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return err
	}

	c.setAuth(req)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return handleError(resp)
}

// AddAttachmentFromFile adds an attachment from a file path
func (c *Client) AddAttachmentFromFile(cardID string, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	filename := filepath.Base(filePath)
	return c.AddAttachment(cardID, filename, data)
}

// DeleteAttachment removes an attachment from a card
func (c *Client) DeleteAttachment(cardID string, filename string) error {
	url := fmt.Sprintf("%s/cards/%s/attachments/%s", baseURL, cardID, filename)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	c.setAuth(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return handleError(resp)
}

// === Deck Operations ===

// ListDecks lists all decks
func (c *Client) ListDecks(bookmark string) (*models.PaginatedResponse, error) {
	params := url.Values{}
	if bookmark != "" {
		params.Set("bookmark", bookmark)
	}

	url := baseURL + "/decks"
	if len(params) > 0 {
		url = url + "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.setAuth(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleError(resp); err != nil {
		return nil, err
	}

	var result models.PaginatedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDeck retrieves a specific deck
func (c *Client) GetDeck(deckID string) (*models.Deck, error) {
	url := baseURL + "/decks/" + deckID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.setAuth(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleError(resp); err != nil {
		return nil, err
	}

	var deck models.Deck
	if err := json.NewDecoder(resp.Body).Decode(&deck); err != nil {
		return nil, err
	}

	return &deck, nil
}

// CreateDeck creates a new deck
func (c *Client) CreateDeck(deck *models.Deck) (*models.Deck, error) {
	url := baseURL + "/decks"

	body, err := json.Marshal(deck)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.setAuth(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleError(resp); err != nil {
		return nil, err
	}

	var created models.Deck
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}

	return &created, nil
}

// UpdateDeck updates an existing deck
func (c *Client) UpdateDeck(deckID string, deck *models.Deck) (*models.Deck, error) {
	url := baseURL + "/decks/" + deckID

	body, err := json.Marshal(deck)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.setAuth(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleError(resp); err != nil {
		return nil, err
	}

	var updated models.Deck
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteDeck permanently deletes a deck
func (c *Client) DeleteDeck(deckID string) error {
	url := baseURL + "/decks/" + deckID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	c.setAuth(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return handleError(resp)
}

// === Template Operations ===

// ListTemplates lists all templates
func (c *Client) ListTemplates(bookmark string) (*models.PaginatedResponse, error) {
	params := url.Values{}
	if bookmark != "" {
		params.Set("bookmark", bookmark)
	}

	url := baseURL + "/templates"
	if len(params) > 0 {
		url = url + "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.setAuth(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleError(resp); err != nil {
		return nil, err
	}

	var result models.PaginatedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTemplate retrieves a specific template
func (c *Client) GetTemplate(templateID string) (*models.Template, error) {
	url := baseURL + "/templates/" + templateID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.setAuth(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleError(resp); err != nil {
		return nil, err
	}

	var template models.Template
	if err := json.NewDecoder(resp.Body).Decode(&template); err != nil {
		return nil, err
	}

	return &template, nil
}

// CreateTemplate creates a new template
func (c *Client) CreateTemplate(template *models.Template) (*models.Template, error) {
	url := baseURL + "/templates"

	body, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.setAuth(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleError(resp); err != nil {
		return nil, err
	}

	var created models.Template
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}

	return &created, nil
}

// === Due Cards Operations ===

// GetDueCards retrieves cards due on a specific date
func (c *Client) GetDueCards(date string, deckID string) (*models.DueResponse, error) {
	params := url.Values{}
	if date != "" {
		params.Set("date", date)
	}

	url := baseURL + "/due"
	if deckID != "" {
		url = url + "/" + deckID
	}
	if len(params) > 0 {
		url = url + "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.setAuth(req)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleError(resp); err != nil {
		return nil, err
	}

	var result models.DueResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAllDueCards retrieves all cards due across all decks
func (c *Client) GetAllDueCards(date string) (*models.DueResponse, error) {
	return c.GetDueCards(date, "")
}

// SearchCards searches for cards matching a query (client-side filtering)
func (c *Client) SearchCards(query string, deckID string) ([]models.Card, error) {
	var allCards []models.Card
	var bookmark string

	for {
		resp, err := c.ListCards(deckID, 100, bookmark)
		if err != nil {
			return nil, err
		}

		docs, ok := resp.Docs.([]interface{})
		if !ok {
			break
		}

		for _, doc := range docs {
			data, _ := json.Marshal(doc)
			var card models.Card
			if err := json.Unmarshal(data, &card); err == nil {
				// Simple case-insensitive search in content and name
				if strings.Contains(strings.ToLower(card.Content), strings.ToLower(query)) ||
					strings.Contains(strings.ToLower(card.Name), strings.ToLower(query)) {
					allCards = append(allCards, card)
				}
			}
		}

		if resp.Bookmark == "" {
			break
		}
		bookmark = resp.Bookmark
	}

	return allCards, nil
}

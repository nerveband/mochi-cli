package models

import "time"

// Card represents a flashcard in Mochi
type Card struct {
	ID            string           `json:"id"`
	Content       string           `json:"content"`
	Name          string           `json:"name,omitempty"`
	DeckID        string           `json:"deck-id"`
	TemplateID    string           `json:"template-id,omitempty"`
	Pos           string           `json:"pos"`
	Fields        map[string]Field `json:"fields,omitempty"`
	ManualTags    []string         `json:"manual-tags,omitempty"`
	Archived      bool             `json:"archived?"`
	ReviewReverse bool             `json:"review-reverse?"`
	Trashed       *time.Time       `json:"trashed?"`
	New           bool             `json:"new?"`
	References    []string         `json:"references"`
	Reviews       []Review         `json:"reviews"`
	CreatedAt     time.Time        `json:"created-at"`
	UpdatedAt     time.Time        `json:"updated-at"`
}

// Field represents a template field value
type Field struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// Review represents a card review record
type Review struct {
	Date       time.Time `json:"date"`
	Due        time.Time `json:"due"`
	Remembered bool      `json:"remembered?"`
}

// Deck represents a deck (collection) in Mochi
type Deck struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	ParentID        string     `json:"parent-id,omitempty"`
	Sort            int        `json:"sort"`
	Archived        bool       `json:"archived?"`
	Trashed         *time.Time `json:"trashed?"`
	SortBy          string     `json:"sort-by,omitempty"`
	CardsView       string     `json:"cards-view,omitempty"`
	ShowSides       bool       `json:"show-sides?"`
	SortByDirection bool       `json:"sort-by-direction"`
	ReviewReverse   bool       `json:"review-reverse?"`
}

// Template represents a card template in Mochi
type Template struct {
	ID      string                   `json:"id"`
	Name    string                   `json:"name"`
	Content string                   `json:"content"`
	Pos     string                   `json:"pos"`
	Fields  map[string]TemplateField `json:"fields"`
	Style   map[string]interface{}   `json:"style,omitempty"`
	Options map[string]interface{}   `json:"options,omitempty"`
}

// TemplateField represents a field definition in a template
type TemplateField struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name,omitempty"`
	Type    string                 `json:"type,omitempty"`
	Pos     string                 `json:"pos,omitempty"`
	Content string                 `json:"content,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Bookmark string      `json:"bookmark,omitempty"`
	Docs     interface{} `json:"docs"`
}

// DueResponse represents the response from the due cards endpoint
type DueResponse struct {
	Cards []Card `json:"cards"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Errors interface{} `json:"errors"`
}

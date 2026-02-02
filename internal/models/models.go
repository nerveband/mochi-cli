package models

import (
	"encoding/json"
	"strings"
	"time"
)

// MochiTime handles Mochi API time formats.
type MochiTime struct {
	time.Time
}

// UnmarshalJSON accepts RFC3339 strings or {"date": "..."} objects.
func (t *MochiTime) UnmarshalJSON(data []byte) error {
	if t == nil {
		return nil
	}
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		return nil
	}
	if strings.HasPrefix(trimmed, "\"") {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		if s == "" {
			return nil
		}
		parsed, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			parsed, err = time.Parse(time.RFC3339, s)
			if err != nil {
				return err
			}
		}
		t.Time = parsed
		return nil
	}
	var obj struct {
		Date string `json:"date"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if obj.Date == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, obj.Date)
	if err != nil {
		parsed, err = time.Parse(time.RFC3339, obj.Date)
		if err != nil {
			return err
		}
	}
	t.Time = parsed
	return nil
}

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
	Trashed       *MochiTime       `json:"trashed?"`
	New           bool             `json:"new?"`
	References    []string         `json:"references"`
	Reviews       []Review         `json:"reviews"`
	CreatedAt     *MochiTime       `json:"created-at"`
	UpdatedAt     *MochiTime       `json:"updated-at"`
}

// Field represents a template field value
type Field struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// Review represents a card review record
type Review struct {
	Date       *MochiTime `json:"date"`
	Due        *MochiTime `json:"due"`
	Remembered bool       `json:"remembered?"`
}

// Deck represents a deck (collection) in Mochi
type Deck struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	ParentID        string     `json:"parent-id,omitempty"`
	Sort            int        `json:"sort"`
	Archived        bool       `json:"archived?"`
	Trashed         *MochiTime `json:"trashed?"`
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

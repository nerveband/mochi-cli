package mochiimport

import (
	"fmt"
	"strings"
	"time"
)

// MochiData represents the top-level structure of a .mochi file
type MochiData struct {
	Version   int             `json:"version" edn:"version"`
	Decks     []MochiDeck     `json:"decks,omitempty" edn:"decks,omitempty"`
	Cards     []MochiCard     `json:"cards,omitempty" edn:"cards,omitempty"`
	Templates []MochiTemplate `json:"templates,omitempty" edn:"templates,omitempty"`
}

// MochiDeck represents a deck in the Mochi format
type MochiDeck struct {
	ID       string      `json:"id,omitempty" edn:"id,omitempty"`
	Name     string      `json:"name" edn:"name"`
	ParentID string      `json:"parent-id,omitempty" edn:"parent-id,omitempty"`
	Cards    []MochiCard `json:"cards,omitempty" edn:"cards,omitempty"`
}

// MochiCard represents a card in the Mochi format
type MochiCard struct {
	ID      string                 `json:"id,omitempty" edn:"id,omitempty"`
	Name    string                 `json:"name,omitempty" edn:"name,omitempty"`
	Content string                 `json:"content" edn:"content"`
	DeckID  string                 `json:"deck-id,omitempty" edn:"deck-id,omitempty"`
	Pos     string                 `json:"pos,omitempty" edn:"pos,omitempty"`
	Fields  map[string]interface{} `json:"fields,omitempty" edn:"fields,omitempty"`
	Reviews []MochiReview          `json:"reviews,omitempty" edn:"reviews,omitempty"`
}

// MochiTemplate represents a template in the Mochi format
type MochiTemplate struct {
	ID      string                `json:"id" edn:"id"`
	Name    string                `json:"name" edn:"name"`
	Content string                `json:"content,omitempty" edn:"content,omitempty"`
	Pos     string                `json:"pos,omitempty" edn:"pos,omitempty"`
	Fields  map[string]MochiField `json:"fields,omitempty" edn:"fields,omitempty"`
}

// MochiField represents a field definition in a template
type MochiField struct {
	ID             string                 `json:"id" edn:"id"`
	Name           string                 `json:"name" edn:"name"`
	Type           string                 `json:"type,omitempty" edn:"type,omitempty"`
	Pos            string                 `json:"pos,omitempty" edn:"pos,omitempty"`
	Options        map[string]interface{} `json:"options,omitempty" edn:"options,omitempty"`
	Lang           string                 `json:"lang,omitempty" edn:"lang,omitempty"`
	From           string                 `json:"from,omitempty" edn:"from,omitempty"`
	To             string                 `json:"to,omitempty" edn:"to,omitempty"`
	BooleanDefault bool                   `json:"boolean-default,omitempty" edn:"boolean-default,omitempty"`
}

// MochiReview represents a review record
type MochiReview struct {
	Date       string `json:"date" edn:"date"`
	Due        string `json:"due" edn:"due"`
	Interval   int    `json:"interval" edn:"interval"`
	Remembered bool   `json:"remembered?" edn:"remembered?"`
}

// ImportOptions contains options for importing
type ImportOptions struct {
	Format     string // "json" or "edn"
	SkipMedia  bool
	DeckID     string // Import into specific deck
	TemplateID string // Use specific template
	DryRun     bool   // Preview only
}

// ExportOptions contains options for exporting
type ExportOptions struct {
	Format         string // "json" or "edn"
	IncludeMedia   bool
	OnlyCards      bool // Export only cards, no decks structure
	IncludeReviews bool // Include review history
}

// Validate validates the MochiData structure
func (m *MochiData) Validate() error {
	if m.Version != 2 {
		return fmt.Errorf("unsupported version: %d (only version 2 is supported)", m.Version)
	}

	// Validate deck IDs are unique
	deckIDs := make(map[string]bool)
	for _, deck := range m.Decks {
		if deck.ID != "" {
			if deckIDs[deck.ID] {
				return fmt.Errorf("duplicate deck ID: %s", deck.ID)
			}
			deckIDs[deck.ID] = true
		}
	}

	// Validate cards have deck references
	for _, card := range m.Cards {
		if card.DeckID == "" {
			return fmt.Errorf("top-level card must have deck-id: %s", card.Name)
		}
		if card.DeckID != "" && !deckIDs[card.DeckID] {
			return fmt.Errorf("card references non-existent deck: %s", card.DeckID)
		}
	}

	// Validate nested cards have proper deck references
	for _, deck := range m.Decks {
		for _, card := range deck.Cards {
			if card.DeckID != "" && card.DeckID != deck.ID {
				return fmt.Errorf("nested card deck-id mismatch: card has %s but parent deck is %s",
					card.DeckID, deck.ID)
			}
		}
	}

	return nil
}

// GenerateID generates a unique ID for Mochi entities
func GenerateID() string {
	// Generate an 8-character alphanumeric ID
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// ToKeyword converts a string to a keyword format (for EDN)
func ToKeyword(s string) string {
	if strings.HasPrefix(s, ":") {
		return s
	}
	return ":" + s
}

// FromKeyword removes the leading colon from a keyword
func FromKeyword(s string) string {
	return strings.TrimPrefix(s, ":")
}

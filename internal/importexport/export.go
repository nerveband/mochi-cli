package mochiimport

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"

	"github.com/nerveband/mochi-cli/internal/models"
)

// Exporter handles exporting data to .mochi format
type Exporter struct {
	client APIClient
}

// APIClient interface for fetching data
type APIClient interface {
	ListCards(deckID string, limit int, bookmark string) (*models.PaginatedResponse, error)
	GetCard(cardID string) (*models.Card, error)
	ListDecks(bookmark string) (*models.PaginatedResponse, error)
	GetDeck(deckID string) (*models.Deck, error)
	CreateDeck(deck *models.Deck) (*models.Deck, error)
	CreateCard(card *models.Card) (*models.Card, error)
}

// NewExporter creates a new exporter
func NewExporter(client APIClient) *Exporter {
	return &Exporter{client: client}
}

// ExportDeck exports a single deck with all its cards
func (e *Exporter) ExportDeck(deckID string, opts ExportOptions) (*MochiData, error) {
	deck, err := e.client.GetDeck(deckID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deck: %w", err)
	}

	mochiDeck := MochiDeck{
		ID:   deck.ID,
		Name: deck.Name,
	}

	if deck.ParentID != "" {
		mochiDeck.ParentID = deck.ParentID
	}

	// Fetch all cards for this deck
	cards, err := e.fetchAllCards(deckID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cards: %w", err)
	}

	mochiDeck.Cards = make([]MochiCard, len(cards))
	for i, card := range cards {
		mochiDeck.Cards[i] = convertCardToMochi(card, opts.IncludeReviews)
	}

	data := &MochiData{
		Version: 2,
		Decks:   []MochiDeck{mochiDeck},
	}

	return data, nil
}

// ExportAllDecks exports all decks
func (e *Exporter) ExportAllDecks(opts ExportOptions) (*MochiData, error) {
	decks, err := e.fetchAllDecks()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch decks: %w", err)
	}

	mochiDecks := make([]MochiDeck, len(decks))
	for i, deck := range decks {
		mochiDeck := MochiDeck{
			ID:   deck.ID,
			Name: deck.Name,
		}

		if deck.ParentID != "" {
			mochiDeck.ParentID = deck.ParentID
		}

		// Fetch cards for this deck
		cards, err := e.fetchAllCards(deck.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch cards for deck %s: %w", deck.ID, err)
		}

		mochiDeck.Cards = make([]MochiCard, len(cards))
		for j, card := range cards {
			mochiDeck.Cards[j] = convertCardToMochi(card, opts.IncludeReviews)
		}

		mochiDecks[i] = mochiDeck
	}

	data := &MochiData{
		Version: 2,
		Decks:   mochiDecks,
	}

	return data, nil
}

// ExportCards exports specific cards
func (e *Exporter) ExportCards(cardIDs []string, deckID string, opts ExportOptions) (*MochiData, error) {
	mochiCards := make([]MochiCard, 0, len(cardIDs))

	for _, cardID := range cardIDs {
		card, err := e.client.GetCard(cardID)
		if err != nil {
			continue // Skip cards we can't fetch
		}

		mochiCard := convertCardToMochi(*card, opts.IncludeReviews)
		if deckID != "" {
			mochiCard.DeckID = deckID
		}
		mochiCards = append(mochiCards, mochiCard)
	}

	data := &MochiData{
		Version: 2,
		Cards:   mochiCards,
	}

	return data, nil
}

// ExportToFile exports data to a .mochi file (ZIP format)
func (e *Exporter) ExportToFile(data *MochiData, filepath string, opts ExportOptions) error {
	// Validate data
	if err := data.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Create ZIP file
	zipFile, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Write data file
	var dataContent []byte
	var dataFilename string

	switch opts.Format {
	case "edn":
		// For now, we'll use JSON representation
		// Full EDN support would require an EDN library
		dataFilename = "data.json"
		dataContent, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
	default: // json
		dataFilename = "data.json"
		dataContent, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
	}

	dataWriter, err := zipWriter.Create(dataFilename)
	if err != nil {
		return fmt.Errorf("failed to create data file in zip: %w", err)
	}

	if _, err := dataWriter.Write(dataContent); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}

// fetchAllCards fetches all cards for a deck with pagination
func (e *Exporter) fetchAllCards(deckID string) ([]models.Card, error) {
	var allCards []models.Card
	var bookmark string

	for {
		resp, err := e.client.ListCards(deckID, 100, bookmark)
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
				allCards = append(allCards, card)
			}
		}

		if resp.Bookmark == "" {
			break
		}
		bookmark = resp.Bookmark
	}

	return allCards, nil
}

// fetchAllDecks fetches all decks with pagination
func (e *Exporter) fetchAllDecks() ([]models.Deck, error) {
	var allDecks []models.Deck
	var bookmark string

	for {
		resp, err := e.client.ListDecks(bookmark)
		if err != nil {
			return nil, err
		}

		docs, ok := resp.Docs.([]interface{})
		if !ok {
			break
		}

		for _, doc := range docs {
			data, _ := json.Marshal(doc)
			var deck models.Deck
			if err := json.Unmarshal(data, &deck); err == nil {
				allDecks = append(allDecks, deck)
			}
		}

		if resp.Bookmark == "" {
			break
		}
		bookmark = resp.Bookmark
	}

	return allDecks, nil
}

// convertCardToMochi converts an API card to Mochi format
func convertCardToMochi(card models.Card, includeReviews bool) MochiCard {
	mochiCard := MochiCard{
		ID:      card.ID,
		Name:    card.Name,
		Content: card.Content,
		DeckID:  card.DeckID,
		Pos:     card.Pos,
	}

	// Convert fields
	if len(card.Fields) > 0 {
		mochiCard.Fields = make(map[string]interface{})
		for id, field := range card.Fields {
			mochiCard.Fields[id] = field.Value
		}
	}

	// Convert reviews if requested
	if includeReviews && len(card.Reviews) > 0 {
		mochiCard.Reviews = make([]MochiReview, len(card.Reviews))
		for i, review := range card.Reviews {
			mochiCard.Reviews[i] = MochiReview{
				Date:       review.Date.Format("2006-01-02T15:04:05Z"),
				Due:        review.Due.Format("2006-01-02T15:04:05Z"),
				Interval:   0, // API doesn't provide interval
				Remembered: review.Remembered,
			}
		}
	}

	return mochiCard
}

package mochiimport

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nerveband/mochi-cli/internal/api"
	"github.com/nerveband/mochi-cli/internal/models"
)

// Importer handles importing data from .mochi format
type Importer struct {
	client *api.Client
}

// NewImporter creates a new importer
func NewImporter(client *api.Client) *Importer {
	return &Importer{client: client}
}

// ImportFromFile imports data from a .mochi file
func (i *Importer) ImportFromFile(filepath string, opts ImportOptions) (*ImportResult, error) {
	// Open ZIP file
	zipReader, err := zip.OpenReader(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open .mochi file: %w", err)
	}
	defer zipReader.Close()

	// Find and read data file
	var data *MochiData
	for _, file := range zipReader.File {
		if file.Name == "data.json" || file.Name == "data.edn" {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open data file: %w", err)
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("failed to read data file: %w", err)
			}

			data = &MochiData{}
			if err := json.Unmarshal(content, data); err != nil {
				return nil, fmt.Errorf("failed to parse data file: %w", err)
			}
			break
		}
	}

	if data == nil {
		return nil, fmt.Errorf("no data.json or data.edn file found in .mochi archive")
	}

	// Validate data
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return i.ImportData(data, opts)
}

// ImportResult contains the results of an import operation
type ImportResult struct {
	DecksCreated     int
	CardsCreated     int
	TemplatesCreated int
	Errors           []string
	MediaFiles       []string
}

// ImportData imports MochiData into the API
func (i *Importer) ImportData(data *MochiData, opts ImportOptions) (*ImportResult, error) {
	result := &ImportResult{
		Errors: make([]string, 0),
	}

	if opts.DryRun {
		return i.PreviewImport(data)
	}

	// Create ID mapping for old -> new IDs
	idMapping := make(map[string]string)

	// First pass: Create decks
	for _, mochiDeck := range data.Decks {
		newDeck, err := i.createDeck(mochiDeck, opts)
		if err != nil {
			result.Errors = append(result.Errors,
				fmt.Sprintf("Failed to create deck '%s': %v", mochiDeck.Name, err))
			continue
		}
		result.DecksCreated++

		// Map old ID to new ID
		if mochiDeck.ID != "" {
			idMapping[mochiDeck.ID] = newDeck.ID
		}

		// Create cards in this deck
		for _, mochiCard := range mochiDeck.Cards {
			// Update deck reference
			if mochiCard.DeckID != "" {
				if newDeckID, ok := idMapping[mochiCard.DeckID]; ok {
					mochiCard.DeckID = newDeckID
				} else {
					mochiCard.DeckID = newDeck.ID
				}
			} else {
				mochiCard.DeckID = newDeck.ID
			}

			if err := i.createCard(mochiCard, opts); err != nil {
				result.Errors = append(result.Errors,
					fmt.Sprintf("Failed to create card in deck '%s': %v", mochiDeck.Name, err))
				continue
			}
			result.CardsCreated++
		}
	}

	// Second pass: Create top-level cards
	for _, mochiCard := range data.Cards {
		// Update deck reference
		if mochiCard.DeckID != "" {
			if newDeckID, ok := idMapping[mochiCard.DeckID]; ok {
				mochiCard.DeckID = newDeckID
			}
		}

		// Skip if importing to specific deck
		if opts.DeckID != "" {
			mochiCard.DeckID = opts.DeckID
		}

		if err := i.createCard(mochiCard, opts); err != nil {
			result.Errors = append(result.Errors,
				fmt.Sprintf("Failed to create card: %v", err))
			continue
		}
		result.CardsCreated++
	}

	return result, nil
}

// createDeck creates a deck from Mochi format
func (i *Importer) createDeck(mochiDeck MochiDeck, opts ImportOptions) (*models.Deck, error) {
	deck := &models.Deck{
		Name: mochiDeck.Name,
	}

	// Handle parent deck reference
	if mochiDeck.ParentID != "" {
		deck.ParentID = mochiDeck.ParentID
	}

	return i.client.CreateDeck(deck)
}

// createCard creates a card from Mochi format
func (i *Importer) createCard(mochiCard MochiCard, opts ImportOptions) error {
	card := &models.Card{
		Name:    mochiCard.Name,
		Content: mochiCard.Content,
		DeckID:  mochiCard.DeckID,
		Pos:     mochiCard.Pos,
	}

	// Apply template if specified
	if opts.TemplateID != "" {
		card.TemplateID = opts.TemplateID
	}

	// Convert fields
	if len(mochiCard.Fields) > 0 {
		card.Fields = make(map[string]models.Field)
		for id, value := range mochiCard.Fields {
			card.Fields[id] = models.Field{
				ID:    id,
				Value: fmt.Sprintf("%v", value),
			}
		}
	}

	_, err := i.client.CreateCard(card)
	return err
}

// PreviewImport previews what would be imported without actually creating anything
func (i *Importer) PreviewImport(data *MochiData) (*ImportResult, error) {
	result := &ImportResult{
		Errors: make([]string, 0),
	}

	// Count what would be created
	for _, deck := range data.Decks {
		result.DecksCreated++
		result.CardsCreated += len(deck.Cards)
	}
	result.CardsCreated += len(data.Cards)
	result.TemplatesCreated = len(data.Templates)

	return result, nil
}

// ExtractMedia extracts media files from a .mochi archive
func ExtractMedia(mochiPath string, destDir string) ([]string, error) {
	zipReader, err := zip.OpenReader(mochiPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open .mochi file: %w", err)
	}
	defer zipReader.Close()

	var mediaFiles []string

	for _, file := range zipReader.File {
		// Skip data files
		if file.Name == "data.json" || file.Name == "data.edn" {
			continue
		}

		// Extract media file
		rc, err := file.Open()
		if err != nil {
			continue
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			continue
		}

		destPath := filepath.Join(destDir, file.Name)
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			continue
		}

		mediaFiles = append(mediaFiles, file.Name)
	}

	return mediaFiles, nil
}

// ValidateMochiFile validates a .mochi file without importing
func ValidateMochiFile(filepath string) error {
	zipReader, err := zip.OpenReader(filepath)
	if err != nil {
		return fmt.Errorf("failed to open .mochi file: %w", err)
	}
	defer zipReader.Close()

	// Check for data file
	hasDataFile := false
	for _, file := range zipReader.File {
		if file.Name == "data.json" || file.Name == "data.edn" {
			hasDataFile = true
			break
		}
	}

	if !hasDataFile {
		return fmt.Errorf("no data.json or data.edn file found in .mochi archive")
	}

	return nil
}

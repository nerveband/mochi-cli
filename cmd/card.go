package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nerveband/mochi-cli/internal/models"
	"github.com/spf13/cobra"
)

// cardCmd represents the card command
var cardCmd = &cobra.Command{
	Use:   "card",
	Short: "Manage cards",
	Long:  `Create, list, get, update, and delete flashcards.`,
}

// cardListCmd lists cards
var cardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cards",
	Long:  `List cards with optional filtering by deck.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		deckID, _ := cmd.Flags().GetString("deck")
		limit, _ := cmd.Flags().GetInt("limit")

		client, err := getClient()
		if err != nil {
			return err
		}

		resp, err := client.ListCards(deckID, limit, "")
		if err != nil {
			return err
		}

		// Convert docs to cards
		docs, ok := resp.Docs.([]interface{})
		if !ok {
			return fmt.Errorf("unexpected response format")
		}

		var cards []models.Card
		for _, doc := range docs {
			data, _ := json.Marshal(doc)
			var card models.Card
			if err := json.Unmarshal(data, &card); err == nil {
				cards = append(cards, card)
			}
		}

		// Handle output formats
		if idOnly || outputOnly == "id" {
			for _, card := range cards {
				fmt.Println(card.ID)
			}
			return nil
		}

		if outputOnly != "" {
			for _, card := range cards {
				val := extractField(card, outputOnly)
				if val != nil {
					fmt.Println(val)
				}
			}
			return nil
		}

		switch format {
		case "json":
			printJSON(map[string]interface{}{
				"cards":    cards,
				"bookmark": resp.Bookmark,
			})
		case "compact":
			printCompactJSON(cards)
		case "table":
			headers := []string{"ID", "NAME", "DECK", "CONTENT"}
			rows := make([][]string, len(cards))
			for i, card := range cards {
				content := truncateString(strings.ReplaceAll(card.Content, "\n", " "), 40)
				rows[i] = []string{
					card.ID,
					card.Name,
					card.DeckID,
					content,
				}
			}
			printTable(headers, rows)
		default:
			for _, card := range cards {
				fmt.Printf("%s: %s\n", card.ID, truncateString(card.Content, 60))
			}
		}

		return nil
	},
}

// cardGetCmd gets a specific card
var cardGetCmd = &cobra.Command{
	Use:   "get <card-id>",
	Short: "Get a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID := args[0]

		client, err := getClient()
		if err != nil {
			return err
		}

		card, err := client.GetCard(cardID)
		if err != nil {
			return err
		}

		if idOnly || outputOnly == "id" {
			fmt.Println(card.ID)
			return nil
		}

		if outputOnly != "" {
			val := extractField(card, outputOnly)
			if val != nil {
				fmt.Println(val)
			}
			return nil
		}

		switch format {
		case "json":
			printJSON(card)
		case "compact":
			printCompactJSON(card)
		case "markdown":
			fmt.Printf("# %s\n\n", card.Name)
			fmt.Println(card.Content)
			fmt.Printf("\n---\n")
			fmt.Printf("**ID:** %s\n", card.ID)
			fmt.Printf("**Deck:** %s\n", card.DeckID)
			fmt.Printf("**Created:** %s\n", formatTime(card.CreatedAt))
		default:
			fmt.Printf("ID: %s\n", card.ID)
			fmt.Printf("Name: %s\n", card.Name)
			fmt.Printf("Deck: %s\n", card.DeckID)
			fmt.Printf("Content:\n%s\n", card.Content)
		}

		return nil
	},
}

// cardCreateCmd creates a new card
var cardCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new card",
	Long:  `Create a new flashcard. Content can be provided via --content, --file, or stdin.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		deckID, _ := cmd.Flags().GetString("deck")
		content, _ := cmd.Flags().GetString("content")
		name, _ := cmd.Flags().GetString("name")
		templateID, _ := cmd.Flags().GetString("template")
		file, _ := cmd.Flags().GetString("file")
		stdin, _ := cmd.Flags().GetBool("stdin")

		if deckID == "" {
			return fmt.Errorf("deck ID is required (use --deck)")
		}

		// Get content from various sources
		if stdin {
			data, err := readStdin()
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			content = data
		} else if file != "" {
			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			content = string(data)
		}

		if content == "" {
			return fmt.Errorf("content is required (use --content, --file, or --stdin)")
		}

		if dryRun {
			printInfo("Dry run - would create card:")
			printInfo(fmt.Sprintf("  Deck: %s", deckID))
			printInfo(fmt.Sprintf("  Name: %s", name))
			printInfo(fmt.Sprintf("  Content: %s", truncateString(content, 50)))
			return nil
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		card := &models.Card{
			DeckID:     deckID,
			Content:    content,
			Name:       name,
			TemplateID: templateID,
		}

		created, err := client.CreateCard(card)
		if err != nil {
			return err
		}

		if idOnly || outputOnly == "id" {
			fmt.Println(created.ID)
			return nil
		}

		if !quiet {
			printSuccess(fmt.Sprintf("Card created: %s", created.ID))
		}

		if format == "json" {
			printJSON(created)
		}

		return nil
	},
}

// cardUpdateCmd updates a card
var cardUpdateCmd = &cobra.Command{
	Use:   "update <card-id>",
	Short: "Update a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID := args[0]

		content, _ := cmd.Flags().GetString("content")
		name, _ := cmd.Flags().GetString("name")
		deckID, _ := cmd.Flags().GetString("deck")
		archived, _ := cmd.Flags().GetBool("archive")
		unarchive, _ := cmd.Flags().GetBool("unarchive")
		file, _ := cmd.Flags().GetString("file")
		stdin, _ := cmd.Flags().GetBool("stdin")

		// Get content from various sources
		if stdin {
			data, err := readStdin()
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			content = data
		} else if file != "" {
			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			content = string(data)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		// Get existing card to update only provided fields
		existing, err := client.GetCard(cardID)
		if err != nil {
			return err
		}

		update := &models.Card{
			ID: existing.ID,
		}

		if content != "" {
			update.Content = content
		}
		if name != "" {
			update.Name = name
		}
		if deckID != "" {
			update.DeckID = deckID
		}
		if archived {
			update.Archived = true
		} else if unarchive {
			update.Archived = false
		}

		if dryRun {
			printInfo("Dry run - would update card:")
			printInfo(fmt.Sprintf("  ID: %s", cardID))
			if content != "" {
				printInfo(fmt.Sprintf("  Content: %s", truncateString(content, 50)))
			}
			if name != "" {
				printInfo(fmt.Sprintf("  Name: %s", name))
			}
			return nil
		}

		updated, err := client.UpdateCard(cardID, update)
		if err != nil {
			return err
		}

		if !quiet {
			printSuccess(fmt.Sprintf("Card updated: %s", updated.ID))
		}

		if format == "json" {
			printJSON(updated)
		}

		return nil
	},
}

// cardDeleteCmd deletes a card
var cardDeleteCmd = &cobra.Command{
	Use:   "delete <card-id>",
	Short: "Delete a card",
	Long:  `Permanently delete a card. This cannot be undone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID := args[0]
		force, _ := cmd.Flags().GetBool("force")

		if !force && !quiet {
			fmt.Printf("Delete card %s? This cannot be undone. [y/N]: ", cardID)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled")
				return nil
			}
		}

		if dryRun {
			printInfo(fmt.Sprintf("Dry run - would delete card: %s", cardID))
			return nil
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		if err := client.DeleteCard(cardID); err != nil {
			return err
		}

		if !quiet {
			printSuccess(fmt.Sprintf("Card deleted: %s", cardID))
		}

		return nil
	},
}

// cardSearchCmd searches cards
var cardSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search cards",
	Long:  `Search for cards by content (client-side search).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		deckID, _ := cmd.Flags().GetString("deck")

		client, err := getClient()
		if err != nil {
			return err
		}

		cards, err := client.SearchCards(query, deckID)
		if err != nil {
			return err
		}

		if idOnly || outputOnly == "id" {
			for _, card := range cards {
				fmt.Println(card.ID)
			}
			return nil
		}

		switch format {
		case "json":
			printJSON(map[string]interface{}{
				"cards": cards,
				"count": len(cards),
			})
		case "compact":
			printCompactJSON(cards)
		case "table":
			headers := []string{"ID", "NAME", "DECK", "CONTENT"}
			rows := make([][]string, len(cards))
			for i, card := range cards {
				content := truncateString(strings.ReplaceAll(card.Content, "\n", " "), 40)
				rows[i] = []string{
					card.ID,
					card.Name,
					card.DeckID,
					content,
				}
			}
			printTable(headers, rows)
		default:
			for _, card := range cards {
				fmt.Printf("%s: %s\n", card.ID, truncateString(card.Content, 60))
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cardCmd)
	cardCmd.AddCommand(cardListCmd)
	cardCmd.AddCommand(cardGetCmd)
	cardCmd.AddCommand(cardCreateCmd)
	cardCmd.AddCommand(cardUpdateCmd)
	cardCmd.AddCommand(cardDeleteCmd)
	cardCmd.AddCommand(cardSearchCmd)

	// List flags
	cardListCmd.Flags().StringP("deck", "d", "", "Filter by deck ID")
	cardListCmd.Flags().IntP("limit", "l", 10, "Number of cards to return")

	// Get flags (none needed)

	// Create flags
	cardCreateCmd.Flags().StringP("deck", "d", "", "Deck ID (required)")
	cardCreateCmd.Flags().StringP("content", "c", "", "Card content (markdown)")
	cardCreateCmd.Flags().StringP("name", "n", "", "Card name")
	cardCreateCmd.Flags().StringP("template", "t", "", "Template ID")
	cardCreateCmd.Flags().StringP("file", "f", "", "Read content from file")
	cardCreateCmd.Flags().Bool("stdin", false, "Read content from stdin")

	// Update flags
	cardUpdateCmd.Flags().StringP("content", "c", "", "New content")
	cardUpdateCmd.Flags().StringP("name", "n", "", "New name")
	cardUpdateCmd.Flags().StringP("deck", "d", "", "Move to different deck")
	cardUpdateCmd.Flags().Bool("archive", false, "Archive the card")
	cardUpdateCmd.Flags().Bool("unarchive", false, "Unarchive the card")
	cardUpdateCmd.Flags().StringP("file", "f", "", "Read content from file")
	cardUpdateCmd.Flags().Bool("stdin", false, "Read content from stdin")

	// Delete flags
	cardDeleteCmd.Flags().Bool("force", false, "Skip confirmation prompt")

	// Search flags
	cardSearchCmd.Flags().StringP("deck", "d", "", "Limit search to specific deck")
}

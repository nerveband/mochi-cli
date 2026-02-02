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

// deckCmd represents the deck command
var deckCmd = &cobra.Command{
	Use:   "deck",
	Short: "Manage decks",
	Long:  `Create, list, get, update, and delete decks.`,
}

// deckListCmd lists decks
var deckListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all decks",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		resp, err := client.ListDecks("")
		if err != nil {
			return err
		}

		// Convert docs to decks
		docs, ok := resp.Docs.([]interface{})
		if !ok {
			return fmt.Errorf("unexpected response format")
		}

		var decks []models.Deck
		for _, doc := range docs {
			data, _ := json.Marshal(doc)
			var deck models.Deck
			if err := json.Unmarshal(data, &deck); err == nil {
				decks = append(decks, deck)
			}
		}

		if idOnly || outputOnly == "id" {
			for _, deck := range decks {
				fmt.Println(deck.ID)
			}
			return nil
		}

		if outputOnly != "" {
			for _, deck := range decks {
				val := extractField(deck, outputOnly)
				if val != nil {
					fmt.Println(val)
				}
			}
			return nil
		}

		switch format {
		case "json":
			printJSON(map[string]interface{}{
				"decks":    decks,
				"bookmark": resp.Bookmark,
			})
		case "compact":
			printCompactJSON(decks)
		case "table":
			headers := []string{"ID", "NAME", "SORT", "ARCHIVED"}
			rows := make([][]string, len(decks))
			for i, deck := range decks {
				archived := ""
				if deck.Archived {
					archived = "yes"
				}
				rows[i] = []string{
					deck.ID,
					deck.Name,
					fmt.Sprintf("%d", deck.Sort),
					archived,
				}
			}
			printTable(headers, rows)
		default:
			for _, deck := range decks {
				fmt.Printf("%s: %s\n", deck.ID, deck.Name)
			}
		}

		return nil
	},
}

// deckGetCmd gets a specific deck
var deckGetCmd = &cobra.Command{
	Use:   "get <deck-id>",
	Short: "Get a deck",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckID := args[0]

		client, err := getClient()
		if err != nil {
			return err
		}

		deck, err := client.GetDeck(deckID)
		if err != nil {
			return err
		}

		if idOnly || outputOnly == "id" {
			fmt.Println(deck.ID)
			return nil
		}

		if outputOnly != "" {
			val := extractField(deck, outputOnly)
			if val != nil {
				fmt.Println(val)
			}
			return nil
		}

		switch format {
		case "json":
			printJSON(deck)
		case "compact":
			printCompactJSON(deck)
		default:
			fmt.Printf("ID: %s\n", deck.ID)
			fmt.Printf("Name: %s\n", deck.Name)
			fmt.Printf("Sort: %d\n", deck.Sort)
			if deck.ParentID != "" {
				fmt.Printf("Parent: %s\n", deck.ParentID)
			}
			fmt.Printf("Archived: %v\n", deck.Archived)
		}

		return nil
	},
}

// deckCreateCmd creates a new deck
var deckCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new deck",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		parentID, _ := cmd.Flags().GetString("parent")
		sort, _ := cmd.Flags().GetInt("sort")

		if dryRun {
			printInfo("Dry run - would create deck:")
			printInfo(fmt.Sprintf("  Name: %s", name))
			if parentID != "" {
				printInfo(fmt.Sprintf("  Parent: %s", parentID))
			}
			return nil
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		deck := &models.Deck{
			Name:     name,
			ParentID: parentID,
			Sort:     sort,
		}

		created, err := client.CreateDeck(deck)
		if err != nil {
			return err
		}

		if idOnly || outputOnly == "id" {
			fmt.Println(created.ID)
			return nil
		}

		if !quiet {
			printSuccess(fmt.Sprintf("Deck created: %s", created.ID))
		}

		if format == "json" {
			printJSON(created)
		}

		return nil
	},
}

// deckUpdateCmd updates a deck
var deckUpdateCmd = &cobra.Command{
	Use:   "update <deck-id>",
	Short: "Update a deck",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckID := args[0]

		name, _ := cmd.Flags().GetString("name")
		parentID, _ := cmd.Flags().GetString("parent")
		sort, _ := cmd.Flags().GetInt("sort")
		archive, _ := cmd.Flags().GetBool("archive")
		unarchive, _ := cmd.Flags().GetBool("unarchive")

		client, err := getClient()
		if err != nil {
			return err
		}

		// Get existing deck
		existing, err := client.GetDeck(deckID)
		if err != nil {
			return err
		}

		update := &models.Deck{
			ID: existing.ID,
		}

		if name != "" {
			update.Name = name
		}
		if cmd.Flags().Changed("parent") {
			update.ParentID = parentID
		}
		if cmd.Flags().Changed("sort") {
			update.Sort = sort
		}
		if archive {
			update.Archived = true
		} else if unarchive {
			update.Archived = false
		}

		if dryRun {
			printInfo("Dry run - would update deck:")
			printInfo(fmt.Sprintf("  ID: %s", deckID))
			if name != "" {
				printInfo(fmt.Sprintf("  Name: %s", name))
			}
			return nil
		}

		updated, err := client.UpdateDeck(deckID, update)
		if err != nil {
			return err
		}

		if !quiet {
			printSuccess(fmt.Sprintf("Deck updated: %s", updated.ID))
		}

		if format == "json" {
			printJSON(updated)
		}

		return nil
	},
}

// deckDeleteCmd deletes a deck
var deckDeleteCmd = &cobra.Command{
	Use:   "delete <deck-id>",
	Short: "Delete a deck",
	Long:  `Permanently delete a deck. This cannot be undone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckID := args[0]
		force, _ := cmd.Flags().GetBool("force")

		if !force && !quiet {
			fmt.Printf("Delete deck %s? This cannot be undone. [y/N]: ", deckID)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled")
				return nil
			}
		}

		if dryRun {
			printInfo(fmt.Sprintf("Dry run - would delete deck: %s", deckID))
			return nil
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		if err := client.DeleteDeck(deckID); err != nil {
			return err
		}

		if !quiet {
			printSuccess(fmt.Sprintf("Deck deleted: %s", deckID))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deckCmd)
	deckCmd.AddCommand(deckListCmd)
	deckCmd.AddCommand(deckGetCmd)
	deckCmd.AddCommand(deckCreateCmd)
	deckCmd.AddCommand(deckUpdateCmd)
	deckCmd.AddCommand(deckDeleteCmd)

	// Create flags
	deckCreateCmd.Flags().StringP("parent", "P", "", "Parent deck ID")
	deckCreateCmd.Flags().IntP("sort", "s", 0, "Sort order")

	// Update flags
	deckUpdateCmd.Flags().StringP("name", "n", "", "New name")
	deckUpdateCmd.Flags().StringP("parent", "P", "", "Parent deck ID")
	deckUpdateCmd.Flags().IntP("sort", "s", 0, "Sort order")
	deckUpdateCmd.Flags().Bool("archive", false, "Archive the deck")
	deckUpdateCmd.Flags().Bool("unarchive", false, "Unarchive the deck")

	// Delete flags
	deckDeleteCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}

package cmd

import (
	"fmt"
	"os"

	importexport "github.com/nerveband/mochi-cli/internal/importexport"
	"github.com/spf13/cobra"
)

// importexportCmd represents the import/export command
var importexportCmd = &cobra.Command{
	Use:     "import-export",
	Short:   "Import and export Mochi data",
	Long:    `Import and export decks, cards, and templates in the .mochi format.`,
	Aliases: []string{"ie"},
}

// exportCmd exports data to a .mochi file
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export decks and cards to .mochi file",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		deckID, _ := cmd.Flags().GetString("deck")
		format, _ := cmd.Flags().GetString("format")
		includeReviews, _ := cmd.Flags().GetBool("include-reviews")

		if output == "" {
			return fmt.Errorf("output file is required (use --output)")
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		exporter := importexport.NewExporter(client)
		opts := importexport.ExportOptions{
			Format:         format,
			IncludeReviews: includeReviews,
		}

		var data *importexport.MochiData

		if deckID != "" {
			// Export single deck
			if !quiet {
				fmt.Printf("Exporting deck %s...\n", deckID)
			}
			data, err = exporter.ExportDeck(deckID, opts)
		} else {
			// Export all decks
			if !quiet {
				fmt.Println("Exporting all decks...")
			}
			data, err = exporter.ExportAllDecks(opts)
		}

		if err != nil {
			return fmt.Errorf("failed to export: %w", err)
		}

		if dryRun {
			printInfo("Dry run - would export:")
			printInfo(fmt.Sprintf("  Decks: %d", len(data.Decks)))

			cardCount := len(data.Cards)
			for _, deck := range data.Decks {
				cardCount += len(deck.Cards)
			}
			printInfo(fmt.Sprintf("  Cards: %d", cardCount))
			return nil
		}

		if err := exporter.ExportToFile(data, output, opts); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		// Calculate total cards
		cardCount := len(data.Cards)
		for _, deck := range data.Decks {
			cardCount += len(deck.Cards)
		}

		if !quiet {
			printSuccess(fmt.Sprintf("Exported %d decks and %d cards to %s", len(data.Decks), cardCount, output))
		}

		if format == "json" {
			printJSON(map[string]interface{}{
				"decks": len(data.Decks),
				"cards": cardCount,
				"file":  output,
			})
		}

		return nil
	},
}

// importCmd imports data from a .mochi file
var importCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import decks and cards from .mochi file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath := args[0]
		deckID, _ := cmd.Flags().GetString("deck")
		templateID, _ := cmd.Flags().GetString("template")
		skipMedia, _ := cmd.Flags().GetBool("skip-media")

		// Validate file exists
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filepath)
		}

		// Validate .mochi file
		if err := importexport.ValidateMochiFile(filepath); err != nil {
			return err
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		importer := importexport.NewImporter(client)
		opts := importexport.ImportOptions{
			SkipMedia:  skipMedia,
			DeckID:     deckID,
			TemplateID: templateID,
			DryRun:     dryRun,
		}

		if !quiet {
			if dryRun {
				fmt.Println("Previewing import...")
			} else {
				fmt.Printf("Importing from %s...\n", filepath)
			}
		}

		result, err := importer.ImportFromFile(filepath, opts)
		if err != nil {
			return fmt.Errorf("import failed: %w", err)
		}

		if dryRun {
			printInfo("Would import:")
			printInfo(fmt.Sprintf("  Decks: %d", result.DecksCreated))
			printInfo(fmt.Sprintf("  Cards: %d", result.CardsCreated))
			printInfo(fmt.Sprintf("  Templates: %d", result.TemplatesCreated))
		} else {
			if !quiet {
				printSuccess(fmt.Sprintf("Imported %d decks, %d cards, %d templates",
					result.DecksCreated, result.CardsCreated, result.TemplatesCreated))
			}
		}

		if len(result.Errors) > 0 {
			printWarning(fmt.Sprintf("Encountered %d errors during import:", len(result.Errors)))
			for _, err := range result.Errors {
				printWarning(fmt.Sprintf("  - %s", err))
			}
		}

		if format == "json" {
			printJSON(map[string]interface{}{
				"decks":     result.DecksCreated,
				"cards":     result.CardsCreated,
				"templates": result.TemplatesCreated,
				"errors":    len(result.Errors),
			})
		}

		return nil
	},
}

// validateCmd validates a .mochi file
var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a .mochi file without importing",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath := args[0]

		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filepath)
		}

		if err := importexport.ValidateMochiFile(filepath); err != nil {
			printError(fmt.Sprintf("Validation failed: %s", err))
			return nil
		}

		if !quiet {
			printSuccess(fmt.Sprintf("File is valid: %s", filepath))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(importexportCmd)
	importexportCmd.AddCommand(exportCmd)
	importexportCmd.AddCommand(importCmd)
	importexportCmd.AddCommand(validateCmd)

	// Export flags
	exportCmd.Flags().StringP("output", "o", "", "Output file path (.mochi)")
	exportCmd.Flags().StringP("deck", "d", "", "Export specific deck only")
	exportCmd.Flags().StringP("format", "f", "json", "Export format: json or edn")
	exportCmd.Flags().Bool("include-reviews", false, "Include review history")

	// Import flags
	importCmd.Flags().StringP("deck", "d", "", "Import into specific deck")
	importCmd.Flags().StringP("template", "t", "", "Apply template to imported cards")
	importCmd.Flags().Bool("skip-media", false, "Skip media files")
}

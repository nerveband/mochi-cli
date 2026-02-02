package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// dueCmd represents the due command
var dueCmd = &cobra.Command{
	Use:   "due",
	Short: "Get cards due for review",
	Long:  `List cards that are due for review on a specific date.`,
}

// dueListCmd lists due cards
var dueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List due cards",
	RunE: func(cmd *cobra.Command, args []string) error {
		date, _ := cmd.Flags().GetString("date")
		deckID, _ := cmd.Flags().GetString("deck")

		// If no date specified, use today
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		resp, err := client.GetDueCards(date, deckID)
		if err != nil {
			return err
		}

		if idOnly || outputOnly == "id" {
			for _, card := range resp.Cards {
				fmt.Println(card.ID)
			}
			return nil
		}

		if outputOnly != "" {
			for _, card := range resp.Cards {
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
				"date":  date,
				"cards": resp.Cards,
				"count": len(resp.Cards),
			})
		case "compact":
			printCompactJSON(resp.Cards)
		case "table":
			headers := []string{"ID", "NAME", "DECK", "CONTENT"}
			rows := make([][]string, len(resp.Cards))
			for i, card := range resp.Cards {
				content := truncateString(card.Content, 40)
				rows[i] = []string{
					card.ID,
					card.Name,
					card.DeckID,
					content,
				}
			}
			printTable(headers, rows)
		default:
			if len(resp.Cards) == 0 {
				if !quiet {
					fmt.Printf("No cards due on %s\n", date)
				}
			} else {
				fmt.Printf("Cards due on %s:\n", date)
				for _, card := range resp.Cards {
					fmt.Printf("  %s: %s\n", card.ID, truncateString(card.Content, 50))
				}
			}
		}

		return nil
	},
}

// dueCountCmd counts due cards
var dueCountCmd = &cobra.Command{
	Use:   "count",
	Short: "Count due cards",
	RunE: func(cmd *cobra.Command, args []string) error {
		date, _ := cmd.Flags().GetString("date")
		deckID, _ := cmd.Flags().GetString("deck")

		// If no date specified, use today
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		resp, err := client.GetDueCards(date, deckID)
		if err != nil {
			return err
		}

		count := len(resp.Cards)

		switch format {
		case "json":
			printJSON(map[string]interface{}{
				"date":  date,
				"count": count,
			})
		case "compact":
			fmt.Println(count)
		default:
			if deckID != "" {
				fmt.Printf("%d cards due on %s in deck %s\n", count, date, deckID)
			} else {
				fmt.Printf("%d cards due on %s\n", count, date)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dueCmd)
	dueCmd.AddCommand(dueListCmd)
	dueCmd.AddCommand(dueCountCmd)

	// Flags
	dueListCmd.Flags().StringP("date", "d", "", "Date to check (YYYY-MM-DD, defaults to today)")
	dueListCmd.Flags().String("deck", "", "Limit to specific deck")

	dueCountCmd.Flags().StringP("date", "d", "", "Date to check (YYYY-MM-DD, defaults to today)")
	dueCountCmd.Flags().String("deck", "", "Limit to specific deck")
}

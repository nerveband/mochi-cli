package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// attachmentCmd represents the attachment command
var attachmentCmd = &cobra.Command{
	Use:   "attachment",
	Short: "Manage card attachments",
	Long:  `Add and delete attachments on cards.`,
}

// attachmentAddCmd adds an attachment
var attachmentAddCmd = &cobra.Command{
	Use:   "add <card-id> <filepath>",
	Short: "Add an attachment to a card",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID := args[0]
		filePath := args[1]

		if dryRun {
			printInfo(fmt.Sprintf("Dry run - would add attachment %s to card %s", filePath, cardID))
			return nil
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		if err := client.AddAttachmentFromFile(cardID, filePath); err != nil {
			return err
		}

		if !quiet {
			printSuccess(fmt.Sprintf("Attachment added to card %s", cardID))
		}

		return nil
	},
}

// attachmentDeleteCmd deletes an attachment
var attachmentDeleteCmd = &cobra.Command{
	Use:   "delete <card-id> <filename>",
	Short: "Delete an attachment from a card",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID := args[0]
		filename := args[1]

		if dryRun {
			printInfo(fmt.Sprintf("Dry run - would delete attachment %s from card %s", filename, cardID))
			return nil
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		if err := client.DeleteAttachment(cardID, filename); err != nil {
			return err
		}

		if !quiet {
			printSuccess(fmt.Sprintf("Attachment deleted from card %s", cardID))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(attachmentCmd)
	attachmentCmd.AddCommand(attachmentAddCmd)
	attachmentCmd.AddCommand(attachmentDeleteCmd)
}

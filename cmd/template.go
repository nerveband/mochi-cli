package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/nerveband/mochi-cli/internal/models"
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage templates",
	Long:  `List and get card templates.`,
}

// templateListCmd lists templates
var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		resp, err := client.ListTemplates("")
		if err != nil {
			return err
		}

		// Convert docs to templates
		docs, ok := resp.Docs.([]interface{})
		if !ok {
			return fmt.Errorf("unexpected response format")
		}

		var templates []models.Template
		for _, doc := range docs {
			data, _ := json.Marshal(doc)
			var template models.Template
			if err := json.Unmarshal(data, &template); err == nil {
				templates = append(templates, template)
			}
		}

		if idOnly || outputOnly == "id" {
			for _, template := range templates {
				fmt.Println(template.ID)
			}
			return nil
		}

		if outputOnly != "" {
			for _, template := range templates {
				val := extractField(template, outputOnly)
				if val != nil {
					fmt.Println(val)
				}
			}
			return nil
		}

		switch format {
		case "json":
			printJSON(map[string]interface{}{
				"templates": templates,
				"bookmark":  resp.Bookmark,
			})
		case "compact":
			printCompactJSON(templates)
		case "table":
			headers := []string{"ID", "NAME", "POS"}
			rows := make([][]string, len(templates))
			for i, template := range templates {
				rows[i] = []string{
					template.ID,
					template.Name,
					template.Pos,
				}
			}
			printTable(headers, rows)
		default:
			for _, template := range templates {
				fmt.Printf("%s: %s\n", template.ID, template.Name)
			}
		}

		return nil
	},
}

// templateGetCmd gets a specific template
var templateGetCmd = &cobra.Command{
	Use:   "get <template-id>",
	Short: "Get a template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateID := args[0]

		client, err := getClient()
		if err != nil {
			return err
		}

		template, err := client.GetTemplate(templateID)
		if err != nil {
			return err
		}

		if idOnly || outputOnly == "id" {
			fmt.Println(template.ID)
			return nil
		}

		if outputOnly != "" {
			val := extractField(template, outputOnly)
			if val != nil {
				fmt.Println(val)
			}
			return nil
		}

		switch format {
		case "json":
			printJSON(template)
		case "compact":
			printCompactJSON(template)
		case "markdown":
			fmt.Printf("# %s\n\n", template.Name)
			fmt.Printf("**ID:** %s\n\n", template.ID)
			fmt.Printf("**Position:** %s\n\n", template.Pos)
			fmt.Println("## Content")
			fmt.Println(template.Content)
			fmt.Println("\n## Fields")
			for id, field := range template.Fields {
				fmt.Printf("- **%s** (%s): %s\n", field.Name, id, field.Type)
			}
		default:
			fmt.Printf("ID: %s\n", template.ID)
			fmt.Printf("Name: %s\n", template.Name)
			fmt.Printf("Content:\n%s\n", template.Content)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateGetCmd)
}

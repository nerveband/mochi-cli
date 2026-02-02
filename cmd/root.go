package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	format     string
	quiet      bool
	jsonErrors bool
	outputOnly string
	idOnly     bool
	noHeaders  bool
	apiKey     string
	profile    string
	dryRun     bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "mochi",
	Short: "A powerful CLI for Mochi.cards",
	Long: `mochi-cli is a command-line interface for managing flashcards, decks,
and templates in Mochi. Built for LLMs and automation workflows.

For LLM and scripting workflows:
  - Use --quiet or -q to suppress status messages
  - Use --format json for structured output (default)
  - Use --output-only to extract specific fields
  - Pipe content to commands: echo "content" | mochi card create --stdin

Get started:
  mochi config add <name> <api-key>    # Set up your API key
  mochi deck list                      # List all your decks
  mochi card list                      # List cards
  mochi --help                         # Show all commands`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Check for updates after command execution
		notifyUpdateAvailable()
	},
}

// Execute runs the root command
func Execute(version string) {
	rootCmd.Version = version

	if err := rootCmd.Execute(); err != nil {
		if jsonErrors {
			fmt.Fprintf(os.Stderr, `{"error": "%s"}\n`, err.Error())
		} else if !quiet {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "json", "Output format: json, table, markdown, compact")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress status messages")
	rootCmd.PersistentFlags().BoolVar(&jsonErrors, "json-errors", false, "Output errors as JSON")
	rootCmd.PersistentFlags().StringVar(&outputOnly, "output-only", "", "Extract only specific field from output")
	rootCmd.PersistentFlags().BoolVar(&idOnly, "id-only", false, "Output only IDs (shorthand for --output-only=id)")
	rootCmd.PersistentFlags().BoolVar(&noHeaders, "no-headers", false, "Suppress table headers")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "api-key", "k", "", "API key (overrides profile and env var)")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "Profile to use (overrides active profile)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Preview changes without executing")
}

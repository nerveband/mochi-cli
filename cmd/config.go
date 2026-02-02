package cmd

import (
	"fmt"
	"strings"

	"github.com/nerveband/mochi-cli/internal/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration profiles",
	Long:  `Add, remove, and switch between Mochi API profiles.`,
}

// configAddCmd adds a new profile
var configAddCmd = &cobra.Command{
	Use:   "add <name> <api-key>",
	Short: "Add a new profile",
	Long:  `Add a new profile with an API key. The first profile added becomes active by default.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		apiKey := args[1]

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := cfg.AddProfile(name, apiKey); err != nil {
			return err
		}

		if !quiet {
			fmt.Printf("Profile '%s' added successfully\n", name)
			if cfg.ActiveProfile == name {
				fmt.Printf("Profile '%s' is now active\n", name)
			}
		}

		return nil
	},
}

// configRemoveCmd removes a profile
var configRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := cfg.RemoveProfile(name); err != nil {
			return err
		}

		if !quiet {
			fmt.Printf("Profile '%s' removed\n", name)
		}

		return nil
	},
}

// configUseCmd sets the active profile
var configUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Switch to a different profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := cfg.UseProfile(name); err != nil {
			return err
		}

		if !quiet {
			fmt.Printf("Now using profile '%s'\n", name)
		}

		return nil
	},
}

// configListCmd lists all profiles
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		names := cfg.ListProfiles()

		if len(names) == 0 {
			if !quiet {
				fmt.Println("No profiles configured")
				fmt.Println("Run 'mochi config add <name> <api-key>' to create one")
			}
			return nil
		}

		if format == "json" {
			output := map[string]interface{}{
				"active_profile": cfg.ActiveProfile,
				"profiles":       names,
			}
			printJSON(output)
		} else if format == "table" {
			if !noHeaders {
				fmt.Printf("%-15s %s\n", "NAME", "ACTIVE")
				fmt.Println(strings.Repeat("-", 25))
			}
			for _, name := range names {
				active := ""
				if name == cfg.ActiveProfile {
					active = "*"
				}
				fmt.Printf("%-15s %s\n", name, active)
			}
		} else {
			for _, name := range names {
				marker := ""
				if name == cfg.ActiveProfile {
					marker = " *"
				}
				fmt.Printf("%s%s\n", name, marker)
			}
		}

		return nil
	},
}

// configResetCmd resets all configuration
var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all configuration",
	Long:  `Remove all profiles and reset configuration to defaults.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !quiet {
			fmt.Print("Are you sure? This will remove all profiles. [y/N]: ")
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				fmt.Println("Cancelled")
				return nil
			}
		}

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := cfg.Reset(); err != nil {
			return err
		}

		if !quiet {
			fmt.Println("Configuration reset")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configRemoveCmd)
	configCmd.AddCommand(configUseCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configResetCmd)
}

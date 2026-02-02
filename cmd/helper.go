package cmd

import (
	"fmt"
	"os"

	"github.com/nerveband/mochi-cli/internal/api"
	"github.com/nerveband/mochi-cli/internal/config"
)

// getClient creates an API client using the active profile or provided credentials
func getClient() (*api.Client, error) {
	// Priority: CLI flag > Environment variable > Config profile

	// 1. Check CLI flag
	key := apiKey

	// 2. Check environment variable
	if key == "" {
		key = os.Getenv("MOCHI_API_KEY")
	}

	// 3. Check profile
	if key == "" {
		cfg, err := config.GetConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}

		// Use specified profile or active profile
		profileName := profile
		if profileName == "" {
			profileName = cfg.ActiveProfile
		}

		if profileName != "" {
			p, err := cfg.GetProfile(profileName)
			if err != nil {
				return nil, err
			}
			key = p.APIKey
		}
	}

	if key == "" {
		return nil, fmt.Errorf("no API key found. Set MOCHI_API_KEY environment variable, use --api-key flag, or run 'mochi config add <name> <api-key>'")
	}

	return api.NewClient(key), nil
}

// getActiveProfileName returns the name of the active profile
func getActiveProfileName() string {
	if profile != "" {
		return profile
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return ""
	}

	return cfg.ActiveProfile
}

// exitWithError exits with an error code
func exitWithError(code int, msg string) {
	if jsonErrors {
		fmt.Fprintf(os.Stderr, `{"error": "%s", "code": %d}`, msg, code)
	} else if !quiet {
		fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	}
	os.Exit(code)
}

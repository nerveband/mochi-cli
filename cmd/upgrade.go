package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"github.com/fatih/color"
	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
)

const (
	repoOwner = "nerveband"
	repoName  = "mochi-cli"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to the latest version",
	Long:  `Check for and install the latest version of mochi-cli.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		checkOnly, _ := cmd.Flags().GetBool("check")

		if !quiet {
			fmt.Println("Checking for updates...")
		}

		latest, err := getLatestRelease()
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		if latest.TagName == versionString {
			if !quiet {
				fmt.Printf("You're already on the latest version: %s\n", versionString)
			}
			return nil
		}

		if checkOnly {
			if !quiet {
				fmt.Printf("A new version is available: %s (current: %s)\n", latest.TagName, versionString)
				fmt.Println("Run 'mochi upgrade' to update")
			}
			return nil
		}

		if !quiet {
			fmt.Printf("New version available: %s (current: %s)\n", latest.TagName, versionString)
			fmt.Println("Downloading update...")
		}

		// Download and apply update
		if err := performUpdate(latest); err != nil {
			return fmt.Errorf("failed to update: %w", err)
		}

		if !quiet {
			color.Green("Successfully upgraded to %s!\n", latest.TagName)
		}

		return nil
	},
}

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func getLatestRelease() (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func performUpdate(release *GitHubRelease) error {
	// Determine correct asset for current platform
	assetName := fmt.Sprintf("mochi-cli_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)

	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no suitable binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	// Download binary
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Apply update
	return update.Apply(resp.Body, update.Options{})
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().Bool("check", false, "Only check for updates, don't install")
}

package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionString string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		if versionString == "" {
			versionString = "dev"
		}

		if format == "json" {
			printJSON(map[string]string{
				"version": versionString,
				"os":      runtime.GOOS,
				"arch":    runtime.GOARCH,
			})
		} else {
			fmt.Printf("mochi-cli version %s\n", versionString)
			fmt.Printf("OS: %s\n", runtime.GOOS)
			fmt.Printf("Arch: %s\n", runtime.GOARCH)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

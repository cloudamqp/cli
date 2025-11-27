package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Version information - these can be set at build time with -ldflags
	Version   = "dev"
	BuildDate = "unknown"
	GitCommit = "none"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show cloudamqp version information",
	Long:  `Display the version number, build date, and release information for cloudamqp CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		if Version == "dev" {
			fmt.Printf("cloudamqp version %s (development build)\n", Version)
		} else {
			fmt.Printf("cloudamqp version %s (%s)\n", Version, BuildDate)
			// Strip leading 'v' if present to avoid double 'v' in URL
			versionTag := strings.TrimPrefix(Version, "v")
			fmt.Printf("https://github.com/cloudamqp/cli/releases/tag/v%s\n", versionTag)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

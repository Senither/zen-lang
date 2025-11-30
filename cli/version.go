package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func init() {
	rootCommand.AddCommand(versionCommand)
}

var versionCommand = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Get the current version of the CLI.",
	Long:    "Displays the current version, build date, and commit hash of the CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Zen %s (built: %s)\n", Version, strings.ReplaceAll(BuildDate, "_", " "))
		fmt.Printf("Commit: %s\n", GitCommit)
	},
}

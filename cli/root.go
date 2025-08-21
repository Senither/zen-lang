package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:        "zen",
	Short:      "Run Zen code",
	Long:       "Runs the Zen interpreter",
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"file"},
	ValidArgs:  []string{"file"},
	Run: func(cmd *cobra.Command, args []string) {
		content, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Printf("Failed to %s\n", err)
		}

		fmt.Println(string(content))
	},
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err)
		os.Exit(1)
	}
}

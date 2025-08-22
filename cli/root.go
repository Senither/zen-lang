package cli

import (
	"fmt"
	"os"

	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/parser"
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

		lexer := lexer.New(string(content))
		parser := parser.New(lexer)

		program := parser.ParseProgram()

		for _, statement := range program.Statements {
			fmt.Printf("%q\n", statement)
		}
	},
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err)
		os.Exit(1)
	}
}

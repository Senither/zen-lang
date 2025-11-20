package cli

import (
	"fmt"
	"strings"

	"github.com/senither/zen-lang/cli/colors"
	"github.com/senither/zen-lang/compiler"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(codeCommand)
	codeCommand.Flags().BoolP("serialize", "s", false, "Compare the serialized/deserialized and the original bytecode")
}

var codeCommand = &cobra.Command{
	Use:   "code",
	Short: "Run code and get the bytecode instructions as output",
	Long:  "Runs the code provided and outputs the bytecode instructions it generates.",
	Run: func(cmd *cobra.Command, args []string) {
		serialize, _ := cmd.Flags().GetBool("serialize")

		table, _, constants := createCompilerParameters()

		createREPLRunner(args, []string{
			"Welcome to the Zen Bytecode generator, type your code below to see the Bytecode output.",
			"Type 'exit' to exit the Bytecode generator or press Ctrl+C.",
		}, func(input string, path any) {
			lexer := inputToLexer(input)
			program := lexerToProgram(lexer, path)

			if bytecode := programToBytecode(path, program, table, constants); bytecode != nil {
				if !serialize {
					fmt.Print(bytecode.String())
				} else {
					series := bytecode.Serialize()
					deserializedBytecode, err := compiler.Deserialize(series)
					if err != nil {
						fmt.Printf("Deserialization Error: %s\n", err)
						return
					}

					printBytecodeComparison(bytecode, deserializedBytecode)
				}
			}
		})
	},
}

func printBytecodeComparison(original, deserialized *compiler.Bytecode) {
	originalStr := strings.Split(original.String(), "\n")
	deserializedStr := strings.Split(deserialized.String(), "\n")

	fmt.Printf("%-35s%-35s\n", "ORIGINAL", "SERIALIZED & DESERIALIZED")
	for i := 0; len(originalStr) > i || len(deserializedStr) > i; i++ {
		originalLine := "~ empty ~"
		deserializedLine := "~ empty ~"
		color := colors.Gray

		if i < len(originalStr) {
			originalLine = originalStr[i]
		}

		if i < len(deserializedStr) {
			deserializedLine = deserializedStr[i]
		}

		if originalLine != deserializedLine {
			color = colors.BgRed
		}

		fmt.Printf("%s%-35s%-35s%s\n", color, originalLine, deserializedLine, colors.Reset)
	}
}

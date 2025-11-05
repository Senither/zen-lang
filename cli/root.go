package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
	"github.com/senither/zen-lang/vm"
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
		path, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		content, err := loadFileContents(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		bytecode, err := compiler.Deserialize(content)
		table, _, constants := createCompilerParameters()
		if err != nil {
			lexer := inputToLexer(string(content))
			program := lexerToProgram(lexer, path)
			bytecode = programToBytecode(program, table, constants)
		}

		vm := vm.New(bytecode)
		if err := vm.Run(); err != nil {
			fmt.Println(err)
			return
		}

		stackTop := vm.LastPoppedStackElem()
		if stackTop != nil {
			fmt.Printf("%s\n", stackTop.Inspect())
		}
	},
}

func createREPLRunner(
	args []string,
	details []string,
	callback func(input string, path interface{}),
) {
	if len(args) > 0 {
		content, err := loadFileContents(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		path, _ := filepath.Abs(args[0])
		callback(string(content), path)

		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	for _, detail := range details {
		fmt.Println(detail)
	}

	for {
		fmt.Printf(">>> ")

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}

		callback(line, nil)
	}
}

func inputToLexer(input string) *lexer.Lexer {
	return lexer.New(input)
}

func lexerToProgram(lex *lexer.Lexer, filePath interface{}) *ast.Program {
	if lex == nil {
		return nil
	}

	parser := parser.New(lex, filePath)

	program := parser.ParseProgram()
	if len(parser.Errors()) > 0 {
		for _, err := range parser.Errors() {
			fmt.Println("Parse error:", err.String())
		}

		return nil
	}

	return program
}

func createCompilerParameters() (*compiler.SymbolTable, []objects.Object, []objects.Object) {
	constants := []objects.Object{}
	globals := make([]objects.Object, vm.GLOBALS_SIZE)
	table := compiler.NewSymbolTable()

	compiler.WriteBuiltinSymbols(table)

	return table, globals, constants
}

func programToBytecode(
	prog *ast.Program,
	table *compiler.SymbolTable,
	constants []objects.Object,
) *compiler.Bytecode {
	if prog == nil {
		return nil
	}

	compile := compiler.NewWithState(table, constants)
	if err := compile.Compile(prog); err != nil {
		fmt.Println("Compilation error:", err)
		return nil
	}

	return compile.Bytecode()
}

func loadFileContents(file string) ([]byte, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", file, err)
	}

	return content, nil
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err)
		os.Exit(1)
	}
}

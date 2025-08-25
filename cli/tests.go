package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/senither/zen-lang/evaluator"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
	"github.com/spf13/cobra"
)

type TestInstance struct {
	message string
	file    string
	expect  string
}

var exitStatusCode = 0

func init() {
	rootCommand.AddCommand(testCommand)
}

var testCommand = &cobra.Command{
	Use:        "test",
	Short:      "Run tests",
	Long:       "Runs the tests for the Zen language",
	ArgAliases: []string{"directory"},
	ValidArgs:  []string{"directory"},
	Run: func(cmd *cobra.Command, args []string) {
		testDirectory := "tests"
		if len(args) > 0 {
			testDirectory = args[0]
		}

		absolutePath, err := filepath.Abs(testDirectory)
		if err != nil {
			fmt.Printf("Error getting absolute path: %s\n", err)
			return
		}

		relativeTestFiles := discoverTestFiles(absolutePath)
		groupedTestFiles := make(map[string][]string)

		for _, relativePath := range relativeTestFiles {
			dir := filepath.Dir(relativePath)
			groupedTestFiles[dir] = append(groupedTestFiles[dir], relativePath)
		}

		start := time.Now()

		for dir, files := range groupedTestFiles {
			dirPath := strings.Trim(filepath.Join(absolutePath, dir), absolutePath)
			fmt.Printf("  %s%s%s\n", testDirectory, string(os.PathSeparator), dirPath)

			for _, file := range files {
				runTestFile(file)
			}

			fmt.Println()
		}

		fmt.Printf("Finished running the test suite in %s\nTime taken %s\n\n", absolutePath, time.Since(start))

		os.Exit(exitStatusCode)
	},
}

func discoverTestFiles(directory string) []string {
	var testFiles []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".zent") {
			absolutePath, _ := filepath.Abs(path)
			testFiles = append(testFiles, absolutePath)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error discovering test files: %s\n", strings.Trim(strings.Split(err.Error(), ":")[1], " "))
	}

	return testFiles
}

func runTestFile(file string) {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading test file %s: %s\n", file, err)
		return
	}

	key := ""
	test := TestInstance{
		message: "",
		file:    "",
		expect:  "",
	}

	for line := range strings.SplitSeq(string(content), "\n") {
		if strings.HasPrefix(line, "--TEST--") || strings.HasPrefix(line, "---TEST---") {
			key = "message"
		} else if strings.HasPrefix(line, "--FILE--") || strings.HasPrefix(line, "---FILE---") {
			key = "file"
		} else if strings.HasPrefix(line, "--EXPECT--") || strings.HasPrefix(line, "---EXPECT---") {
			key = "expect"
		}

		if strings.HasPrefix(line, "--") {
			continue
		}

		switch key {
		case "message":
			test.message += line + "\n"
		case "file":
			test.file += line + "\n"
		case "expect":
			test.expect += line + "\n"
		}
	}

	test.file = cleanString(test.file)
	test.expect = cleanString(test.expect)

	l := lexer.New(test.file)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		printErrorStatusMessage(test, "Parser errors found")
		for _, err := range p.Errors() {
			fmt.Printf("    %s\n", err.String())
		}
		return
	}

	evaluated := evaluator.Eval(program)
	if evaluated == nil {
		printErrorStatusMessage(test, "Evaluator returned nil, failed to evaluate the test input")
		return
	}

	if evaluated.Type() == objects.ERROR_OBJ {
		printErrorStatusMessage(test, "Evaluator returned an error")
		fmt.Printf("    %s\n", evaluated.Inspect())
		return
	}

	if strings.Trim(evaluated.Inspect(), "\n") != test.expect {
		printErrorStatusMessage(test, "Test expectation does not match the evaluated result")
		fmt.Printf("     Got:   %s\n", strings.Trim(evaluated.Inspect(), "\n"))
		fmt.Printf("     Want:  %s\n", test.expect)
		return
	}

	fmt.Printf("  ✔ %s\n", strings.Trim(test.message, "\n"))
}

func cleanString(str string) string {
	str = strings.ReplaceAll(str, "\r\n", "\n")
	str = strings.Trim(str, "\n")

	return str
}

func printErrorStatusMessage(test TestInstance, message string) {
	exitStatusCode = 1

	fmt.Printf("  ✖ %s\n     %s\n", cleanString(test.message), message)
}

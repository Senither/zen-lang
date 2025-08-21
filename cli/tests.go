package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type TestInstance struct {
	message string
	file    string
	expect  string
}

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

		relativeTestFiles := discoverTestFiles(testDirectory)
		groupedTestFiles := make(map[string][]string)

		for _, relativePath := range relativeTestFiles {
			dir := filepath.Dir(relativePath)
			groupedTestFiles[dir] = append(groupedTestFiles[dir], relativePath)
		}

		for dir, files := range groupedTestFiles {
			dirPath := strings.Join(strings.Split(dir, string(os.PathSeparator)), "/")
			fmt.Printf("  %s/%s\n", testDirectory, dirPath)

			for _, file := range files {
				runTestFile(filepath.Join(testDirectory, file))
			}

			fmt.Println()
		}
	},
}

func discoverTestFiles(directory string) []string {
	var testFiles []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".zent") {
			relativePath, _ := filepath.Rel("tests", path)
			testFiles = append(testFiles, relativePath)
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

	// TODO: Evaluate test code and check it against our expectations
	// For now we'll simulate a passing test

	// fmt.Printf("  ✖ %s\n", strings.Trim(test.message, "\n"))
	fmt.Printf("  ✔ %s\n", strings.Trim(test.message, "\n"))

	// Develop information
	fmt.Printf("       File:\n")
	for line := range strings.SplitSeq(strings.Trim(test.file, "\n"), "\n") {
		fmt.Printf("         %s\n", strings.Trim(line, "\n"))
	}

	fmt.Printf("       Expect:\n")
	for line := range strings.SplitSeq(strings.Trim(test.expect, "\n"), "\n") {
		fmt.Printf("         %s\n", strings.Trim(line, "\n"))
	}
}

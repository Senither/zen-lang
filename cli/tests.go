package cli

import (
	"fmt"
	"path/filepath"

	"github.com/senither/zen-lang/cli/tester"
	"github.com/spf13/cobra"
)

func init() {
	testCommand.Flags().BoolP("verbose", "v", false, "Adds more verbose output to the test runner")
	testCommand.Flags().BoolP("compact", "c", false, "Replace default result output with Compact format")
	testCommand.Flags().StringP("engine", "e", "all", "Specifies which engine to run the tests against. Options are: all, eval, vm")

	rootCommand.AddCommand(testCommand)
}

var testCommand = &cobra.Command{
	Use:        "test",
	Short:      "Run tests",
	Long:       "Runs the tests for the Zen language",
	ArgAliases: []string{"directory"},
	ValidArgs:  []string{"directory"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()

		testDirectory := "tests"
		if len(args) > 0 {
			testDirectory = args[0]
		}

		absolutePath, err := filepath.Abs(testDirectory)
		if err != nil {
			fmt.Printf("Error getting absolute path: %s\n", err)
			return
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		compact, _ := cmd.Flags().GetBool("compact")

		runner := tester.NewTestRunner(
			absolutePath, testDirectory,
			tester.RunnerOptions{
				Engine:  getTestRunnerEngine(cmd),
				Verbose: verbose,
				Compact: compact,
			},
		)

		err = runner.RunTests()
		if err != nil {
			fmt.Printf("Error running tests: %s\n", err)
		}
	},
}

func getTestRunnerEngine(cmd *cobra.Command) tester.EngineType {
	engine, _ := cmd.Flags().GetString("engine")
	switch engine {
	case "eval":
		return tester.EvaluatorEngine
	case "vm":
		return tester.VirtualMachineEngine

	default:
		return tester.AllEngines
	}
}

// type TestInstance struct {
// 	message string
// 	file    string
// 	expect  string
// 	errors  string
// }

// var (
// 	PASSED          = BgGreen + Gray + " PASS " + Reset
// 	FAILED          = BgRed + Gray + " FAIL " + Reset
// 	messages        = []string{}
// 	collectedErrors = []string{}
// 	exitStatusCode  = 0
// )

// var (
// 	totalTimeTakenForParsing    = time.Duration(0)
// 	totalTimeTakenForEvaluation = time.Duration(0)
// )

// func init() {
// 	rootCommand.AddCommand(testCommand)
// }

// var testCommand = &cobra.Command{
// 	Use:        "test",
// 	Short:      "Run tests",
// 	Long:       "Runs the tests for the Zen language",
// 	ArgAliases: []string{"directory"},
// 	ValidArgs:  []string{"directory"},
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println()

// 		testDirectory := "tests"
// 		if len(args) > 0 {
// 			testDirectory = args[0]
// 		}

// 		absolutePath, err := filepath.Abs(testDirectory)
// 		if err != nil {
// 			fmt.Printf("Error getting absolute path: %s\n", err)
// 			return
// 		}

// 		relativeTestFiles := discoverTestFiles(absolutePath)
// 		groupedTestFiles := make(map[string][]string)

// 		for _, relativePath := range relativeTestFiles {
// 			dir := filepath.Dir(relativePath)
// 			groupedTestFiles[dir] = append(groupedTestFiles[dir], relativePath)
// 		}

// 		start := time.Now()
// 		for dir, files := range groupedTestFiles {
// 			messages = []string{}

// 			fullPath := fmt.Sprintf("%s%s%s", testDirectory, string(os.PathSeparator), dir[len(absolutePath)+1:])

// 			errorsCount := len(collectedErrors)
// 			for _, file := range files {
// 				runTestFile(fullPath, file)
// 			}

// 			if errorsCount == len(collectedErrors) {
// 				fmt.Printf("  %s %s\n", PASSED, fullPath)
// 			} else {
// 				fmt.Printf("  %s %s\n", FAILED, fullPath)
// 			}

// 			fmt.Println(strings.Join(messages, ""))
// 		}

// 		if len(collectedErrors) > 0 {
// 			fmt.Println("Test suite failed with the following errors:")

// 			for _, err := range collectedErrors {
// 				parts := strings.Split(err, "\n")

// 				fmt.Println()
// 				fmt.Printf(" %s- %s%s\n%s\n", Red, parts[0], Reset, strings.Join(parts[1:], "\n"))
// 			}

// 			fmt.Println()
// 		}

// 		timeTaken := time.Since(start)

// 		fmt.Printf("Finished running the test suite in %s\n\n", absolutePath)
// 		fmt.Printf("   Reading files: %s\n", timeTaken-totalTimeTakenForParsing-totalTimeTakenForEvaluation)
// 		fmt.Printf("  Lexer + Parser: %s\n", totalTimeTakenForParsing)
// 		fmt.Printf("      Evaluation: %s\n", totalTimeTakenForEvaluation)
// 		fmt.Printf("           Total: %s\n", timeTaken)
// 		fmt.Printf("\n")

// 		os.Exit(exitStatusCode)
// 	},
// }

// func discoverTestFiles(directory string) []string {
// 	var testFiles []string

// 	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if !info.IsDir() && strings.HasSuffix(info.Name(), ".zent") {
// 			absolutePath, _ := filepath.Abs(path)
// 			testFiles = append(testFiles, absolutePath)
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		fmt.Printf("Error discovering test files: %s\n", strings.Trim(strings.Split(err.Error(), ":")[1], " "))
// 	}

// 	return testFiles
// }

// func runTestFile(fullPath, file string) {
// 	evaluator.Stdout.Clear()

// 	content, err := os.ReadFile(file)
// 	if err != nil {
// 		messages = append(messages, fmt.Sprintf("Error reading test file %s: %s\n", file, err))
// 		return
// 	}

// 	key := ""
// 	test := TestInstance{
// 		message: "",
// 		file:    "",
// 		expect:  "",
// 		errors:  "",
// 	}

// 	for line := range strings.SplitSeq(string(content), "\n") {
// 		if strings.HasPrefix(line, "--TEST--") || strings.HasPrefix(line, "---TEST---") {
// 			key = "message"
// 		} else if strings.HasPrefix(line, "--FILE--") || strings.HasPrefix(line, "---FILE---") {
// 			key = "file"
// 		} else if strings.HasPrefix(line, "--EXPECT--") || strings.HasPrefix(line, "---EXPECT---") {
// 			key = "expect"
// 		} else if strings.HasPrefix(line, "--ERROR--") || strings.HasPrefix(line, "---ERROR---") {
// 			key = "errors"
// 		}

// 		if strings.HasPrefix(line, "--") {
// 			continue
// 		}

// 		switch key {
// 		case "message":
// 			test.message += line + "\n"
// 		case "file":
// 			test.file += line + "\n"
// 		case "expect":
// 			test.expect += line + "\n"
// 		case "errors":
// 			test.errors += line + "\n"
// 		}
// 	}

// 	test.file = cleanString(test.file)
// 	test.expect = cleanString(test.expect)
// 	test.errors = cleanString(test.errors)

// 	startParserTimer := time.Now()
// 	l := lexer.New(test.file)
// 	p := parser.New(l, file)

// 	program := p.ParseProgram()
// 	totalTimeTakenForParsing += time.Since(startParserTimer)

// 	if len(p.Errors()) > 0 {
// 		msg := []string{"Parser errors found"}
// 		for _, err := range p.Errors() {
// 			msg = append(msg, fmt.Sprintf("     %s", err.String()))
// 		}
// 		printErrorStatusMessage(test, fullPath, strings.Join(msg, "\n"))
// 		return
// 	}

// 	runTestWithEvaluator(test, fullPath, file, program)
// }

// func runTestWithEvaluator(test TestInstance, fullPath, file string, program *ast.Program) {
// 	startEvaluatorTimer := time.Now()
// 	evaluated := evaluator.Stdout.Mute(func() objects.Object {
// 		env := objects.NewEnvironment(file)
// 		return evaluator.Eval(program, env)
// 	})
// 	totalTimeTakenForEvaluation += time.Since(startEvaluatorTimer)

// 	if evaluated == nil {
// 		printErrorStatusMessage(test, fullPath, "Evaluator returned nil, failed to evaluate the test input")
// 		return
// 	}

// 	if objects.IsError(evaluated) {
// 		stripFileLocationsAndPrintError(test, fullPath, evaluated)
// 	} else if evaluated.Type() != objects.NULL_OBJ {
// 		compareEvaluatedWithExpected(test, fullPath, evaluated)
// 	} else {
// 		compareStandardOutputWithExpected(test, fullPath)
// 	}
// }

// func compareEvaluatedWithExpected(test TestInstance, fullPath string, evaluated objects.Object) {
// 	if strings.Trim(evaluated.Inspect(), "\n") != test.expect {
// 		printErrorStatusMessage(
// 			test,
// 			fullPath,
// 			fmt.Sprintf(
// 				"%s\n     -----------------[ RESULT ]-----------------\n%s\n     ----------------[ EXPECTED ]-----------------\n%s",
// 				"Test expectation does not match the evaluated result",
// 				strings.Trim(evaluated.Inspect(), "\n"),
// 				test.expect,
// 			),
// 		)
// 		return
// 	}

// 	printSuccessStatusMessage(test)
// }

// func stripFileLocationsAndPrintError(test TestInstance, fullPath string, evaluated objects.Object) {
// 	err := evaluated.Inspect()
// 	lines := strings.Split(err, "\n")

// 	for i, line := range lines {
// 		if strings.Contains(line, "at ") && strings.Contains(line, ".zent:") {
// 			fileInfo := strings.Split(line, ".zent:")[1]

// 			lines[i] = fmt.Sprintf("    at <unknown>:%s", fileInfo)
// 		}
// 	}

// 	err = strings.Trim(strings.Join(lines, "\n"), "\n")
// 	if err != test.errors {
// 		var message = "Test expectation does not match the evaluated result"
// 		if len(test.errors) == 0 {
// 			message = "No error expectation were provided, despite the result being *objects.Error"
// 		}

// 		printErrorStatusMessage(
// 			test,
// 			fullPath,
// 			fmt.Sprintf(
// 				"%s\n     -----------------[ RESULT ]-----------------\n%s\n     ----------------[ EXPECTED ]-----------------\n%s",
// 				message, err, test.errors,
// 			),
// 		)
// 		return
// 	}

// 	printSuccessStatusMessage(test)
// }

// func compareStandardOutputWithExpected(test TestInstance, fullPath string) {
// 	messages := evaluator.Stdout.ReadAll()
// 	if len(messages) == 0 {
// 		printErrorStatusMessage(test, fullPath, "No output captured from standard output")
// 		return
// 	}

// 	out := strings.Trim(strings.Join(messages, ""), "\n")
// 	if out != test.expect {
// 		printErrorStatusMessage(
// 			test,
// 			fullPath,
// 			fmt.Sprintf(
// 				"%s\n     -----------------[ RESULT ]-----------------\n%s\n     ----------------[ EXPECTED ]-----------------\n%s",
// 				"Test expectation does not match the standard output",
// 				out,
// 				test.expect,
// 			),
// 		)
// 		return
// 	}

// 	printSuccessStatusMessage(test)
// }

// func cleanString(str string) string {
// 	str = strings.ReplaceAll(str, "\r\n", "\n")
// 	str = strings.Trim(str, "\n")

// 	return str
// }

// func printSuccessStatusMessage(test TestInstance) {
// 	messages = append(messages, fmt.Sprintf("  %s✔%s %s\n", Green, Reset, cleanString(test.message)))
// }

// func printErrorStatusMessage(test TestInstance, fullPath, message string) {
// 	exitStatusCode = 1

// 	errorMessage := fmt.Sprintf("%s\n     %s", cleanString(test.message), message)
// 	collectedErrors = append(collectedErrors, fmt.Sprintf("%s: %s", fullPath, errorMessage))

// 	messages = append(messages, fmt.Sprintf("  %s✖%s %s\n", Red, Reset, errorMessage))
// }

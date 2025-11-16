package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/senither/zen-lang/cli/tester"
	"github.com/spf13/cobra"
)

func init() {
	testCommand.Flags().BoolP("verbose", "v", false, "Adds more verbose output to the test runner")
	testCommand.Flags().BoolP("compact", "c", false, "Replace default result output with Compact format")
	testCommand.Flags().StringP("filter", "f", "", "Filter tests to run by name (substring match)")
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

		filter, _ := cmd.Flags().GetString("filter")
		verbose, _ := cmd.Flags().GetBool("verbose")
		compact, _ := cmd.Flags().GetBool("compact")

		runner := tester.NewTestRunner(
			absolutePath, testDirectory,
			tester.RunnerOptions{
				Engine:  getTestRunnerEngine(cmd),
				Filter:  filter,
				Verbose: verbose,
				Compact: compact,
			},
		)

		if err = runner.RunTests(); err != nil {
			fmt.Printf("Error running tests: %s\n", err)
		}
	},
}

func getTestRunnerEngine(cmd *cobra.Command) tester.EngineType {
	engine, _ := cmd.Flags().GetString("engine")
	switch strings.ToLower(engine) {
	case "evaluator", "eval", "e", "interpreter", "int", "inter":
		return tester.EvaluatorEngine
	case "virtual", "machine", "vm", "v", "virt", "mach":
		return tester.VirtualMachineEngine
	case "all", "a", "both", "b":
		return tester.AllEngines

	default:
		fmt.Printf("Unknown engine value: %s\n", engine)
		os.Exit(0)
	}

	return tester.AllEngines
}

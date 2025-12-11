package tester

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/senither/zen-lang/cli/colors"
	"github.com/senither/zen-lang/evaluator"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/objects/process"
	"github.com/senither/zen-lang/objects/timer"
	"github.com/senither/zen-lang/vm"
)

var closureRegex = regexp.MustCompile(`Closure\[0x[a-fA-F0-9]+\]`)

func (tr *TestRunner) printSuccessStatusMessage(test *Test, engineType EngineType) {
	tr.passedTests++

	if tr.options.Compact {
		fmt.Print(".")
		return
	}

	timings := ""
	if tr.options.Verbose {
		switch engineType {
		case EvaluatorEngine:
			timings = fmt.Sprintf(" %sT:%s%s",
				colors.Gray, test.metadata[EvaluatorExecutionTiming], colors.Reset,
			)
		case VirtualMachineEngine, _VirtualMachineEngineUnprocessed, _VirtualMachineEngineSerialized:
			vmTiming := test.metadata[VMExecutionTiming]
			if vmTiming == nil {
				vmTiming = "n/a"
			}

			timings = fmt.Sprintf(" %sC:%s VM:%s%s",
				colors.Gray, test.metadata[CompilationTiming], vmTiming, colors.Reset,
			)
		}
	}

	message := fmt.Sprintf("  %s✔%s %s %s[%s%s%s]%s",
		colors.Green, colors.Reset, tr.normalizeLineEndings(test.message),
		colors.Gray, engineType.GetTag(), timings, colors.Gray, colors.Reset,
	)

	messages = append(messages, message+"\n")
}

func (tr *TestRunner) printErrorDoesNotMatchExpectation(
	test *Test,
	fullPath,
	message, result, expected string,
	engineType EngineType,
) {
	tr.printErrorStatusMessage(
		test,
		fullPath,
		fmt.Sprintf(
			"%s\n     -----------------[ RESULT ]-----------------\n%s\n     ----------------[ EXPECTED ]-----------------\n%s",
			message, result, expected,
		),
		engineType,
	)
}

func (tr *TestRunner) printErrorStatusMessage(test *Test, fullPath, message string, engineType EngineType) {
	exitStatusCode = 1

	errorMessage := fmt.Sprintf("%s %s[%s%s]%s\n     %s",
		tr.normalizeLineEndings(test.message),
		colors.Gray, engineType.GetTag(), colors.Gray, colors.Reset,
		message,
	)

	if tr.options.Compact {
		fmt.Print(colors.Red + "x" + colors.Reset)
	} else {
		messages = append(messages, fmt.Sprintf("  %s✖%s %s\n", colors.Red, colors.Reset, errorMessage))
	}

	collectedErrors = append(collectedErrors, fmt.Sprintf("%s: %s", fullPath, errorMessage))
}

func (tr *TestRunner) normalizeLineEndings(str string) string {
	str = strings.ReplaceAll(str, "\r\n", "\n")
	str = strings.Trim(str, "\n")

	return str
}

func (tr *TestRunner) normalizeFileLocations(input string) string {
	lines := strings.Split(input, "\n")

	for i, line := range lines {
		if strings.Contains(line, "at ") && strings.Contains(line, ".zent:") {
			fileInfo := strings.Split(line, ".zent:")[1]

			lines[i] = fmt.Sprintf("    at <unknown>:%s", fileInfo)
		}
	}

	return strings.Trim(strings.Join(lines, "\n"), "\n")
}

func (tr *TestRunner) normalizeClosurePointers(content string) string {
	return closureRegex.ReplaceAllString(content, "Closure[<pointer>]")
}

func (tr *TestRunner) applyTestEnvVariables(test *Test) {
	evaluator.Stdout.Clear()
	vm.Stdout.Clear()

	for key, value := range test.envs {
		switch key {
		case "time":
			time, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				fmt.Printf("Invalid time env variable: %s\n", value)
				continue
			}

			timer.Freeze(time)
		case "timezone":
			err := timer.SetTimezone(value)
			if err != nil {
				fmt.Printf("Invalid timezone env variable: %s\n", value)
				continue
			}
		case "process":
			process.Fake()
		case "process.args":
			process.FakeArgs(strings.Split(value, ","))
		case "process.envs":
			process.Fake()

			envPairs := strings.Split(value, ";")
			for _, pair := range envPairs {
				kv := strings.SplitN(pair, "=", 2)
				if len(kv) != 2 {
					fmt.Printf("Invalid process.envs variable: %s\n", pair)
					continue
				}
				process.FakeEnv(kv[0], kv[1])
			}
		}
	}
}

func (tr *TestRunner) clearTestEnvVariables() {
	objects.RestoreObjectsState()
}

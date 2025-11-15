package tester

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/senither/zen-lang/cli/colors"
)

var closureRegex = regexp.MustCompile(`Closure\[0x[a-fA-F0-9]+\]`)

func (tr *TestRunner) printSuccessStatusMessage(test *Test, engineType EngineType) {
	tr.passedTests++

	if tr.options.Compact {
		fmt.Print(".")
		return
	}

	messages = append(messages, fmt.Sprintf("  %s✔%s %s %s[%s%s]%s\n",
		colors.Green, colors.Reset, tr.normalizeLineEndings(test.message),
		colors.Gray, engineType.GetTag(), colors.Gray, colors.Reset,
	))
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

func (tr *TestRunner) normalizeFileLocations(err string) string {
	lines := strings.Split(err, "\n")

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

package tester

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/senither/zen-lang/cli/colors"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/parser"
)

type TestRunner struct {
	directory string
	path      string
	options   RunnerOptions
	timings   map[RunnerTimings]time.Duration

	passedTests int
}

type RunnerOptions struct {
	Engine  EngineType
	Verbose bool
}

type EngineType int

const (
	AllEngines EngineType = iota
	EvaluatorEngine
	VirtualMachineEngine
)

func (et EngineType) GetTag() string {
	switch et {
	case AllEngines:
		return colors.Green + "All"
	case EvaluatorEngine:
		return colors.Cyan + "Eval"
	case VirtualMachineEngine:
		return colors.Magenta + "VM"

	default:
		return "unknown"
	}
}

type RunnerTimings string

const (
	FileDiscoveryTiming      RunnerTimings = "File Discovery"
	ReadingFilesTiming       RunnerTimings = "Reading Files"
	LexingAndParsingTiming   RunnerTimings = "Lexing + Parsing"
	CompilationTiming        RunnerTimings = "Compilation"
	EvaluatorExecutionTiming RunnerTimings = "Evaluator Execution"
	VMExecutionTiming        RunnerTimings = "VM Execution"
)

type Test struct {
	message         string
	file            string
	expect          string
	errors          string
	supportedEngine EngineType
}

var (
	PASSED          = colors.BgGreen + colors.Gray + " PASS " + colors.Reset
	FAILED          = colors.BgRed + colors.Gray + " FAIL " + colors.Reset
	messages        = []string{}
	collectedErrors = []string{}
	exitStatusCode  = 0
)

func NewTestRunner(directory, path string, options RunnerOptions) *TestRunner {
	return &TestRunner{
		directory: directory,
		path:      path,
		options:   options,
		timings:   make(map[RunnerTimings]time.Duration),
	}
}

func (tr *TestRunner) setTiming(timingType RunnerTimings, duration time.Duration) {
	tr.timings[timingType] = duration
}

func (tr *TestRunner) getTiming(timingType RunnerTimings) time.Duration {
	timing := tr.timings[timingType]
	if timing <= 0 {
		return 0
	}
	return timing
}

func (tr *TestRunner) RunTests() error {
	groupedFiles, err := tr.discoverAndGroupTestFiles()
	if err != nil {
		return err
	}

	var groups []string
	for dir := range groupedFiles {
		groups = append(groups, dir)
	}

	sort.Strings(groups)

	for _, dir := range groups {
		messages = []string{}

		fullPath := fmt.Sprintf("%s%s%s", tr.path, string(os.PathSeparator), dir[len(tr.path)+1:])
		errorsCount := len(collectedErrors)

		for _, file := range groupedFiles[dir] {
			tr.runTestFile(fullPath, file)
		}

		if errorsCount == len(collectedErrors) {
			fmt.Printf("  %s %s\n", PASSED, fullPath)
		} else {
			fmt.Printf("  %s %s\n", FAILED, fullPath)
		}

		fmt.Println(strings.Join(messages, ""))
	}

	if len(collectedErrors) > 0 {
		fmt.Println("Test suite failed with the following errors:")

		for _, err := range collectedErrors {
			parts := strings.Split(err, "\n")

			fmt.Println()
			fmt.Printf(" %s- %s%s\n%s\n", colors.Red, parts[0], colors.Reset, strings.Join(parts[1:], "\n"))
		}

		fmt.Println()
	}

	fmt.Printf("  Finished running the test suite in %s\n", tr.directory)
	fmt.Printf("  Summary: %s\n\n", tr.getStatusSummary())

	fmt.Printf("     Tests discovery: %s\n", tr.getTiming(FileDiscoveryTiming))
	fmt.Printf("       Reading files: %s\n", tr.getTiming(ReadingFilesTiming))
	fmt.Printf("      Lexer + Parser: %s\n", tr.getTiming(LexingAndParsingTiming))

	if tr.options.Engine == AllEngines || tr.options.Engine == EvaluatorEngine {
		fmt.Printf(" -----------------------------------\n")
		fmt.Printf("          Evaluation: %s\n", tr.getTiming(EvaluatorExecutionTiming))
	}

	if tr.options.Engine == AllEngines || tr.options.Engine == VirtualMachineEngine {
		fmt.Printf(" -----------------------------------\n")
		fmt.Printf("  Compile + Optimize: %s\n", tr.getTiming(CompilationTiming))
		fmt.Printf("          VM Runtime: %s\n", tr.getTiming(VMExecutionTiming))
	}

	var totalTimeTake time.Duration
	for _, timing := range tr.timings {
		totalTimeTake += timing
	}

	fmt.Printf(" -----------------------------------\n")
	fmt.Printf("               Total: %s\n", totalTimeTake)
	fmt.Printf("\n")

	os.Exit(exitStatusCode)
	return nil
}

func (tr *TestRunner) getStatusSummary() string {
	var parts []string

	if len(collectedErrors) > 0 {
		parts = append(parts, fmt.Sprintf("%s%d failed%s", colors.Red, len(collectedErrors), colors.Reset))
	}

	parts = append(parts, fmt.Sprintf("%s%d passed%s", colors.Green, tr.passedTests, colors.Reset))

	return strings.Join(parts, ", ")
}

func (tr *TestRunner) discoverAndGroupTestFiles() (map[string][]string, error) {
	start := time.Now()
	var relativeTestFiles []string

	err := filepath.Walk(tr.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".zent") {
			absolutePath, _ := filepath.Abs(path)
			relativeTestFiles = append(relativeTestFiles, absolutePath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	groupedTestFiles := make(map[string][]string)

	for _, relativePath := range relativeTestFiles {
		dir := filepath.Dir(relativePath)
		groupedTestFiles[dir] = append(groupedTestFiles[dir], relativePath)
	}

	tr.setTiming(FileDiscoveryTiming, time.Since(start))

	return groupedTestFiles, nil
}

func (tr *TestRunner) runTestFile(fullPath, file string) {
	test, err := tr.parseTestFile(file)
	if err != nil {
		collectedErrors = append(collectedErrors, fmt.Sprintf("Error parsing test file %s: %s", file, err.Error()))
		return
	}

	startLexingAndParsing := time.Now()
	l := lexer.New(test.file)
	p := parser.New(l, file)

	program := p.ParseProgram()
	tr.setTiming(LexingAndParsingTiming, tr.getTiming(LexingAndParsingTiming)+time.Since(startLexingAndParsing))

	if len(p.Errors()) > 0 {
		msg := []string{"Parser errors found"}
		for _, err := range p.Errors() {
			msg = append(msg, fmt.Sprintf("     %s", err.String()))
		}

		tr.printErrorStatusMessage(test, fullPath, strings.Join(msg, "\n"), AllEngines)
		return
	}

	if tr.ShouldRunTest(test, EvaluatorEngine) {
		tr.runEvaluatorTest(test, program, fullPath, file)
	}

	if tr.ShouldRunTest(test, VirtualMachineEngine) {
		tr.runVMTest(test, program, fullPath, file)
	}
}

func (tr *TestRunner) parseTestFile(file string) (*Test, error) {
	start := time.Now()

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key := ""
	test := &Test{
		message:         "",
		file:            "",
		expect:          "",
		errors:          "",
		supportedEngine: AllEngines,
	}

	for line := range strings.SplitSeq(string(content), "\n") {
		if strings.HasPrefix(line, "--TEST--") || strings.HasPrefix(line, "---TEST---") {
			key = "message"
		} else if strings.HasPrefix(line, "--FILE--") || strings.HasPrefix(line, "---FILE---") {
			key = "file"
		} else if strings.HasPrefix(line, "--EXPECT--") || strings.HasPrefix(line, "---EXPECT---") {
			key = "expect"
		} else if strings.HasPrefix(line, "--ERROR--") || strings.HasPrefix(line, "---ERROR---") {
			key = "errors"
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
		case "errors":
			test.errors += line + "\n"
		}
	}

	test.file = tr.cleanString(test.file)
	test.expect = tr.cleanString(test.expect)
	test.errors = tr.cleanString(test.errors)

	if strings.HasSuffix(file, ".vm.zent") {
		test.supportedEngine = VirtualMachineEngine
	} else if strings.HasSuffix(file, ".eval.zent") {
		test.supportedEngine = EvaluatorEngine
	}

	tr.setTiming(ReadingFilesTiming, tr.getTiming(ReadingFilesTiming)+time.Since(start))

	return test, nil
}

func (tr *TestRunner) printSuccessStatusMessage(test *Test, engineType EngineType) {
	tr.passedTests++

	messages = append(messages, fmt.Sprintf("  %s✔%s %s %s[%s%s]%s\n",
		colors.Green, colors.Reset, tr.cleanString(test.message),
		colors.Gray, engineType.GetTag(), colors.Gray, colors.Reset,
	))
}

func (tr *TestRunner) printErrorStatusMessage(test *Test, fullPath, message string, engineType EngineType) {
	exitStatusCode = 1

	errorMessage := fmt.Sprintf("%s %s[%s%s]%s\n     %s",
		tr.cleanString(test.message),
		colors.Gray, engineType.GetTag(), colors.Gray, colors.Reset,
		message,
	)
	collectedErrors = append(collectedErrors, fmt.Sprintf("%s: %s", fullPath, errorMessage))

	messages = append(messages, fmt.Sprintf("  %s✖%s %s\n", colors.Red, colors.Reset, errorMessage))
}

func (tr *TestRunner) cleanString(str string) string {
	str = strings.ReplaceAll(str, "\r\n", "\n")
	str = strings.Trim(str, "\n")

	return str
}

func (tr *TestRunner) stripFileLocationsFromError(err string) string {
	lines := strings.Split(err, "\n")

	for i, line := range lines {
		if strings.Contains(line, "at ") && strings.Contains(line, ".zent:") {
			fileInfo := strings.Split(line, ".zent:")[1]

			lines[i] = fmt.Sprintf("    at <unknown>:%s", fileInfo)
		}
	}

	return strings.Trim(strings.Join(lines, "\n"), "\n")
}

func (tr *TestRunner) ShouldRunTest(test *Test, engine EngineType) bool {
	if tr.options.Engine != AllEngines {
		if tr.options.Engine != engine {
			return false
		}

		if test.supportedEngine != AllEngines && test.supportedEngine != tr.options.Engine {
			return false
		}

		return true
	}

	if test.supportedEngine != AllEngines && test.supportedEngine != engine {
		return false
	}

	return true
}

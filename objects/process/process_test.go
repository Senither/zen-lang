package process

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestFakeArgs(t *testing.T) {
	defer RestoreFromFake()

	expectedArgs := []string{"zen", "is", "awesome"}
	FakeArgs(expectedArgs)

	if len(os.Args) != len(expectedArgs) {
		t.Fatalf("Expected %d args, but got %d", len(expectedArgs), len(os.Args))
	}

	for i, expected := range expectedArgs {
		if os.Args[i] != expected {
			t.Errorf("Expected arg %d to be '%s', but got '%s'", i, expected, os.Args[i])
		}
	}
}

func TestExit(t *testing.T) {
	defer RestoreFromFake()
	Fake()

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Exit(0)
	Exit(1)
	Exit(42)
	Exit(255)

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)

	messages := strings.Split(buf.String(), "\n")
	expectedMessages := []string{
		"INTERNAL_FAKE_PROCESS_EXIT(0)",
		"INTERNAL_FAKE_PROCESS_EXIT(1)",
		"INTERNAL_FAKE_PROCESS_EXIT(42)",
		"INTERNAL_FAKE_PROCESS_EXIT(255)",
	}

	for i, expected := range expectedMessages {
		if messages[i] != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, messages[i])
		}
	}
}

func TestLookupEnv(t *testing.T) {
	defer RestoreFromFake()
	Fake()

	FakeEnv("ZEN_TEST_ENV", "test_value")

	value, exists := LookupEnv("ZEN_TEST_ENV")
	if !exists {
		t.Fatalf("Expected ZEN_TEST_ENV to exist")
	}

	if value != "test_value" {
		t.Errorf("Expected ZEN_TEST_ENV to be 'test_value', but got '%s'", value)
	}

	_, exists = LookupEnv("NON_EXISTENT_ENV")
	if exists {
		t.Fatalf("Expected NON_EXISTENT_ENV to not exist")
	}
}

package vm

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/objects/process"
)

func TestStandardOutWriteAndReadAll(t *testing.T) {
	s := &StandardOut{messages: []string{}, muted: true}

	s.Write("hello")
	s.Write(" world")

	messages := s.ReadAll()
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}

	if strings.Join(messages, "") != "hello world" {
		t.Fatalf("unexpected combined messages: %q", strings.Join(messages, ""))
	}
}

func TestStandardOutClear(t *testing.T) {
	s := &StandardOut{messages: []string{"a", "b"}, muted: true}

	s.Clear()

	if len(s.ReadAll()) != 0 {
		t.Fatalf("expected messages to be cleared, got %d", len(s.ReadAll()))
	}
}

func TestStandardOutMute(t *testing.T) {
	s := &StandardOut{messages: []string{}, muted: false}

	wasMutedInFn := false
	result := s.Mute(func() objects.Object {
		wasMutedInFn = s.muted
		s.Write("muted-message")
		return objects.TRUE
	})

	if !wasMutedInFn {
		t.Fatal("expected stdout to be muted while mute callback is running")
	}

	if s.muted {
		t.Fatal("expected stdout mute state to be restored after callback")
	}

	if result != objects.TRUE {
		t.Fatalf("expected mute callback result to be returned, got %T", result)
	}

	if strings.Join(s.ReadAll(), "") != "muted-message" {
		t.Fatalf("expected muted callback writes to still be buffered, got %q", strings.Join(s.ReadAll(), ""))
	}
}

func TestCaptureStdoutForBuiltin(t *testing.T) {
	Stdout.Clear()
	t.Cleanup(func() {
		Stdout.Clear()
	})

	t.Run("captures builtin output", func(t *testing.T) {
		Stdout.Clear()

		result := captureStdoutForBuiltin(
			objects.GetBuiltinByName("print"),
			[]objects.Object{&objects.String{Value: "captured"}},
		)

		if result != objects.NULL {
			t.Fatalf("expected NULL result from print builtin, got %T", result)
		}

		if strings.Join(Stdout.ReadAll(), "") != "captured" {
			t.Fatalf("expected captured output to be buffered, got %q", strings.Join(Stdout.ReadAll(), ""))
		}
	})

	t.Run("does not buffer newline-only output", func(t *testing.T) {
		Stdout.Clear()

		result := captureStdoutForBuiltin(
			&objects.Builtin{Fn: func(args ...objects.Object) (objects.Object, error) {
				fmt.Fprint(os.Stdout, "\n")
				return objects.NULL, nil
			}},
			nil,
		)

		if result != objects.NULL {
			t.Fatalf("expected NULL result, got %T", result)
		}

		if len(Stdout.ReadAll()) != 0 {
			t.Fatalf("expected newline-only output to be ignored, got %q", strings.Join(Stdout.ReadAll(), ""))
		}
	})

	t.Run("converts builtin errors to error objects", func(t *testing.T) {
		Stdout.Clear()

		result := captureStdoutForBuiltin(
			&objects.Builtin{Fn: func(args ...objects.Object) (objects.Object, error) {
				fmt.Fprint(os.Stdout, "before-error")
				return nil, errors.New("builtin failed")
			}},
			nil,
		)

		if result == nil || result.Type() != objects.ERROR_OBJ {
			t.Fatalf("expected error object result, got %T", result)
		}

		if strings.Join(Stdout.ReadAll(), "") != "before-error" {
			t.Fatalf("expected output emitted before error to be captured, got %q", strings.Join(Stdout.ReadAll(), ""))
		}
	})
}

func TestCaptureStdoutForCaptureEnabledBuiltins(t *testing.T) {
	Stdout.Clear()
	t.Cleanup(func() {
		Stdout.Clear()
		process.RestoreFromFake()
	})

	t.Run("print builtin writes to stdout buffer", func(t *testing.T) {
		Stdout.Clear()

		builtin := objects.GetBuiltinByName("print")
		if builtin == nil {
			t.Fatal("expected print builtin to exist")
		}

		if !builtin.CaptureStdout {
			t.Fatal("expected print builtin to be marked as capture-enabled")
		}

		captureStdoutForBuiltin(
			builtin,
			[]objects.Object{&objects.String{Value: "print-output"}},
		)

		if strings.Join(Stdout.ReadAll(), "") != "print-output" {
			t.Fatalf("expected print output in stdout buffer, got %q", strings.Join(Stdout.ReadAll(), ""))
		}
	})

	t.Run("println builtin writes to stdout buffer", func(t *testing.T) {
		Stdout.Clear()

		builtin := objects.GetBuiltinByName("println")
		if builtin == nil {
			t.Fatal("expected println builtin to exist")
		}

		if !builtin.CaptureStdout {
			t.Fatal("expected println builtin to be marked as capture-enabled")
		}

		captureStdoutForBuiltin(
			builtin,
			[]objects.Object{&objects.String{Value: "println-output"}},
		)

		if strings.Join(Stdout.ReadAll(), "") != "println-output\n" {
			t.Fatalf("expected println output in stdout buffer, got %q", strings.Join(Stdout.ReadAll(), ""))
		}
	})

	t.Run("process.exit global builtin writes fake exit output to stdout buffer", func(t *testing.T) {
		Stdout.Clear()
		process.Fake()
		defer process.RestoreFromFake()

		builtin := objects.GetGlobalBuiltinByName("process", "exit")
		if builtin == nil {
			t.Fatal("expected process.exit global builtin to exist")
		}

		if !builtin.CaptureStdout {
			t.Fatal("expected process.exit builtin to be marked as capture-enabled")
		}

		captureStdoutForBuiltin(
			builtin,
			[]objects.Object{&objects.Integer{Value: 42}},
		)

		if strings.Join(Stdout.ReadAll(), "") != "INTERNAL_FAKE_PROCESS_EXIT(42)\n" {
			t.Fatalf("expected process.exit fake output in stdout buffer, got %q", strings.Join(Stdout.ReadAll(), ""))
		}
	})
}

package process

import (
	"fmt"
	"os"
)

var faking bool = false
var args []string = nil
var envs map[string]string = nil

func Fake() {
	faking = true
}

func RestoreFromFake() {
	faking = false

	if args != nil {
		os.Args = args
		args = nil
	}

	if envs != nil {
		for key, value := range envs {
			os.Setenv(key, value)
		}

		envs = nil
	}
}

func FakeArgs(fakeArgs []string) {
	args = os.Args
	os.Args = fakeArgs
}

func FakeEnv(key, value string) {
	if envs == nil {
		envs = make(map[string]string)
	}

	envs[key] = os.Getenv(key)
	os.Setenv(key, value)
}

func Exit(code int) {
	if faking {
		fmt.Fprintf(os.Stdout, "INTERNAL_FAKE_PROCESS_EXIT(%d)\n", code)
		return
	}

	os.Exit(code)
}

func LookupEnv(key string) (string, bool) {
	if faking {
		if value, exists := envs[key]; exists {
			return value, true
		}
		return "", false
	}

	return os.LookupEnv(key)
}

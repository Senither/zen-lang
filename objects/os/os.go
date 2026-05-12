package os

import (
	"os"
	"runtime"
	"strconv"
)

var faking bool = false
var fakeState map[string]string = nil

func Fake() {
	faking = true
}

func SetFakeState(key string, value string) {
	if !faking {
		return
	}

	if fakeState == nil {
		fakeState = make(map[string]string)
	}

	fakeState[key] = value
}

func getFakeState(key string, defaultValue string) string {
	if value, ok := fakeState[key]; ok {
		return value
	}

	return defaultValue
}

func RestoreFromFake() {
	faking = false
	fakeState = nil
}

func Hostname() (string, error) {
	if faking {
		return getFakeState("hostname", "OS_FAKE_HOSTNAME"), nil
	}

	return os.Hostname()
}

func Platform() (string, error) {
	if faking {
		return getFakeState("platform", "OS_FAKE_PLATFORM"), nil
	}

	return runtime.GOOS, nil
}

func Arch() (string, error) {
	if faking {
		return getFakeState("arch", "OS_FAKE_ARCH"), nil
	}

	return runtime.GOARCH, nil
}

func CPUs() int64 {
	if faking {
		val, err := strconv.ParseInt(getFakeState("cpus", "8"), 10, 64)
		if err != nil {
			return 8
		}

		return val
	}

	return int64(runtime.NumCPU())
}

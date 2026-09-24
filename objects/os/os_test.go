package os

import (
	"runtime"
	"testing"
)

func TestFakeStateDefaults(t *testing.T) {
	defer RestoreFromFake()

	Fake()

	hostname, err := Hostname()
	if err != nil {
		t.Fatalf("Did not expect error getting fake hostname, got %v", err)
	}

	if hostname != "OS_FAKE_HOSTNAME" {
		t.Errorf("Expected fake hostname to be 'OS_FAKE_HOSTNAME', got '%s'", hostname)
	}

	platform, err := Platform()
	if err != nil {
		t.Fatalf("Did not expect error getting fake platform, got %v", err)
	}

	if platform != "OS_FAKE_PLATFORM" {
		t.Errorf("Expected fake platform to be 'OS_FAKE_PLATFORM', got '%s'", platform)
	}

	arch, err := Arch()
	if err != nil {
		t.Fatalf("Did not expect error getting fake architecture, got %v", err)
	}

	if arch != "OS_FAKE_ARCH" {
		t.Errorf("Expected fake architecture to be 'OS_FAKE_ARCH', got '%s'", arch)
	}

	if cpus := CPUs(); cpus != 8 {
		t.Errorf("Expected fake CPU count to be 8, got %d", cpus)
	}
}

func TestSetFakeState(t *testing.T) {
	defer RestoreFromFake()

	Fake()

	SetFakeState("hostname", "test-host")
	SetFakeState("platform", "test-platform")
	SetFakeState("arch", "test-arch")
	SetFakeState("cpus", "16")

	hostname, _ := Hostname()
	if hostname != "test-host" {
		t.Errorf("Expected fake hostname to be 'test-host', got '%s'", hostname)
	}

	platform, _ := Platform()
	if platform != "test-platform" {
		t.Errorf("Expected fake platform to be 'test-platform', got '%s'", platform)
	}

	arch, _ := Arch()
	if arch != "test-arch" {
		t.Errorf("Expected fake architecture to be 'test-arch', got '%s'", arch)
	}

	if cpus := CPUs(); cpus != 16 {
		t.Errorf("Expected fake CPU count to be 16, got %d", cpus)
	}
}

func TestSetFakeStateWithoutFaking(t *testing.T) {
	defer RestoreFromFake()

	SetFakeState("hostname", "ignored-host")

	if fakeState != nil {
		t.Errorf("Expected fake state to remain nil when faking is disabled")
	}
}

func TestFakeCPUsWithInvalidValue(t *testing.T) {
	defer RestoreFromFake()

	Fake()
	SetFakeState("cpus", "not-a-number")

	if cpus := CPUs(); cpus != 8 {
		t.Errorf("Expected invalid fake CPU count to fall back to 8, got %d", cpus)
	}
}

func TestRestoreFromFake(t *testing.T) {
	Fake()
	SetFakeState("hostname", "test-host")
	RestoreFromFake()

	if faking {
		t.Errorf("Expected faking to be disabled after RestoreFromFake")
	}

	if fakeState != nil {
		t.Errorf("Expected fake state to be cleared after RestoreFromFake")
	}

	platform, err := Platform()
	if err != nil {
		t.Fatalf("Did not expect error getting real platform, got %v", err)
	}

	if platform != runtime.GOOS {
		t.Errorf("Expected real platform to be '%s', got '%s'", runtime.GOOS, platform)
	}

	arch, err := Arch()
	if err != nil {
		t.Fatalf("Did not expect error getting real architecture, got %v", err)
	}

	if arch != runtime.GOARCH {
		t.Errorf("Expected real architecture to be '%s', got '%s'", runtime.GOARCH, arch)
	}
}

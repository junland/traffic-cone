package daemon

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestEnsureParentDirCreatesNestedDir(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	targetFile := filepath.Join(root, "a", "b", "c", "daemon.pid")

	if err := ensureParentDir(targetFile); err != nil {
		t.Fatalf("ensureParentDir returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Dir(targetFile)); err != nil {
		t.Fatalf("expected parent dir to exist: %v", err)
	}
}

func TestStatusFileNotFound(t *testing.T) {
	t.Parallel()

	pidFile := filepath.Join(t.TempDir(), "missing.pid")
	pid, running, err := Status(pidFile)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pid != 0 || running {
		t.Fatalf("expected pid=0,running=false; got pid=%d running=%t", pid, running)
	}
}

func TestStatusRejectsEmptyAndInvalidPID(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	tests := []struct {
		name    string
		content string
		errSub  string
	}{
		{name: "empty", content: "", errSub: "pid file is empty"},
		{name: "nonnumeric", content: "abc\n", errSub: "invalid pid"},
		{name: "nonpositive", content: "0\n", errSub: "invalid pid"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pidFile := filepath.Join(root, tc.name+".pid")
			if err := os.WriteFile(pidFile, []byte(tc.content), 0o644); err != nil {
				t.Fatalf("write pid file: %v", err)
			}

			_, _, err := Status(pidFile)
			if err == nil {
				t.Fatalf("expected error for %q", tc.name)
			}
			if !strings.Contains(err.Error(), tc.errSub) {
				t.Fatalf("expected error containing %q, got %q", tc.errSub, err.Error())
			}
		})
	}
}

func TestAcquirePIDFileWritesAndReleases(t *testing.T) {
	t.Parallel()

	pidFile := filepath.Join(t.TempDir(), "subdir", "daemon.pid")

	release, err := acquirePIDFile(pidFile)
	if err != nil {
		t.Fatalf("acquirePIDFile returned error: %v", err)
	}

	data, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatalf("read pid file: %v", err)
	}

	gotPID, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatalf("pid file does not contain an integer: %v", err)
	}
	if gotPID != os.Getpid() {
		t.Fatalf("expected pid %d, got %d", os.Getpid(), gotPID)
	}

	release()
	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		t.Fatalf("expected pid file to be removed, got err=%v", err)
	}
}

func TestDockerHostFromSocket(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "already scheme",
			input:  "unix:///var/run/docker.sock",
			output: "unix:///var/run/docker.sock",
		},
		{
			name:   "unix path",
			input:  "/var/run/docker.sock",
			output: "unix:///var/run/docker.sock",
		},
		{
			name:   "windows npipe",
			input:  "//./pipe/docker_engine",
			output: "npipe:////./pipe/docker_engine",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := dockerHostFromSocket(tc.input)
			if got != tc.output {
				t.Fatalf("expected %q, got %q", tc.output, got)
			}
		})
	}
}

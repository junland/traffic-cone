package daemon

import (
	"strings"
	"testing"
)

func TestStartRequiresPIDFile(t *testing.T) {
	t.Parallel()

	cfg := RunConfig{
		DockerSocket: "/var/run/docker.sock",
	}
	err := Start(cfg)
	if err == nil {
		t.Fatal("expected error when PIDFile is empty, got nil")
	}
	if !strings.Contains(err.Error(), "pid file is required") {
		t.Fatalf("expected 'pid file is required' error, got %q", err.Error())
	}
}

func TestStartRequiresDockerSocket(t *testing.T) {
	t.Parallel()

	cfg := RunConfig{
		PIDFile: "/tmp/test.pid",
	}
	err := Start(cfg)
	if err == nil {
		t.Fatal("expected error when DockerSocket is empty, got nil")
	}
	if !strings.Contains(err.Error(), "docker socket is required") {
		t.Fatalf("expected 'docker socket is required' error, got %q", err.Error())
	}
}

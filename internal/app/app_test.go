package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunNoArgs(t *testing.T) {
	t.Parallel()

	if got := Run([]string{"-invalid-flag"}); got != 1 {
		t.Fatalf("expected exit code 1 for invalid flag, got %d", got)
	}
	if got := Run([]string{"-pid-file"}); got != 1 {
		t.Fatalf("expected exit code 1 for missing flag value, got %d", got)
	}
}

func TestResolveHAProxyDataPlanePasswordFromFile(t *testing.T) {
	t.Parallel()

	passwordFile := filepath.Join(t.TempDir(), "password.txt")
	if err := os.WriteFile(passwordFile, []byte(" secret \n"), 0o600); err != nil {
		t.Fatalf("write password file: %v", err)
	}

	got, err := resolveHAProxyDataPlanePassword(passwordFile)
	if err != nil {
		t.Fatalf("resolve password from file: %v", err)
	}
	if got != "secret" {
		t.Fatalf("expected trimmed password %q, got %q", "secret", got)
	}
}

func TestResolveHAProxyDataPlanePasswordFromEnv(t *testing.T) {
	t.Setenv("HAPROXY_DATA_PLANE_API_PASSWORD", "env-secret")

	got, err := resolveHAProxyDataPlanePassword("")
	if err != nil {
		t.Fatalf("resolve password from env: %v", err)
	}
	if got != "env-secret" {
		t.Fatalf("expected env password %q, got %q", "env-secret", got)
	}
}

func TestResolveHAProxyDataPlanePasswordFileTooLarge(t *testing.T) {
	t.Parallel()

	passwordFile := filepath.Join(t.TempDir(), "password.txt")
	oversized := strings.Repeat("x", maxHAProxyDataPlanePasswordFileSize+1)
	if err := os.WriteFile(passwordFile, []byte(oversized), 0o600); err != nil {
		t.Fatalf("write oversized password file: %v", err)
	}

	_, err := resolveHAProxyDataPlanePassword(passwordFile)
	if err == nil {
		t.Fatal("expected error for oversized password file, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds maximum allowed size") {
		t.Fatalf("expected oversized error, got %q", err.Error())
	}
}

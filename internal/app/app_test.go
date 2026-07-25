package app

import (
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

package app

import (
	"testing"
)

func TestRunNoArgs(t *testing.T) {
	t.Parallel()

	if got := Run(nil); got != 1 {
		t.Fatalf("expected exit code 1 for no args, got %d", got)
	}
	if got := Run([]string{}); got != 1 {
		t.Fatalf("expected exit code 1 for empty args, got %d", got)
	}
}

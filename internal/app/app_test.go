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

func TestRunParsesFlagsAfterDaemonName(t *testing.T) {
	t.Parallel()

	if got := Run([]string{"traffic-cone", "-docker-socket="}); got != 1 {
		t.Fatalf("expected exit code 1 when parsed flags produce invalid config, got %d", got)
	}
}

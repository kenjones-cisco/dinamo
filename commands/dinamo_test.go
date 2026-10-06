package commands

import (
	"context"
	"log/slog"
	"strings"
	"testing"
)

func isLogLevelEnabled(level slog.Level) bool {
	return slog.Default().Enabled(context.Background(), level)
}

func TestRootCmdDebug(t *testing.T) {
	cmd := getRootCommand()

	if result := runCmd(cmd, "-D"); result.Error != nil {
		t.Error(result.Error)
	}

	if !isLogLevelEnabled(slog.LevelDebug) {
		t.Error("expected debug logging to be enabled")
	}
}

func TestRootCmdBadLogLevel(t *testing.T) {
	cmd := getRootCommand()

	result := runCmd(cmd, "-l=fake")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !strings.Contains(result.Output, "Unknown log-level provided:") {
		t.Error("expected an error message to be printed out, but the message was not found.")
	}

	if isLogLevelEnabled(slog.LevelDebug) || !isLogLevelEnabled(slog.LevelInfo) {
		t.Error("expected info logging to be enabled and debug logging to be disabled")
	}
}

func TestRootCmdLogLevel(t *testing.T) {
	cmd := getRootCommand()

	result := runCmd(cmd, "--log-level warn")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if isLogLevelEnabled(slog.LevelInfo) || !isLogLevelEnabled(slog.LevelWarn) {
		t.Error("expected warn logging to be enabled and info logging to be disabled")
	}
}

func TestRootCmdDisplayVersion(t *testing.T) {
	cmd := getRootCommand()
	// short flag
	result := runCmd(cmd, "-v")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !strings.Contains(result.Output, "Dynamic Generator\n version") {
		t.Error("expected version message to be printed out, but the message was not found.")
	}

	result = runCmd(cmd, "--version")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !strings.Contains(result.Output, "Dynamic Generator\n version") {
		t.Error("expected version message to be printed out, but the message was not found.")
	}
}

package commands

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestRootCmdDebug(t *testing.T) {
	cmd := getRootCommand()

	if result := runCmd(cmd, "-D"); result.Error != nil {
		t.Error(result.Error)
	}

	if logrus.GetLevel() != logrus.DebugLevel {
		t.Errorf("expected: %v, got: %v", logrus.DebugLevel, logrus.GetLevel())
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

	if logrus.GetLevel() != logrus.InfoLevel {
		t.Errorf("expected: %v, got: %v", logrus.InfoLevel, logrus.GetLevel())
	}
}

func TestRootCmdLogLevel(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "--log-level warn")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if logrus.GetLevel() != logrus.WarnLevel {
		t.Errorf("expected: %v, got: %v", logrus.WarnLevel, logrus.GetLevel())
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

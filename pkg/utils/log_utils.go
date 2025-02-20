package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"

	log "github.com/charmbracelet/log"
	"github.com/cloudposse/atmos/pkg/ui/theme"
	"github.com/fatih/color"
)

const (
	LogLevelTrace   = "Trace"
	LogLevelDebug   = "Debug"
	LogLevelInfo    = "Info"
	LogLevelWarning = "Warning"
)

// PrintMessage prints the message to the console.
func PrintMessage(message string) {
	fmt.Println(message)
}

// PrintMessageInColor prints the message to the console using the provided color.
func PrintMessageInColor(message string, messageColor *color.Color) {
	_, _ = messageColor.Fprint(os.Stdout, message)
}

// Deprecated: Use `log.Error` instead. This function will be removed in a future release.
func PrintErrorInColor(message string) {
	messageColor := theme.Colors.Error
	_, _ = messageColor.Fprint(os.Stderr, message)
}

// LogErrorAndExit logs errors to std.Error and exits with an error code.
func LogErrorAndExit(err error) {
	log.Error(err)

	// Find the executed command's exit code from the error
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		os.Exit(exitError.ExitCode())
	}

	os.Exit(1)
}

// LogError logs errors to std.Error.
func LogError(err error) {
	if err != nil {
		log.Error(err)
		// Print stack trace
		if log.GetLevel() == log.DebugLevel {
			debug.PrintStack()
		}
	}
}

// LogTrace logs the provided trace message.
func LogTrace(message string) {
	LogDebug(message)
}

// LogDebug logs the provided debug message.
func LogDebug(message string) {
	log.Debug(message)
}

// LogInfo logs the provided info message.
func LogInfo(message string) {
	log.Info(message)
}

// LogWarning logs the provided warning message.
func LogWarning(message string) {
	log.Warn(message)
}

package tests

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath" // For resolving absolute paths
	"regexp"
	"strings"
	"testing"

	"github.com/creack/pty"
	"github.com/fatih/color"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

// Command-line flag for regenerating snapshots
var regenerateSnapshots = flag.Bool("regenerate-snapshots", false, "Regenerate all golden snapshots")
var startingDir string
var snapshotBaseDir string

type DiffResult struct {
	HasDiff  bool
	DiffText string
}

type Expectation struct {
	Stdout       []string            `yaml:"stdout"`
	Stderr       []string            `yaml:"stderr"`
	ExitCode     int                 `yaml:"exit_code"`
	FileExists   []string            `yaml:"file_exists"`
	FileContains map[string][]string `yaml:"file_contains"`
	IgnoreDiffs  []string            `yaml:"ignore_diffs"`
}

type TestCase struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Enabled     bool              `yaml:"enabled"`
	Workdir     string            `yaml:"workdir"`
	Command     string            `yaml:"command"`
	Args        []string          `yaml:"args"`
	Env         map[string]string `yaml:"env"`
	Expect      Expectation       `yaml:"expect"`
	Tty         bool              `yaml:"tty"`
}

type TestSuite struct {
	Tests []TestCase `yaml:"tests"`
}

func loadTestSuite(filePath string) (*TestSuite, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var suite TestSuite
	err = yaml.Unmarshal(data, &suite)
	if err != nil {
		return nil, err
	}

	return &suite, nil
}

type PathManager struct {
	OriginalPath string
	Prepended    []string
}

// NewPathManager initializes a PathManager with the current PATH.
func NewPathManager() *PathManager {
	return &PathManager{
		OriginalPath: os.Getenv("PATH"),
		Prepended:    []string{},
	}
}

// Prepend adds directories to the PATH with precedence.
func (pm *PathManager) Prepend(dirs ...string) {
	for _, dir := range dirs {
		absPath, err := filepath.Abs(dir)
		if err != nil {
			fmt.Printf("Failed to resolve absolute path for %q: %v\n", dir, err)
			continue
		}
		pm.Prepended = append(pm.Prepended, absPath)
	}
}

// GetPath returns the updated PATH.
func (pm *PathManager) GetPath() string {
	return fmt.Sprintf("%s%c%s",
		strings.Join(pm.Prepended, string(os.PathListSeparator)),
		os.PathListSeparator,
		pm.OriginalPath,
	)
}

// Apply updates the PATH environment variable globally.
func (pm *PathManager) Apply() error {
	return os.Setenv("PATH", pm.GetPath())
}

// Apply regex ignore patterns to text
func applyIgnorePatterns(input string, patterns []string) string {
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		input = re.ReplaceAllString(input, "")
	}
	return input
}

// Simulate TTY command execution
func simulateTtyCommand(cmd *exec.Cmd) (string, error) {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to start TTY: %v", err)
	}
	defer func() { _ = ptmx.Close() }() // Best effort cleanup

	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(ptmx)
	if err != nil {
		return "", fmt.Errorf("failed to read TTY output: %v", err)
	}
	return buffer.String(), nil
}

// Execute the command and return the exit code
func executeCommand(t *testing.T, cmd *exec.Cmd) int {
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		t.Fatalf("Command execution failed: %v", err)
	}
	return 0
}

func colorizeDiff(diff string) string {
	var result strings.Builder
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "+") {
			// Additions are green
			result.WriteString(color.GreenString(line) + "\n")
		} else if strings.HasPrefix(line, "-") {
			// Deletions are red
			result.WriteString(color.RedString(line) + "\n")
		} else if strings.HasPrefix(line, "@") {
			// Diff context is yellow
			result.WriteString(color.YellowString(line) + "\n")
		} else {
			// Unchanged lines are plain
			result.WriteString(line + "\n")
		}
	}
	return result.String()
}

// loadTestSuites loads and merges all .yaml files from the test-cases directory
func loadTestSuites(testCasesDir string) (*TestSuite, error) {
	var mergedSuite TestSuite

	entries, err := os.ReadDir(testCasesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read test-cases directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			filePath := filepath.Join(testCasesDir, entry.Name())
			suite, err := loadTestSuite(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to load %s: %v", filePath, err)
			}
			mergedSuite.Tests = append(mergedSuite.Tests, suite.Tests...)
		}
	}

	return &mergedSuite, nil
}

// Entry point for tests to parse flags and handle setup/teardown
func TestMain(m *testing.M) {
	// Declare err in the function's scope
	var err error

	// Capture the starting working directory
	startingDir, err = os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get the current working directory: %v\n", err)
		os.Exit(1) // Exit with a non-zero code to indicate failure
	}

	// Define the base directory for snapshots relative to startingDir
	snapshotBaseDir = filepath.Join(startingDir, "snapshots")

	flag.Parse() // Parse command-line flags
	os.Exit(m.Run())
}

func TestCLICommands(t *testing.T) {
	// Declare err in the function's scope
	var err error

	// Initialize PathManager and update PATH
	pathManager := NewPathManager()
	pathManager.Prepend("../build", "..")
	err = pathManager.Apply()
	if err != nil {
		t.Fatalf("Failed to apply updated PATH: %v", err)
	}
	fmt.Printf("Updated PATH: %s\n", pathManager.GetPath())

	// Update the test suite loading
	testSuite, err := loadTestSuites("test-cases")
	if err != nil {
		t.Fatalf("Failed to load test suites: %v", err)
	}

	for _, tc := range testSuite.Tests {

		if !tc.Enabled {
			t.Logf("Skipping disabled test: %s", tc.Name)
			continue
		}

		t.Run(tc.Name, func(t *testing.T) {
			defer func() {
				// Change back to the original working directory after the test
				if err := os.Chdir(startingDir); err != nil {
					t.Fatalf("Failed to change back to the starting directory: %v", err)
				}
			}()

			// Change to the specified working directory
			if tc.Workdir != "" {
				err := os.Chdir(tc.Workdir)
				if err != nil {
					t.Fatalf("Failed to change directory to %q: %v", tc.Workdir, err)
				}
			}

			// Check if the binary exists
			binaryPath, err := exec.LookPath(tc.Command)
			if err != nil {
				t.Fatalf("Binary not found: %s. Current PATH: %s", tc.Command, pathManager.GetPath())
			}

			// Prepare the command
			cmd := exec.Command(binaryPath, tc.Args...)

			// Set environment variables
			envVars := os.Environ()
			for key, value := range tc.Env {
				envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
			}
			cmd.Env = envVars

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			var exitCode int
			if tc.Tty {
				// Run the command in a pseudo-terminal
				ptyOutput, err := simulateTtyCommand(cmd)
				if err != nil {
					if exitErr, ok := err.(*exec.ExitError); ok {
						exitCode = exitErr.ExitCode()
					} else {
						t.Fatalf("Failed to run TTY command: %v", err)
					}
				}
				stdout.WriteString(ptyOutput)
			} else {
				// Run the command directly and capture the exit code
				err := cmd.Run()
				if err != nil {
					if exitErr, ok := err.(*exec.ExitError); ok {
						exitCode = exitErr.ExitCode()
					} else {
						t.Fatalf("Failed to run command: %v", err)
					}
				}
			}

			// Validate exit code
			if !verifyExitCode(t, tc.Expect.ExitCode, exitCode) {
				t.Errorf("Description: %s", tc.Description)
			}

			// Validate stdout
			if !verifyOutput(t, "stdout", stdout.String(), tc.Expect.Stdout) {
				t.Errorf("Description: %s", tc.Description)
			}

			// Validate stderr
			if !verifyOutput(t, "stderr", stderr.String(), tc.Expect.Stderr) {
				t.Errorf("Description: %s", tc.Description)
			}

			// Validate file existence
			if !verifyFileExists(t, tc.Expect.FileExists) {
				t.Errorf("Description: %s", tc.Description)
			}

			// Validate file contents
			if !verifyFileContains(t, tc.Expect.FileContains) {
				t.Errorf("Description: %s", tc.Description)
			}
		})
	}
}

func verifyExitCode(t *testing.T, expected, actual int) bool {
	success := true
	if expected != actual {
		t.Errorf("Reason: Expected exit code %d, got %d", expected, actual)
		success = false
	}
	return success
}

func verifyOutput(t *testing.T, outputType, output string, patterns []string) bool {
	success := true
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			t.Errorf("Invalid %s regex: %q, error: %v", outputType, pattern, err)
			success = false
			continue
		}
		if !re.MatchString(output) {
			t.Errorf("Reason: %s did not match pattern %q.", outputType, pattern)
			t.Errorf("Output: %q", output)
			success = false
		}
	}
	return success
}

func verifyFileExists(t *testing.T, files []string) bool {
	success := true
	for _, file := range files {
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			t.Errorf("Reason: Expected file does not exist: %q", file)
			success = false
		}
	}
	return success
}

func verifyFileContains(t *testing.T, filePatterns map[string][]string) bool {
	success := true
	for file, patterns := range filePatterns {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("Reason: Failed to read file %q: %v", file, err)
			success = false
			continue
		}
		for _, pattern := range patterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				t.Errorf("Invalid regex for file %q: %q, error: %v", file, pattern, err)
				success = false
				continue
			}
			if !re.Match(content) {
				t.Errorf("Reason: File %q did not match pattern %q.", file, pattern)
				t.Errorf("Content: %q", string(content))
				success = false
			}
		}
	}
	return success
}

func updateSnapshot(path, output string) {
	fullPath := filepath.Join(snapshotBaseDir, path)
	err := os.MkdirAll(filepath.Dir(fullPath), 0755) // Ensure parent directories exist
	if err != nil {
		panic(fmt.Sprintf("Failed to create snapshot directory: %v", err))
	}
	err = os.WriteFile(fullPath, []byte(output), 0644) // Write snapshot
	if err != nil {
		panic(fmt.Sprintf("Failed to write snapshot file: %v", err))
	}
}

func readSnapshot(t *testing.T, path string) string {
	fullPath := filepath.Join(snapshotBaseDir, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Error reading snapshot file %q: %v", fullPath, err)
	}
	return string(data)
}

func verifySnapshot(t *testing.T, goldenPath, actualOutput string, ignorePatterns []string, regenerate bool) DiffResult {
	// Filter ignored diffs from the actual output
	filteredActual := applyIgnorePatterns(actualOutput, ignorePatterns)

	if regenerate {
		t.Logf("Regenerating snapshot at %q", filepath.Join(snapshotBaseDir, goldenPath))
		updateSnapshot(goldenPath, filteredActual)
		return DiffResult{HasDiff: false}
	}

	var filteredExpected string
	if _, err := os.Stat(filepath.Join(snapshotBaseDir, goldenPath)); errors.Is(err, os.ErrNotExist) {
		t.Logf("Snapshot not found at %q. Creating new one.", filepath.Join(snapshotBaseDir, goldenPath))
		updateSnapshot(goldenPath, filteredActual)
		return DiffResult{HasDiff: false}
	} else {
		filteredExpected = applyIgnorePatterns(readSnapshot(t, goldenPath), ignorePatterns)
	}

	// Compare outputs and generate a colorized diff
	if diff := cmp.Diff(filteredExpected, filteredActual); diff != "" {
		return DiffResult{
			HasDiff:  true,
			DiffText: colorizeDiff(diff),
		}
	}

	return DiffResult{HasDiff: false}
}

// filterIgnoredDiffs removes ignored patterns from the output
func filterIgnoredDiffs(output string, ignores []string) string {
	for _, pattern := range ignores {
		re := regexp.MustCompile(pattern)
		output = re.ReplaceAllString(output, "")
	}
	return output
}

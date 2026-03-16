package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// buildBinary compiles the calculator binary into a temp directory and returns
// the path to the executable. The test is skipped if the build fails.
func buildBinary(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	binaryPath := filepath.Join(dir, "calculator"+ext)

	// Determine module root by walking up from this file's location.
	// We rely on the go.mod being at the repo root.
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/calculator")
	cmd.Dir = moduleRoot(t)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build calculator binary: %v\nstderr: %s", err, stderr.String())
	}
	return binaryPath
}

// moduleRoot locates the directory that contains go.mod by starting from the
// test binary's working directory and walking upward.
func moduleRoot(t *testing.T) string {
	t.Helper()
	// os.Getwd() during tests returns the package directory.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	// Walk up until we find go.mod.
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find go.mod in any parent directory")
		}
		dir = parent
	}
}

// TestHelpFlag_ExitCodeZero verifies that 'calculator --help' exits with code 0.
func TestHelpFlag_ExitCodeZero(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "--help")
	err := cmd.Run()
	if err != nil {
		// exec.ExitError means the process exited non-zero.
		if exitErr, ok := err.(*exec.ExitError); ok {
			t.Errorf("'calculator --help' exited with code %d; want 0", exitErr.ExitCode())
		} else {
			t.Errorf("'calculator --help' failed to run: %v", err)
		}
	}
}

// TestHelpFlag_OutputContainsUsage verifies that 'calculator --help' prints "Usage" to stdout.
func TestHelpFlag_OutputContainsUsage(t *testing.T) {
	binary := buildBinary(t)

	var stdout bytes.Buffer
	cmd := exec.Command(binary, "--help")
	cmd.Stdout = &stdout
	_ = cmd.Run() // ignore exit error; we only check output here

	output := stdout.String()
	if !strings.Contains(output, "Usage") {
		t.Errorf("'calculator --help' stdout does not contain 'Usage'; got:\n%s", output)
	}
}

// TestHelpFlag_OutputContainsOperations verifies the help output mentions key operations.
func TestHelpFlag_OutputContainsOperations(t *testing.T) {
	binary := buildBinary(t)

	var stdout bytes.Buffer
	cmd := exec.Command(binary, "--help")
	cmd.Stdout = &stdout
	_ = cmd.Run()

	output := stdout.String()
	required := []string{"add", "subtract", "multiply", "divide"}
	for _, op := range required {
		if !strings.Contains(output, op) {
			t.Errorf("'calculator --help' output missing operation %q; got:\n%s", op, output)
		}
	}
}

// TestNoHelpFlag_ProcessStarts verifies that 'calculator' (no flags) starts without
// immediately crashing (we send EOF to stdin so it exits cleanly).
func TestNoHelpFlag_ProcessStarts(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary)
	// Provide an empty stdin so the REPL reads EOF and exits.
	cmd.Stdin = strings.NewReader("")
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Exit code 1 is acceptable if the REPL returns an error on EOF.
			// We just want to verify the process started (not a crash / signal).
			if exitErr.ExitCode() < 0 {
				t.Errorf("calculator exited with signal/crash: %v", err)
			}
		} else {
			t.Errorf("failed to run calculator: %v", err)
		}
	}
}

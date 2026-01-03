package testutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// AssertNoError fails the test if err is not nil
func AssertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: unexpected error: %v", msg, err)
	}
}

// AssertError fails the test if err is nil
func AssertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error but got nil", msg)
	}
}

// AssertEqual fails the test if expected != actual
func AssertEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: expected %v, got %v", msg, expected, actual)
	}
}

// AssertNotEqual fails the test if expected == actual
func AssertNotEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: expected values to be different, but both are %v", msg, expected)
	}
}

// AssertNotNil fails the test if value is nil
func AssertNotNil(t *testing.T, value interface{}, msg string) {
	t.Helper()
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		t.Fatalf("%s: expected non-nil value", msg)
	}
}

// AssertNil fails the test if value is not nil
func AssertNil(t *testing.T, value interface{}, msg string) {
	t.Helper()
	if value != nil && !(reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		t.Fatalf("%s: expected nil but got %v", msg, value)
	}
}

// AssertTrue fails the test if condition is false
func AssertTrue(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Errorf("%s: expected true but got false", msg)
	}
}

// AssertFalse fails the test if condition is true
func AssertFalse(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Errorf("%s: expected false but got true", msg)
	}
}

// AssertContains fails the test if str doesn't contain substr
func AssertContains(t *testing.T, str, substr string, msg string) {
	t.Helper()
	if !contains(str, substr) {
		t.Errorf("%s: expected string to contain %q, but got %q", msg, substr, str)
	}
}

// AssertErrorContains fails the test if error message doesn't contain expected string
func AssertErrorContains(t *testing.T, err error, expectedMsg string, msg string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error but got nil", msg)
	}
	if !contains(err.Error(), expectedMsg) {
		t.Errorf("%s: expected error to contain %q, but got %q", msg, expectedMsg, err.Error())
	}
}

// AssertPanics fails the test if fn doesn't panic
func AssertPanics(t *testing.T, fn func(), msg string) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: expected panic but function completed normally", msg)
		}
	}()
	fn()
}

// AssertNoPanic fails the test if fn panics
func AssertNoPanic(t *testing.T, fn func(), msg string) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%s: unexpected panic: %v", msg, r)
		}
	}()
	fn()
}

// LoadTestFile loads a file from the testdata directory
func LoadTestFile(t *testing.T, filename string) []byte {
	t.Helper()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to load test file %s: %v", filename, err)
	}

	return data
}

// CreateTempFile creates a temporary file with the given content
func CreateTempFile(t *testing.T, content []byte) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "catfetch_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.Write(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	return tmpFile.Name()
}

// CreateTempDir creates a temporary directory
func CreateTempDir(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "catfetch_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return tmpDir
}

// WriteTestFile writes content to a file in the given directory
func WriteTestFile(t *testing.T, dir, filename string, content []byte) string {
	t.Helper()

	filePath := filepath.Join(dir, filename)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("Failed to write test file %s: %v", filePath, err)
	}

	return filePath
}

// CopyFile copies a file from src to dst
func CopyFile(t *testing.T, src, dst string) {
	t.Helper()

	sourceFile, err := os.Open(src)
	if err != nil {
		t.Fatalf("Failed to open source file %s: %v", src, err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		t.Fatalf("Failed to create destination file %s: %v", dst, err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		t.Fatalf("Failed to copy file from %s to %s: %v", src, dst, err)
	}
}

// SkipIfShort skips the test if running in short mode
func SkipIfShort(t *testing.T, msg string) {
	t.Helper()
	if testing.Short() {
		t.Skipf("Skipping in short mode: %s", msg)
	}
}

// Parallel marks the test as capable of running in parallel
func Parallel(t *testing.T) {
	t.Helper()
	t.Parallel()
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

// findSubstring performs a simple substring search
func findSubstring(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// FormatTestName formats a test name for table-driven tests
func FormatTestName(parts ...interface{}) string {
	return fmt.Sprint(parts...)
}

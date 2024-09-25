package flags

import (
	"testing"
)

func TestOutputFileFlag(t *testing.T) {
	tests := []struct {
		args     []string
		expected bool
	}{
		{[]string{"-O=/path/to/file"}, true},
		{[]string{"-P=/path/to/file"}, false},
		{[]string{"-O=output.txt"}, true},
		{[]string{"--rate-limit=10"}, false},
		{[]string{"-O="}, true}, // Edge case: empty value
	}

	for _, test := range tests {
		result := OutputFileFlag(test.args)
		if result != test.expected {
			t.Errorf("OutputFileFlag(%v) = %v; expected %v", test.args, result, test.expected)
		}
	}
}

func TestGetFlagInput(t *testing.T) {
	tests := []struct {
		flagInput string
		expected  string
	}{
		{"-O=/path/to/file", "/path/to/file"},
		{"-P=/path/to/dir", "/path/to/dir"},
		{"-O=", ""},
		{"-P=", ""},
	}

	for _, test := range tests {
		result := GetFlagInput(test.flagInput)
		if result != test.expected {
			t.Errorf("GetFlagInput(%q) = %q; expected %q", test.flagInput, result, test.expected)
		}
	}
}

func TestFlagType(t *testing.T) {
	tests := []struct {
		args     []string
		expected []string
	}{
		{[]string{"-O=/path/to/file", "-P=/path/to/dir", "--rate-limit=10"}, []string{"-O=", "-P=", "--rate-limit="}},
		{[]string{"-P=/path/to/dir", "--rate-limit=10"}, []string{"-P=", "--rate-limit="}},
		{[]string{"--rate-limit=10", "random_arg"}, []string{"--rate-limit="}},
		{[]string{"-O=", "-P="}, []string{"-O=", "-P="}}, // Edge case: empty values
		{[]string{"random_arg"}, []string{}},             // Edge case: no flags
	}

	for _, test := range tests {
		result := FlagType(test.args)
		if !stringSlicesEqual(result, test.expected) {
			t.Errorf("FlagType(%v) = %v; expected %v", test.args, result, test.expected)
		}
	}
}

// Helper function to compare slices of strings
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

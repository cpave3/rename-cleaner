package main

import (
	"testing"
)

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect bool
	}{
		{"valid simple name", "file.txt", true},
		{"invalid name with space", "file name.txt", false},
		{"invalid name with brackets", "file[name].txt", false},
		{"name with hyphen", "file-name.txt", true},
		{"name with period", "file.name.txt", true},
		{"name with underscore", "file_name.txt", true},
		{"invalid name with special chars", "file!name.txt", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := isValidName(test.input)
			if got != test.expect {
				t.Errorf("isValidName(%q) = %v, want %v", test.input, got, test.expect)
			}
		})
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"valid simple name", "file.txt", "file.txt"},
		{"invalid name with space", "file name.txt", "file-name.txt"},
		{"invalid name with brackets", "file[name].txt", "file-name.txt"},
		{"name with hyphen", "file-name.txt", "file-name.txt"},
		{"name with period", "file.name.txt", "file.name.txt"},
		{"name with underscore", "file_name.txt", "file_name.txt"},
		{"invalid name with special chars", "file!name.txt", "filename.txt"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := sanitizeName(test.input)
			if got != test.expect {
				t.Errorf("sanitizeName(%q) = %v, want %v", test.input, got, test.expect)
			}
		})
	}
}

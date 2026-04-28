package main

import (
	"reflect"
	"testing"
)

func TestParseFzfExpectOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		wantKey  string
		wantSel  []string
	}{
		{
			"empty input",
			[]byte{},
			"",
			nil,
		},
		{
			"no key with single file",
			[]byte("\nfile.txt\n"),
			"",
			[]string{"file.txt"},
		},
		{
			"ctrl-e key with single file",
			[]byte("ctrl-e\nfile.txt\n"),
			"ctrl-e",
			[]string{"file.txt"},
		},
		{
			"no key with multiple files",
			[]byte("\nfile1.txt\nfile2.go\nfile3.md"),
			"",
			[]string{"file1.txt", "file2.go", "file3.md"},
		},
		{
			"CRLF line endings",
			[]byte("ctrl-e\r\nfile.txt\r\n"),
			"ctrl-e",
			[]string{"file.txt"},
		},
		{
			"blank lines in selection filtered",
			[]byte("\nfile1.txt\n\nfile2.go\n"),
			"",
			[]string{"file1.txt", "file2.go"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotSel := parseFzfExpectOutput(tt.input)
			if gotKey != tt.wantKey {
				t.Errorf("key: got %q, want %q", gotKey, tt.wantKey)
			}
			if !reflect.DeepEqual(gotSel, tt.wantSel) {
				t.Errorf("selected: got %v, want %v", gotSel, tt.wantSel)
			}
		})
	}
}

func TestFirstCommandToken(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"bat", "bat"},
		{"bat --color=always", "bat"},
		{`"C:\tools\bat.exe" --color`, `C:\tools\bat.exe`},
		{`"bat"`, "bat"},
		{`"bat" --color`, "bat"},
		{"  bat  ", "bat"},
		{"cmd /C type {1}", "cmd"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := firstCommandToken(tt.input); got != tt.want {
				t.Errorf("firstCommandToken(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// IsBinaryFile reports whether path is binary by scanning the first 512 bytes
// for null bytes — the same heuristic git uses.
func IsBinaryFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && n == 0 {
		return false
	}
	return bytes.IndexByte(buf[:n], 0) >= 0
}

// OpenFileWithDefault launches path using the Windows default program association.
func OpenFileWithDefault(path string) error {
	cmd := exec.Command("cmd.exe", "/C", "start", "", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}

// OpenInExplorer opens path in Windows Explorer with the file selected.
func OpenInExplorer(path string) error {
	cmd := exec.Command("cmd.exe", "/C", "start", "explorer.exe", fmt.Sprintf(`/select,%s`, path))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}

// RunEditorCommand runs editor with the given args in dir (empty dir = inherit CWD).
func RunEditorCommand(editor, dir string, args ...string) error {
	cmd := exec.Command(editor, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	return nil
}

// OpenMixedFiles splits files into binary vs text, launches binary files with
// the OS default program, and opens text files together in the editor.
func OpenMixedFiles(editor string, files []string) error {
	var textFiles []string
	for _, f := range files {
		if IsBinaryFile(f) {
			if err := OpenFileWithDefault(f); err != nil {
				return fmt.Errorf("open default %s: %w", f, err)
			}
		} else {
			textFiles = append(textFiles, f)
		}
	}
	if len(textFiles) == 0 {
		return nil
	}
	return RunEditorCommand(editor, "", textFiles...)
}

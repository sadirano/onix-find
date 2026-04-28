package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func runFF(target, query, editor string, vis *Config) error {
	if _, err := exec.LookPath("fzf"); err != nil {
		return fmt.Errorf("ff requires fzf in PATH")
	}
	if _, err := exec.LookPath("es"); err == nil {
		return runFFWithEverythingStream(target, query, editor, vis)
	}
	return runFFWithWalkFallback(target, query, editor, vis)
}

func buildFZFArgs(query string, vis *Config) []string {
	preview := resolvePreviewCommand(vis.FZF.FF.Preview, "")
	args := []string{
		"--multi",
		"--expect=ctrl-e",
		"--prompt", vis.FZF.FF.Prompt,
		"--query", query,
	}
	args = appendLayoutArg(args, vis.FZF.FF.Layout)
	if preview != "" {
		args = append(args, "--preview", preview)
		if window := strings.TrimSpace(vis.FZF.FF.PreviewWindow); window != "" {
			args = append(args, "--preview-window", window)
		}
	}
	return args
}

func runFFWithEverythingStream(target, query, editor string, vis *Config) error {
	esArgs := []string{"-p", "-path", target}
	if query != "" {
		esArgs = append(esArgs, query)
	}
	esCmd := exec.Command("es", esArgs...)
	esCmd.Stderr = os.Stderr
	esOut, err := esCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("run es: %w", err)
	}
	if err := esCmd.Start(); err != nil {
		return fmt.Errorf("run es: %w", err)
	}
	err = runFFWithInput(esOut, query, editor, vis)
	_ = esCmd.Wait()
	return err
}

func runFFWithWalkFallback(target, query, editor string, vis *Config) error {
	files, err := gatherFilesWithWalk(target, query)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Println("No files found.")
		return nil
	}
	var input bytes.Buffer
	for _, file := range files {
		input.WriteString(file)
		input.WriteByte('\n')
	}
	return runFFWithInput(&input, query, editor, vis)
}

// runFFWithInput runs fzf with the given stdin and handles the result.
func runFFWithInput(stdin io.Reader, query, editor string, vis *Config) error {
	fzfCmd := exec.Command("fzf", buildFZFArgs(query, vis)...)
	fzfCmd.Stdin = stdin
	fzfCmd.Stderr = os.Stderr
	out, err := fzfCmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil
		}
		return fmt.Errorf("run fzf: %w", err)
	}
	key, selected := parseFzfExpectOutput(out)
	if len(selected) == 0 {
		return nil
	}
	if key == "ctrl-e" {
		for _, file := range selected {
			if err := OpenInExplorer(file); err != nil {
				return err
			}
		}
		return nil
	}
	return OpenMixedFiles(editor, selected)
}

func gatherFilesWithWalk(root, query string) ([]string, error) {
	queryLower := strings.ToLower(query)
	files := make([]string, 0, 256)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if queryLower != "" && !strings.Contains(strings.ToLower(filepath.Base(path)), queryLower) {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk files: %w", err)
	}
	sort.Strings(files)
	return files, nil
}

func parseFzfExpectOutput(out []byte) (string, []string) {
	lines := strings.Split(strings.ReplaceAll(string(out), "\r\n", "\n"), "\n")
	if len(lines) == 0 {
		return "", nil
	}
	key := strings.TrimSpace(lines[0])
	var selected []string
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		selected = append(selected, line)
	}
	return key, selected
}

func resolvePreviewCommand(configured, fallback string) string {
	preview := strings.TrimSpace(configured)
	if preview == "" {
		return fallback
	}
	token := firstCommandToken(preview)
	if token == "" {
		return fallback
	}
	base := strings.ToLower(strings.TrimSuffix(filepath.Base(token), filepath.Ext(token)))
	if base == "bat" {
		if _, err := exec.LookPath("bat"); err != nil {
			return fallback
		}
	}
	return preview
}

func firstCommandToken(command string) string {
	command = strings.TrimSpace(command)
	if command == "" {
		return ""
	}
	if command[0] == '"' {
		rest := command[1:]
		if idx := strings.Index(rest, `"`); idx >= 0 {
			return rest[:idx]
		}
		return rest
	}
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func appendLayoutArg(args []string, layout string) []string {
	layout = strings.TrimSpace(layout)
	if layout == "" || strings.EqualFold(layout, "default") {
		return args
	}
	return append(args, "--layout", layout)
}

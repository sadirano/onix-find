package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sadirano/onix-find/internal/config"
	"github.com/sadirano/onix-find/internal/search"
)

func main() {
	target := strings.TrimSpace(os.Getenv("ONIX_TARGET"))
	if target == "" {
		var err error
		target, err = os.Getwd()
		if err != nil {
			fatal("no ONIX_TARGET set and cannot determine working directory: %v", err)
		}
	}

	editor := resolveEditor()
	vis := config.LoadVisualConfig(os.Getenv("ONIX_HOME"))

	query := strings.TrimSpace(strings.Join(os.Args[1:], " "))

	if err := search.RunFF(target, query, editor, &vis); err != nil {
		fatal("%v", err)
	}
}

func resolveEditor() string {
	for _, env := range []string{"ONIX_EDITOR", "EDITOR"} {
		if e := strings.TrimSpace(os.Getenv(env)); e != "" {
			return e
		}
	}
	return "nvim"
}

func fatal(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "onix-ff: "+format+"\n", a...)
	os.Exit(1)
}

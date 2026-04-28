package main

// onix-ff — file finder module for onix.
//
// Invoked by the onix dispatch mechanism:
//
//	ONIX_MODULE=ff onix <alias> [query]   find files, open in editor (ctrl-e = open in Explorer)
//
// Environment variables (set by onix core at dispatch time):
//
//	ONIX_TARGET   resolved absolute path of the target directory
//	ONIX_HOME     onix home directory (~/.onix), used to locate onix.visual.toml
//	ONIX_EDITOR   preferred editor (resolved by onix core from config + EDITOR env)

import (
	"fmt"
	"os"
	"strings"
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
	vis := loadConfig(os.Getenv("ONIX_HOME"))

	query := strings.TrimSpace(strings.Join(os.Args[1:], " "))

	if err := runFF(target, query, editor, &vis); err != nil {
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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const configFileName = "onix.visual.toml"

type Config struct {
	FZF FZFConfig `toml:"fzf"`
}

type FZFConfig struct {
	FF PickerConfig `toml:"ff"`
}

type PickerConfig struct {
	Prompt        string `toml:"prompt"`
	Layout        string `toml:"layout"`
	Preview       string `toml:"preview"`
	PreviewWindow string `toml:"preview_window"`
}

func defaultConfig() Config {
	return Config{
		FZF: FZFConfig{
			FF: PickerConfig{
				Prompt:        "> ",
				Layout:        "default",
				Preview:       `bat --color=always {}`,
				PreviewWindow: "right,55%,border-left",
			},
		},
	}
}

func loadConfig(onixHome string) Config {
	cfg := defaultConfig()
	if onixHome == "" {
		return cfg
	}
	p := filepath.Join(onixHome, configFileName)
	if _, err := os.Stat(p); err != nil {
		return cfg
	}
	if _, err := toml.DecodeFile(p, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to parse %s: %v\n", p, err)
		return cfg
	}
	applyDefaults(&cfg)
	return cfg
}

func applyDefaults(v *Config) {
	def := defaultConfig()
	v.FZF.FF.Prompt = fallback(v.FZF.FF.Prompt, def.FZF.FF.Prompt)
	v.FZF.FF.Layout = fallback(v.FZF.FF.Layout, def.FZF.FF.Layout)
	v.FZF.FF.Preview = fallback(v.FZF.FF.Preview, def.FZF.FF.Preview)
	v.FZF.FF.PreviewWindow = fallback(v.FZF.FF.PreviewWindow, def.FZF.FF.PreviewWindow)
}

func appendLayoutArg(args []string, layout string) []string {
	layout = strings.TrimSpace(layout)
	if layout == "" || strings.EqualFold(layout, "default") {
		return args
	}
	return append(args, "--layout", layout)
}

func fallback(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}

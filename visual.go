package main

import (
	"fmt"
	"os"
	"path/filepath"

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
		return defaultConfig()
	}
	return cfg
}

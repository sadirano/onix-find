package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const configFileName = "onix.visual.toml"

type VisualConfig struct {
	FZF FZFConfig `toml:"fzf"`
	RG  RGConfig  `toml:"rg"`
}

type FZFConfig struct {
	FF PickerConfig `toml:"ff"`
	SG PickerConfig `toml:"sg"`
}

type PickerConfig struct {
	Prompt        string `toml:"prompt"`
	Layout        string `toml:"layout"`
	Color         string `toml:"color"`
	Preview       string `toml:"preview"`
	PreviewWindow string `toml:"preview_window"`
}

type RGConfig struct {
	Color string `toml:"color"`
	Case  string `toml:"case"`
}

func DefaultVisualConfig() VisualConfig {
	return VisualConfig{
		FZF: FZFConfig{
			FF: PickerConfig{
				Prompt: "> ",
				Layout: "default",
			},
			SG: PickerConfig{
				Prompt:        "> ",
				Layout:        "default",
				Color:         "hl:-1:underline,hl+:-1:underline:reverse",
				Preview:       `bat --color=always {1} --highlight-line {2}`,
				PreviewWindow: "up,60%,border-bottom,+{2}+3/3,~3",
			},
		},
		RG: RGConfig{
			Color: "always",
			Case:  "smart",
		},
	}
}

func LoadVisualConfig(onixHome string) VisualConfig {
	cfg := DefaultVisualConfig()
	if onixHome == "" {
		return cfg
	}
	p := filepath.Join(onixHome, configFileName)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return cfg
	}
	if _, err := toml.DecodeFile(p, &cfg); err != nil {
		// Silent fail
	}
	applyDefaults(&cfg)
	return cfg
}

func applyDefaults(v *VisualConfig) {
	def := DefaultVisualConfig()
	v.RG.Color = fallback(v.RG.Color, def.RG.Color)
	v.RG.Case = fallback(v.RG.Case, def.RG.Case)

	v.FZF.FF.Prompt = fallback(v.FZF.FF.Prompt, def.FZF.FF.Prompt)
	v.FZF.FF.Layout = fallback(v.FZF.FF.Layout, def.FZF.FF.Layout)
	v.FZF.FF.Color = fallback(v.FZF.FF.Color, def.FZF.FF.Color)
	v.FZF.FF.Preview = fallback(v.FZF.FF.Preview, def.FZF.FF.Preview)
	v.FZF.FF.PreviewWindow = fallback(v.FZF.FF.PreviewWindow, def.FZF.FF.PreviewWindow)

	v.FZF.SG.Prompt = fallback(v.FZF.SG.Prompt, def.FZF.SG.Prompt)
	v.FZF.SG.Layout = fallback(v.FZF.SG.Layout, def.FZF.SG.Layout)
	v.FZF.SG.Color = fallback(v.FZF.SG.Color, def.FZF.SG.Color)
	v.FZF.SG.Preview = fallback(v.FZF.SG.Preview, def.FZF.SG.Preview)
	v.FZF.SG.PreviewWindow = fallback(v.FZF.SG.PreviewWindow, def.FZF.SG.PreviewWindow)
}

func AppendLayoutArg(args []string, layout string) []string {
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

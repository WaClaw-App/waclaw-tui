package tui

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/WaClaw-App/waclaw/internal/tui/i18n"
	"github.com/WaClaw-App/waclaw/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

// Config represents the user-customizable config.yaml structure.
// This is the runtime configuration (locale, etc.) that is NOT visual theme data.
// Per doc/00-philosophy-and-design.md, locale goes in ~/.waclaw/config.yaml,
// separate from the visual palette in theme.yaml.
type Config struct {
	// Locale is the UI display language: "id" (Indonesian, default) or "en" (English).
	Locale string `yaml:"locale"`
}

// DefaultConfig returns the built-in config with Indonesian as default locale.
func DefaultConfig() *Config {
	return &Config{Locale: "id"}
}

// LoadConfig reads config.yaml from ~/.waclaw/ and applies overrides.
// Falls back to DefaultConfig() if file is missing or invalid.
func LoadConfig() *Config {
	cfg := DefaultConfig()

	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		// File missing or unreadable — use defaults.
		return cfg
	}

	// Parse user overrides on top of defaults.
	var user Config
	if err := yaml.Unmarshal(data, &user); err != nil {
		// Invalid YAML — use defaults.
		return cfg
	}

	if user.Locale != "" {
		cfg.Locale = user.Locale
	}
	return cfg
}

// ApplyConfig applies runtime configuration (locale, etc.).
func ApplyConfig(cfg *Config) {
	if cfg.Locale != "" {
		locale := i18n.Locale(cfg.Locale)
		if i18n.IsValidLocale(locale) {
			i18n.SetLocale(locale)
		}
	}
}

// ConfigPath returns the expected path of config.yaml.
func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".waclaw", "config.yaml")
	}
	return filepath.Join(home, ".waclaw", "config.yaml")
}

// ThemeConfig represents the user-customizable theme.yaml structure.
// Contains ONLY visual palette data — no locale or other non-visual settings.
type ThemeConfig struct {
	Palette struct {
		Bg          string `yaml:"bg"`
		BgRaised    string `yaml:"bg_raised"`
		BgActive    string `yaml:"bg_active"`
		Text        string `yaml:"text"`
		TextMuted   string `yaml:"text_muted"`
		TextDim     string `yaml:"text_dim"`
		Success     string `yaml:"success"`
		Warning     string `yaml:"warning"`
		Danger      string `yaml:"danger"`
		Accent      string `yaml:"accent"`
		Pulse       string `yaml:"pulse"`
		Highlight   string `yaml:"highlight"`
		Gold        string `yaml:"gold"`
		Celebration string `yaml:"celebration"`
	} `yaml:"palette"`
}

// paletteFieldOrder defines the mapping between ThemeConfig.Palette field names
// and the corresponding style package color variable pointers.
// This is the single source of truth for the field→token mapping — no more
// repetitive if-chains.
var paletteFieldOrder = []struct {
	YAMLName string // matches the yaml tag (and struct field name in PascalCase)
	// ColorPtr points to the corresponding lipgloss.Color variable in the style package.
	// We store these as *lipgloss.Color so ApplyTheme can update them in-place.
	ColorPtr *lipgloss.Color
}{
	{"Bg", &style.Bg},
	{"BgRaised", &style.BgRaised},
	{"BgActive", &style.BgActive},
	{"Text", &style.Text},
	{"TextMuted", &style.TextMuted},
	{"TextDim", &style.TextDim},
	{"Success", &style.Success},
	{"Warning", &style.Warning},
	{"Danger", &style.Danger},
	{"Accent", &style.Accent},
	{"Pulse", &style.Pulse},
	{"Highlight", &style.Highlight},
	{"Gold", &style.Gold},
	{"Celebration", &style.Celebration},
}

// DefaultTheme returns the built-in theme matching doc/16-design-system.md.
func DefaultTheme() *ThemeConfig {
	cfg := &ThemeConfig{}
	cfg.Palette.Bg = "#0A0A0B"
	cfg.Palette.BgRaised = "#141416"
	cfg.Palette.BgActive = "#1A1A1E"
	cfg.Palette.Text = "#E8E8EC"
	cfg.Palette.TextMuted = "#6B6B76"
	cfg.Palette.TextDim = "#3D3D44"
	cfg.Palette.Success = "#34D399"
	cfg.Palette.Warning = "#FBBF24"
	cfg.Palette.Danger = "#F87171"
	cfg.Palette.Accent = "#818CF8"
	cfg.Palette.Pulse = "#818CF866"
	cfg.Palette.Highlight = "#FFFFFF22"
	cfg.Palette.Gold = "#FFD700"
	cfg.Palette.Celebration = "#FFFFFF"
	return cfg
}

// LoadTheme reads theme.yaml from ~/.waclaw/ and applies overrides.
// Falls back to DefaultTheme() if file is missing or invalid.
func LoadTheme() *ThemeConfig {
	cfg := DefaultTheme()

	path := ThemePath()
	data, err := os.ReadFile(path)
	if err != nil {
		// File missing or unreadable — use defaults.
		return cfg
	}

	// Parse user overrides on top of defaults.
	var user ThemeConfig
	if err := yaml.Unmarshal(data, &user); err != nil {
		// Invalid YAML — use defaults.
		return cfg
	}

	// Apply non-empty overrides.
	applyOverrides(cfg, &user)
	return cfg
}

// ApplyTheme updates the runtime color tokens from a ThemeConfig.
// This allows hot-reloading without restart.
func ApplyTheme(cfg *ThemeConfig) {
	paletteVal := reflect.ValueOf(cfg.Palette)

	for _, entry := range paletteFieldOrder {
		field := paletteVal.FieldByName(entry.YAMLName)
		if !field.IsValid() {
			continue
		}
		hex := field.String()
		if hex == "" {
			continue
		}
		*entry.ColorPtr = lipgloss.Color(hex)
	}

	// Rebuild all pre-built styles to pick up new colors.
	style.RebuildStyles()
}

// ThemePath returns the expected path of theme.yaml.
func ThemePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback — should not happen in practice.
		return filepath.Join(".", ".waclaw", "theme.yaml")
	}
	return filepath.Join(home, ".waclaw", "theme.yaml")
}

// applyOverrides merges non-empty user overrides into the default config.
// Uses reflect to iterate over palette fields — single loop, no repetitive
// per-field if-chains.
func applyOverrides(base, override *ThemeConfig) {
	baseVal := reflect.ValueOf(&base.Palette).Elem()
	overVal := reflect.ValueOf(&override.Palette).Elem()

	for i := 0; i < baseVal.NumField(); i++ {
		overField := overVal.Field(i)
		if overField.Kind() == reflect.String && overField.String() != "" {
			baseVal.Field(i).Set(overField)
		}
	}
}

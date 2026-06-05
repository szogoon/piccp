package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Target TargetConfig `toml:"target"`
	SDCard SDCardConfig `toml:"sdcard"`
	Import ImportConfig `toml:"import"`
	UI     UIConfig     `toml:"ui"`
}

type TargetConfig struct {
	OutputDir string `toml:"output_dir"`
}

type SDCardConfig struct {
	MinSizeGB   float64 `toml:"min_size_gb"`
	MaxSizeGB   float64 `toml:"max_size_gb"`
	AutoUnmount bool    `toml:"auto_unmount"`
}

type ImportConfig struct {
	GroupingGapDays    int      `toml:"grouping_gap_days"`
	AllowedExtensions []string `toml:"allowed_extensions"`
	PreserveDir        bool     `toml:"preserve_directory_structure"`
	OverwriteExisting  bool     `toml:"overwrite_existing"`
}

type UIConfig struct {
	ProgressBar bool `toml:"progress_bar"`
	Verbose     bool `toml:"verbose"`
	DryRun      bool `toml:"dry_run"`
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, ".config", "piccp", "config.toml")
	}

	var cfg Config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		Target: TargetConfig{
			OutputDir: filepath.Join(home, "Pictures", "Imports"),
		},
		SDCard: SDCardConfig{
			MinSizeGB:   4,
			MaxSizeGB:   1024,
			AutoUnmount: true,
		},
		Import: ImportConfig{
			GroupingGapDays:    1,
			AllowedExtensions: []string{"jpg", "jpeg", "png", "mp4", "mov"},
			PreserveDir:        false,
			OverwriteExisting:  false,
		},
		UI: UIConfig{
			ProgressBar: true,
			Verbose:     true,
			DryRun:      false,
		},
	}
}

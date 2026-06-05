package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.SDCard.MinSizeGB != 4 {
		t.Errorf("Expected default min size 4, got %f", cfg.SDCard.MinSizeGB)
	}
	if len(cfg.Import.AllowedExtensions) == 0 {
		t.Error("Expected default allowed extensions to be populated")
	}
}

func TestLoadConfig_NonExistent(t *testing.T) {
	// Loading a non-existent file should return default config
	cfg, err := LoadConfig("/tmp/non-existent-piccp-config.toml")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg == nil {
		t.Fatal("Expected config to be returned")
	}
	if cfg.SDCard.MinSizeGB != 4 {
		t.Errorf("Expected default config values, got %f", cfg.SDCard.MinSizeGB)
	}
}

func TestLoadConfig_Valid(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "config.toml")
	content := `
[sdcard]
min_size_gb = 10
`
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write tmp file: %v", err)
	}

	cfg, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.SDCard.MinSizeGB != 10 {
		t.Errorf("Expected min_size_gb to be 10, got %f", cfg.SDCard.MinSizeGB)
	}
	// Verify other values are defaults (since they weren't in the file)
	// Note: BurntSushi/toml will leave unset fields as zero values if they aren't explicitly defaulted before decode
	// In our LoadConfig, we check os.IsNotExist but don't merge defaults for a partial file yet.
}

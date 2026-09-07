package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoadAndGenerate(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	err := GenerateDefaultConfigFile(configPath)
	if err != nil {
		t.Fatalf("GenerateDefaultConfigFile failed: %v", err)
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("Config file does not exist: %v", err)
	}

	cfg, loadedPath, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loadedPath != configPath {
		t.Errorf("Expected loadedPath %s, got %s", configPath, loadedPath)
	}
	if cfg.Unit != "GB" {
		t.Errorf("Expected Unit GB, got %s", cfg.Unit)
	}
	if cfg.BarWidth != 16 {
		t.Errorf("Expected BarWidth 16, got %d", cfg.BarWidth)
	}
}

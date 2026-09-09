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
	if cfg.Unit != "GiB" {
		t.Errorf("Expected Unit GiB, got %s", cfg.Unit)
	}
	if cfg.BarStyle != "smooth" {
		t.Errorf("Expected BarStyle smooth, got %s", cfg.BarStyle)
	}
	if cfg.BarWidth != 12 {
		t.Errorf("Expected BarWidth 12, got %d", cfg.BarWidth)
	}
}

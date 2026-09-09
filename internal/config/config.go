package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds user configuration for gdf
type Config struct {
	// DefaultMounts: list of mountpoints to display by default.
	// If empty, all real local disks are shown.
	DefaultMounts []string `yaml:"default_mounts"`

	// ExcludeMounts: list of mount points or prefixes to ignore.
	ExcludeMounts []string `yaml:"exclude_mounts"`

	// ExcludeFSTypes: list of filesystem types to ignore.
	ExcludeFSTypes []string `yaml:"exclude_fstypes"`

	// OnlyRealDisks: if true, ignore pseudo/virtual filesystems (default: true).
	OnlyRealDisks bool `yaml:"only_real_disks"`

	// Unit: measurement unit ("GiB", "GB", "human-si", "human-iec"). Default: "GiB".
	Unit string `yaml:"unit"`

	// BarWidth: width of the progress bar in columns (default: 10).
	BarWidth int `yaml:"bar_width"`

	// BarStyle: "smooth", "unicode", "blocks", or "ascii". Default: "smooth".
	BarStyle string `yaml:"bar_style"`

	// Colors: enable ANSI colors (default: true).
	Colors bool `yaml:"colors"`

	// ShowTotal: show total summary row (default: true).
	ShowTotal bool `yaml:"show_total"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		DefaultMounts: []string{}, // empty means show all detected real disks
		ExcludeMounts: []string{
			"/var/lib/docker",
			"/var/lib/containers",
			"/run/credentials",
		},
		ExcludeFSTypes: []string{
			"sysfs", "proc", "devtmpfs", "devpts", "tmpfs", "cgroup", "cgroup2",
			"pstore", "bpf", "securityfs", "selinuxfs", "fusectl", "mqueue",
			"hugetlbfs", "tracefs", "debugfs", "configfs", "autofs", "nsfs",
			"overlay", "efivarfs", "binfmt_misc", "squashfs", "iso9660",
			"rpc_pipefs", "nfsd",
		},
		OnlyRealDisks: true,
		Unit:          "GiB",
		BarWidth:      10,
		BarStyle:      "smooth",
		Colors:        true,
		ShowTotal:     true,
	}
}

// DefaultConfigPath returns the path to ~/.config/gdf/config.yaml
func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig != "" {
		return filepath.Join(xdgConfig, "gdf", "config.yaml"), nil
	}
	return filepath.Join(home, ".config", "gdf", "config.yaml"), nil
}

// LoadConfig attempts to load the config from a custom path, or the default paths
func LoadConfig(customPath string) (*Config, string, error) {
	cfg := DefaultConfig()

	var searchPaths []string
	if customPath != "" {
		searchPaths = append(searchPaths, customPath)
	} else {
		if def, err := DefaultConfigPath(); err == nil {
			searchPaths = append(searchPaths, def)
		}
		if home, err := os.UserHomeDir(); err == nil {
			searchPaths = append(searchPaths, filepath.Join(home, ".gdfrc.yaml"))
			searchPaths = append(searchPaths, filepath.Join(home, ".gdfrc"))
		}
	}

	for _, path := range searchPaths {
		data, err := os.ReadFile(path)
		if err == nil {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, path, fmt.Errorf("error parsing config file %s: %w", path, err)
			}
			return cfg, path, nil
		}
	}

	return cfg, "", nil
}

// GenerateDefaultConfigFile generates a commented YAML configuration file
func GenerateDefaultConfigFile(targetPath string) error {
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	template := `# ==========================================================
# gdf (Go Disk Free) Configuration File
# ==========================================================

# 1. Default mount points to display.
# If empty ([]), gdf automatically discovers and displays all real physical disks.
# If populated, gdf will ONLY display the specified mount points by default.
# Example:
# default_mounts:
#   - /
#   - /home
#   - /boot
default_mounts: []

# 2. Excluded mount points (exact paths or prefixes)
exclude_mounts:
  - /var/lib/docker
  - /var/lib/containers
  - /run/credentials

# 3. Excluded filesystem types
exclude_fstypes:
  - sysfs
  - proc
  - devtmpfs
  - devpts
  - tmpfs
  - cgroup
  - cgroup2
  - pstore
  - bpf
  - securityfs
  - selinuxfs
  - fusectl
  - mqueue
  - hugetlbfs
  - tracefs
  - debugfs
  - configfs
  - autofs
  - nsfs
  - overlay
  - efivarfs
  - binfmt_misc
  - squashfs

# 4. Only show real physical storage devices (skip pseudo/virtual filesystems)
only_real_disks: true

# 5. Measurement unit:
# "GiB"       - Binary Gibibytes (1 GiB = 1,073,741,824 bytes, like df -h) - default
# "GB"        - Decimal Gigabytes (1 GB = 1,000,000,000 bytes)
# "human-si"  - Auto-scale decimal (kB, MB, GB, TB)
# "human-iec" - Auto-scale binary (KiB, MiB, GiB, TiB)
unit: "GiB"

# 6. Progress bar width (in characters)
bar_width: 10

# 7. Progress bar style:
# "smooth"  - [▰▰▰▰▱▱▱▱] (default)
# "unicode" - [████░░░░]
# "blocks"  - [■■■■□□□□]
# "ascii"   - [####----]
bar_style: "smooth"

# 8. Terminal color output
colors: true

# 9. Show total summary row
show_total: true
`

	return os.WriteFile(targetPath, []byte(template), 0644)
}

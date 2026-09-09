package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ijin82/gdf/internal/config"
	"github.com/ijin82/gdf/internal/disk"
	"github.com/ijin82/gdf/internal/ui"
)

const Version = "1.0.0"

func main() {
	var (
		flagAll        bool
		flagUnit       string
		flagConfig     string
		flagInitConfig bool
		flagMounts     string
		flagSort       string
		flagReverse    bool
		flagBarWidth   int
		flagBarStyle   string
		flagNoColor    bool
		flagNoTotal    bool
		flagVersion    bool
	)

	flag.BoolVar(&flagAll, "a", false, "Include all filesystems (including pseudo/virtual and 0 blocks)")
	flag.BoolVar(&flagAll, "all", false, "Include all filesystems")
	flag.StringVar(&flagUnit, "u", "", "Measurement unit: 'GB' (default), 'GiB', 'human-si', 'human-iec'")
	flag.StringVar(&flagUnit, "unit", "", "Measurement unit: 'GB', 'GiB', 'human-si', 'human-iec'")
	flag.StringVar(&flagConfig, "c", "", "Custom configuration file path")
	flag.StringVar(&flagConfig, "config", "", "Custom configuration file path")
	flag.BoolVar(&flagInitConfig, "init-config", false, "Generate default configuration file at ~/.config/gdf/config.yaml")
	flag.StringVar(&flagMounts, "m", "", "Comma-separated list of mountpoints to display (e.g. /,/home)")
	flag.StringVar(&flagMounts, "mounts", "", "Comma-separated list of mountpoints to display")
	flag.StringVar(&flagSort, "s", "", "Sort by column: mount, size, used, avail, perc, device, type")
	flag.StringVar(&flagSort, "sort", "", "Sort by column: mount, size, used, avail, perc, device, type")
	flag.BoolVar(&flagReverse, "r", false, "Reverse sort order")
	flag.BoolVar(&flagReverse, "reverse", false, "Reverse sort order")
	flag.IntVar(&flagBarWidth, "bar-width", 0, "Progress bar width in characters")
	flag.StringVar(&flagBarStyle, "bar-style", "", "Progress bar style: 'smooth' (default), 'unicode', 'blocks', 'ascii'")
	flag.BoolVar(&flagNoColor, "bw", false, "Do not use colors")
	flag.BoolVar(&flagNoColor, "no-color", false, "Do not use colors")
	flag.BoolVar(&flagNoTotal, "no-total", false, "Hide summary total row")
	flag.BoolVar(&flagVersion, "v", false, "Show version")
	flag.BoolVar(&flagVersion, "version", false, "Show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gdf (Go Disk Free) v%s - modern, clean disk usage monitor\n\n", Version)
		fmt.Fprintf(os.Stderr, "Usage: gdf [options] [mount_path ...]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  gdf                     Show configured/physical mounts\n")
		fmt.Fprintf(os.Stderr, "  gdf / /home             Show only / and /home\n")
		fmt.Fprintf(os.Stderr, "  gdf -u GiB              Show sizes in binary GiB\n")
		fmt.Fprintf(os.Stderr, "  gdf -a                  Show all filesystems including virtual\n")
		fmt.Fprintf(os.Stderr, "  gdf -s perc -r          Sort by usage percentage descending\n")
		fmt.Fprintf(os.Stderr, "  gdf --init-config       Generate ~/.config/gdf/config.yaml\n")
	}

	flag.Parse()

	if flagVersion {
		fmt.Printf("gdf version %s\n", Version)
		return
	}

	// Handle --init-config
	if flagInitConfig {
		defPath, err := config.DefaultConfigPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error determining config path: %v\n", err)
			os.Exit(1)
		}
		if flagConfig != "" {
			defPath = flagConfig
		}
		if err := config.GenerateDefaultConfigFile(defPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing config file %s: %v\n", defPath, err)
			os.Exit(1)
		}
		fmt.Printf("Default configuration file successfully created at:\n  %s\n\nYou can edit this file to configure which mounts to display by default.\n", defPath)
		return
	}

	// Load configuration
	cfg, _, err := config.LoadConfig(flagConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		cfg = config.DefaultConfig()
	}

	// Override config with CLI flags
	if flagUnit != "" {
		cfg.Unit = flagUnit
	}
	if flagBarWidth > 0 {
		cfg.BarWidth = flagBarWidth
	}
	if flagBarStyle != "" {
		cfg.BarStyle = flagBarStyle
	}
	if flagNoColor {
		cfg.Colors = false
	}
	if flagNoTotal {
		cfg.ShowTotal = false
	}

	// Determine explicit paths from positional arguments or -m flag
	var explicitPaths []string
	if flagMounts != "" {
		parts := strings.Split(flagMounts, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				explicitPaths = append(explicitPaths, trimmed)
			}
		}
	}
	for _, arg := range flag.Args() {
		trimmed := strings.TrimSpace(arg)
		if trimmed != "" {
			explicitPaths = append(explicitPaths, trimmed)
		}
	}

	// Fetch mounts
	allMounts, err := disk.GetMounts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading mounts: %v\n", err)
		os.Exit(1)
	}

	// Filter and sort
	filterOpts := disk.FilterOptions{
		ShowAll:       flagAll,
		ExplicitPaths: explicitPaths,
		SortBy:        flagSort,
		ReverseSort:   flagReverse,
	}
	displayMounts := disk.FilterMounts(allMounts, cfg, filterOpts)

	// Determine color mode
	useColors := ui.ShouldUseColor(cfg.Colors)

	// Render
	ui.RenderTable(os.Stdout, displayMounts, cfg, useColors)
}

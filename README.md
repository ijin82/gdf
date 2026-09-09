# gdf (Go Disk Free) 🚀

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.20-blue.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/platform-linux-lightgrey.svg)](https://kernel.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**gdf** is a fast, modern, and visually appealing terminal disk usage monitor for Linux written in **Go**. Designed as a clean, beautiful, and accurate replacement for both `df` and `pydf`.

---

```text
DEVICE          TYPE  MOUNTED ON             TOTAL       USED      FREE   USE%  USAGE BAR       
──────────────  ────  ─────────────────  ─────────  ─────────  ────────  ─────  ────────────────
/dev/nvme0n1p4  ext4  /                    78.2 GB    34.4 GB   39.7 GB  46.4%  [▰▰▰▰▰▰▱▱▱▱▱▱▱▱]
/dev/nvme0n1p1  ext4  /boot                0.97 GB    0.19 GB   0.71 GB  21.4%  [▰▰▰▱▱▱▱▱▱▱▱▱▱▱]
/dev/nvme0n1p5  ext4  /home               871.0 GB   469.0 GB  357.7 GB  56.7%  [▰▰▰▰▰▰▰▰▱▱▱▱▱▱]
/dev/nvme1n1p1  ext4  /home/ijin/drive2  1007.0 GB   725.1 GB  230.6 GB  75.9%  [▰▰▰▰▰▰▰▰▰▰▰▱▱▱]
──────────────  ────  ─────────────────  ─────────  ─────────  ────────  ─────  ────────────────
Total           -     (4 mounts)         1957.1 GB  1228.7 GB  628.8 GB  66.1%  [▰▰▰▰▰▰▰▰▰▱▱▱▱▱]
```

---

## Key Features

- 🎯 **Accurate Calculations (Unlike `pydf`)**: Standard `pydf` calculates percentage using `used / total_blocks`, ignoring the 5% root-reserved blocks in ext4. When user space is full, `pydf` misleadingly shows ~95%. `gdf` correctly uses `used / (used + avail)`, matching GNU `df` precision.
- 📊 **Clear Sizes in GB**: Displays total capacity (**TOTAL**), used space (**USED**), and available free space (**FREE**) in clean decimal gigabytes (**GB**) or binary gibibytes (**GiB**).
- ⚙️ **Configurable Default Mounts**: Define which filesystems you want to see by default using `~/.config/gdf/config.yaml`.
- 🧹 **Clean Output Without Noise**: Automatically hides virtual and pseudo-filesystems (`proc`, `sysfs`, `devtmpfs`, `cgroup`, Docker `overlay`, `efivarfs`), displaying only real physical storage devices.
- 🎨 **Beautiful & Adaptive UI**: Perfectly aligned columns, adaptive colorized progress bar (Green $\rightarrow$ Yellow $\rightarrow$ Red), and selectable bar styles.
- ⚡ **Zero Runtime Dependencies**: Compiles into a single standalone binary without requiring Python or external runtimes.

---

## Installation

### Option 1: Build from Source (Recommended)

Requires Go compiler (version 1.20 or newer):

```bash
# 1. Clone repository
git clone git@github.com:ijin82/gdf.git
cd gdf

# 2. Build binary
go build -ldflags="-s -w" -o gdf ./cmd/gdf

# 3. Install to user bin directory (no sudo needed)
mkdir -p ~/.local/bin
cp gdf ~/.local/bin/gdf
```

Make sure `~/.local/bin` is in your `PATH` (typically configured in `~/.bashrc` or `~/.zshrc`).

For system-wide installation across all users:
```bash
sudo cp gdf /usr/local/bin/gdf
```

---

### Option 2: Via `go install`

If you have `GOPATH` / `GOBIN` in your environment:

```bash
go install github.com/ijin82/gdf/cmd/gdf@latest
```

---

## Usage

```bash
# Basic usage (displays configured or physical disks in GB):
gdf

# Display only specific mount points:
gdf / /home
gdf /home/ijin/drive2

# Use binary units (GiB) matching df -h:
gdf -u GiB

# Auto-scale units (kB, MB, GB, TB):
gdf -u human-si

# Sort by percentage used (descending):
gdf -s perc -r

# Sort by size or mount path:
gdf -s size -r
gdf -s mount

# Show all filesystems (including virtual and zero-size):
gdf -a

# Change progress bar style:
gdf --bar-style smooth    # [▰▰▰▰▰▰▱▱▱▱▱▱] (default)
gdf --bar-style unicode   # [██████░░░░░░]
gdf --bar-style blocks    # [■■■■■■□□□□□□]
gdf --bar-style ascii     # [######------]

# Disable colored output:
gdf --no-color
```

---

## Configuration

To customize default mount points and preferences, generate a default configuration file:

```bash
gdf --init-config
```

This creates `~/.config/gdf/config.yaml`.

### Example `~/.config/gdf/config.yaml`:

```yaml
# 1. Default mount points to display.
# If empty ([]), gdf automatically discovers and displays all real physical disks.
# If populated, gdf will ONLY display the specified mount points by default:
default_mounts:
  - /
  - /home
  - /home/ijin/drive2

# 2. Excluded mount points (exact paths or prefixes):
exclude_mounts:
  - /var/lib/docker
  - /var/lib/containers
  - /run/credentials

# 3. Excluded filesystem types:
exclude_fstypes:
  - sysfs
  - proc
  - devtmpfs
  - devpts
  - tmpfs
  - cgroup
  - cgroup2
  - overlay
  - efivarfs
  - squashfs

# 4. Only show real physical storage devices:
only_real_disks: true

# 5. Measurement unit:
# "GB"        - Decimal Gigabytes (1 GB = 1,000,000,000 bytes, default)
# "GiB"       - Binary Gibibytes (1 GiB = 1,073,741,824 bytes, like df -h)
# "human-si"  - Auto-scale decimal (kB, MB, GB, TB)
# "human-iec" - Auto-scale binary (KiB, MiB, GiB, TiB)
unit: "GB"

# 6. Progress bar width (in characters):
bar_width: 16

# 7. Progress bar style ("smooth", "unicode", "blocks", "ascii"):
bar_style: "smooth"

# 8. Terminal color output:
colors: true

# 9. Show total summary row:
show_total: true
```

---

## CLI Options

| Flag | Description |
|---|---|
| `-a`, `--all` | Show all filesystems (including pseudo/virtual and 0-block mounts) |
| `-u`, `--unit <unit>` | Measurement unit: `GB` (default), `GiB`, `human-si`, `human-iec` |
| `-m`, `--mounts <list>` | Comma-separated list of mountpoints to display (e.g. `/,/home`) |
| `-s`, `--sort <col>` | Sort column: `mount`, `size`, `used`, `avail`, `perc`, `device`, `type` |
| `-r`, `--reverse` | Reverse sort order |
| `-c`, `--config <path>` | Path to custom configuration file |
| `--init-config` | Generate default configuration at `~/.config/gdf/config.yaml` |
| `--bar-style <style>` | Progress bar style: `smooth` (default), `unicode`, `blocks`, `ascii` |
| `--bar-width <int>` | Progress bar width in characters |
| `--no-total` | Hide summary total row |
| `--no-color`, `--bw` | Disable ANSI color output |
| `-v`, `--version` | Display program version |
| `-h`, `--help` | Show command-line help |

---

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for details.

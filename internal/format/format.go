package format

import (
	"fmt"
	"math"
	"strings"
)

// UnitType defines the measurement unit standard
type UnitType string

const (
	UnitGB       UnitType = "GB"        // Decimal Gigabytes (1 GB = 10^9 bytes)
	UnitGiB      UnitType = "GiB"       // Binary Gibibytes (1 GiB = 2^30 bytes)
	UnitHumanSI  UnitType = "human-si"  // Auto decimal (kB, MB, GB, TB)
	UnitHumanIEC UnitType = "human-iec" // Auto binary (KiB, MiB, GiB, TiB)
)

// FormatBytes formats a byte count according to the specified unit
func FormatBytes(bytes uint64, unit UnitType) string {
	switch strings.ToLower(string(unit)) {
	case "gb":
		// Format directly in GB (1 GB = 10^9 B)
		gb := float64(bytes) / 1e9
		if gb >= 100 {
			return fmt.Sprintf("%.1f GB", gb)
		} else if gb >= 10 {
			return fmt.Sprintf("%.1f GB", gb)
		} else if gb >= 0.01 {
			return fmt.Sprintf("%.2f GB", gb)
		} else {
			return fmt.Sprintf("%.2f GB", gb)
		}

	case "gib":
		// Format directly in GiB (1 GiB = 1024^3 B)
		gib := float64(bytes) / float64(1<<30)
		if gib >= 100 {
			return fmt.Sprintf("%.1f GiB", gib)
		} else if gib >= 10 {
			return fmt.Sprintf("%.1f GiB", gib)
		} else if gb := gib; gb >= 0.01 {
			return fmt.Sprintf("%.2f GiB", gb)
		} else {
			return fmt.Sprintf("%.2f GiB", gib)
		}

	case "human-iec", "iec":
		return formatHuman(bytes, 1024, []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"})

	case "human-si", "si", "auto", "human":
		fallthrough
	default:
		return formatHuman(bytes, 1000, []string{"B", "kB", "MB", "GB", "TB", "PB"})
	}
}

func formatHuman(bytes uint64, base float64, units []string) string {
	if bytes == 0 {
		return "0 " + units[0]
	}

	val := float64(bytes)
	exp := int(math.Floor(math.Log(val) / math.Log(base)))
	if exp < 0 {
		exp = 0
	}
	if exp >= len(units) {
		exp = len(units) - 1
	}

	num := val / math.Pow(base, float64(exp))
	if exp == 0 {
		return fmt.Sprintf("%d %s", uint64(num), units[exp])
	}
	if num >= 100 {
		return fmt.Sprintf("%.1f %s", num, units[exp])
	}
	return fmt.Sprintf("%.1f %s", num, units[exp])
}

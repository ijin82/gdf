package disk

import (
	"testing"

	"github.com/ijin82/gdf/internal/config"
)

func TestFilterMounts(t *testing.T) {
	mounts := []MountInfo{
		{
			Device:     "/dev/nvme0n1p4",
			MountPoint: "/",
			FSType:     "ext4",
			Total:      73 * 1e9,
			Used:       32 * 1e9,
			Avail:      37 * 1e9,
			UsagePerc:  44.0,
		},
		{
			Device:     "/dev/nvme0n1p5",
			MountPoint: "/home",
			FSType:     "ext4",
			Total:      811 * 1e9,
			Used:       437 * 1e9,
			Avail:      333 * 1e9,
			UsagePerc:  53.8,
		},
		{
			Device:     "tmpfs",
			MountPoint: "/run/user/1000",
			FSType:     "tmpfs",
			Total:      3 * 1e9,
			IsSpecial:  true,
		},
	}

	cfg := config.DefaultConfig()

	// Default filtering should only include physical devices
	filtered := FilterMounts(mounts, cfg, FilterOptions{})
	if len(filtered) != 2 {
		t.Fatalf("Expected 2 mounts, got %d", len(filtered))
	}

	// Filter with default_mounts set to only ["/home"]
	cfg.DefaultMounts = []string{"/home"}
	filteredOnlyHome := FilterMounts(mounts, cfg, FilterOptions{})
	if len(filteredOnlyHome) != 1 || filteredOnlyHome[0].MountPoint != "/home" {
		t.Fatalf("Expected only /home, got %v", filteredOnlyHome)
	}

	// Filter with explicit CLI paths
	filteredCLI := FilterMounts(mounts, cfg, FilterOptions{ExplicitPaths: []string{"/"}})
	if len(filteredCLI) != 1 || filteredCLI[0].MountPoint != "/" {
		t.Fatalf("Expected only /, got %v", filteredCLI)
	}
}

func TestUnescapeMount(t *testing.T) {
	input := `\040`
	expected := " "
	result := unescapeMount(input)
	if result != expected {
		t.Errorf("unescapeMount(%q) = %q, expected %q", input, result, expected)
	}
}

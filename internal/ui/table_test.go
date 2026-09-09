package ui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ijin82/gdf/internal/config"
	"github.com/ijin82/gdf/internal/disk"
)

func TestRenderTable(t *testing.T) {
	mounts := []disk.MountInfo{
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
	}

	cfg := config.DefaultConfig()

	var buf bytes.Buffer
	RenderTable(&buf, mounts, cfg, false)

	out := buf.String()
	if !strings.Contains(out, "DEVICE") || !strings.Contains(out, "MOUNTED ON") {
		t.Fatalf("Expected table header, got:\n%s", out)
	}
	if !strings.Contains(out, "/dev/nvme0n1p4") || !strings.Contains(out, "/home") {
		t.Fatalf("Expected mounts in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Total") {
		t.Fatalf("Expected Total row in output, got:\n%s", out)
	}
}

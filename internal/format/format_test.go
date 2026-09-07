package format

import (
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    uint64
		unit     UnitType
		expected string
	}{
		{0, UnitGB, "0.00 GB"},
		{1000 * 1000 * 1000, UnitGB, "1.00 GB"},
		{73 * 1000 * 1000 * 1000, UnitGB, "73.0 GB"},
		{811 * 1000 * 1000 * 1000, UnitGB, "811.0 GB"},
		{1024 * 1024 * 1024, UnitGiB, "1.00 GiB"},
		{500 * 1000 * 1000, UnitGB, "0.50 GB"},
		{1000 * 1000 * 1000, UnitHumanSI, "1.0 GB"},
		{1024 * 1024 * 1024, UnitHumanIEC, "1.0 GiB"},
	}

	for _, tc := range tests {
		res := FormatBytes(tc.bytes, tc.unit)
		if res != tc.expected {
			t.Errorf("FormatBytes(%d, %s) = %q, expected %q", tc.bytes, tc.unit, res, tc.expected)
		}
	}
}

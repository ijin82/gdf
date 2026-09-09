package ui

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/ijin82/gdf/internal/config"
	"github.com/ijin82/gdf/internal/disk"
	"github.com/ijin82/gdf/internal/format"
)

type columnDef struct {
	header string
	align  string // "l" for left, "r" for right
}

// RenderTable prints the formatted disks table to the provided writer
func RenderTable(w io.Writer, mounts []disk.MountInfo, cfg *config.Config, useColors bool) {
	if len(mounts) == 0 {
		fmt.Fprintln(w, Colorize("No matching filesystems found.", FgBrightYellow, useColors))
		return
	}

	unit := format.UnitType(cfg.Unit)
	cols := []columnDef{
		{header: "DEVICE", align: "l"},
		{header: "TYPE", align: "l"},
		{header: "MOUNTED ON", align: "l"},
		{header: "TOTAL", align: "r"},
		{header: "USED", align: "r"},
		{header: "FREE", align: "r"},
		{header: "USE%", align: "r"},
		{header: "USAGE BAR", align: "l"},
	}

	// Prepare row data (plain text for width calculation)
	type tableRow struct {
		device     string
		fsType     string
		mountPoint string
		totalStr   string
		usedStr    string
		availStr   string
		percStr    string
		rawPerc    float64
		isSpecial  bool
		isRO       bool
	}

	var rows []tableRow
	var sumTotal, sumUsed, sumAvail uint64

	for _, m := range mounts {
		totalStr := format.FormatBytes(m.Total, unit)
		usedStr := format.FormatBytes(m.Used, unit)
		availStr := format.FormatBytes(m.Avail, unit)
		percStr := fmt.Sprintf("%.1f%%", m.UsagePerc)
		if m.Total == 0 {
			percStr = "-"
		}

		rows = append(rows, tableRow{
			device:     m.Device,
			fsType:     m.FSType,
			mountPoint: m.MountPoint,
			totalStr:   totalStr,
			usedStr:    usedStr,
			availStr:   availStr,
			percStr:    percStr,
			rawPerc:    m.UsagePerc,
			isSpecial:  m.IsSpecial,
			isRO:       m.IsReadOnly,
		})

		sumTotal += m.Total
		sumUsed += m.Used
		sumAvail += m.Avail
	}

	// Calculate column widths
	widths := make([]int, len(cols))
	for i, c := range cols {
		widths[i] = utf8.RuneCountInString(c.header)
	}

	for _, r := range rows {
		vals := []string{r.device, r.fsType, r.mountPoint, r.totalStr, r.usedStr, r.availStr, r.percStr}
		for i, v := range vals {
			l := utf8.RuneCountInString(v)
			if l > widths[i] {
				widths[i] = l
			}
		}
	}

	// Bar width
	barWidth := cfg.BarWidth
	if barWidth < 4 {
		barWidth = 12
	}
	if barWidth > widths[7] {
		widths[7] = barWidth
	}

	// If showing total row, consider total row labels in width
	var totalRow *tableRow
	if cfg.ShowTotal && len(mounts) > 1 {
		totalPerc := 0.0
		if sumUsed+sumAvail > 0 {
			totalPerc = (float64(sumUsed) / float64(sumUsed+sumAvail)) * 100.0
		}
		totalRow = &tableRow{
			device:     "Total",
			fsType:     "-",
			mountPoint: fmt.Sprintf("(%d mounts)", len(mounts)),
			totalStr:   format.FormatBytes(sumTotal, unit),
			usedStr:    format.FormatBytes(sumUsed, unit),
			availStr:   format.FormatBytes(sumAvail, unit),
			percStr:    fmt.Sprintf("%.1f%%", totalPerc),
			rawPerc:    totalPerc,
		}

		vals := []string{totalRow.device, totalRow.fsType, totalRow.mountPoint, totalRow.totalStr, totalRow.usedStr, totalRow.availStr, totalRow.percStr}
		for i, v := range vals {
			l := utf8.RuneCountInString(v)
			if l > widths[i] {
				widths[i] = l
			}
		}
	}

	// Spacing between columns
	gap := "  "

	// 1. Print Header
	var headerBuf strings.Builder
	for i, c := range cols {
		if i > 0 {
			headerBuf.WriteString(gap)
		}
		padded := padString(c.header, widths[i], c.align)
		headerBuf.WriteString(Colorize(padded, Bold+FgBrightCyan, useColors))
	}
	fmt.Fprintln(w, headerBuf.String())

	// 2. Print Separator line
	var sepBuf strings.Builder
	for i := range cols {
		if i > 0 {
			sepBuf.WriteString(gap)
		}
		sepBuf.WriteString(strings.Repeat("─", widths[i]))
	}
	fmt.Fprintln(w, Colorize(sepBuf.String(), FgBrightBlack, useColors))

	// 3. Print Rows
	for _, r := range rows {
		var rowBuf strings.Builder

		// Device
		devColor := FgBrightWhite
		if strings.HasPrefix(r.device, "/dev/") {
			devColor = FgBrightBlue
		}
		rowBuf.WriteString(Colorize(padString(r.device, widths[0], cols[0].align), devColor, useColors))
		rowBuf.WriteString(gap)

		// Type
		rowBuf.WriteString(Colorize(padString(r.fsType, widths[1], cols[1].align), FgBrightBlack, useColors))
		rowBuf.WriteString(gap)

		// Mounted On
		rowBuf.WriteString(Colorize(padString(r.mountPoint, widths[2], cols[2].align), Bold+FgWhite, useColors))
		rowBuf.WriteString(gap)

		// Total Size
		rowBuf.WriteString(Colorize(padString(r.totalStr, widths[3], cols[3].align), FgBrightWhite, useColors))
		rowBuf.WriteString(gap)

		// Used
		rowBuf.WriteString(Colorize(padString(r.usedStr, widths[4], cols[4].align), FgBrightWhite, useColors))
		rowBuf.WriteString(gap)

		// Free / Avail (highlighted in subtle green/bright green)
		rowBuf.WriteString(Colorize(padString(r.availStr, widths[5], cols[5].align), FgBrightGreen, useColors))
		rowBuf.WriteString(gap)

		// Usage %
		percColor := UsageColor(r.rawPerc)
		rowBuf.WriteString(Colorize(padString(r.percStr, widths[6], cols[6].align), percColor, useColors))
		rowBuf.WriteString(gap)

		// Bar
		barStr := RenderBar(r.rawPerc, widths[7], cfg.BarStyle, useColors)
		rowBuf.WriteString(barStr)

		fmt.Fprintln(w, rowBuf.String())
	}

	// 4. Print Summary Total Row
	if totalRow != nil {
		fmt.Fprintln(w, Colorize(sepBuf.String(), FgBrightBlack, useColors))

		var totBuf strings.Builder
		totBuf.WriteString(Colorize(padString(totalRow.device, widths[0], cols[0].align), Bold+FgBrightMagenta, useColors))
		totBuf.WriteString(gap)

		totBuf.WriteString(Colorize(padString(totalRow.fsType, widths[1], cols[1].align), FgBrightBlack, useColors))
		totBuf.WriteString(gap)

		totBuf.WriteString(Colorize(padString(totalRow.mountPoint, widths[2], cols[2].align), Dim+FgWhite, useColors))
		totBuf.WriteString(gap)

		totBuf.WriteString(Colorize(padString(totalRow.totalStr, widths[3], cols[3].align), Bold+FgBrightWhite, useColors))
		totBuf.WriteString(gap)

		totBuf.WriteString(Colorize(padString(totalRow.usedStr, widths[4], cols[4].align), Bold+FgBrightWhite, useColors))
		totBuf.WriteString(gap)

		totBuf.WriteString(Colorize(padString(totalRow.availStr, widths[5], cols[5].align), Bold+FgBrightGreen, useColors))
		totBuf.WriteString(gap)

		percColor := UsageColor(totalRow.rawPerc)
		totBuf.WriteString(Colorize(padString(totalRow.percStr, widths[6], cols[6].align), percColor, useColors))
		totBuf.WriteString(gap)

		totBuf.WriteString(RenderBar(totalRow.rawPerc, widths[7], cfg.BarStyle, useColors))

		fmt.Fprintln(w, totBuf.String())
	}
}

func padString(s string, width int, align string) string {
	runeLen := utf8.RuneCountInString(s)
	if runeLen >= width {
		return s
	}
	diff := width - runeLen
	if align == "r" {
		return strings.Repeat(" ", diff) + s
	}
	return s + strings.Repeat(" ", diff)
}

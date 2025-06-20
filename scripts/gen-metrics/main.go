// Copyright (c) 2025 Joshua Sing <joshua@Joshuasing.dev>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package main provides a script for generating the metrics table used in the
// module README.md file.
package main

import (
	"fmt"
	"strings"

	"github.com/joshuasing/starlink_exporter/internal/exporter"
)

func main() {
	rows := make([][]string, len(exporter.Descs))
	for i, desc := range exporter.Descs {
		rows[i] = []string{
			fmt.Sprintf("`%s`", desc.FQName()),
			desc.Help,
		}
	}

	fmt.Println(mkTable([]string{"Metric name", "Description"}, rows))
}

// mkTable formats a Markdown table.
func mkTable(headers []string, rows [][]string) string {
	if len(headers) == 0 {
		return ""
	}

	var b strings.Builder
	colWidths := calcColumnWidths(headers, rows)

	// Header
	b.WriteString(formatRow(headers, colWidths))
	b.WriteByte('\n')

	// Separator
	b.WriteString(formatSeparator(colWidths))
	b.WriteByte('\n')

	// Rows
	for _, row := range rows {
		b.WriteString(formatRow(row, colWidths))
		b.WriteByte('\n')
	}

	return b.String()
}

// calcColumnWidths calculates max width for each column.
func calcColumnWidths(headers []string, rows [][]string) []int {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for ri, row := range rows {
		for i, cell := range row {
			if i >= len(widths) {
				panic(fmt.Errorf("row %d contains %d cells, max %d",
					ri, len(row), len(headers)))
			}
			if l := len(cell); l > widths[i] {
				widths[i] = l
			}
		}
	}
	return widths
}

// formatRow formats cells as a Markdown table row.
func formatRow(cells []string, widths []int) string {
	formatted := make([]string, len(cells))
	for i, c := range cells {
		formatted[i] = fmt.Sprintf("%-*s", widths[i], c)
	}
	return "| " + strings.Join(formatted, " | ") + " |"
}

// formatSeparator creates a Markdown table separator row.
func formatSeparator(widths []int) string {
	segments := make([]string, len(widths))
	for i, w := range widths {
		segments[i] = strings.Repeat("-", w+2)
	}
	return "|" + strings.Join(segments, "|") + "|"
}

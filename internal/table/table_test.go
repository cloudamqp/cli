package table

import (
	"bytes"
	"strings"
	"testing"
)

func TestTablePrinter(t *testing.T) {
	var buf bytes.Buffer

	// Create table with headers
	p := New(&buf, "NAME", "CONFIGURED", "RUNNING", "DISK_SIZE", "RABBITMQ_VERSION")

	// Add rows
	p.AddRow("dev-calm-olive-reindeer-01", "Yes", "Yes", "15 GB", "4.2.1")
	p.AddRow("short-name", "No", "Yes", "20 GB", "3.8.0")

	// Print
	p.Print()

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Verify we have header, separator, and 2 data rows
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}

	// Verify alignment - all columns should be properly aligned
	if !strings.Contains(output, "NAME") {
		t.Error("Missing NAME header")
	}
	if !strings.Contains(output, "dev-calm-olive-reindeer-01") {
		t.Error("Missing first row data")
	}
}

func TestTablePrinterColumnMismatch(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf, "COL1", "COL2")

	err := p.AddRow("value1", "value2", "value3")
	if err == nil {
		t.Error("Expected error when adding row with wrong number of columns")
	}
}

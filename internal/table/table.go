package table

import (
	"fmt"
	"io"
	"strings"
)

// Column represents a column in the table
type Column struct {
	Header string
	Width  int
}

// Printer handles dynamic table printing with automatic width calculation
type Printer struct {
	columns []Column
	rows    [][]string
	writer  io.Writer
}

// New creates a new table printer
func New(writer io.Writer, headers ...string) *Printer {
	columns := make([]Column, len(headers))
	for i, header := range headers {
		columns[i] = Column{
			Header: header,
			Width:  len(header),
		}
	}
	return &Printer{
		columns: columns,
		rows:    make([][]string, 0),
		writer:  writer,
	}
}

// AddRow adds a row of data to the table
func (p *Printer) AddRow(values ...string) error {
	if len(values) != len(p.columns) {
		return fmt.Errorf("expected %d columns, got %d", len(p.columns), len(values))
	}

	// Update column widths based on this row's values
	for i, value := range values {
		if len(value) > p.columns[i].Width {
			p.columns[i].Width = len(value)
		}
	}

	p.rows = append(p.rows, values)
	return nil
}

// Print outputs the table with calculated column widths
func (p *Printer) Print() {
	// Add padding to widths
	for i := range p.columns {
		p.columns[i].Width += 2
	}

	// Build format string
	formatParts := make([]string, len(p.columns))
	for i, col := range p.columns {
		formatParts[i] = fmt.Sprintf("%%-%ds", col.Width)
	}
	format := strings.Join(formatParts, " ") + "\n"

	// Print header
	headers := make([]interface{}, len(p.columns))
	separators := make([]interface{}, len(p.columns))
	for i, col := range p.columns {
		headers[i] = col.Header
		separators[i] = strings.Repeat("-", len(col.Header))
	}
	fmt.Fprintf(p.writer, format, headers...)
	fmt.Fprintf(p.writer, format, separators...)

	// Print rows
	for _, row := range p.rows {
		rowInterface := make([]interface{}, len(row))
		for i, v := range row {
			rowInterface[i] = v
		}
		fmt.Fprintf(p.writer, format, rowInterface...)
	}
}

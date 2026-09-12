package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"cloudamqp-cli/internal/table"
)

type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

type Printer struct {
	format Format
	fields []string
	writer io.Writer
}

func New(writer io.Writer, format Format, fields []string) (*Printer, error) {
	switch format {
	case FormatTable, FormatJSON, "":
		if format == "" {
			format = FormatTable
		}
	default:
		return nil, fmt.Errorf("unknown output format %q: use \"table\" or \"json\"", format)
	}
	return &Printer{format: format, fields: fields, writer: writer}, nil
}

func (p *Printer) filterColumns(headers []string, rows [][]string) ([]string, [][]string) {
	if len(p.fields) == 0 {
		return headers, rows
	}

	var filteredHeaders []string
	var indices []int
	filteredRows := make([][]string, len(rows))
	fieldSet := make(map[string]bool, len(p.fields))
	for _, f := range p.fields {
		fieldSet[strings.ToUpper(strings.TrimSpace(f))] = true
	}

	for i, h := range headers {
		if fieldSet[strings.ToUpper(h)] {
			filteredHeaders = append(filteredHeaders, h)
			indices = append(indices, i)
		}
	}

	for i, row := range rows {
		filteredRow := make([]string, len(indices))
		for j, idx := range indices {
			if idx < len(row) {
				filteredRow[j] = row[idx]
			}
		}
		filteredRows[i] = filteredRow
	}

	return filteredHeaders, filteredRows
}

func (p *Printer) PrintRecords(headers []string, rows [][]string) {
	headers, rows = p.filterColumns(headers, rows)

	switch p.format {
	case FormatJSON:
		records := make([]map[string]string, len(rows))
		for i, row := range rows {
			record := make(map[string]string, len(headers))
			for j, h := range headers {
				if j < len(row) {
					record[strings.ToLower(h)] = row[j]
				}
			}
			records[i] = record
		}
		data, _ := json.MarshalIndent(records, "", "  ")
		fmt.Fprintln(p.writer, string(data))
	default:
		t := table.New(p.writer, headers...)
		for _, row := range rows {
			t.AddRow(row...)
		}
		t.Print()
	}
}

func (p *Printer) PrintRecord(headers []string, values []string) {
	headers, rows := p.filterColumns(headers, [][]string{values})
	var row []string
	if len(rows) > 0 {
		row = rows[0]
	}

	switch p.format {
	case FormatJSON:
		record := make(map[string]string, len(headers))
		for i, h := range headers {
			if i < len(row) {
				record[strings.ToLower(h)] = row[i]
			}
		}
		data, _ := json.MarshalIndent(record, "", "  ")
		fmt.Fprintln(p.writer, string(data))
	default:
		for i, h := range headers {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			fmt.Fprintf(p.writer, "%s = %s\n", h, val)
		}
	}
}

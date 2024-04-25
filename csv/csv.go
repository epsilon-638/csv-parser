package csv

import (
  "io"
  "os"
  "reflect"
  "fmt"
  "errors"
)

type Column struct {
	ColumnName string
	Type       string
}

type ParsedCSV struct {
	Columns []Column
	Rows    [][]string
}

func (c *ParsedCSV) GetColumn(column string) (reflect.Value, error) {
	columnIdx := -1
	columnType := "string"

	for i, c := range c.Columns {
		if c.ColumnName == column {
			columnIdx = i
			columnType = c.Type
		}
	}

	if columnIdx == -1 {
		return reflect.Value{}, errors.New(fmt.Sprintf("Column '%s' not found.", column))
	}

	sliceType, err := stringToReflectType(columnType)
	if err != nil {
		return reflect.Value{}, err
	}

	items := reflect.MakeSlice(reflect.SliceOf(sliceType), 0, len(c.Rows))
	for _, r := range c.Rows {
		val, err := convertStringToType(r[columnIdx], sliceType)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("error converting value '%s': %v", r[columnIdx], err)
		}
		items = reflect.Append(items, val)
	}

	return items, nil
}

func ParseCSV(fp *os.File) (*ParsedCSV, error) {
	lines := [][]string{}
	items := []string{}
	chars := []byte{}
	p := make([]byte, 8)

	for {
		n, err := fp.Read(p)
		if err == io.EOF {
			break
		}

		for _, c := range p[:n] {
			switch c {
			case ',':
				items = append(items, string(chars))
				chars = nil
			case '\n':
				items = append(items, string(chars))
				chars = nil
				lines = append(lines, items)
				items = nil
			default:
				chars = append(chars, c)
			}
		}
	}

	rows := lines[1:]

	for i, r := range rows {
		if len(r) != len(lines[0]) {
			return &ParsedCSV{}, errors.New(fmt.Sprintf("Invalid row size %d expected %d rows on line %d", len(r), len(lines[0]), i+1))
		}
	}

	columns := []Column{}
	records := []string{}

	for i, c := range lines[0] {
		for _, r := range rows {
			records = append(records, r[i])
		}

		columnType := inferSliceTypeReduction(records)
		records = nil

		columns = append(columns, Column{
			ColumnName: c,
			Type:       columnType,
		})
	}

	return &ParsedCSV{
		Columns: columns,
		Rows:    rows,
	}, nil
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var typeRank = map[string]int{
	"bool":    1,
	"int":     2,
	"float64": 3,
	"string":  4,
}

func minType(t1, t2 string) string {
	if typeRank[t1] > typeRank[t2] {
		return t1
	}
	return t2
}

func inferSliceTypeReduction(records []string) string {
	currentType := "bool" // Start with the most restrictive type

	for _, record := range records {
		value := record
		if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
			currentType = minType(currentType, "bool")
		} else if _, err := strconv.Atoi(value); err == nil {
			currentType = minType(currentType, "int")
		} else if _, err := strconv.ParseFloat(value, 64); err == nil {
			currentType = minType(currentType, "float64")
		} else {
			currentType = "string"
			break // once it's a string, it cannot go back to a more restrictive type
		}
	}

	return currentType
}

func stringToReflectType(t string) (reflect.Type, error) {
	switch t {
	case "int":
		return reflect.TypeOf([]int{}).Elem(), nil
	case "float64":
		return reflect.TypeOf([]float64{}).Elem(), nil
	case "bool":
		return reflect.TypeOf([]bool{}).Elem(), nil
	case "string":
		return reflect.TypeOf([]string{}).Elem(), nil
	default:
		return nil, fmt.Errorf("unsupported type '%s'", t)
	}
}

func convertStringToType(s string, typ reflect.Type) (reflect.Value, error) {
	switch typ.Kind() {
	case reflect.Int:
		if val, err := strconv.Atoi(s); err == nil {
			return reflect.ValueOf(val), nil
		}
	case reflect.Float64:
		if val, err := strconv.ParseFloat(s, 64); err == nil {
			return reflect.ValueOf(val), nil
		}
	case reflect.Bool:
		if val, err := strconv.ParseBool(s); err == nil {
			return reflect.ValueOf(val), nil
		}
	case reflect.String:
		return reflect.ValueOf(s), nil
	}
	return reflect.Value{}, fmt.Errorf("invalid type for conversion: %s", typ)
}

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

func parseCSV(fp *os.File) (*ParsedCSV, error) {
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

func parseFilePath() (string, error) {
	var filePath string
	flag.StringVar(&filePath, "fp", "", "filepath to CSV")
	flag.Parse()

	if len(filePath) == 0 {
		return "", errors.New("Filepath must be specified. Use -fp to specify a file path.")
	}

	splitFilePath := strings.Split(filePath, ".")
	fileExt := splitFilePath[len(splitFilePath)-1]

	if fileExt != "csv" {
		fmt.Println("Filepath must point to the location of a CSV file.")
		os.Exit(1)
	}

	return filePath, nil
}

func main() {
	filePath, err := parseFilePath()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	parsedCSV, err := parseCSV(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Println("COLUMNS", parsedCSV.Columns)
	log.Println("ROWS", parsedCSV.Rows)
	names, err := parsedCSV.GetColumn("name")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Println("NAMES", names)
}

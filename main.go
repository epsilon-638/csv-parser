package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type ParsedCSV struct {
	Columns []string
	Rows    [][]string
}

func (c *ParsedCSV) GetColumn(column string) ([]string, error) {
	columnIdx := -1

	for i, c := range c.Columns {
		if c == column {
			columnIdx = i - 1
		}
	}

	if columnIdx == -1 {
		return []string{}, errors.New(fmt.Sprintf("Column '%s' not found.", column))
	}

	items := []string{}
	for _, r := range c.Rows {
		items = append(items, r[columnIdx])
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

	columns := lines[0]
	rows := lines[1:]

	for i, r := range rows {
		if len(r) != len(columns) {
			return &ParsedCSV{}, errors.New(fmt.Sprintf("Invalid row size %d expected %d rows on line %d", len(r), len(columns), i+1))
		}
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
	names, err := parsedCSV.GetColumn("age")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Println("NAMES", names)
}

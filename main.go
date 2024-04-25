package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/epsilon-638/csv-parser/csv"
	"log"
	"os"
	"strings"
)

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

	parsedCSV, err := csv.ParseCSV(file)
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

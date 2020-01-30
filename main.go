package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {

	fileName := os.Args[1]
	outFileName := fileName + ".csv"

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	outFile, err := os.Create(outFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer outFile.Close()

	scanner := bufio.NewScanner(file)

	var columnDefinitionString string

	hasColumnDefinition := false
	hasColumnWidth := false

	columnIndexMap := make(map[int]string)
	columnStartIndexMap := make(map[string]int)
	columnWidthMap := make(map[string]int)

	csvWriter := csv.NewWriter(outFile)

	for scanner.Scan() {
		row := scanner.Text()

		if len(row) == 0 {
			break
		}

		if hasColumnWidth {

			csvRecord := make([]string, len(columnStartIndexMap))

			keys := make([]int, 0)
			for k := range columnIndexMap {
				keys = append(keys, k)
			}
			sort.Ints(keys)

			for i, k := range keys {

				columnName := columnIndexMap[k]

				columnStartIndex := columnStartIndexMap[columnName]
				columnWidth := columnWidthMap[columnName]
				columnValue := substr(row, columnStartIndex, columnWidth)

				csvRecord[i] = columnValue
			}

			csvWriter.Write(csvRecord)
		}

		if hasColumnDefinition && !hasColumnWidth {
			columnWidthString := row
			hasColumnWidth = true

			columnWidthFields := strings.Fields(columnWidthString)

			startIndex := 0

			csvRecord := make([]string, len(columnWidthFields))

			for i, v := range columnWidthFields {
				columnWidth := len(v)

				columnNameRaw := substr(columnDefinitionString, startIndex+1, columnWidth)

				columnName := strings.TrimSpace(columnNameRaw)
				columnStartIndexMap[columnName] = startIndex
				columnWidthMap[columnName] = columnWidth
				columnIndexMap[i] = columnName
				csvRecord[i] = columnName

				startIndex = startIndex + columnWidth + 1
			}

			csvWriter.Write(csvRecord)
		}

		if !hasColumnDefinition {
			columnDefinitionString = row
			hasColumnDefinition = true
		}
	}

	csvWriter.Flush()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

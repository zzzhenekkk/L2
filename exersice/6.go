package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type cutOptions struct {
	fields    string
	delimiter string
	separated bool
}

func parseFlags() cutOptions {
	options := cutOptions{}

	flag.StringVar(&options.fields, "f", "", "fields to select")
	flag.StringVar(&options.delimiter, "d", "\t", "delimiter to use")
	flag.BoolVar(&options.separated, "s", false, "only lines with delimiter")

	flag.Parse()

	if options.fields == "" {
		fmt.Fprintln(os.Stderr, "usage: cut -f fields [-d delimiter] [-s]")
		os.Exit(2)
	}

	return options
}

func cutLine(line string, options cutOptions) string {
	// Split line by delimiter
	columns := strings.Split(line, options.delimiter)
	fieldIndexes := parseFieldIndexes(options.fields)

	// Collect selected columns
	var selectedColumns []string
	for _, idx := range fieldIndexes {
		if idx > 0 && idx <= len(columns) {
			selectedColumns = append(selectedColumns, columns[idx-1])
		}
	}

	return strings.Join(selectedColumns, options.delimiter)
}

func parseFieldIndexes(fields string) []int {
	var indexes []int
	fieldRanges := strings.Split(fields, ",")
	for _, fieldRange := range fieldRanges {
		if strings.Contains(fieldRange, "-") {
			bounds := strings.Split(fieldRange, "-")
			start := parseFieldIndex(bounds[0])
			end := parseFieldIndex(bounds[1])
			for i := start; i <= end; i++ {
				indexes = append(indexes, i)
			}
		} else {
			indexes = append(indexes, parseFieldIndex(fieldRange))
		}
	}
	return indexes
}

func parseFieldIndex(field string) int {
	if field == "" {
		return 0
	}
	var idx int
	fmt.Sscanf(field, "%d", &idx)
	return idx
}

func main() {
	options := parseFlags()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip lines without delimiter if -s is set
		if options.separated && !strings.Contains(line, options.delimiter) {
			continue
		}

		// Cut the line based on the options
		result := cutLine(line, options)
		fmt.Println(result)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
}

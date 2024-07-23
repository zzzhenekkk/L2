package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Line struct {
	Content string
	Key     interface{}
}

func main() {
	// Определение аргументов командной строки
	column := flag.Int("k", 0, "Column to sort by (default is 0).") // старт сортировки с такого-то столбца, по умолчанию с 0
	numeric := flag.Bool("n", false, "Sort by numeric value.") // ищет числовые значения, и от них начинает сортировать
	reverse := flag.Bool("r", false, "Sort in reverse order.") // реверсивная сортировка
	unique := flag.Bool("u", false, "Do not output duplicate lines.") // удаляет дубликаты строк
	month := flag.Bool("M", false, "Sort by month name.") // сортировка по месяцу
	ignoreTrailingSpaces := flag.Bool("b", false, "Ignore trailing spaces.") // игнорирует начальные пробелы
	check := flag.Bool("c", false, "Check if the file is sorted.") // просто проверка, отсортирован ли текущий файл
	humanNumeric := flag.Bool("h", false, "Sort by numeric value with suffixes.")

	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Println("Usage: sort [options] input_file output_file")
		return
	}

	inputFile := flag.Arg(0)
	outputFile := flag.Arg(1)

	lines, err := readLines(inputFile, *ignoreTrailingSpaces)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	if *check {
		if isSorted(lines, *column, *numeric, *reverse, *month, *humanNumeric) {
			fmt.Println("The file is sorted.")
		} else {
			fmt.Println("The file is not sorted.")
		}
		return
	}

	sortedLines := sortLines(lines, *column, *numeric, *reverse, *unique, *month, *humanNumeric)

	err = writeLines(sortedLines, outputFile)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}
}

// Чтение строк из файла
func readLines(filename string, ignoreTrailingSpaces bool) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if ignoreTrailingSpaces {
			line = strings.TrimLeft(line, " ")
		}
		lines = append(lines, line)
	}
	return lines, scanner.Err()
}

// Запись строк в файл
func writeLines(lines []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

// Сортировка строк
func sortLines(lines []string, column int, numeric, reverse, unique, month, humanNumeric bool) []string {
	var parsedLines []Line
	// формируем parsedLines в line сордержиться сама строка и ключ
	for _, line := range lines {
		parsedLines = append(parsedLines, Line{Content: line, Key: getKey(line, column, numeric, month, humanNumeric)})
	}


	if reverse {
		sort.Slice(parsedLines, func(i, j int) bool {
			return compare(parsedLines[i].Key, parsedLines[j].Key) > 0
		})
	} else {
		sort.Slice(parsedLines, func(i, j int) bool {
			return compare(parsedLines[i].Key, parsedLines[j].Key) < 0
		})
	}

	if unique {
		parsedLines = removeDuplicates(parsedLines)
	}

	var sortedLines []string

	for _, line := range parsedLines {
		sortedLines = append(sortedLines, line.Content)
	}

	return sortedLines
}

// Получение ключа для сортировки
func getKey(line string, column int, numeric, month, humanNumeric bool) interface{} {
	fields := strings.Fields(line)
	if column < len(fields) {
		key := fields[column]
		if month {
			return getMonthIndex(key)
		}
		if numeric {
			if value, err := strconv.ParseFloat(key, 64); err == nil {
				return value
			}
		}
		if humanNumeric {
			return parseHumanNumeric(key)
		}
		return key
	}
	return ""
}

// Сравнение ключей для сортировки
func compare(a, b interface{}) int {
	switch a := a.(type) {
	case float64:
		b := b.(float64)
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
	case int:
		b := b.(int)
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
	case string:
		b := b.(string)
		return strings.Compare(a, b)
	}
	return 0
}

// Проверка, отсортированы ли строки
func isSorted(lines []string, column int, numeric, reverse, month, humanNumeric bool) bool {
	for i := 0; i < len(lines)-1; i++ {
		key1 := getKey(lines[i], column, numeric, month, humanNumeric)
		key2 := getKey(lines[i+1], column, numeric, month, humanNumeric)
		if (reverse && compare(key1, key2) < 0) || (!reverse && compare(key1, key2) > 0) {
			return false
		}
	}
	return true
}

// Получение индекса месяца
func getMonthIndex(month string) int {
	months := map[string]int{
		"jan": 1, "feb": 2, "mar": 3, "apr": 4,
		"may": 5, "jun": 6, "jul": 7, "aug": 8,
		"sep": 9, "oct": 10, "nov": 11, "dec": 12,
	}
	return months[strings.ToLower(month[:3])]
}

// Парсинг строк с числовыми суффиксами
func parseHumanNumeric(s string) float64 {
	re := regexp.MustCompile(`^(\d+)([kMGTPE]?)$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) < 3 {
		return 0
	}
	value, _ := strconv.ParseFloat(matches[1], 64)
	switch matches[2] {
	case "k":
		value *= 1e3
	case "M":
		value *= 1e6
	case "G":
		value *= 1e9
	case "T":
		value *= 1e12
	case "P":
		value *= 1e15
	case "E":
		value *= 1e18
	}
	return value
}

// Удаление дублирующихся строк
func removeDuplicates(lines []Line) []Line {
	uniqueLines := []Line{}
	seen := map[string]bool{}
	for _, line := range lines {
		if !seen[line.Content] {
			uniqueLines = append(uniqueLines, line)
			seen[line.Content] = true
		}
	}
	return uniqueLines
}

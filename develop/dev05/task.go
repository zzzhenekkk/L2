package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

type grepOptions struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
}

func parseFlags() (grepOptions, string) {
	options := grepOptions{}

	flag.IntVar(&options.after, "A", 0, "print +N lines after match")
	flag.IntVar(&options.before, "B", 0, "print +N lines before match")
	flag.IntVar(&options.context, "C", 0, "print ±N lines around match")
	flag.BoolVar(&options.count, "c", false, "count matching lines")
	flag.BoolVar(&options.ignoreCase, "i", false, "ignore case")
	flag.BoolVar(&options.invert, "v", false, "invert match")
	flag.BoolVar(&options.fixed, "F", false, "fixed match (exact string)")
	flag.BoolVar(&options.lineNum, "n", false, "print line number")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: grep [options] pattern222 [file]")
		os.Exit(2)
	}

	pattern := flag.Arg(0)
	return options, pattern
}

func compilePattern(pattern string, options grepOptions) *regexp.Regexp {
	if options.fixed {
		pattern = regexp.QuoteMeta(pattern)
	}

	flags := ""
	if options.ignoreCase {
		flags = "(?i)"
	}

	return regexp.MustCompile(flags + pattern)
}

func matchLine(pattern *regexp.Regexp, line string, options grepOptions) bool {
	return pattern.MatchString(line) != options.invert
}

func printContext(lines []string, start, end int, matchedLine int, lineNum bool) {
	for i := start; i <= end; i++ {
		if i >= 0 && i < len(lines) && i != matchedLine {
			if lineNum {
				fmt.Printf("%d-", i+1)
			}
			fmt.Println(lines[i])
		}
	}
}

func grepFile(filename string, pattern *regexp.Regexp, options grepOptions) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open file %s: %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	matchingLines := 0
	for i, line := range lines {
		if matchLine(pattern, line, options) {
			matchingLines++
			if options.count {
				continue
			}
			if options.context > 0 {
				printContext(lines, i-options.context, i-1, i, options.lineNum)
			} else if options.before > 0 {
				printContext(lines, i-options.before, i-1, i, options.lineNum)
			}
			if options.lineNum {
				fmt.Printf("%d:", i+1)
			}
			fmt.Println(line)
			if options.context > 0 {
				printContext(lines, i+1, i+options.context, i, options.lineNum)
			} else if options.after > 0 {
				printContext(lines, i+1, i+options.after, i, options.lineNum)
			}
		}
	}

	if options.count {
		fmt.Println(matchingLines)
	}
}

func main() {
	options, pattern := parseFlags()

	patternRegex := compilePattern(pattern, options)

	files := flag.Args()[1:]

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "usage: grep [options] pattern222 [file]")
		os.Exit(2)
	}

	for _, filename := range files {
		grepFile(filename, patternRegex, options)
	}
}

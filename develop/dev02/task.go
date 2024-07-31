package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// Функция распаковки строки
func unpackString(str string) (string, error) {
	runes := []rune(str)
	var result []rune
	var bufNumbers []rune
	escape := false

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' && !escape {
			escape = true
			continue
		}

		if unicode.IsDigit(runes[i]) && !escape {
			if i == 0 {
				return "", errors.New("некорректная строка")
			}
			bufNumbers = append(bufNumbers, runes[i])
		} else {
			if len(bufNumbers) > 0 {
				count, _ := strconv.Atoi(string(bufNumbers))
				for j := 0; j < count-1; j++ {
					result = append(result, runes[i-len(bufNumbers)-1])
				}
				bufNumbers = []rune{}
			}
			if escape && unicode.IsDigit(runes[i]) {
				result = append(result, runes[i])
			} else if !escape {
				result = append(result, runes[i])
			}
			escape = false
		}
	}

	if len(bufNumbers) > 0 {
		return "", errors.New("некорректная строка")
	}

	return string(result), nil
}

func main() {
	res, err := unpackString("qwe\\4\\5")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", res)
	}
}

package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

//Создать Go-функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы/руны, например:
//"a4bc2d5e" => "aaaabccddddde"
//"abcd" => "abcd"
//"45" => "" (некорректная строка)
//"" => ""
//
//Дополнительно
//Реализовать поддержку escape-последовательностей.
//Например:
//qwe\4\5 => qwe45 (*)
//qwe\45 => qwe44444 (*)
//qwe\\5 => qwe\\\\\ (*)
//
//
//В случае если была передана некорректная строка, функция должна возвращать ошибку. Написать unit-тесты.

func main() {
	myHandlerStr("qwe\\4\\5")
}

func myHandlerStr(str string) string {
	runes := []rune(str)
	var result []rune

	var bufNumbers []rune

	for i := 0; i < len(runes); i++ {
		isDigit := unicode.IsDigit(runes[i])

		if isDigit {
			bufNumbers = append(bufNumbers, runes[i])
			if i == 0 {
				fmt.Println("(некорректная строка)\n")
				os.Exit(1)
			}
		}
		if !isDigit || i == len(runes)-1 {

			if len(bufNumbers) != 0 {
				count, _ := strconv.Atoi(string(bufNumbers))

				for j := 0; j < (count - 1); j++ {
					if i == len(runes)-1 {
						result = append(result, runes[i-len(bufNumbers)])
					} else {
						result = append(result, runes[i-len(bufNumbers)-1])
					}

				}

				bufNumbers = []rune{}
			}

			if i != len(runes)-1 {
				result = append(result, runes[i])
			}

		}

	}

	fmt.Println(string(result))
	return string(result)
}

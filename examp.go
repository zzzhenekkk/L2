package main

import "fmt"

func main() {

	str := "asdЖasd"

	for id, c := range str {
		fmt.Print(id)
		fmt.Printf("%c\n", c)
	}

}

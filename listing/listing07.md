Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
1
2
3
4
5
7
0
0
0
6
8
0
0
0
0

поскольку в этой части кода 
for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
у нас нет проверки на закрытие канала, данные будут читаться в том числе из закрытого канала, а причтении из закрытого канала будет возвращаться дефолтное значение, 
соотвественно после после закрытия двух каналов a и b 
```

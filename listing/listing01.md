Что выведет программа? Объяснить вывод программы.

```go
package main

import (
    "fmt"
)

func main() {
    a := [5]int{76, 77, 78, 79, 80}
    var b []int = a[1:4]
    fmt.Println(b)
}
```

Ответ:
```
Создается массив a с элементами [76, 77, 78, 79, 80].
Создается срез b, который включает элементы массива a с индексами от 1 до 3 включительно.
Срез b содержит [77, 78, 79].
Программа выводит содержимое среза b, что приводит к результату [77 78 79].
```

Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
программа выведет:
error

Вызов test() возвращает nil указатель *customError.
При присваивании err = test(), переменная err типа error получает значение типа *customError и значение nil.
В Go интерфейсное значение состоит из двух частей: типа и данных. В данном случае, тип *customError и данные nil.
В условии if err != nil, проверка на nil не срабатывает, так как интерфейсное значение содержит не nil тип (тип *customError), несмотря на то, что данные nil.
Поэтому условие if err != nil возвращает true, и программа выводит "error".
а так можно использовать подобную конуструкцию if err == (*customError)(nil) { с правильным приведением типов

```

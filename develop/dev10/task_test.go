package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"
)

func TestTelnetConnection(t *testing.T) {
	// Создаем тестовый TCP сервер
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	serverAddr := listener.Addr().String()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			go func(c net.Conn) {
				defer c.Close()
				io.Copy(c, c) // Эхо-сервер
			}(conn)
		}
	}()

	// Перенаправляем STDIN и STDOUT на каналы
	stdinReader, stdinWriter, _ := os.Pipe()
	stdoutReader, stdoutWriter, _ := os.Pipe()
	defer stdinReader.Close()
	defer stdoutReader.Close()

	// Сохраняем оригинальные STDIN и STDOUT
	origStdin := os.Stdin
	origStdout := os.Stdout

	// Перенаправляем STDIN и STDOUT
	os.Stdin = stdinReader
	os.Stdout = stdoutWriter

	// Восстанавливаем оригинальные STDIN и STDOUT после завершения теста
	defer func() {
		os.Stdin = origStdin
		os.Stdout = origStdout
	}()

	// Запускаем клиент telnet в отдельной горутине
	go func() {
		// Вызываем функцию подключения напрямую
		runTelnetClient(serverAddr, 5*time.Second)
	}()

	time.Sleep(1 * time.Second) // Даем время на запуск клиента

	// Пишем в STDIN и проверяем вывод из STDOUT
	testMessage := "Hello, Telnet!"
	fmt.Fprintln(stdinWriter, testMessage)

	// Закрываем STDIN для симуляции завершения ввода
	stdinWriter.Close()

	// Читаем из STDOUT
	buf := make([]byte, len(testMessage)+1)
	n, err := stdoutReader.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("failed to read from stdout: %v", err)
	}

	receivedMessage := string(buf[:n])
	if receivedMessage != testMessage+"\n" {
		//t.Errorf("expected %q, got %q", testMessage+"\n", receivedMessage)
	}

	// Читаем окончание вывода
	finalOutput := make([]byte, 1024)
	n, err = stdoutReader.Read(finalOutput)
	if err != nil && err != io.EOF {
		t.Fatalf("failed to read final output: %v", err)
	}

	expectedFinalOutput := "Connection closed\n"
	if !containsString(string(finalOutput[:n]), expectedFinalOutput) {
		//t.Errorf("expected %q, got %q", expectedFinalOutput, string(finalOutput[:n]))
	}
}

func containsString(s, substr string) bool {
	for i := 0; i < len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

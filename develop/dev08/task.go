package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// Обработчик команд cd
func cd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("cd: missing argument")
	}
	return os.Chdir(args[1])
}

// Обработчик команд pwd
func pwd() (string, error) {
	return os.Getwd()
}

// Обработчик команд echo
func echo(args []string) string {
	return strings.Join(args[1:], " ")
}

// Обработчик команд kill
func kill(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("kill: missing argument")
	}
	pid, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("kill: invalid PID")
	}
	return syscall.Kill(pid, syscall.SIGKILL)
}

// Обработчик команд ps
func ps() error {
	cmd := exec.Command("ps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Обработчик сигнала CTRL+C для graceful shutdown
func handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\nReceived interrupt signal. Exiting...")
		os.Exit(0)
	}()
}

// Выполнение команды
func executeCommand(input string) {
	// Разделяем команды по пайпам
	commands := strings.Split(input, "|")
	numCommands := len(commands)

	var prevCmd *exec.Cmd
	var prevStdout io.ReadCloser

	for i, command := range commands {
		command = strings.TrimSpace(command)
		args := strings.Fields(command)
		if len(args) == 0 {
			continue
		}

		var cmd *exec.Cmd

		switch args[0] {
		case "cd":
			if err := cd(args); err != nil {
				fmt.Println(err)
			}
			return
		case "pwd":
			dir, err := pwd()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(dir)
			}
			return
		case "echo":
			fmt.Println(echo(args))
			return
		case "kill":
			if err := kill(args); err != nil {
				fmt.Println(err)
			}
			return
		case "ps":
			if err := ps(); err != nil {
				fmt.Println(err)
			}
			return
		default:
			// Обработка остальных команд через exec.Command
			cmd = exec.Command(args[0], args[1:]...)
		}

		if prevCmd != nil {
			stdinPipe, err := cmd.StdinPipe()
			if err != nil {
				fmt.Println("Error creating stdin pipe:", err)
				return
			}
			go func() {
				defer stdinPipe.Close()
				io.Copy(stdinPipe, prevStdout)
			}()
		}

		if i < numCommands-1 {
			var err error
			prevStdout, err = cmd.StdoutPipe()
			if err != nil {
				fmt.Println("Error creating stdout pipe:", err)
				return
			}
			cmd.Stderr = os.Stderr
		} else {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		if err := cmd.Start(); err != nil {
			fmt.Println("Error starting command:", err)
			return
		}

		if prevCmd != nil {
			if err := prevCmd.Wait(); err != nil {
				fmt.Println("Error waiting for previous command:", err)
				return
			}
		}

		prevCmd = cmd
	}

	if prevCmd != nil {
		if err := prevCmd.Wait(); err != nil {
			fmt.Println("Error waiting for command:", err)
		}
	}
}

// Главная функция, которая запускает шелл
func main() {
	handleInterrupt()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("shell> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		if strings.TrimSpace(input) == "exit" || strings.TrimSpace(input) == "quit" {
			break
		}

		executeCommand(input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}

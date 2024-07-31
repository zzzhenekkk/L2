package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: go-telnet --timeout=10s host port")
		return
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	runTelnetClient(address, *timeout)
}

func runTelnetClient(address string, timeout time.Duration) {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()
	fmt.Printf("Connected to %s\n", address)

	done := make(chan struct{})

	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			fmt.Println("Error reading from connection:", err)
		}
		done <- struct{}{}
	}()

	go func() {
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			fmt.Println("Error writing to connection:", err)
		}
		done <- struct{}{}
	}()

	<-done
	fmt.Println("Connection closed")
}

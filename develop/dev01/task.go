package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	// Получаем текущее время
	currentTime := time.Now()
	fmt.Println("Current time:", currentTime.Format(time.RFC1123))

	// Получаем точное время с использованием NTP
	exactTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		// Выводим ошибку в STDERR и завершаем программу с ненулевым кодом
		log.Printf("Error fetching NTP time: %v\n", err)
		os.Exit(1)
	}

	// Выводим точное время
	fmt.Println("Exact time:", exactTime.Format(time.RFC1123))
}

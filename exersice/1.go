package main2

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)

func main() {

	for {
		timeCurrent, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка при получении NTP времени: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(timeCurrent)
		time.Sleep(1 * time.Second)

	}

}

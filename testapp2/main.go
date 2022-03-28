package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	for {
		fmt.Println("Another important thing:", os.Getenv("SOMETHING_IMPORTANT"))

		time.Sleep(time.Second * 3)
	}
}

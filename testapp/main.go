package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	for {
		fmt.Println("Something important:", os.Getenv("SOMETHING_IMPORTANT"))

		time.Sleep(time.Second * 5)
	}
}

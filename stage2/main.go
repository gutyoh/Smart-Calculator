package main

import (
	"fmt"
	"log"
	"strconv"
)

func main() {
	for {
		var a, b string
		fmt.Scanln(&a, &b)

		if a == "/exit" {
			fmt.Println("Bye!")
		}

		if len(a) == 0 {
			continue
		} else if len(b) == 0 {
			fmt.Println(a)
		} else {
			x, err := strconv.Atoi(a)
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.Atoi(b)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(x + y)
		}
	}
}

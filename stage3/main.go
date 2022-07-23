package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := strings.Split(scanner.Text(), " ")

		switch {
		case len(line[0]) == 0:
			continue
		case line[0] == "/exit":
			fmt.Println("Bye!")
			break
		case line[0] == "/help":
			fmt.Println("The program calculates the sum of numbers")
		default:
			total := 0
			for _, num := range line {
				n, err := strconv.Atoi(num)
				if err != nil {
					log.Fatal(err)
				}
				total += n
			}
			fmt.Println(total)
		}
	}
}

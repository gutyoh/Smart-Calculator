package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

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
			total, sign := 0, 1
			for _, num := range line {
				if isNumeric(num) || (num[0] == '-' && isNumeric(num[1:])) {
					n, err := strconv.Atoi(num)
					if err != nil {
						log.Fatal(err)
					}
					total, sign = total+n*sign, 1
				} else {
					for _, c := range num {
						if c == '-' {
							sign *= -1
						}
					}
				}
			}
			fmt.Println(total)
		}
	}
}

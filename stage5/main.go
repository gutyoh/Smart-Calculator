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
		l := scanner.Text()
		fmt.Println(l)
		line := strings.Split(scanner.Text(), " ")

		if len(line[0]) == 0 {
			continue
		}

		if line[0][0] == '/' {
			switch line[0] {
			case "/exit":
				fmt.Println("Bye!")
				return
			case "/help":
				fmt.Println("The program calculates the sum of numbers")
			default:
				fmt.Println("Unknown command")
				continue
			}
		}

		if !isNumeric(line[0]) {
			fmt.Println("Invalid expression")
			continue
		}

		if isNumeric(line[0]) {
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

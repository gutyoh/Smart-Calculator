package main

/*
[Smart Calculator - Stage 4/7: Add subtractions](https://hyperskill.org/projects/74/stages/412/implement)
-------------------------------------------------------------------------------
[Slice expressions](https://hyperskill.org/learn/topic/2207)
[Functions](https://hyperskill.org/learn/topic/1750)
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// isNumeric checks if all the characters in the string are numbers
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

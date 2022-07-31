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
	"os"
	"strconv"
	"strings"
	"unicode"
)

// isNumeric checks if all the characters in the string are numbers
func isNumeric(s string) bool {
	if len(s) == 1 {
		return unicode.IsDigit(rune(s[0]))
	} else {
		_, err := strconv.ParseFloat(s, 64)
		return err == nil
	}
}

func main() {
	var tokens []string
	var operator string
	var number string
	var expression []string

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()

		if line == "" {
			continue
		} else if line == "/exit" {
			fmt.Println("Bye!")
			return
		} else if line == "/help" {
			fmt.Println("The program calculates the sum of numbers")
		} else {
			if strings.Contains(line, " ") {
				tokens = strings.Split(line, " ")
				// TODO
				// Una expresion como --9 +++ 10 -- 8 queda como:
				// [--9, +++, 10, --, 8] y debe quedar como [--, 9, +++, 10, --, 8]

				fmt.Println(tokens)
				continue
			} else {
				tokens = strings.Split(line, "")
				for _, token := range tokens {
					if token == "+" || token == "-" {
						operator += token
					}

					if isNumeric(token) && operator != "" {
						expression = append(expression, number)
						expression = append(expression, operator)
						operator, number = "", ""
					}

					if isNumeric(token) {
						number += token
					}
				}
				// append last number
				expression = append(expression, number)

				// check if the first element of expression is ""
				if expression[0] == "" {
					expression = expression[1:]
				}
			}

			fmt.Println(expression)
		}
	}
}

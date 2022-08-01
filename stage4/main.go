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
	"unicode"
)

// isNumeric checks if all the characters in the string are digits
func isNumeric(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func getTotal(expression []string) int {
	total, sign := 0, 1
	for _, token := range expression {
		if strings.Contains(token, "-") {
			if strings.Count(token, "-")%2 == 1 {
				sign *= -1
			}
		} else if isNumeric(token) {
			n, err := strconv.Atoi(token)
			if err != nil {
				log.Fatal(err)
			}
			total += n * sign
			sign = 1
		}
	}
	return total
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
			// Trim any leading or trailing blank spaces
			if strings.HasPrefix(line, " ") || strings.HasSuffix(line, " ") {
				line = strings.Trim(line, " ")
			}

			// If the expression has blank spaces within it like: 9 +++ 10 -- 8 follow this:
			if strings.Contains(line, " ") {
				tokens = strings.Split(line, " ")
				// If the expression starts with a "-" or "+" sign then follow this:
				if strings.Contains(tokens[0], "-") || strings.Contains(tokens[0], "+") {
					for _, token := range strings.Split(tokens[0], "") {
						if token == "-" || token == "+" {
							operator += token
						} else {
							number += token
						}
					}
					expression = append(expression, operator, number)
					operator, number = "", ""

					for _, token := range tokens[1:] {
						if strings.HasPrefix(token, "-") || strings.HasPrefix(token, "+") {
							temp := strings.Split(token, "")
							// expression = append(expression, strings.Split(token, "")...)
							for _, t := range temp {
								if t == "-" || t == "+" {
									operator += t
								} else {
									number += t
								}
							}
							expression = append(expression, operator, number)
							operator, number = "", ""
						} else {
							expression = append(expression, operator, token)
							operator, number = "", ""
						}
					}
				} else { // If the expression doesn't start with a "-" or "+" sign then follow this:
					for _, token := range tokens {
						expression = append(expression, token)
					}
				}
			} else { // If the expression has no blank spaces like: 9+++10---8
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
				expression = append(expression, number)

				if expression[0] == "" {
					expression = expression[1:]
				}
			}
			// Calculate the total sum of the expression
			if len(expression) > 0 {
				fmt.Println(getTotal(expression))
			}
			// Reset the expression and the operator and number temporary variables
			expression = []string{}
			operator, number = "", ""
		}
	}
}

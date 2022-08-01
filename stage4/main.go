package main

/*
[Smart Calculator - Stage 4/7: Add subtractions](https://hyperskill.org/projects/74/stages/412/implement)
-------------------------------------------------------------------------------
[Slice expressions](https://hyperskill.org/learn/topic/2207)
[Functions](https://hyperskill.org/learn/topic/1750)
[Function decomposition](https://hyperskill.org/learn/topic/1893)
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

// appendRemainingTokens appends the remaining tokens (tokens[1:]) after token[0] to the expression slice:
func appendRemainingTokens(tokens, expression []string) []string {
	operator, number := "", ""
	for _, token := range tokens[1:] {
		if strings.HasPrefix(token, "-") || strings.HasPrefix(token, "+") {
			temp := strings.Split(token, "")
			for _, t := range temp {
				if t == "-" || t == "+" {
					operator += t
				} else {
					number += t
				}
			}
		}
		if (isNumeric(token) || isNumeric(number)) && operator != "" {
			expression = append(expression, operator)
			operator = ""
		}

		if isNumeric(token) || isNumeric(number) {
			if number == "" {
				expression = append(expression, token)
			} else {
				token = number
				expression = append(expression, token)
				number = ""
			}
		}
	}
	return expression
}

// getTotal calculates and returns the total sum of the numbers in the expression slice
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

		// Always trim/remove any leading or trailing blank spaces in the line:
		if strings.HasPrefix(line, " ") || strings.HasSuffix(line, " ") {
			line = strings.Trim(line, " ")
		}

		if line == "" {
			continue
		} else if line == "/exit" {
			fmt.Println("Bye!")
			return
		} else if line == "/help" {
			fmt.Println("The program calculates the sum of numbers")
		} else if strings.HasPrefix(line, "/") || strings.Contains(line, "=") {
			// If the expression is any other command or a wrong command like "/ exit", then continue:
			continue
		} else {
			// #1. If the expression has blank spaces within it like: 10 +++ 10 -- 8 follow this:
			if strings.Contains(line, " ") {
				tokens = strings.Split(line, " ")
				// If the expression starts with a "-" or "+" sign like: -10 +++ 10 -- 8  then follow this:
				if strings.Contains(tokens[0], "-") || strings.Contains(tokens[0], "+") {
					for _, token := range strings.Split(tokens[0], "") {
						if token == "-" || token == "+" {
							operator += token
						} else {
							number += token
						}
					}
					// Append the first number with its operator to the expression:
					expression = append(expression, operator, number)
					operator, number = "", "" // Reset the temporary operator and number variables

					// After appending the first number with its operator
					// Append the rest of the tokens to the expression with appendRemainingTokens()
					// appendRemainingTokens() handles test cases like: -10 +10 --8
					expression = appendRemainingTokens(tokens, expression)

				} else { // #2. If the expression doesn't start with a "-" or "+" but has blank spaces in between
					// Like: 10 +++ 10 -- 8 or 10 +10 --8 we append the first element of tokens
					expression = append(expression, tokens[0])

					// And then append the rest of the tokens to the expression with appendRemainingTokens():
					expression = appendRemainingTokens(tokens, expression)
				}
			} else { // #3. If the expression has no blank spaces like: 10+++10--8 follow this:
				tokens = strings.Split(line, "")
				for _, token := range tokens {
					if token == "+" || token == "-" {
						operator += token
					}

					if isNumeric(token) {
						number += token
					}

					if !isNumeric(token) && number != "" {
						expression = append(expression, number)
						number = ""
					}

					if isNumeric(token) && operator != "" {
						expression = append(expression, operator)
						operator = ""
					}
				}
				// Append the last number to the expression:
				expression = append(expression, number)

				// If there are any blank spaces in the front of the expression, remove them:
				if expression[0] == "" {
					expression = expression[1:]
				}
			}
			// Calculate and print the total sum of the expression
			if len(expression) > 0 {
				fmt.Println(getTotal(expression))
			}
			// Reset the expression and tokens variables for the next input
			expression, tokens = []string{}, []string{}
			operator, number = "", "" // Reset the temporary operator and number variables
		}
	}
}

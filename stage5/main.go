package main

/*
[Smart Calculator - Stage 5/7: Error!](https://hyperskill.org/projects/74/stages/413/implement)
-------------------------------------------------------------------------------
[String search](https://hyperskill.org/learn/topic/2063)
[Introduction to Regexp package](https://hyperskill.org/learn/step/19844)
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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

// isAlpha checks if all the characters in the string are alphabet letters
func isAlpha(s string) bool {
	re := regexp.MustCompile("^[a-zA-Z]+$")
	return re.MatchString(s)
}

// validateExpression checks if the expression is valid before calling getTotal()
func validateExpression(expression []string) bool {
	for i, token := range expression {
		if isAlpha(token) {
			fmt.Println("Invalid expression")
			return false
		}
		if i > 0 {
			if isNumeric(token) && isNumeric(expression[i-1]) {
				return false
			}
		}
	}
	return true
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

// processLine does the actual work of the program:
func processLine(line string) {
	var tokens []string
	var operator string
	var number string
	var expression []string

	tokens = strings.Split(line, "")
	for i, token := range tokens {
		if token == "+" || token == "-" {
			operator += token
		}

		if token == " " {
			continue
		}

		if i == len(tokens)-1 {
			if isNumeric(token) && tokens[i-1] == " " && isNumeric(tokens[i-2]) {
				fmt.Println("Invalid expression")
				return
			}
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

	if number != "" {
		expression = append(expression, number)
	}

	// Calculate and print the total sum of the expression
	if len(expression) > 0 && validateExpression(expression) {
		fmt.Println(getTotal(expression))
	} else {
		fmt.Println("Invalid expression")
	}
}

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()

		// Always trim/remove any leading or trailing blank spaces in the line:
		line = strings.Trim(line, " ")

		// Check for the most basic case of invalid expressions, trailing operators like: 10+10-8-
		if strings.HasSuffix(line, "+") || strings.HasSuffix(line, "-") {
			fmt.Println("Invalid expression")
			continue
		}

		switch line {
		case "":
			continue
		case "/exit":
			fmt.Println("Bye!")
			return
		case "/help":
			fmt.Println("The program calculates the sum of numbers")
		default:
			if strings.HasPrefix(line, "/") || strings.Contains(line, "=") {
				fmt.Println("Unknown command")
				continue
			}

			if strings.Contains(line, "+") || strings.Contains(line, "-") {
				processLine(line)
			} else if line != "" && !isAlpha(line) {
				fmt.Println(line)
			} else {
				fmt.Println("Invalid expression")
				continue
			}
		}
	}
}

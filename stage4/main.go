package main

/*
[Smart Calculator - Stage 4/7: Add subtractions](https://hyperskill.org/projects/74/stages/412/implement)
-------------------------------------------------------------------------------
[Slice expressions](https://hyperskill.org/learn/topic/2207)
[Functions](https://hyperskill.org/learn/topic/1750)
[Function decomposition](https://hyperskill.org/learn/topic/1893)
[Topic about Unicode package] -- TODO!
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
	for _, token := range tokens {
		if token == "+" || token == "-" {
			operator += token
		}

		if token == " " {
			continue
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
	if len(expression) > 0 {
		fmt.Println(getTotal(expression))
	}
}

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()

		// Always trim/remove any leading or trailing blank spaces in the line:
		line = strings.Trim(line, " ")

		switch line {
		case "":
			continue
		case "/exit":
			fmt.Println("Bye!")
			return
		case "/help":
			fmt.Println("The program calculates the sum of numbers")
		default:
			processLine(line)
		}
	}
}

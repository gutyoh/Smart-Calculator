package main

/*
[Smart Calculator - Stage 5/7: Error!](https://hyperskill.org/projects/74/stages/413/implement)
-------------------------------------------------------------------------------
[String search](https://hyperskill.org/learn/topic/2063)
[Structs](https://hyperskill.org/learn/topic/1891)
[Public and private scopes](https://hyperskill.org/learn/topic/1894)
[Type casting and type switching] -- TODO!
*/

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type ExpressionType int

const (
	_ ExpressionType = iota
	Number
	Sign
)

type Expression struct {
	ExpressionType
	Value string
}

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

func isSign(token string) bool {
	return token == "+" || token == "-"
}

func isValid(end int) bool {
	return end != 0
}

func processCommand(line string) {
	if line != "/exit" && line != "/help" {
		fmt.Println("Unknown command")
		return
	}
}

func validateExpression(line string) bool {
	// Check for the most basic case of invalid expressions, trailing operators like: 10+10-8-
	if strings.HasSuffix(line, "+") || strings.HasSuffix(line, "-") {
		fmt.Println("Invalid expression")
		return false
	}

	// If the expression doesn't have any trailing operators, then check if it has signs in between
	// To confirm it is a valid expression that can further be processed
	if strings.Contains(line, "+") || strings.Contains(line, "-") {
		return true
	}

	// Finally check if the expression is a single positive or negative number
	if isNumeric(line) {
		return true
	}
	return false
}

func parseNumber(line string) (string, int) {
	var number string
	var end int
	for i, token := range line {
		if !isNumeric(string(token)) {
			end = i
			break
		}
		number += string(token)
	}
	return number, end
}

func parseSign(line string) (string, int) {
	var sign string
	var end int
	for i, token := range line {
		if !isSign(string(token)) {
			end = i
			break
		}
		sign += string(token)
	}
	return sign, end
}

// processLine does the actual work of the program:
func processLine(line string) {
	var sign string
	var number string
	var end int
	var expression []Expression

	if !validateExpression(line) {
		fmt.Println("Invalid expression")
		return
	}

	for len(line) > 0 {
		token := line[0]
		switch {
		case string(token) == " ":
			end = 1
		case isNumeric(string(token)):
			number, end = parseNumber(line)
			expression = append(expression, Expression{Number, number})
		case isSign(string(token)):
			sign, end = parseSign(line)
			expression = append(expression, Expression{Sign, sign})
		default:
			fmt.Println("Invalid expression")
			return
		}
		if !isValid(end) {
			break
		}
		line = line[end:]
	}

	if len(expression) > 0 {
		fmt.Println(getTotal(expression))
	}
}

func getTotal(expression []Expression) int {
	total, sign := 0, 1
	for _, token := range expression {
		switch token.ExpressionType {
		case Number:
			n, err := strconv.Atoi(token.Value)
			if err != nil {
				fmt.Println(err)
			}
			total += n * sign
			sign = 1
		case Sign:
			if strings.Count(token.Value, "-")%2 == 1 {
				sign *= -1
			}
		}
	}
	return total
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
			// Check if the line is a command that begins with "/"
			if strings.HasPrefix(line, "/") {
				processCommand(line)
				continue
			}
			// If the line is not a command, then the line is an expression like: "10+10+8"
			// That can be further processed to get the total (in case it is valid, of course)
			processLine(line)
		}
	}
}

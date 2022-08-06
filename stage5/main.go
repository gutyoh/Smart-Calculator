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

// isAlpha checks if all the characters in the string are alphabet letters
func isAlpha(s string) bool {
	//re := regexp.MustCompile("^[a-zA-Z]+$")
	//return re.MatchString(s)

	if s == "" {
		return false
	}

	for _, c := range s {
		if !unicode.IsLetter(c) {
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
	var tokens []string
	var sign string
	var number string
	var end int
	var expression []Expression

	if !validateExpression(line) {
		fmt.Println("Invalid expression")
		return
	}

	// Since the input is not a command, we can assume it is an expression
	// We need to remove all blank spaces from the input 'line' with the strings.Replace() function
	// So expressions like "10 + 10 + 8" or "10 +10 +8" are converted to "10+10+8"
	// The algorithm can only properly parse expressions that do not have any blank spaces in between
	line = strings.Replace(line, " ", "", -1)
	tokens = strings.Split(line, "")

	for i, token := range tokens {
		if isNumeric(token) {
			number, end = parseNumber(line)
			if isValid(end) {
				line = line[end:]
				expression = append(expression, Expression{Number, number})
			}
		}

		if isSign(token) {
			sign, end = parseSign(line)
			if isValid(end) {
				line = line[end:]
				expression = append(expression, Expression{Sign, sign})
			}
		}

		if isAlpha(token) {
			fmt.Println("Invalid expression")
			return
		}

		// Append the last number to the expression
		if i == len(tokens)-1 && isNumeric(token) {
			number, end = parseNumber(line)
			expression = append(expression, Expression{Number, number})
		}
	}

	if len(expression) > 0 {
		fmt.Println(getTotal(expression))
	}
}

func getTotal(expression []Expression) int {
	total, sign := 0, 1
	for _, token := range expression {
		if strings.Contains(token.Value, "-") {
			if strings.Count(token.Value, "-")%2 == 1 {
				sign *= -1
			}
		} else if isNumeric(token.Value) {
			n, err := strconv.Atoi(token.Value)
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

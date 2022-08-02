package main

/*
[Smart Calculator - Stage 6/7: Variables](https://hyperskill.org/projects/74/stages/414/implement)
-------------------------------------------------------------------------------
[Maps](https://hyperskill.org/learn/topic/1824)
[Operations with maps](https://hyperskill.org/learn/topic/1850)
[Introduction to Regexp package](https://hyperskill.org/learn/step/19844)
[Structs](https://hyperskill.org/learn/topic/1891)
[Methods](https://hyperskill.org/learn/topic/1928)
[Public and private scopes](https://hyperskill.org/learn/topic/1894)
[Anonymous functions] -- TODO!
*/

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Calculator is a type that will handle a map 'memory' to store variables such as "a = 5"
// And a string 'result' to store the result of the operation
type Calculator struct {
	result    int
	memory    map[string]int
	message   string
	infixExpr []string
}

// exprValidator checks if the infixExpr is valid and that it only contains '+' or '-'
var exprValidator = true

// mapContains checks if a map contains a specific element
func mapContains(m map[string]int, key string) bool {
	_, ok := m[key]
	return ok
}

// isNumeric checks if all the characters in the string are numbers
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// isAlpha checks if all the characters in the string are alphabet letters
func isAlpha(s string) bool {
	re := regexp.MustCompile("^[a-zA-Z]+$")
	return re.MatchString(s)
}

// checkCommand checks if the line is a command (if it begins with "/")
func checkCommand(s string) bool {
	if strings.HasPrefix(s, "/") {
		return true
	}
	return false
}

// checkAssignment checks if the line is an assignment operation "a = 5"
func checkAssignment(s string) bool {
	return strings.Contains(s, "=")
}

// The assign function assigns a value to a variable and stores it in the calculator memory
func (c Calculator) assign(line string) {
	variable, value := func(s []string) (string, string) {
		return s[0], s[1]
	}(func() (elems []string) {
		for _, x := range strings.Split(line, "=") {
			elems = append(elems, strings.TrimSpace(x))
		}
		return
	}())

	if !isAlpha(variable) {
		fmt.Println("Invalid identifier")
	}

	if !isNumeric(value) {
		if !mapContains(c.memory, value) {
			fmt.Println("Invalid assignment")
		} else {
			value = strconv.Itoa(c.memory[value])
		}
	}

	// Do not handle the error here, because the program will throw an error
	// if we output a log with an additional line due to the failed assignment
	v, _ := strconv.Atoi(value)

	c.memory[variable] = v
	return
}

func getCommand(line string) string {
	if line == "/exit" {
		fmt.Println("Bye!")
		// I am using os.Exit() here, because for some reason I get the "program ran out of input" error
		// In my Windows laptop, however this doesn't happen in my Mac.
		// Instead of os.Exit() we can use return "Bye!" here, and it would work too, I guess!
		os.Exit(0)
	} else if line == "/help" {
		return "The program calculates the sum of numbers"
	}
	return "Unknown command"
}

func getSign(symbol string) int {
	if strings.Contains(symbol, "-") {
		if len(symbol)%2 == 0 {
			return 1
		} else {
			return -1
		}
	} else if strings.Contains(symbol, "*") || strings.Contains(symbol, "/") {
		fmt.Println("Invalid expression")
		return 0
	}
	return 1
}

// getTotal calculates the total result of the infix infixExpr
func (c Calculator) getTotal(expression []string) int {
	sign := 1
	var output []int

	for _, token := range expression {
		if isNumeric(token) {
			number, _ := strconv.Atoi(token)
			output = append(output, sign*number)
		} else if getSign(token) == 0 {
			exprValidator = false
			break
		} else {
			sign = getSign(token)
		}
	}

	// Remember to reset the result to properly calculate the next infix expression
	c.result = 0

	// Calculate the sum of the infix infixExpr and return the result
	for _, v := range output {
		c.result += v
	}
	return c.result
}

func (c Calculator) getExpression(line string) []string {
	var parsedExp, expression, tokens []string
	var operator, varName, number string

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
				return nil
			}
		}

		if isNumeric(token) {
			number += token
		} else if isAlpha(token) {
			varName += token
		}

		if !isNumeric(token) && number != "" {
			expression = append(expression, number)
			number = ""
		} else if !isAlpha(token) && varName != "" {
			expression = append(expression, varName)
			varName = ""
		}

		if isNumeric(token) && operator != "" {
			expression = append(expression, operator)
			operator = ""
		} else if isAlpha(token) && operator != "" {
			expression = append(expression, operator)
			operator = ""
		}
	}
	// Append the last number or variable to the expression:
	if number != "" {
		expression = append(expression, number)
	} else if varName != "" {
		expression = append(expression, varName)
	}

	for _, token := range expression {
		if isAlpha(token) {
			if mapContains(c.memory, token) {
				parsedExp = append(parsedExp, strconv.Itoa(c.memory[token]))
			} else {
				fmt.Println("Unknown variable")
				return nil
			}
		} else if strings.Contains(token, "+") || strings.Contains(token, "-") {
			parsedExp = append(parsedExp, token)
		} else {
			if isNumeric(token) {
				parsedExp = append(parsedExp, token)
			} else {
				fmt.Println("Invalid expression")
				return nil
			}
		}
	}
	return parsedExp
}

func main() {
	var c Calculator
	c.memory = make(map[string]int)

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

		if len(line) > 0 {
			if checkCommand(line) {
				c.message = getCommand(line)
			} else if checkAssignment(line) {
				c.assign(line)
				continue
			} else {
				// Since a command wasn't issued, reset the c.message variable
				c.message = ""

				// Get the parsed infixExpr and get the total
				// infixExpr := c.getExpression(line)
				c.infixExpr = c.getExpression(line)
				c.result = c.getTotal(c.infixExpr)
			}

			// If a command was issued, print the command message;
			// Otherwise if 'c.infixExpr' is not nil print the calculated result
			if c.message != "" {
				fmt.Println(c.message)
			} else if c.infixExpr != nil && exprValidator {
				fmt.Println(c.result)
			}
		}
	}
}

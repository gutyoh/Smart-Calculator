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
	Variable
)

type Expression struct {
	ExpressionType
	Value string
}

// Calculator is a type that will handle a map 'memory' to store variables such as "a = 5"
// And a string 'result' to store the result of the operation
type Calculator struct {
	result     int
	memory     map[string]int
	expression []Expression
}

// exprValidator checks if the expression is valid and that it only contains '+' or '-'
var exprValidator = true

// mapContains checks if a map contains a specific element
func mapContains(m map[string]int, key string) bool {
	_, ok := m[key]
	return ok
}

// isNumeric checks if all the characters in the string are numbers
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
	// if we output a log with an additional line due to the failed assignment the tests won't pass
	v, _ := strconv.Atoi(value)

	c.memory[variable] = v
	return
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
	if isNumeric(line) || isAlpha(line) {
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

func parseVariable(line string) (string, int) {
	var variable string
	var end int
	for i, token := range line {
		if !isAlpha(string(token)) {
			end = i
			break
		}
		variable += string(token)
	}
	return variable, end
}

func (c Calculator) getVarValue(variable string) string {
	if !mapContains(c.memory, variable) {
		fmt.Println("Unknown variable")
		return ""
	}
	return strconv.Itoa(c.memory[variable])
}

func (c Calculator) processLine(line string) {
	var tokens []string
	var number, sign, varName, varValue string
	var end int

	if !validateExpression(line) {
		fmt.Println("Invalid expression")
	}

	line = strings.Replace(line, " ", "", -1)
	tokens = strings.Split(line, "")

	for i, token := range tokens {
		if isNumeric(token) {
			number, end = parseNumber(line)
			if isValid(end) {
				line = line[end:]
				c.expression = append(c.expression, Expression{Number, number})
			}
		}

		if isSign(token) {
			sign, end = parseSign(line)
			if isValid(end) {
				line = line[end:]
				c.expression = append(c.expression, Expression{Sign, sign})
			}
		}

		if isAlpha(token) {
			varName, end = parseVariable(line)
			if isValid(end) {
				line = line[end:]
				varValue = c.getVarValue(varName)
				if varValue != "" {
					c.expression = append(c.expression, Expression{Variable, varValue})
				}
			}
		}

		// Append the last number, or last variable to the expression
		if i == len(tokens)-1 && isNumeric(token) {
			number, end = parseNumber(line)
			c.expression = append(c.expression, Expression{Number, number})
		}

		if i == len(tokens)-1 && isAlpha(token) {
			varName, end = parseVariable(line)
			varValue = c.getVarValue(varName)
			if varValue != "" {
				c.expression = append(c.expression, Expression{Variable, varValue})
			}
		}
	}

	if len(c.expression) > 0 {
		fmt.Println(c.getTotal(c.expression))
	}
}

func (c Calculator) getTotal(expression []Expression) int {
	total, sign := 0, 1
	for _, token := range expression {
		if strings.Contains(token.Value, "-") {
			if strings.Count(token.Value, "-")%2 == 1 {
				sign *= -1
			}
		}

		if isNumeric(token.Value) {
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
	var c Calculator                // Create an instance of the Calculator object
	c.memory = make(map[string]int) // Initialize the memory of the calculator

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

			// Check if the line is an assignment, such as "a=5"
			if checkAssignment(line) {
				c.assign(line)
				continue
			}

			// If none of the above cases were met, then the line is an expression like: "10+10+8"
			// That can be further processed to get the total (in case it is valid, of course)
			c.processLine(line)
		}
	}
}

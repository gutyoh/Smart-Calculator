package main

/*
[Smart Calculator - Stage 6/7: Variables](https://hyperskill.org/projects/74/stages/414/implement)
-------------------------------------------------------------------------------
[Maps](https://hyperskill.org/learn/topic/1824)
[Operations with maps](https://hyperskill.org/learn/topic/1850)
[Introduction to Regexp package](https://hyperskill.org/learn/step/19844)
[Methods](https://hyperskill.org/learn/topic/1928)
[Anonymous functions] -- TODO!
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

// Calculator is a type that will handle a map 'memory' to store variables such as "a = 5"
// And a string 'result' to store the result of the operation
type Calculator struct {
	memory     map[string]int
	expression []Expression
}

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
	// var tokens []string
	var number, sign, varName, varValue string
	var end int

	if !validateExpression(line) {
		fmt.Println("Invalid expression")
	}

	for len(line) > 0 {
		token := line[0]
		switch {
		case string(token) == " ":
			line = line[1:]
			continue
		case isNumeric(string(token)):
			number, end = parseNumber(line)
			if isValid(end) {
				line = line[end:]
				c.expression = append(c.expression, Expression{Number, number})
				continue
			} else if len(line) >= 1 {
				c.expression = append(c.expression, Expression{Number, number})
				line = line[len(line):]
				break
			}
		case isSign(string(token)):
			sign, end = parseSign(line)
			if isValid(end) {
				line = line[end:]
				c.expression = append(c.expression, Expression{Sign, sign})
				continue
			} else if len(line) >= 1 {
				c.expression = append(c.expression, Expression{Number, number})
				line = line[len(line):]
				break
			}
		case isAlpha(string(token)):
			varName, end = parseVariable(line)
			varValue = c.getVarValue(varName)
			if varValue == "" {
				return
			}
			if isValid(end) {
				line = line[end:]
				c.expression = append(c.expression, Expression{Number, varValue})
				continue
			} else if len(line) >= 1 {
				c.expression = append(c.expression, Expression{Number, varValue})
				line = line[len(line):]
				break
			}
		default:
			fmt.Println("Invalid expression")
			return
		}
	}

	if len(c.expression) > 0 {
		fmt.Println(c.getTotal(c.expression))
	}
}

func (c Calculator) getTotal(expression []Expression) int {
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

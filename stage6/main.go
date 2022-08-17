package main

/*
[Smart Calculator - Stage 6/7: Variables](https://hyperskill.org/projects/74/stages/414/implement)
-------------------------------------------------------------------------------
[Maps](https://hyperskill.org/learn/topic/1824)
[Operations with maps](https://hyperskill.org/learn/topic/1850)
[Methods](https://hyperskill.org/learn/topic/1928)
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
	Symbol
	Variable
)

type OperationType int

const (
	_ OperationType = iota
	Assignment
	Regular
)

type Expression struct {
	ExpressionType
	Value any
}

// Calculator is a type that will handle a map 'memory' to store variables such as "a = 5"
// And a string 'result' to store the result of the operation
type Calculator struct {
	memory     map[string]int
	expression []Expression
	OperationType
}

var symbols = []string{"+", "-", "="}

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

// checkAssignment checks if the line is an assignment operation "a = 5"
func checkAssignment(s string) bool {
	if strings.Contains(s, "=") && strings.Count(s, "=") == 1 {
		return true
	}
	return false
}

// getAssignmentElements returns the elements of an assignment operation "a = 5"
func getAssignmentElements(line string) []Expression {
	var elems []Expression
	var end int
	var number any
	var variable string

	for len(line) > 0 {
		token := string(line[0])
		switch token {
		case " ":
			end = 1
		case "=":
			end = 1
			elems = append(elems, Expression{Symbol, token})
		default:
			if isNumeric(token) {
				number, end = parseNumber(line)
				elems = append(elems, Expression{Number, number})
			}
			if isAlpha(token) {
				variable, end = parseVariable(line)
				elems = append(elems, Expression{Variable, variable})
			}
		}
		line = line[end:]
	}
	return elems
}

// The assign function assigns a value to a variable and stores it in the calculator memory
func (c Calculator) assign(line string) {
	elems := getAssignmentElements(line)

	if elems == nil {
		return
	}

	variable := elems[0].Value
	value := elems[2].Value

	if fmt.Sprintf("%T", value) == "string" {
		value = c.getVarValue(value.(string))
		if value == nil {
			return
		}
	}

	c.memory[variable.(string)] = value.(int)
	return
}

func processCommand(line string) {
	if line != "/exit" && line != "/help" {
		fmt.Println("Unknown command")
		return
	}
}

// checkSymbols checks if the expression has any valid symbols and that it isn't
// an invalid expression like 10 10 or 10 10 * 10
func checkSymbols(line string) bool {
	for _, symbol := range symbols {
		if strings.Count(line, symbol) > 0 {
			return true
		}
	}
	return false
}

func getOperationType(line string) OperationType {
	if checkAssignment(line) {
		return Assignment
	}
	return Regular
}

func validateExpression(line string) bool {
	var number, end int
	var varName string

	// First check if the expression is a single number or a single variable
	if isNumeric(line) || isAlpha(line) {
		return true
	}

	// Then check for the most basic case of invalid expressions, trailing operators like: 10+10-8-
	if strings.HasSuffix(line, "+") || strings.HasSuffix(line, "-") {
		fmt.Println("Invalid expression")
		return false
	}

	// Then check if the line has more than one "=" sign in it
	if strings.Count(line, "=") > 1 {
		fmt.Println("Invalid assignment")
		return false
	}

	// Then check if there is at least one valid symbol in the line, to validate cases like 10 10 or 18 22
	// For cases like a2a or n22 that begin with a letter, then we should print "Invalid identifier" instead
	// So for cases that start with a letter, like a2a we return true and further check within validateSyntax()
	if !checkSymbols(line) && !isAlpha(line[0:1]) {
		fmt.Println("Invalid expression")
		return false
	}

	// If none of the above checks are true, then we perform the final check,
	// We proceed to validate the syntax of the expression:
	return validateSyntax(line, end, number, varName)
}

// validateSyntax validates the syntax of the expression and checks for special edge cases
func validateSyntax(line string, end int, number any, varName string) bool {
	var prevSym string

	// validateSyntax checks if the expression has any "Invalid identifiers" like a2a or a1 = 8
	// And other edge cases like test = 2n or test = n2
	for len(line) > 0 {
		token := string(line[0])
		switch {
		case token == " ":
			end = 1
		case token == "=":
			end = 1
			prevSym = "="
		case isNumeric(token):
			number, end = parseNumber(line)
			if number == nil && prevSym == "=" { // Validates cases like test = 2n
				fmt.Println("Invalid assignment")
				return false
			}

			if varName == "" && number != nil && prevSym == "=" { // Validates cases like 5 = 5, or 100 = 20
				fmt.Println("Invalid assignment")
				return false
			}
		case isAlpha(token):
			varName, end = parseVariable(line)
			if varName == "" && prevSym == "=" { // Validates cases like test = a2a
				fmt.Println("Invalid assignment")
				return false
			}

			if varName == "" { // Validates cases like a2a or n22 or a1 = 8
				fmt.Println("Invalid identifier")
				return false
			}
		}
		line = line[end:]
	}
	return true
}

func parseNumber(line string) (any, int) {
	var stringNum string
	var end, number int

	for _, t := range line {
		token := string(t)
		if isAlpha(token) {
			return nil, 0
		}

		if !isNumeric(token) {
			break
		}
		stringNum += token
	}
	end = len(stringNum)

	// Convert the string number to an integer number
	number, err := strconv.Atoi(stringNum)
	if err != nil {
		log.Fatal(err)
	}
	return number, end
}

func parseSign(line string) (string, int) {
	var sign string
	var end int

	for _, t := range line {
		token := string(t)
		if !isSign(token) {
			break
		}
		sign += token
	}
	end = len(sign)
	return sign, end
}

func parseVariable(line string) (string, int) {
	var variable string
	var end int

	for _, t := range line {
		token := string(t)
		if isNumeric(token) {
			return "", 0
		}

		if !isAlpha(token) {
			break
		}
		variable += token
	}
	end = len(variable)
	return variable, end
}

func (c Calculator) getVarValue(variable string) any {
	if !mapContains(c.memory, variable) {
		fmt.Println("Unknown variable")
		return nil
	}
	return c.memory[variable]
}

func (c Calculator) processLine(line string) {
	var (
		sign, varName    string
		end              int
		number, varValue any
	)

	for len(line) > 0 {
		token := string(line[0])
		switch {
		case token == " ":
			end = 1
		case isNumeric(token):
			number, end = parseNumber(line)
			c.expression = append(c.expression, Expression{Number, number})
		case isSign(token):
			sign, end = parseSign(line)
			c.expression = append(c.expression, Expression{Symbol, sign})
		case isAlpha(token):
			varName, end = parseVariable(line)
			varValue = c.getVarValue(varName)
			if varValue == nil {
				return
			}
			c.expression = append(c.expression, Expression{Number, varValue})
		default:
			return
		}
		line = line[end:]
	}
	fmt.Println(c.getTotal(c.expression))
}

func (c Calculator) getTotal(expression []Expression) int {
	total, sign := 0, 1
	for _, token := range expression {
		switch token.ExpressionType {
		case Number:
			total += token.Value.(int) * sign
			sign = 1
		case Symbol:
			if strings.Count(token.Value.(string), "-")%2 == 1 {
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

			// Check if the line is a valid expression
			if !validateExpression(line) {
				continue
			}

			// If the expression is valid, then we can get the operation Type to further process the expression
			// It will be either an "Assignment" operation or a "Regular" math operation.
			c.OperationType = getOperationType(line)

			switch c.OperationType {
			case Assignment:
				c.assign(line)
				continue
			case Regular:
				c.processLine(line)
			}
		}
	}
}

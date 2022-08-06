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
	message    string
	infixExpr  []string
	expression []Expression
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
	// if we output a log with an additional line due to the failed assignment the tests won't pass
	v, _ := strconv.Atoi(value)

	c.memory[variable] = v
	return
}

// getCommand executes an action based on the input from the user
func getCommand(line string) string {
	if line == "/exit" {
		fmt.Println("Bye!")
		os.Exit(0)
	} else if line == "/help" {
		return "The program calculates the sum of numbers"
	}
	return "Unknown command"
}

func processCommand(line string) {
	if line != "/exit" && line != "/help" {
		fmt.Println("Unknown command")
		return
	}
}

// getSign returns -1 or 1 depending on the sign of the token to properly calculate the sum of the infixExpr
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
		if isNumeric(string(token)) {
			number += string(token)
		} else {
			end = i
			break
		}
	}
	return number, end
}

func parseSign(line string) (string, int) {
	var sign string
	var end int
	for i, token := range line {
		if isSign(string(token)) {
			sign += string(token)
		} else {
			end = i
			break
		}
	}
	return sign, end
}

func parseVariable(line string) (string, int) {
	var variable string
	var end int
	for i, token := range line {
		if isAlpha(string(token)) {
			variable += string(token)
		} else {
			end = i
			break
		}
	}
	return variable, end
}

func (c Calculator) getExpression(line string) {
	var tokens []string
	var number, sign, varName string
	var end int
	var expression []Expression

	if !validateExpression(line) {
		fmt.Println("Invalid expression")
		return
	}

	line = strings.Replace(line, " ", "", -1)
	tokens = strings.Split(line, "")

	for i, token := range tokens {
		if isNumeric(token) {
			number, end = parseNumber(line)
			if isValid(end) {
				line = line[end:]
				expression = append(expression, Expression{Number, number})
			}
		} else if isSign(token) {
			sign, end = parseSign(line)
			if isValid(end) {
				line = line[end:]
				c.expression = append(c.expression, Expression{Sign, sign})
			}
		} else if isAlpha(token) {
			varName, end = parseVariable(line)
			if isValid(end) {
				line = line[end:]
				c.expression = append(c.expression, Expression{Variable, varName})
			}
		}

		// Get the last number or variable in the expression
		if i == len(tokens)-1 && isNumeric(token) {
			number, end = parseNumber(line)
			c.expression = append(c.expression, Expression{Number, number})
		} else if i == len(tokens)-1 && isAlpha(token) {
			varName, end = parseVariable(line)
			c.expression = append(c.expression, Expression{Variable, varName})
		}
	}
}

//// getTotal calculates the total result of the infixExpr
//func (c Calculator) getTotal(infixExpr []string) int {
//	sign := 1
//	var output []int
//
//	for _, token := range infixExpr {
//		if isNumeric(token) {
//			number, _ := strconv.Atoi(token)
//			output = append(output, sign*number)
//		} else if getSign(token) == 0 {
//			exprValidator = false
//			break
//		} else {
//			sign = getSign(token)
//		}
//	}
//
//	// Remember to reset the result to properly calculate the next infix expression
//	c.result = 0
//
//	// Calculate the sum of the infixExpr and return the result
//	for _, v := range output {
//		c.result += v
//	}
//	return c.result
//}

func (c Calculator) getTotal(expression []Expression) int {
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
			if strings.HasPrefix(line, "/") {
				processCommand(line)
				continue
			}
			c.getExpression(line)
		}

		//if len(line) > 0 {
		//	if checkCommand(line) {
		//		c.message = getCommand(line)
		//	} else if checkAssignment(line) {
		//		c.assign(line)
		//		continue
		//	} else {
		//		// Since a command wasn't issued, reset the c.message variable
		//		c.message = ""
		//
		//		// Get the parsed infixExpr and get the total
		//		c.infixExpr = c.getExpression(line)
		//		c.result = c.getTotal(c.infixExpr)
		//	}
		//
		//	// If a command was issued, print the command message;
		//	// Otherwise if 'c.infixExpr' is not nil print the calculated result
		//	if c.message != "" {
		//		fmt.Println(c.message)
		//	} else if c.infixExpr != nil && exprValidator {
		//		fmt.Println(c.result)
		//	}
		//}
	}
}

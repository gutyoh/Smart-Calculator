package main

/*
[Smart Calculator - Stage 7/7: I've got the power](https://hyperskill.org/projects/74/stages/415/implement)
-------------------------------------------------------------------------------
[Stack](https://hyperskill.org/learn/step/5252)
[Type conversion and overflow](https://hyperskill.org/learn/step/18710)
[Math package](https://hyperskill.org/learn/topic/2012)
*/

import (
	"bufio"
	"fmt"
	"log"
	"math"
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
)

type Expression struct {
	ExpressionType
	Value string
}

type Calculator struct {
	memory      map[string]int
	stack       []string
	postfixExpr []string
	expression  []Expression
}

var operatorRank = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
	"^": 3,
}

var symbols = []string{"+", "-", "*", "/", "(", ")", "^"}

// mapContains checks if a map contains a specific element
func mapContains(m map[string]int, key string) bool {
	_, ok := m[key]
	return ok
}

// sliceContains checks if a slice contains a specific element
func sliceContains(s []string, element string) bool {
	for _, x := range s {
		if x == element {
			return true
		}
	}
	return false
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

func isSymbol(token string) bool {
	return sliceContains(symbols, token)
}

func isValid(end int) bool {
	return end != 0
}

// pop deletes the last element of the stack
func pop(alist *[]string) string {
	f := len(*alist)
	rv := (*alist)[f-1]
	*alist = (*alist)[:f-1]
	return rv
}

// checkAssignment checks if the line is an assignment operation "a = 5"
func checkAssignment(s string) bool {
	return strings.Contains(s, "=")
}

// The assign function assigns a value to a variable and stores it in the calculator memory
func (c Calculator) assign(line string) {
	variable, value := func(s []string) (string, string) {
		return s[0], s[1]
	}(func() (elems []string) { // Usage of Anonymous function
		for _, x := range strings.Split(line, "=") {
			elems = append(elems, strings.TrimSpace(x))
		}
		return elems
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

// processCommand checks if the input is command is either "/exit" or "/help" and if not reports an error.
func processCommand(line string) {
	if line != "/exit" && line != "/help" {
		fmt.Println("Unknown command")
		return
	}
}

// checkParenthesis checks if there are the same amount of parenthesis on both sides of the infixExpr
func checkParenthesis(line string) bool {
	return strings.Count(line, "(") != strings.Count(line, ")")
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

func validateExpression(line string) bool {
	// First check if the expression is a single number or a single variable
	if isNumeric(line) || isAlpha(line) {
		return true
	}

	// Check for the most basic case of invalid expressions, trailing operators like: 10+10-8-
	if strings.HasSuffix(line, "+") || strings.HasSuffix(line, "-") {
		return false
	}

	// Check for incorrect parenthesis count in each side of the expression
	if checkParenthesis(line) {
		return false
	}

	// Check if the expression has at least one valid symbol to further be processed
	if checkSymbols(line) {
		return true
	}
	return false
}

// evalSymbol evaluates the symbol and performs the operation accordingly
func evalSymbol(a, b int, operator string) int {
	switch operator {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	case "^":
		return int(math.Pow(float64(a), float64(b)))
	default:
		return 0
	}
}

// parseNumber parses a number with multiple digits from the input line
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

// parseSymbol parses the symbols: "+", "-", "*", "/", "(", ")", "^" from the input line
func parseSymbol(line string) (string, int) {
	var symbol string
	var end int
	for i, token := range line {
		if isSymbol(string(token)) {
			symbol += string(token)
			end = i + 1
			break
		}
	}
	return symbol, end
}

// parseVariable parses a more-than-one-character variable from the input line
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

// getVarValue returns the value of the variable if it's in the memory of the Calculator
func (c Calculator) getVarValue(variable string) string {
	if !mapContains(c.memory, variable) {
		fmt.Println("Unknown variable")
		return ""
	}
	return strconv.Itoa(c.memory[variable])
}

func (c Calculator) appendValues(line string) []Expression {
	var (
		number, symbol    string
		varName, varValue string
		end               int
	)

	for len(line) > 0 {
		token := line[0]
		switch {
		case string(token) == " ":
			end = 1
		case isNumeric(string(token)):
			number, end = parseNumber(line)
			c.expression = append(c.expression, Expression{Number, number})
		case isSymbol(string(token)):
			symbol, end = parseSymbol(line)
			c.expression = append(c.expression, Expression{Symbol, symbol})
		case isAlpha(string(token)):
			varName, end = parseVariable(line)
			varValue = c.getVarValue(varName)
			if varValue == "" {
				return nil
			}
			c.expression = append(c.expression, Expression{Number, varValue})
		default:
			fmt.Println("Invalid expression")
			return nil
		}
		if !isValid(end) {
			break
		}
		line = line[end:]
	}
	return c.expression
}

// stackOperator performs the operation on the stack
func (c Calculator) stackOperator(token string) ([]string, []string) {
	if len(c.stack) == 0 || c.stack[len(c.stack)-1] == "(" || token == "(" {
		c.stack = append(c.stack, token)
		return c.stack, c.postfixExpr
	}

	if token == ")" {
		for {
			if c.stack[len(c.stack)-1] == "(" {
				pop(&c.stack)
				break
			}
			c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
		}
		return c.stack, c.postfixExpr
	}

	if higherPrecedence(c.stack[len(c.stack)-1], token) {
		c.stack = append(c.stack, token)
	} else {
		for len(c.stack) > 0 && !higherPrecedence(c.stack[len(c.stack)-1], token) {
			c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
		}
		c.stack = append(c.stack, token)
	}
	return c.stack, c.postfixExpr
}

// getPostfix converts expression to a postfixExpr
func (c Calculator) getPostfix(expression []Expression) []string {
	var prevSym string

	for i, token := range expression {
		if isNumeric(token.Value) {
			if prevSym != "" && isNumeric(prevSym) {
				c.postfixExpr = append(c.postfixExpr, pop(&c.postfixExpr)+token.Value)
			} else {
				if i == 1 && (c.stack[0] == "-" || c.stack[0] == "+") {
					c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
				}
				c.postfixExpr = append(c.postfixExpr, token.Value)
				prevSym = token.Value
			}
		}

		if sliceContains(symbols, token.Value) {
			if prevSym == "" {
				prevSym = token.Value
				c.stack = append(c.stack, token.Value)
				continue
			}

			if prevSym == token.Value {
				switch token.Value {
				case "+":
					c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
				case "-":
					pop(&c.stack)
					c.stack = append(c.stack, "-")
				case "*":
					fmt.Println("Invalid expression")
					return nil
				case "/":
					fmt.Println("Invalid expression")
					return nil
				}
			} else {
				prevSym = token.Value
			}
			c.stack, c.postfixExpr = c.stackOperator(token.Value)
		}
	}

	// if the stack has any remaining parenthesis "(" or ")", remove them:
	for i := len(c.stack) - 1; i >= 0; i-- {
		if c.stack[i] == "(" || c.stack[i] == ")" {
			c.stack = append(c.stack[:i], c.stack[i+1:]...)
		}
	}

	// Append to postfixExpr until c.stack is empty
	for len(c.stack) > 0 {
		c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
	}

	return c.postfixExpr
}

// higherPrecedence returns true if the first symbol has higher precedence than the second
func higherPrecedence(stackPop, token string) bool {
	if stackPop == "(" || operatorRank[token] > operatorRank[stackPop] {
		return true
	}
	return false
}

// checkSingleNum() checks if the expression is a single number or variable like: --10 or -a or 100
func (c Calculator) checkSingleNum() int {
	var singleNum []string
	sign := 1

	// Iterate over the expression and confirm that it is a single number with signs in front or not:
	for i, token := range c.expression {
		if isNumeric(token.Value) && i < len(c.expression)-1 {
			if isSymbol(c.expression[i+1].Value) {
				return 0
			}
		} else {
			singleNum = append(singleNum, token.Value)
		}
	}

	// Calculate the correct "sign" either positive or negative for the single number
	for _, token := range singleNum {
		if token == "-" {
			sign *= -1
		}
		if isNumeric(token) {
			x, err := strconv.Atoi(token)
			if err != nil {
				log.Fatal(err)
			}
			return x * sign
		}
	}
	return 0
}

// getTotal calculates the total result of postfixExpr
func (c Calculator) getTotal() int {
	var (
		end, minusCount int
	)

	// If the expression is a single number, return it and then ask the user for the next expression
	if c.checkSingleNum() != 0 {
		return c.checkSingleNum()
	}

	// If the expression is not a single number, then begin processing the postfixExpr:
	for i, token := range c.postfixExpr {
		// If the token is a positive or negative sign, increment the respective counters:
		switch token {
		case "-":
			end += 1
			minusCount += 1
		case "+":
			end += 1
		}

		// If the first element of the postfixExpr is not a number, break the loop
		if isNumeric(c.postfixExpr[0]) {
			break
		}

		// If the first element of postfixExpr is a number, we need to apply a series of validations
		// To properly check positive and negative signs in front of the number
		if isNumeric(token) {
			// Remove the first sign either positive or negative from the postfixExpr
			c.postfixExpr = c.postfixExpr[end:]

			// Check for cases with only one negative sign in front: -10-12--8
			if minusCount == 1 && isNumeric(c.postfixExpr[i]) {
				c.postfixExpr[0] = "-" + c.postfixExpr[0]
				break
			}

			// Check for cases with two negatives sign in front: --10-12--8
			if minusCount == 1 && !isNumeric(c.postfixExpr[i]) {
				break
			}

			// Check for cases with more than 3 negatives sign in front, like: ---10--12--8
			if minusCount > 1 {
				if minusCount%2 == 1 {
					break
				}
				if minusCount%2 == 0 && c.postfixExpr[i-2] == "+" {
					break
				}
				c.postfixExpr[0] = "-" + c.postfixExpr[0]
				break
			}
		}
	}

	// After checking for multiple negative signs, finally start calculating the c.postfixExpr
	for i, token := range c.postfixExpr {
		if token == "-" && i <= len(c.postfixExpr)-1 {
			minusCount += 1
		}
		if isNumeric(token) || isNumeric(token[1:]) {
			c.stack = append(c.stack, token)
		} else if len(c.stack) > 1 {
			if minusCount%2 == 0 && minusCount != 0 {
				token = "+"
			}
			// Get the two last elements of the stack and perform the math operation according to the 'token'
			b, a := pop(&c.stack), pop(&c.stack)

			// Remember to convert the 'b' and 'a' from string to int before performing the math operation
			x, err := strconv.Atoi(a)
			if err != nil {
				log.Fatal(err)
			}

			y, err := strconv.Atoi(b)
			if err != nil {
				log.Fatal(err)
			}

			// Perform the math operation and push the result to the stack
			c.stack = append(c.stack, strconv.Itoa(evalSymbol(x, y, token)))
			// Reset the minusCount to 0 to properly check for multiple negative signs for the next iteration
			minusCount = 0
		}
	}

	// Finally return the calculated result:
	if len(c.stack) > 0 {
		x, _ := strconv.Atoi(pop(&c.stack))
		return x
	}
	return 0
}

// processLine is the main function that processes the input line and returns the result
func (c Calculator) processLine(line string) {
	if !validateExpression(line) {
		fmt.Println("Invalid expression")
		return
	}

	// If the expression is valid, proceed to append each operator and number to it:
	c.expression = c.appendValues(line)

	// If the expression is not blank, then get its postfix form:
	if len(c.expression) > 0 {
		c.postfixExpr = c.getPostfix(c.expression)
		// If the postfix form is not blank, then proceed to calculate the result:
		if c.postfixExpr != nil {
			fmt.Println(c.getTotal())
		}
		return
	}
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

			c.processLine(line)
		}
	}
}

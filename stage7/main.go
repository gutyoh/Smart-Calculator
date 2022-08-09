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

func isParenthesis(token string) bool {
	return token == "(" || token == ")"
}

func isValid(end int) bool {
	return end != 0
}

// splitParenthesis separates tokens like "(3" or "1)" and returns a slice with the separated tokens
func splitParenthesis(tokens []string) []string {
	var newLine []string
	for _, token := range tokens {
		if strings.HasPrefix(token, "(") {
			newLine = append(newLine, "(")
			newLine = append(newLine, token[1:])
		} else if strings.HasSuffix(token, ")") {
			newLine = append(newLine, token[:len(token)-1])
			newLine = append(newLine, ")")
		} else {
			newLine = append(newLine, token)
		}
	}
	return newLine
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

func (c Calculator) checkStackElements() bool {
	counter := 0
	// check if the stack contains exactly two numeric elements
	for _, element := range c.stack {
		if isNumeric(element) {
			counter++
		}
	}

	if counter == 2 {
		return true
	}
	return false
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

func processCommand(line string) {
	if line != "/exit" && line != "/help" {
		fmt.Println("Unknown command")
		return
	}
}

// checkParentheses checks if there are the same amount of parenthesis in the infixExpr
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

// getTotal calculates the total result of the postfixExpr infixExpr
func (c Calculator) getTotal() int {
	var a, b string
	// var sign = 1
	var end, minusCount = 0, 0
	var sign string

	// Check if c.postfixExpr starts with a negative sign to validate cases like: --10++10--8 or -10+10+8
	for _, token := range c.postfixExpr {
		if isNumeric(token) {
			c.postfixExpr = c.postfixExpr[end:]
			break
		}
		if token == "-" {
			sign = "-"
			end += 1
			minusCount += 1
		} else if token == "+" {
			end += 1
			continue
		}
	}

	if c.postfixExpr[len(c.postfixExpr)-1] == "-" {
		minusCount += 1
	}

	// If c.postfixExpr is only a single number, return the number
	if len(c.postfixExpr) == 1 {
		x, _ := strconv.Atoi(sign + c.postfixExpr[0])
		return x
	} else if len(c.postfixExpr) == 2 && minusCount%2 == 1 {
		x, _ := strconv.Atoi("-" + c.postfixExpr[0])
		return x
	} else if len(c.postfixExpr) == 2 && minusCount%2 == 0 {
		x, _ := strconv.Atoi(c.postfixExpr[0])
		return x
	}

	// TODO - Fix this logic to handle multiple "-" or "+" symbols as the first - DONE
	// symbols in the expression, like: ++10++10--8 or --10--10--8
	for i, token := range c.postfixExpr {
		if sign == "-" && i == len(c.postfixExpr)-1 && c.postfixExpr[i] == "-" {
			c.stack[1] = "-" + c.stack[1]
		}

		if isNumeric(token) && i < len(c.postfixExpr)-1 {
			if sign != "" && isNumeric(c.postfixExpr[i+1]) || end >= 2 {
				c.stack = append(c.stack, sign+token)
				sign = ""
			} else {
				c.stack = append(c.stack, token)
			}
		} else if len(c.stack) > 1 {
			b, a = pop(&c.stack), pop(&c.stack)
			x, _ := strconv.Atoi(a)
			y, _ := strconv.Atoi(b)

			c.stack = append(c.stack, strconv.Itoa(evalSymbol(x, y, token)))
		}
	}

	if len(c.stack) > 0 {
		x, _ := strconv.Atoi(pop(&c.stack))
		return x
	}
	return 0
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

func parseParenthesis(line string) (string, int) {
	var parenthesis string
	var end int
	for i, token := range line {
		if isParenthesis(string(token)) {
			parenthesis += string(token)
			end = i + 1
			break
		}
	}
	return parenthesis, end
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
	var number, parenthesis, symbol, varName, varValue string
	var end int

	if !validateExpression(line) {
		fmt.Println("Invalid expression")
		return
	}

	line = strings.Replace(line, " ", "", -1)
	tokens = strings.Split(line, "")
	tokens = splitParenthesis(tokens)

	// remove any "" in between the tokens
	tokens = strings.Split(strings.Join(tokens, ""), "")

	for i, token := range tokens {
		if isNumeric(token) {
			number, end = parseNumber(line)
			if isValid(end) {
				line = line[end:]
				c.expression = append(c.expression, Expression{Number, number})
			}
		}

		if isSymbol(token) {
			symbol, end = parseSymbol(line)
			if isValid(end) && symbol != "(" && symbol != ")" {
				line = line[end:]
				c.expression = append(c.expression, Expression{Symbol, symbol})
			}
		}

		if isParenthesis(token) {
			parenthesis, end = parseParenthesis(line)
			if isValid(end) {
				line = line[end:]
				c.expression = append(c.expression, Expression{Symbol, parenthesis})
			}
		}

		if isAlpha(token) {
			varName, end = parseVariable(line)
			if isValid(end) {
				line = line[end:]
				varValue = c.getVarValue(varName)
				if varValue != "" {
					c.expression = append(c.expression, Expression{Number, varValue})
				}
			}
		}

		// Append the last number, or last variable to the expression
		if i == len(tokens)-1 && isNumeric(token) {
			number, end = parseNumber(line)
			if number != "" {
				c.expression = append(c.expression, Expression{Number, number})
			}
		}

		if i == len(tokens)-1 && isSymbol(token) {
			symbol, end = parseSymbol(line)
			if symbol != "" {
				c.expression = append(c.expression, Expression{Symbol, symbol})
			}
		}

		if i == len(tokens)-1 && isParenthesis(token) {
			parenthesis, end = parseParenthesis(line)
			if parenthesis != "" {
				c.expression = append(c.expression, Expression{Symbol, parenthesis})
			}
		}

		if i == len(tokens)-1 && isAlpha(token) {
			varName, end = parseVariable(line)
			varValue = c.getVarValue(varName)
			if varValue != "" {
				c.expression = append(c.expression, Expression{Number, varValue})
			}
		}
	}

	if len(c.expression) > 0 && c.getPostfix(c.expression) != nil {
		c.postfixExpr = c.getPostfix(c.expression)
		fmt.Println(c.getTotal())
	}
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
				if "+" == token.Value {
					c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
				} else if "-" == token.Value {
					pop(&c.stack)
					c.stack = append(c.stack, "-")
				} else if "*" == token.Value || "/" == token.Value {
					fmt.Println("Invalid expression")
					return nil
				}
			} else {
				prevSym = token.Value
			}
			c.stack, c.postfixExpr = c.stackOperator(token.Value)
		}
	}

	// if the stack has any "(" or ")" remaining, remove them
	for i := len(c.stack) - 1; i >= 0; i-- {
		if c.stack[i] == "(" || c.stack[i] == ")" {
			c.stack = append(c.stack[:i], c.stack[i+1:]...)
		}
	}

	// TODO fix this logic to properly handle parenthesis operations like: 4*2+5*3+6*(2+3) or 4*2+5*3+6*((2+3))
	// or ((10+10)) * 8 and (10+10) * 8 or ((10+10) / (3 + 5) * 8) + 5
	for _, token := range c.stack {
		if len(c.stack) == 0 {
			break
		} else {
			c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
			fmt.Sprintln(token)
		}

		//if len(c.stack) == 0 {
		//	break
		//} else if token != "(" && token != ")" {
		//	if pop(&c.stack) != "(" || pop(&c.stack) != ")" {
		//		c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
		//	}
		//} else if token == "(" || token == ")" {
		//	pop(&c.stack)
		//} else {
		//	// if the last element of the stack is not "(" or ")" then append it to the postfixExpr
		//	if len(c.stack) > 0 && c.stack[len(c.stack)-1] != "(" && c.stack[len(c.stack)-1] != ")" {
		//		c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
		//	}
		//	pop(&c.stack)
		//}
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
		for {
			if len(c.stack) > 0 && !higherPrecedence(c.stack[len(c.stack)-1], token) {
				c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
			} else {
				break
			}
		}
		c.stack = append(c.stack, token)
	}
	return c.stack, c.postfixExpr
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

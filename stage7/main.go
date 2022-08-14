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
	"reflect"
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

type Calculator struct {
	memory      map[string]int
	stack       []Expression
	postfixExpr []Expression
	expression  []Expression
	OperationType
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

// pop deletes the last element of the stack []Expression and returns it
func pop(alist *[]Expression) Expression {
	f := len(*alist)
	rv := (*alist)[f-1]
	*alist = (*alist)[:f-1]
	return rv
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

	if reflect.TypeOf(value).Kind() == reflect.String {
		value = c.getVarValue(value.(string))
		if value == nil {
			return
		}
	}
	c.memory[variable.(string)] = value.(int)
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

func getOperationType(line string) OperationType {
	if checkAssignment(line) {
		return Assignment
	}
	return Regular
}

func validateExpression(line string) bool {
	var number, end int
	var varName string
	var valid = true

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

	// Then check if there are the same amount of parenthesis on both sides of the line
	if checkParenthesis(line) {
		fmt.Println("Invalid expression")
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
	valid = validateSyntax(line, end, number, valid, varName)
	return valid
}

// validateSyntax validates the syntax of the expression and checks for special edge cases
func validateSyntax(line string, end int, number any, valid bool, varName string) bool {
	var prevSym string

	// validateSyntax checks if the expression has any "Invalid identifiers" like a2a or a1 = 8
	// And other edge cases like test = 2n or test = n2
Loop:
	for len(line) > 0 {
		token := string(line[0])
		switch {
		case token == " ":
			end = 1
		case token == "=":
			end = 1
			prevSym = "="
		case token == "*":
			if prevSym == "*" {
				fmt.Println("Invalid expression")
				valid = false
				break Loop
			}
			end = 1
			prevSym = token
		case token == "/":
			if prevSym == "/" {
				fmt.Println("Invalid expression")
				valid = false
				break Loop
			}
			end = 1
			prevSym = token
		case isNumeric(token):
			number, end = parseNumber(line)
			if number == nil && prevSym == "=" { // Validates cases like test = 2n
				fmt.Println("Invalid assignment")
				valid = false
				break Loop
			}

			if varName == "" { // Validates cases like 5 = 5, or 100 = 20
				if number != nil && prevSym == "=" {
					fmt.Println("Invalid assignment")
					valid = false
					break Loop
				}
			}
		case isAlpha(token):
			varName, end = parseVariable(line)
			if varName == "" && prevSym == "=" { // Validates cases like test = a2a
				fmt.Println("Invalid assignment")
				valid = false
				break Loop
			}

			if varName == "" { // Validates cases like a2a or n22 or a1 = 8
				fmt.Println("Invalid identifier")
				valid = false
				break Loop
			}
		case isSymbol(token):
			_, end = parseSymbol(line)
			prevSym = token
		}
		line = line[end:]
	}
	return valid
}

// evalSymbol evaluates the symbol and performs the operation accordingly
func evalSymbol(a, b int, operator any) int {
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

// parseSymbol parses the symbols: "+", "-", "*", "/", "(", ")", "^" from the input line
func parseSymbol(line string) (string, int) {
	var symbol string
	var end int

	for i, t := range line {
		token := string(t)
		if isSymbol(token) {
			symbol += token
			end = i + 1
			break
		}
	}
	end = len(symbol)
	return symbol, end
}

// parseVariable parses a more-than-one-character variable from the input line
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

// getVarValue returns the value of the variable if it's in the memory of the Calculator
func (c Calculator) getVarValue(variable string) any {
	if !mapContains(c.memory, variable) {
		fmt.Println("Unknown variable")
		return nil
	}
	return c.memory[variable]
}

func (c Calculator) appendValues(line string) []Expression {
	var (
		symbol, varName  string
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
		case isSymbol(token):
			symbol, end = parseSymbol(line)
			c.expression = append(c.expression, Expression{Symbol, symbol})
		case isAlpha(token):
			varName, end = parseVariable(line)
			varValue = c.getVarValue(varName)
			if varValue == nil {
				return nil
			}
			c.expression = append(c.expression, Expression{Number, varValue.(int)})
		default:
			return nil
		}
		line = line[end:]
	}
	return c.expression
}

// stackOperator performs the operation on the stack
func (c Calculator) stackOperator(token string) ([]Expression, []Expression) {
	if len(c.stack) == 0 || token == "(" {
		c.stack = append(c.stack, Expression{Symbol, token})
		return c.stack, c.postfixExpr
	}

	if token == ")" {
		for {
			if c.stack[len(c.stack)-1].Value.(string) == "(" {
				pop(&c.stack)
				break
			}
			c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
		}
		return c.stack, c.postfixExpr
	}

	if higherPrecedence(c.stack[len(c.stack)-1].Value.(string), token) {
		c.stack = append(c.stack, Expression{Symbol, token})
	} else {
		for len(c.stack) > 0 && !higherPrecedence(c.stack[len(c.stack)-1].Value.(string), token) {
			c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
		}
		c.stack = append(c.stack, Expression{Symbol, token})
	}
	return c.stack, c.postfixExpr
}

// getPostfix converts expression to a postfixExpr
func (c Calculator) getPostfix(expression []Expression) []Expression {
	var prevSym any

	for i, token := range expression {
		switch token.ExpressionType {
		case Number:
			if i == 1 && (c.stack[0].Value.(string) == "-" || c.stack[0].Value.(string) == "+") {
				c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
			}
			c.postfixExpr = append(c.postfixExpr, Expression{Number, token.Value.(int)})
			prevSym = token.Value
		case Symbol:
			if sliceContains(symbols, token.Value.(string)) {
				if prevSym == "" || prevSym == nil {
					prevSym = token.Value.(string)
					c.stack = append(c.stack, Expression{Symbol, token.Value.(string)})
					continue
				}

				if prevSym == token.Value {
					switch token.Value {
					case "+":
						c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
					case "-":
						pop(&c.stack)
						c.stack = append(c.stack, Expression{Symbol, "-"})
					}
				} else {
					prevSym = token.Value.(string)
				}
				c.stack, c.postfixExpr = c.stackOperator(token.Value.(string))
			}
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
	minusCount := 0
	number := 0
	numCount := 0

	for i, token := range c.postfixExpr {
		switch token.ExpressionType {
		case Symbol:
			if token.Value.(string) == "-" {
				minusCount += 1
			}
		case Number:
			if numCount > 0 {
				return 0
			}
			number = token.Value.(int)
			numCount += 1
		}

		if i == len(c.postfixExpr)-1 {
			if minusCount%2 != 1 {
				return number
			}
			return number * -1
		}
	}

	return 0
}

// getTotal calculates the total result of postfixExpr
func (c Calculator) getTotal() int {
	var (
		end, minusCount, singleNum int
	)

	var mc2 int

	// If the expression is a single number, return it and then ask the user for the next expression
	singleNum = c.checkSingleNum()
	if singleNum != 0 {
		return singleNum
	}

	// TODO -- Check the logic of "mc" and "minusCount"
	// Get the first symbol of the original expression
	for _, token := range c.expression {
		if token.ExpressionType == Symbol {
			if token.Value.(string) == "-" {
				mc2 += 1
			}
		}
		if token.ExpressionType == Number {
			break
		}
	}

	// If the expression is not a single number, then begin processing the postfixExpr:
	for i, token := range c.postfixExpr {
		if reflect.TypeOf(c.postfixExpr[0].Value).Kind() == reflect.Int {
			break
		}

		switch token.ExpressionType {
		case Symbol:
			if token.Value.(string) == "-" {
				end += 1
				minusCount += 1
			}
			if token.Value.(string) == "+" {
				end += 1
			}
		case Number:
			// Remove the first sign either positive or negative from the postfixExpr
			c.postfixExpr = c.postfixExpr[end:]

			// Check for cases with only one negative sign in front: -10-12--8
			if minusCount == 1 && reflect.TypeOf(c.postfixExpr[i].Value).Kind() == reflect.Int {
				c.postfixExpr[0].Value = c.postfixExpr[0].Value.(int) * -1
				break
			}

			// Check for cases with two negatives sign in front: --10-12--8
			if minusCount == 1 && reflect.TypeOf(c.postfixExpr[i].Value).Kind() != reflect.Int {
				break
			}

			// Check for cases with more than 3 negatives sign in front, like: ---10--12--8
			if mc2 > 1 {
				if mc2%2 == 0 {
					break
				}
				c.postfixExpr[0].Value = c.postfixExpr[0].Value.(int) * -1
				break
			}
			break
		}
	}

	// After checking for multiple negative signs, finally start calculating the c.postfixExpr
	for i, token := range c.postfixExpr {
		switch token.ExpressionType {
		case Symbol:
			if token.Value.(string) == "-" && i <= len(c.postfixExpr)-1 {
				minusCount += 1
			}
		case Number:
			n := strconv.Itoa(token.Value.(int))
			c.stack = append(c.stack, Expression{Number, n})
			continue
		}
		if len(c.stack) > 1 {
			if minusCount%2 == 0 && minusCount != 0 {
				token.Value = "+"
			}

			if mc2%2 == 1 && token.Value.(string) != "^" {
				token.Value = "+"
			}

			// Get the two last elements of the stack and perform the math operation according to the 'token'
			b, a := pop(&c.stack), pop(&c.stack)

			// Remember to convert the 'b' and 'a' from string to int before performing the math operation
			x, err := strconv.Atoi(a.Value.(string))
			if err != nil {
				log.Fatal(err)
			}

			y, err := strconv.Atoi(b.Value.(string))
			if err != nil {
				log.Fatal(err)
			}

			// Perform the math operation and push the result to the stack
			c.stack = append(c.stack, Expression{Number, strconv.Itoa(evalSymbol(x, y, token.Value.(string)))})
			// Reset the minusCount to 0 to properly check for multiple negative signs for the next iteration
			minusCount = 0
			mc2 = 0
		}
	}

	// Finally return the calculated result:
	if len(c.stack) > 0 {
		x, _ := strconv.Atoi(pop(&c.stack).Value.(string))
		return x
	}
	return 0
}

// processLine is the main function that processes the input line and returns the result
func (c Calculator) processLine(line string) {
	//if !validateExpression(line) {
	//	fmt.Println("Invalid expression")
	//	return
	//}

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

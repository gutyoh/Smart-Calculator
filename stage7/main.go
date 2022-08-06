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
	Sign
	Variable
)

type Expression struct {
	ExpressionType
	Value string
}

type Calculator struct {
	result      int
	memory      map[string]int
	message     string
	infixExpr   []Expression
	stack       []string
	postfixExpr []string
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

func isValid(end int) bool {
	return end != 0
}

// FOR FUN ONLY TO HANDLE FLOAT INPUTS -- THIS CAN BE REMOVED FROM FINAL SOLUTION
// floatStringToInt converts a "float string" like: "3.50" to an int: "3"
func floatStringToInt(a string) int {
	f, err := strconv.ParseFloat(a, 64)
	if err != nil {
		log.Fatal(err)
	}
	return int(f)
}

// splitParenthesis separates tokens like "(3" or "1)" and returns a slice with the separated tokens
func splitParenthesis(tokens []string) []Expression {
	var newLine []Expression
	for _, token := range tokens {
		if strings.HasPrefix(token, "(") {
			newLine = append(newLine, Expression{Sign, "("})
			newLine = append(newLine, Expression{Number, token[1:]})
		}

		if strings.HasSuffix(token, ")") {
			newLine = append(newLine, Expression{Number, token[:len(token)-1]})
			newLine = append(newLine, Expression{Sign, ")"})
		}

		newLine = append(newLine, Expression{Number, token})

		//if strings.HasPrefix(token, "(") {
		//	newLine = append(newLine, "(")
		//	newLine = append(newLine, token[1:])
		//} else if strings.HasSuffix(token, ")") {
		//	newLine = append(newLine, token[:len(token)-1])
		//	newLine = append(newLine, ")")
		//} else {
		//	newLine = append(newLine, token)
		//}
	}
	return newLine
}

// repeatedSymbol checks if there is more than one "*" or "/" symbol in the line
func repeatedSymbol(line string) bool {
	return strings.Count(line, "*") > 1 || strings.Count(line, "/") > 1
}

// pop deletes the last element of the stack
func pop(alist *[]string) string {
	f := len(*alist)
	rv := (*alist)[f-1]
	*alist = (*alist)[:f-1]
	return rv
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

func validateExpression(line string) bool {
	// Check for the most basic case of invalid expressions, trailing operators like: 10+10-8-
	if strings.HasSuffix(line, "+") || strings.HasSuffix(line, "-") {
		fmt.Println("Invalid expression")
		return false
	}

	if checkParenthesis(line) {
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

// getTotal calculates the total result of the postfixExpr infixExpr
func (c Calculator) getTotal() int {
	var a, b string

	for i, token := range c.postfixExpr {
		if isNumeric(token) {
			c.stack = append(c.stack, token)
		}

		if len(c.stack) == 1 {
			b = pop(&c.stack)
			x, _ := strconv.Atoi(b)
			// TODO FIX THIS LOGIC
			if i == len(c.postfixExpr)-1 && token == "-" && x <= 0 {
				return x
			} else if i == len(c.postfixExpr)-1 && token == "+" {
				return -x
			} else if token == "-" {
				c.stack = append(c.stack, strconv.Itoa(evalSymbol(0, x, token)))
			} else {
				c.stack = append(c.stack, strconv.Itoa(evalSymbol(x, 0, token)))
			}
		}

		if len(c.stack) > 1 {
			b, a = pop(&c.stack), pop(&c.stack)

			// Finally, convert 'a' and 'b' to int type:
			x, err := strconv.Atoi(a)
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.Atoi(b)
			if err != nil {
				log.Fatal(err)
			}
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

// getPostfix converts the infix infixExpr to postfixExpr
func (c Calculator) getPostfix(line string) []string {
	var prevSym, temp, varName string
	// varValue string
	// var end int
	var tokens []string

	if !validateExpression(line) {
		fmt.Println("Invalid expression")
	}

	line = strings.Replace(line, " ", "", -1)
	tokens = strings.Split(line, "")

	if strings.Contains(line, " ") {
		c.infixExpr = splitParenthesis(tokens)
	} else {
		c.infixExpr = splitParenthesis(tokens)
	}

	for i, token := range tokens {
		if repeatedSymbol(token) {
			fmt.Println("Invalid expression")
			break
		}

		// TODO -- Add the logic to parse long variable names, like "test"
		if isAlpha(token) {
			if mapContains(c.memory, varName) {
				prevSym = token
				token = strconv.Itoa(c.memory[token])
				c.postfixExpr = append(c.postfixExpr, token)
			} else if mapContains(c.memory, strings.Join(tokens, "")) {
				prevSym = strings.Join(tokens, "")
				token = strconv.Itoa(c.memory[strings.Join(tokens, "")])
				c.postfixExpr = append(c.postfixExpr, token)
				break
			} else {
				fmt.Println("Unknown variable")
				return nil
			}
		} else if isNumeric(token) {
			if len(c.stack) > 1 && prevSym != "" && isNumeric(token) {
				c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
				temp = token
			} else if prevSym != "" && isNumeric(prevSym) {
				c.postfixExpr = append(c.postfixExpr, pop(&c.postfixExpr)+token)
			} else {
				if sliceContains(c.stack, "-") && i == 1 {
					c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
				}
				c.postfixExpr = append(c.postfixExpr, temp+token)
				prevSym = token
				temp = ""
			}
		} else if sliceContains(symbols, token) {
			if prevSym == "" {
				prevSym = token
				c.stack = append(c.stack, token)
				continue
			}

			if prevSym == token {
				switch token {
				case "+":
					if len(c.stack) == 1 && c.stack[0] == "+" {
						c.stack = append(c.stack, token)
						continue
					}

					if c.stack[0] == "+" && c.stack[1] == "+" {
						c.stack, c.postfixExpr = c.stackOperator(token)
					}

					continue
				case "-":
					prevSym = "+"
					pop(&c.stack)
					c.stack = append(c.stack, token)
				case "*":
					fmt.Println("Invalid expression")
					break
				case "/":
					fmt.Println("Invalid expression")
					break
				}
			} else {
				prevSym = token
			}
			c.stack, c.postfixExpr = c.stackOperator(token)
		}
	}

	for {
		if len(c.stack) == 0 {
			break
		}
		c.postfixExpr = append(c.postfixExpr, pop(&c.stack))
	}
	return c.postfixExpr
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
				c.stack = c.stack[:len(c.stack)-1]
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

// higherPrecedence returns true if the first symbol has higher precedence than the second
func higherPrecedence(stackPop, token string) bool {
	if stackPop == "(" || operatorRank[token] > operatorRank[stackPop] {
		return true
	}
	return false
}

// checkParentheses checks if there are the same amount of parenthesis in the infixExpr
func checkParenthesis(line string) bool {
	return strings.Count(line, "(") != strings.Count(line, ")")
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

			// If none of the above cases were met, then the line is an expression like: "10+10+8"
			// That can be further processed to get the total (in case it is valid, of course)
			c.postfixExpr = c.getPostfix(line)
		}

		//if len(line) > 0 {
		//	if checkCommand(line) {
		//		c.message = getCommand(line)
		//	} else if checkAssignment(line) {
		//		c.assign(line)
		//		continue
		//	} else {
		//		if checkParenthesis(line) {
		//			c.message = "Invalid expression"
		//		} else {
		//			// Since a command wasn't issued, reset the c.message variable
		//			c.message = ""
		//
		//			// Use splitParenthesis to split any tokens like "(3" or "1)"
		//			if strings.Contains(line, " ") {
		//				c.infixExpr = splitParenthesis(strings.Split(line, " "))
		//			} else {
		//				c.infixExpr = splitParenthesis(strings.Split(line, ""))
		//			}
		//
		//			// Get the postfixExpr and then calculate the result
		//			joinedExpr := strings.Join(c.infixExpr, "")
		//			c.postfixExpr = c.getPostfix(joinedExpr)
		//			c.result = c.getTotal()
		//		}
		//	}
		//
		//	// If a command was issued, print the command message;
		//	// Otherwise if 'c.postfixExpr' is not nil print the calculated result
		//	if c.message != "" {
		//		fmt.Println(c.message)
		//	} else if c.postfixExpr != nil {
		//		fmt.Println(c.result)
		//		c.postfixExpr = nil // Reset the postfixExpr
		//	}
		//}
	}
}

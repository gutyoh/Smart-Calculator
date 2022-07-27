package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Calculator struct {
	stack   []string
	postfix []string
	result  string
	memory  map[string]int
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
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// isAlpha checks if all the characters in the string are alphabet letters
func isAlpha(s string) bool {
	re := regexp.MustCompile("^[a-zA-Z]+$")
	return re.MatchString(s)
}

// stringToFloat converts a string to a float number
func stringToFloat(a string) int {
	f, err := strconv.ParseFloat(a, 64)
	if err != nil {
		log.Fatal(err)
	}
	return int(f)
}

// floatToString converts a float number to a string
func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// removeSpaces removes blank spaces from the line
func removeSpaces(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		}
	}
	return string(rr)
}

// splitParenthesis separates tokens like "(3" or "1)" and returns a slice with the separated tokens
func splitParenthesis(line []string) []string {
	var newLine []string
	for _, token := range line {
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
func (c Calculator) assign(line string) string {
	variable, value := func(s []string) (string, string) {
		return s[0], s[1]
	}(func() (elems []string) {
		for _, x := range strings.Split(line, "=") {
			elems = append(elems, strings.TrimSpace(x))
		}
		return elems
	}())

	if !isAlpha(variable) {
		return "Invalid identifier"
	}

	if !isNumeric(value) {
		if !mapContains(c.memory, value) {
			return "Invalid assignment"
		} else {
			value = strconv.Itoa(c.memory[value])
		}
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}
	c.memory[variable] = v

	return ""
}

func getCommand(line string) string {
	if line == "/exit" {
		return "Bye!"
	} else if line == "/help" {
		return "The program calculates the sum of numbers"
	}
	return "Unknown command"
}

// getTotal calculates the total result of the postfix expression
func (c Calculator) getTotal() string {
	for _, val := range c.postfix {
		if isNumeric(val) {
			c.stack = append(c.stack, val)
		} else {
			b, a := pop(&c.stack), pop(&c.stack)

			//if 'b' and 'a' are float strings, convert them to float numbers:
			if stringToFloat(a) != 0 && stringToFloat(b) != 0 {
				a, b = floatToString(float64(stringToFloat(a))), floatToString(float64(stringToFloat(b)))
			}

			// Finally, convert 'a' and 'b' to int type:
			x, err := strconv.Atoi(a)
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.Atoi(b)
			if err != nil {
				log.Fatal(err)
			}
			c.stack = append(c.stack, strconv.Itoa(evalSymbol(x, y, val)))
		}
	}
	if len(c.stack) > 0 {
		return c.stack[len(c.stack)-1]
	}
	return ""
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

// getPostfix converts the infix expression to postfix
func (c Calculator) getPostfix(line string) []string {
	var prevSym string
	var tokens []string

	if strings.Contains(line, " ") {
		tokens = strings.Split(line, " ")
	} else {
		tokens = strings.Split(line, "")
	}

	for _, token := range tokens {
		if repeatedSymbol(token) {
			fmt.Println("Invalid expression")
			break
		}

		if token == " " {
			continue
		} else if isAlpha(token) {
			if mapContains(c.memory, token) {
				prevSym = token
				token = strconv.Itoa(c.memory[token])
				c.postfix = append(c.postfix, token)
			} else if mapContains(c.memory, strings.Join(tokens, "")) {
				prevSym = strings.Join(tokens, "")
				token = strconv.Itoa(c.memory[strings.Join(tokens, "")])
				c.postfix = append(c.postfix, token)
				break
			} else {
				fmt.Println("Unknown variable")
				break
			}
		} else if isNumeric(token) {
			if prevSym != "" && isNumeric(prevSym) {
				c.postfix = append(c.postfix, pop(&c.postfix)+token)
			} else {
				c.postfix = append(c.postfix, token)
				prevSym = token
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
					continue
				case "-":
					prevSym = "+"
					pop(&c.stack)
					c.stack = append(c.stack, prevSym)
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
			c.stack, c.postfix = c.stackOperator(token)
		}
	}

	for {
		if len(c.stack) == 0 {
			break
		}
		c.postfix = append(c.postfix, pop(&c.stack))
	}
	return c.postfix
}

// stackOperator performs the operation on the stack
func (c Calculator) stackOperator(token string) ([]string, []string) {
	if len(c.stack) == 0 || c.stack[len(c.stack)-1] == "(" || token == "(" {
		c.stack = append(c.stack, token)
		return c.stack, c.postfix
	}

	if token == ")" {
		for {
			if c.stack[len(c.stack)-1] == "(" {
				c.stack = c.stack[:len(c.stack)-1]
				break
			}
			c.postfix = append(c.postfix, pop(&c.stack))
		}
		return c.stack, c.postfix
	}

	if higherPrecedence(c.stack[len(c.stack)-1], token) {
		c.stack = append(c.stack, token)
	} else {
		for {
			if len(c.stack) > 0 && !higherPrecedence(c.stack[len(c.stack)-1], token) {
				c.postfix = append(c.postfix, pop(&c.stack))
			} else {
				break
			}
		}
		c.stack = append(c.stack, token)
	}
	return c.stack, c.postfix
}

// higherPrecedence returns true if the first symbol has higher precedence than the second
func higherPrecedence(stackPop, token string) bool {
	if stackPop == "(" || operatorRank[token] > operatorRank[stackPop] {
		return true
	}
	return false
}

// checkParentheses checks if there are the same amount of parenthesis in the expression
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

		// Check if the entered line has any preceding blank spaces
		if strings.HasPrefix(line, " ") {
			line = removeSpaces(line)
		}

		if len(line) > 0 {
			if checkCommand(line) {
				c.result = getCommand(line)
			} else if checkAssignment(line) {
				c.result = c.assign(line)
			} else {
				if checkParenthesis(line) {
					c.result = "Invalid expression"
				} else {
					// Use splitParenthesis to split any tokens like "(3" or "1)"
					expression := splitParenthesis(strings.Split(line, " "))
					// Get the postfix expression and then calculate the result
					c.postfix = c.getPostfix(strings.Join(expression, " "))
					c.result = c.getTotal()
				}
			}

			if c.result != "" {
				fmt.Println(c.result)
			}
		}
	}
}

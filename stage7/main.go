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
	"regexp"
	"strconv"
	"strings"
)

type Calculator struct {
	result      int
	memory      map[string]int
	message     string
	infixExpr   []string
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
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// isAlpha checks if all the characters in the string are alphabet letters
func isAlpha(s string) bool {
	re := regexp.MustCompile("^[a-zA-Z]+$")
	return re.MatchString(s)
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
		// I am using os.Exit() here, because for some reason I get the "program ran out of input" error
		// In my Windows laptop, however this doesn't happen in my Mac.
		// Instead of os.Exit() we can use return "Bye!" here, and it would work too, I guess!
		os.Exit(0)
	} else if line == "/help" {
		return "The program calculates the sum of numbers"
	}
	return "Unknown command"
}

// getTotal calculates the total result of the postfixExpr infixExpr
func (c Calculator) getTotal() int {
	for _, val := range c.postfixExpr {
		if isNumeric(val) {
			c.stack = append(c.stack, val)
		} else {
			b, a := pop(&c.stack), pop(&c.stack)

			// FOR FUN ONLY TO HANDLE FLOAT INPUTS -- THIS CAN BE REMOVED FROM FINAL SOLUTION
			// Check if 'b' and 'a' are "float strings" like "3.50",
			// Then round them down to the lowest integer value, and then convert them back to strings
			if floatStringToInt(a) != 0 && floatStringToInt(b) != 0 {
				a, b = strconv.Itoa(floatStringToInt(a)), strconv.Itoa(floatStringToInt(b))
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

// getPostfix converts the infix infixExpr to postfixExpr
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
			if prevSym != "" && isNumeric(prevSym) {
				c.postfixExpr = append(c.postfixExpr, pop(&c.postfixExpr)+token)
			} else {
				c.postfixExpr = append(c.postfixExpr, token)
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

		// Check if the entered line has any preceding blank spaces
		if strings.HasPrefix(line, " ") {
			line = strings.TrimSpace(line)
		}

		if len(line) > 0 {
			if checkCommand(line) {
				c.message = getCommand(line)
			} else if checkAssignment(line) {
				c.assign(line)
				continue
			} else {
				if checkParenthesis(line) {
					c.message = "Invalid expression"
				} else {
					// Since a command wasn't issued, reset the c.message variable
					c.message = ""

					// Use splitParenthesis to split any tokens like "(3" or "1)"
					if strings.Contains(line, " ") {
						c.infixExpr = splitParenthesis(strings.Split(line, " "))
					} else {
						c.infixExpr = splitParenthesis(strings.Split(line, ""))
					}

					// Get the postfixExpr and then calculate the result
					c.postfixExpr = c.getPostfix(strings.Join(c.infixExpr, " "))
					c.result = c.getTotal()
				}
			}

			// If a command was issued, print the command message;
			// Otherwise if 'c.postfixExpr' is not nil print the calculated result
			if c.message != "" {
				fmt.Println(c.message)
			} else if c.postfixExpr != nil {
				fmt.Println(c.result)
			}
		}
	}
}

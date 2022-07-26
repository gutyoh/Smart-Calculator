package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var store = make(map[string]int)

var operatorRank = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
	"^": 3,
}

var symbols = []string{"+", "-", "*", "/", "(", ")"}

// count returns the number of times the value x appears in a slice
func count(s []string, x string) int {
	var c int
	for _, a := range s {
		if a == x {
			c++
		}
	}
	return c
}

func split(s, sep string) []string {
	return strings.Split(s, sep)
}

func mapContains(m map[string]int, key string) bool {
	_, ok := m[key]
	return ok
}

func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isAlpha(s string) bool {
	for _, c := range s {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
			return false
		}
	}
	return true
}

func removeSpace(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		}
	}
	return string(rr)
}

// End of Python mock-up functions

func isCommand(s []string) bool {
	if startsWith(s[0], "/") {
		return true
	}
	return false
}

func isAssignment(s string) bool {
	return strings.Contains(s, "=")
}

func assign(line string) string {
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
		if !mapContains(store, value) {
			return "Invalid assignment"
		} else {
			value = strconv.Itoa(store[value])
		}
	}

	store[variable], _ = strconv.Atoi(value)
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

func getSign(symbol string) int {
	if strings.Contains(symbol, "-") {
		if len(symbol)%2 == 0 {
			return 1
		} else {
			return -1
		}
	}
	return 1
}

func getTotal(postfix []string) string {
	var stack []string
	for _, val := range postfix {
		if isNumeric(val) {
			stack = append(stack, val)
		} else {
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			x, _ := strconv.Atoi(a)
			y, _ := strconv.Atoi(b)
			stack = append(stack, strconv.Itoa(evalBinary(x, y, val)))
		}
	}
	if len(stack) > 0 {
		return stack[len(stack)-1]
	}
	return ""
}

func evalBinary(a, b int, op string) int {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	case "^":
		return int(float64(a) * float64(b))
	default:
		return 0
	}
}

func getExpression(line []string) []string {
	var parsedExp []string
	for _, val := range line {
		if isAlpha(val) {
			if mapContains(store, val) {
				val = strconv.Itoa(store[val])
			} else {
				fmt.Println("Unknown variable")
				break
			}
		}
		parsedExp = append(parsedExp, val)
	}
	return parsedExp
}

func getPostfix(line []string) []string {
	var stack []string
	prevSymbol := ""
	postfix := []string{}

	for _, sym := range line {
		if sym == " " {
			continue
		} else if isAlpha(sym) {
			if mapContains(store, sym) {
				prevSymbol = sym
				sym = strconv.Itoa(store[sym])
				postfix = append(postfix, sym)
			} else {
				fmt.Println("Unknown variable")
				break
			}
		} else if isNumeric(sym) {
			if prevSymbol != "" && isNumeric(prevSymbol) {
				// save the last element of 'postfix'
				temp := postfix[len(postfix)-1]
				// remove the last element of 'postfix'
				postfix = postfix[:len(postfix)-1]
				// append the last element of 'postfix' to 'postfix'
				postfix = append(postfix, temp+sym)
			} else {
				postfix = append(postfix, sym)
				prevSymbol = sym
			}
		} else if sliceContains(symbols, sym) {
			if prevSymbol == "" {
				prevSymbol = sym
				// append symbol to stack
				stack = append(stack, sym)
				continue
			}

			if prevSymbol == sym {
				if sym == "+" {
					continue
				} else if sym == "-" {
					prevSymbol = "+"
					// delete the last element of 'stack'
					stack = stack[:len(stack)-1]
					// append symbol to stack
					stack = append(stack, sym)
				} else if sym == "*" || sym == "/" {
					fmt.Println("Invalid expression")
					break
				}
			} else {
				prevSymbol = sym
			}
			stack, postfix = stackOperator(stack, postfix, sym)
		}
	}
	for {
		if len(stack) == 0 {
			break
		}
		postfix = append(postfix, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return postfix
}

func stackOperator(stack, postfix []string, sym string) ([]string, []string) {
	if len(stack) == 0 || stack[len(stack)-1] == "(" || sym == "(" {
		stack = append(stack, sym)
		return stack, postfix
	}
	if sym == ")" {
		for {
			if stack[len(stack)-1] == "(" {
				stack = stack[:len(stack)-1]
				break
			}
			postfix = append(postfix, stack[len(stack)-1])
			stack = stack[:len(stack)-1]
		}
		return stack, postfix
	}

	if higherPrecedence(stack[len(stack)-1], sym) {
		stack = append(stack, sym)
	} else {
		for {
			if len(stack) > 0 && !higherPrecedence(stack[len(stack)-1], sym) {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			} else {
				break
			}
		}
		stack = append(stack, sym)
	}
	return stack, postfix
}

func higherPrecedence(stackPop, sym string) bool {
	if stackPop == "(" || operatorRank[sym] > operatorRank[stackPop] {
		return true
	}
	return false
}

func checkParenthesis(line []string) bool {
	return count(line, "(") != count(line, ")")
}

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		userInput := scanner.Text()

		var output string
		var postfix []string

		if startsWith(userInput, " ") {
			userInput = removeSpace(userInput)
		}

		if len(userInput) > 0 {
			if isCommand(split(userInput, " ")) {
				output = getCommand(userInput)
			} else if isAssignment(userInput) {
				output = assign(userInput)
			} else {
				if checkParenthesis(split(userInput, " ")) {
					output = "Invalid expression"
				} else {
					postfix = getPostfix(getExpression(split(userInput, " ")))
					output = getTotal(postfix)
				}
			}

			if output != "" {
				fmt.Println(output)
			}
		}
	}
}

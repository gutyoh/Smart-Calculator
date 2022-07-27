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

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
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

// repeatedSymbol checks there is more than one "*" or "/" in the line
func repeatedSymbol(line string) bool {
	return strings.Count(line, "*") > 1 || strings.Count(line, "/") > 1
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
		fmt.Println("Bye!")
		os.Exit(0)
	} else if line == "/help" {
		return "The program calculates the sum of numbers"
	}
	return "Unknown command"
}

func getTotal(postfix []string) string {
	var stack []string
	var a string
	for _, val := range postfix {
		if isNumeric(val) {
			stack = append(stack, val)
		} else {
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if len(stack) > 0 {
				a = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
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

func getPostfix(line []string) []string {
	var stack []string
	prevSymbol := ""
	var postfix []string

	for _, sym := range line {
		if repeatedSymbol(sym) {
			fmt.Println("Invalid expression")
			break
		}
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

func checkParenthesis(line string) bool {
	return strings.Count(line, "(") != strings.Count(line, ")")
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
				if checkParenthesis(userInput) {
					output = "Invalid expression"
				} else {

					line := split(userInput, " ")

					newLine := []string{}
					for _, v := range line {
						if strings.HasPrefix(v, "(") {
							newLine = append(newLine, "(")
							newLine = append(newLine, v[1:])
						} else if strings.HasSuffix(v, ")") {
							newLine = append(newLine, v[:len(v)-1])
							newLine = append(newLine, ")")
						} else {
							newLine = append(newLine, v)
						}
					}

					postfix = getPostfix(newLine)
					output = getTotal(postfix)
				}
			}

			if output != "" {
				fmt.Println(output)
			}
		}
	}
}

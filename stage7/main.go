package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var store = make(map[string]int)

func mapContains(m map[string]int, key string) bool {
	_, ok := m[key]
	return ok
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func isCommand(s []string) bool {
	if startsWith(s[0], "/") {
		return true
	}
	return false
}

func isAssignment(s string) bool {
	return strings.Contains(s, "=")
}

func split(s, sep string) []string {
	return strings.Split(s, sep)
}

func assign(line string) string {
	variable, value := func(s []string) (string, string) {
		return s[0], s[1]
	}(func() (elems []string) {
		for _, x := range strings.Split(line, "=") {
			elems = append(elems, strings.TrimSpace(x))
		}
		return
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

func getCommand(s string) string {
	if s == "/exit" {
		return "Bye!"
	} else if s == "/help" {
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

func getTotal(line []string) int {
	sign := 1
	var output []int
	for idx, symbol := range line {
		if idx%2 == 0 {
			symb, _ := strconv.Atoi(symbol)
			output = append(output, sign*symb)
		} else {
			sign = getSign(symbol)
		}
	}

	var sum int
	for _, val := range output {
		sum += val
	}
	return sum
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

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		userInput := scanner.Text()

		var output string

		if len(userInput) > 0 {
			if isCommand(split(userInput, " ")) {
				output = getCommand(userInput)
			} else if isAssignment(userInput) {
				output = assign(userInput)
			} else {
				expression := getExpression(split(userInput, " "))
				output = strconv.Itoa(getTotal(expression))
			}

			if output != "" {
				fmt.Println(output)
			}
		}
	}
}

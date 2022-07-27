package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Calculator struct {
	stack  []string
	result string
	memory map[string]int
}

// mapContains checks if a map contains a specific element
func mapContains(m map[string]int, key string) bool {
	_, ok := m[key]
	return ok
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
		return
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

func (c Calculator) getValue(val string) int {
	if isNumeric(val) {
		return getSign(val) * getSign(val) * getSign(val)
	} else {
		return c.memory[val]
	}
}

func (c Calculator) getExpression(line string) []string {
	var parsedExp []string
	var tokens []string

	if strings.Contains(line, " ") {
		tokens = strings.Split(line, " ")
	} else {
		tokens = strings.Split(line, "")
	}

	for _, token := range tokens {
		if isAlpha(token) {
			if mapContains(c.memory, token) {
				token = strconv.Itoa(c.memory[token])
			} else if mapContains(c.memory, strings.Join(tokens, "")) {
				token = strconv.Itoa(c.memory[strings.Join(tokens, "")])
				parsedExp = append(parsedExp, token)
				break
			} else {
				fmt.Println("Unknown variable")
				break
			}
		}
		parsedExp = append(parsedExp, token)
	}
	return parsedExp
}

func main() {
	var c Calculator
	c.memory = make(map[string]int)

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()

		if len(line) > 0 {
			if checkCommand(line) {
				c.result = getCommand(line)
			} else if checkAssignment(line) {
				c.result = c.assign(line)
			} else {
				expression := c.getExpression(line)
				c.result = strconv.Itoa(getTotal(expression))
			}

			if c.result != "" {
				fmt.Println(c.result)
			}
		}
	}
}

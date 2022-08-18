package main

/*
[Smart Calculator - Stage 4/7: Add subtractions](https://hyperskill.org/projects/74/stages/412/implement)
-------------------------------------------------------------------------------
[Intro to computational thinking](https://hyperskill.org/learn/step/8742)
[Components of computational thinking](https://hyperskill.org/learn/step/8745)
[Functions](https://hyperskill.org/learn/topic/1750)
[Function decomposition](https://hyperskill.org/learn/topic/1893)
[Structs](https://hyperskill.org/learn/topic/1891)
[Public and private scopes](https://hyperskill.org/learn/topic/1894)
[Design principles](https://hyperskill.org/learn/step/8956)
[Type casting and type switching] -- TODO!
[Unicode package] -- TODO!
*/

import (
	"bufio"
	"fmt"
	"log"
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
)

type Expression struct {
	ExpressionType
	Value any
}

// isNumeric checks if all the characters in the string are digits
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

func isSign(token string) bool {
	return token == "+" || token == "-"
}

func parseNumber(line string) (int, int) {
	var stringNum string
	var end, number int

	for _, t := range line {
		token := string(t)
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

func parseSign(line string) (string, int) {
	var sign string
	var end int

	for _, token := range line {
		if !isSign(string(token)) {
			break
		}
		sign += string(token)
	}
	end = len(sign)
	return sign, end
}

// processLine does the actual work of the program:
func processLine(line string) {
	var sign string
	var number, end int
	var expression []Expression

	for len(line) > 0 {
		token := string(line[0])
		switch {
		case token == " ":
			end = 1
		case isNumeric(token):
			number, end = parseNumber(line)
			expression = append(expression, Expression{Number, number})
		case isSign(token):
			sign, end = parseSign(line)
			expression = append(expression, Expression{Sign, sign})
		default:
			return
		}
		line = line[end:]
	}

	// Calculate the expression and output the final result
	fmt.Println(getTotal(expression))
}

// getTotal calculates and returns the total sum of the numbers in the expression slice
func getTotal(expression []Expression) int {
	total, sign := 0, 1
	for _, token := range expression {
		switch token.ExpressionType {
		case Number:
			total += token.Value.(int) * sign
			sign = 1
		case Sign:
			if strings.Count(token.Value.(string), "-")%2 == 1 {
				sign *= -1
			}
		}
	}
	return total
}

func main() {
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
			processLine(line)
		}
	}
}

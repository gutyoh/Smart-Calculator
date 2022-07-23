package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var store = make(map[string]int)

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

type Calculator struct {
}

func isCommand(s string) bool {
	return s[0] == '/'
}

func isAssignment(s string) bool {
	return strings.Contains(s, "=")
}

//func assign(s string) {
//	if
//}

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
	total := 0
	for _, num := range line {
		if isNumeric(num) || (num[0] == '-' && isNumeric(num[1:])) {
			n, err := strconv.Atoi(num)
			if err != nil {
				log.Fatal(err)
			}
			total += n
		} else {
			for _, c := range num {
				if c == '-' {
					total *= -1
				}
			}
		}
	}
	return total
}

func getExpression(line []string) string {
	var parsedExp []string
	tokens := strings.Split(line[0], " ")
	for _, val := range tokens {
		if isAlpha(val) {
			if _, ok := store[val]; ok {
				parsedExp = append(parsedExp, strconv.Itoa(store[val]))
			} else {
				return "Unknown variable"
			}
		}
		parsedExp = append(parsedExp, val)
	}
	return strings.Join(parsedExp, " ")
}

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := strings.Split(scanner.Text(), " ")

		if len(line[0]) == 0 {
			continue
		}

		if isCommand(line[0]) {
			switch line[0] {
			case "/exit":
				fmt.Println("Bye!")
				return
			case "/help":
				fmt.Println("The program calculates the sum of numbers")
			default:
				fmt.Println("Unknown command")
				continue
			}
		}

		if isAssignment(line[0]) {
			fmt.Println("Invalid expression")
			continue
		}

	}
}

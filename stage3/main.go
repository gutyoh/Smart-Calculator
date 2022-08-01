package main

/*
[Smart Calculator - Stage 3/7: Count them all](https://hyperskill.org/projects/74/stages/411/implement)
-------------------------------------------------------------------------------
[Working with slices](https://hyperskill.org/learn/topic/1701)
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()

		if line == "" {
			continue
		} else if line == "/exit" {
			fmt.Println("Bye!")
			return
		} else if line == "/help" {
			fmt.Println("The program calculates the sum of numbers")
		} else if strings.HasPrefix(line, "/") || strings.Contains(line, "=") {
			// If the expression is any other command or a wrong command like "/ exit", then continue:
			continue
		} else {
			total := 0
			tokens := strings.Split(line, " ")
			for _, num := range tokens {
				n, err := strconv.Atoi(num)
				if err != nil {
					log.Fatal(err)
				}
				total += n
			}
			fmt.Println(total)
		}

		//switch line[0] {
		//case "":
		//	continue
		//case "/exit":
		//	fmt.Println("Bye!")
		//	return
		//case "/help":
		//	fmt.Println("The program calculates the sum of numbers")
		//default:
		//	total := 0
		//	for _, num := range line {
		//		n, err := strconv.Atoi(num)
		//		if err != nil {
		//			log.Fatal(err)
		//		}
		//		total += n
		//	}
		//	fmt.Println(total)
		//}
	}
}

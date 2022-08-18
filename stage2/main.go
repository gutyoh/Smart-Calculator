package main

/*
[Smart Calculator - Stage 2/7: 2+2+](https://hyperskill.org/projects/74/stages/410/implement)
-------------------------------------------------------------------------------
[Slices](https://hyperskill.org/learn/topic/1672)
[Slice expressions](https://hyperskill.org/learn/topic/2207)
[Control statements](https://hyperskill.org/learn/topic/1728)
[Loops](https://hyperskill.org/learn/topic/1531)
[Advanced Input](https://hyperskill.org/learn/topic/2027)
[Errors](https://hyperskill.org/learn/topic/1795)
[Operations with strings](https://hyperskill.org/learn/topic/2023)
[Parsing data from strings](https://hyperskill.org/learn/topic/1955)
[Type conversion and overflow](https://hyperskill.org/learn/step/18710)
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
		} else {
			tokens := strings.Split(line, " ")
			if len(tokens) > 1 {
				x, err := strconv.Atoi(tokens[0])
				if err != nil {
					log.Fatal(err)
				}
				y, err := strconv.Atoi(tokens[1])
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(x + y)
			} else {
				x, err := strconv.Atoi(tokens[0])
				if err != nil {
					log.Fatal(err)
				}
				y := 0
				fmt.Println(x + y)
			}
		}
	}
}

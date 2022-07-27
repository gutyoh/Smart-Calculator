package main

/*
[Smart Calculator - Stage 2/7: 2+2+](https://hyperskill.org/projects/74/stages/410/implement)
-------------------------------------------------------------------------------
[Slices](https://hyperskill.org/learn/topic/1672)
[Control statements](https://hyperskill.org/learn/topic/1728)
[Loops](https://hyperskill.org/learn/topic/1531)
[Advanced Input](https://hyperskill.org/learn/topic/2027)
[Errors](https://hyperskill.org/learn/topic/1795)
[Operations with strings](https://hyperskill.org/learn/topic/2023)
[Parsing data from strings](https://hyperskill.org/learn/topic/1955)
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
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "/exit" {
			fmt.Println("Bye!")
			break
		}
		line := strings.Split(scanner.Text(), " ")

		if len(line) > 1 {
			x, _ := strconv.Atoi(line[0])
			y, _ := strconv.Atoi(line[1])
			fmt.Println(x + y)
		} else if line[0] == "" {
			continue
		} else {
			x, err := strconv.Atoi(line[0])
			if err != nil {
				log.Fatal(err)
			}
			y := 0
			fmt.Println(x + y)
		}
	}
}

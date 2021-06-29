package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter numbers")
	// when you hit enter - the data is read
	scanner.Scan()
	// this will read data from the scanner
	scanParams := scanner.Text()

	evenSum := EvenSum(scanParams)
	fmt.Println(evenSum)
}

func EvenSum(scanParams string) (evenSum int) {
	numberSlice := strings.Split(scanParams, ",")
	for _, v := range numberSlice {
		number, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println("You has some error in your input data")
			continue
		}
		//Case where we check even for number
		if number > 0 && number%2 == 0 {
			evenSum += number
		}
	}
	return
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type CardNumber []int

func ValidateCardNumber(numberString string) (CardNumber, error) {
	numberString = strings.ReplaceAll(numberString, " ", "")
	numberStringSlice := strings.Split(numberString, "")
	if len(numberStringSlice) != 16 {
		return nil, fmt.Errorf("Number not valid\n")
	}

	var sum int
	var cardNumbers CardNumber
	for i, v := range numberStringSlice {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("Input data not valid\n")
		}
		cardNumbers = append(cardNumbers, num)
		if i%2 == 0 {
			num *= 2
		}

		sum += num / 10
		sum += num % 10
	}
	if sum%10 != 0 {
		return nil, fmt.Errorf("Card number not valid\n")
	}
	return cardNumbers, nil
}

func (c CardNumber) Print(wtiter io.Writer) {
	printBuilder := strings.Builder{}
	for i, n := range c {
		if i%4 == 0 && i != 0 {
			printBuilder.WriteString(" ")
		}
		if i+4 < len(c) {
			printBuilder.WriteRune('*')
			continue
		}
		printBuilder.WriteString(strconv.Itoa(n))
	}
	fmt.Fprintln(wtiter, printBuilder.String())
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter credit card")
	// when you hit enter - the data is read
	scanner.Scan()
	// this will read data from the scanner
	scanData := scanner.Text()

	cardNumbers, err := ValidateCardNumber(scanData)
	if err != nil {
		log.Fatal(err)
	}
	cardNumbers.Print(os.Stdout)
}

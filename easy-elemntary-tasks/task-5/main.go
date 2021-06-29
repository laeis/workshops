package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	MinDigit = 100000
	MaxDigit = 999999
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter min and max six-digit numbers")
	// when you hit enter - the data is read
	scanner.Scan()
	// this will read data from the scanner
	scanData := scanner.Text()

	min, max, err := ValidateInput(scanData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("EasyFormula: %d \n", EasyFormula(min, max))
	fmt.Printf("HardFormula: %d \n", HardFormula(min, max))
}

func ValidateInput(data string) (min, max int, err error) {
	scanSlice := strings.Split(data, ",")
	for _, v := range scanSlice {
		num, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return min, max, err
		}
		if num < MinDigit || num > MaxDigit {
			return min, max, fmt.Errorf("Input numbers contains exactly 6 digits, got: %d\n", num)
		}
		if num > max {
			min, max = max, num
		}
		if num < max {
			min = num
		}
	}
	return
}

func EasyFormula(min, max int) (amount int) {
	for i := min; i <= max; i++ {
		t := []rune(strconv.Itoa(i))
		var first, second int
		for i := 0; i < len(t)/2; i++ {
			f, err := strconv.Atoi(string(t[i]))
			if err != nil {
				log.Fatal(err)
			}
			first += f
		}
		for i := len(t) / 2; i < len(t); i++ {
			s, err := strconv.Atoi(string(t[i]))
			if err != nil {
				log.Fatal(err)
			}
			second += s
		}
		if first == second {
			amount += 1
		}
	}
	return amount
}

func HardFormula(min, max int) (amount int) {
	for i := min; i <= max; i++ {
		t := []rune(strconv.Itoa(i))
		var even, odd int
		for _, r := range t {
			n, err := strconv.Atoi(string(r))
			if err != nil {
				log.Fatal(err)
			}
			if n%2 == 0 {
				even += n
			} else {
				odd += n
			}

		}
		if even == odd {
			amount += 1
		}
	}
	return amount
}

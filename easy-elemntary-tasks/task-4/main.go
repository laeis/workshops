package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const Min = 10

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter number")
	// when you hit enter - the data is read
	scanner.Scan()
	// this will read data from the scanner
	scanData := scanner.Text()

	palindromes := FindPalindromes(scanData)
	fmt.Println(palindromes)
}

type possiblePalindrome []rune

func reverseSliceRune(r []rune) []rune {
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return r
}

func FindPalindromes(s string) string {
	emptyResult := "0"
	number, err := strconv.Atoi(s)
	if err != nil && number < Min {
		return emptyResult
	}
	positionMap := make(map[rune]int)
	pCollection := []string{}
	r := []rune(s)
	for i, v := range r {
		if _, ok := positionMap[v]; !ok {
			positionMap[v] = i
		} else {
			pPalindrome := possiblePalindrome(r[positionMap[v] : i+1])
			if string(pPalindrome) == string(reverseSliceRune(pPalindrome)) {
				pCollection = append(pCollection, string(pPalindrome))
			} else {
				positionMap[v] = i
			}
		}
	}
	if len(pCollection) == 0 {
		return emptyResult
	}
	return strings.Join(pCollection, ",")
}

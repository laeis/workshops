package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const ParamsLength = 3
const Space = "  "
const (
	Height = iota
	Width
	Symbol
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("For print a chessboard with the specified dimensions use instruction below")
	fmt.Println("Input parameters: height, width, symbol to print and press enter")
	// when you hit enter - the data is read
	scanner.Scan()
	// this will read data from the scanner
	scanParams := scanner.Text()
	paramsSlice := strings.Split(scanParams, ",")
	if len(paramsSlice) != ParamsLength {
		errorMessage := fmt.Errorf("Wrong parameters! Waiting for: height, width, symbol got %s \n", scanParams)
		printError(errorMessage)
	}
	height, convertErr := strconv.Atoi(paramsSlice[Height])
	if convertErr != nil {
		printError(formatParameterTypeError("height", "number", paramsSlice[Height]))
	}
	width, convertErr := strconv.Atoi(paramsSlice[Width])
	if convertErr != nil {
		printError(formatParameterTypeError("width", "number", paramsSlice[Width]))
	}
	if len(paramsSlice[Symbol]) > 1 {
		printError(formatParameterTypeError("symbol", "one symbol", paramsSlice[Symbol]))
	}
	symbol := paramsSlice[Symbol]

	chessboard := prepareChessBoard(height, width, symbol)
	// will print it to stdin
	fmt.Println(chessboard.String())
}

func prepareChessBoard(height, width int, symbol string) fmt.Stringer {
	chessboardBuf := &strings.Builder{}
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if h%2 != 0 {
				chessboardBuf.WriteString(Space)
			}
			chessboardBuf.WriteString(symbol)
			if h%2 == 0 {
				chessboardBuf.WriteString(Space)
			}
		}
		if h+1 != height {
			chessboardBuf.WriteString("\n")
		}
	}
	return chessboardBuf
}

func formatParameterTypeError(name, pType, got string) error {
	return fmt.Errorf("Wrong %s parameters! %s must be %s got '%s' \n", name, strings.Title(name), pType, got)
}

func printError(errorMessage error) {
	fmt.Print(errorMessage)
	os.Exit(1)
}

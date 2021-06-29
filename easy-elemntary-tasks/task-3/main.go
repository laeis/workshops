package main

import (
	"bytes"
	"fmt"
	"strconv"
)

func main() {
	fi := fibonacci()
	fmt.Println(fi(100))
}

func fibonacci() func(max int) string {
	current, next := 0, 1
	return func(max int) string {
		buf := bytes.Buffer{}
		for current <= max {
			buf.WriteString(strconv.Itoa(current))
			if next <= max {
				buf.WriteRune(',')
			}
			current, next = next, current+next
		}
		return buf.String()
	}
}

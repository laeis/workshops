package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestCardNumber_Print(t *testing.T) {
	cases := []struct {
		name string
		card CardNumber
		want string
	}{
		{
			"Number with 16 numbers",
			CardNumber{4, 5, 3, 9, 1, 4, 8, 8, 0, 3, 4, 3, 6, 4, 6, 7},
			"**** **** **** 6467\n",
		},
		{
			"Number with 14 numbers",
			CardNumber{4, 5, 3, 9, 1, 4, 8, 8, 0, 3, 4, 3, 6, 4},
			"**** **** **43 64\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			c.card.Print(buffer)
			got := buffer.String()
			if c.want != got {
				t.Errorf("got %q want %q\n", got, c.want)
			}
		})
	}
}

func TestValidateCardNumber(t *testing.T) {
	cases := []struct {
		name    string
		card    string
		isError bool
		result  CardNumber
	}{
		{
			"Send short number",
			"539 1488 0343 646",
			true,
			nil,
		},
		{
			"Send wrong data in string short number",
			"539 1488 034r 6467",
			true,
			nil,
		},
		{
			"Send wrong number",
			"539 1488 034r 6468",
			true,
			nil,
		},
		{
			"Send valid card number",
			"4539 1488 0343 6467",
			false,
			CardNumber{4, 5, 3, 9, 1, 4, 8, 8, 0, 3, 4, 3, 6, 4, 6, 7},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := ValidateCardNumber(c.card)
			if c.isError && err == nil {
				t.Error("Assert Error and dont have")
			}
			if !reflect.DeepEqual(result, c.result) {
				t.Errorf("got %q want %q\n", result, c.result)
			}
		})
	}
}

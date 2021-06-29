package main

import "testing"

func TestValidateInput(t *testing.T) {
	cases := []struct {
		name    string
		isError bool
		input   string
	}{
		{
			"6 digits error for input",
			true,
			"12012, 120123",
		},
		{
			"6 digits error for second number in input",
			true,
			"120121, 1201113",
		},
		{
			"Not a number in input",
			true,
			"120d121, 120111",
		},
		{
			"Valid input",
			false,
			"120123, 320320",
		},
	}
	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			_, _, err := ValidateInput(v.input)

			if v.isError && err == nil {
				t.Error("Waiting for error and dont have it")
			}
			if !v.isError && err != nil {
				t.Error("Dont Waiting for error and have it")
			}
		})
	}
}

func TestHardFormula(t *testing.T) {
	want := 5790
	got := HardFormula(120123, 320320)
	if got != want {
		t.Errorf("Waiting for %d got %d", want, got)
	}
}

func TestEasyFormula(t *testing.T) {
	want := 1
	got := EasyFormula(100000, 100002)
	if got != want {
		t.Errorf("Waiting for %d got %d", want, got)
	}
}

package main

import "testing"

func TestFindPalindromes(t *testing.T) {
	cases := []struct {
		name  string
		want  string
		input string
	}{
		{
			"Send not a number",
			"0",
			"text",
		},
		{
			"Send less than 10",
			"0",
			"9",
		},
		{
			"Send not palindrome",
			"0",
			"10",
		},
		{
			"Send string with palindrome",
			"44,3443",
			"1234437",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := FindPalindromes(c.input)
			if got != c.want {
				t.Errorf("Waiting for %s got %s", c.want, got)
			}
		})
	}
}

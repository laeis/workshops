package main

import "testing"

func TestFindEvenSum(t *testing.T) {
	t.Run("Just sum even numbers", func(t *testing.T) {
		want := 60
		got := EvenSum("30,30")
		if got != want {
			t.Errorf("got %d want %d\n", got, want)
		}
	})

	t.Run("Sum even numbers with has negative value", func(t *testing.T) {
		want := 30
		got := EvenSum("30,-30")
		if got != want {
			t.Errorf("got %d want %d\n", got, want)
		}
	})
}

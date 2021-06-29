package main

import "testing"

func TestFibonacci(t *testing.T) {
	want := "0,1,1,2,3,5,8,13,21,34,55,89"
	got := fibonacci()(100)
	if got != want {
		t.Errorf("got %q want %q\n", got, want)
	}
}

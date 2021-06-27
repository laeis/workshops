package main

import (
	"testing"
)

func TestPrepareChessBoard(t *testing.T) {
	want := `^  ^  ^  ^  ^  ^
  ^  ^  ^  ^  ^  ^
^  ^  ^  ^  ^  ^
  ^  ^  ^  ^  ^  ^`
	got := prepareChessBoard(4, 6, "^")

	if want != got.String() {
		t.Errorf("Wrong chessboard, want: \n%s \n got: \n%s \n", want, got)
	}

}

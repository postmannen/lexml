package main

import "testing"

const succeed = "*** \u2713 *** "
const failed = "*** \u2717 **** "

func TestTokenSendConsole(t *testing.T) {
	t.Log(succeed, tokenSendConsole(tokenEOF, "EOF"))
	t.Log(succeed, tokenSendConsole(tokenArgumentFound, "some-argument"))
}

func TestFindLettersBetween(t *testing.T) {
	s := "THIS IS SOME STRING TO DO SOME CHECKING ON!"
	w := findLettersBetween(s, 0, 4)
	if w != "THIS" {
		t.Errorf("%v Failed finding THIS, found : %v\n", failed, w)
	} else {
		t.Logf("%v Succeeded in finding THIS \n", succeed)
	}
}

func TestFindChrPositionAfter(t *testing.T) {
	s := "THIS IS SOME STRING TO DO SOME CHECKING ON!"
	n := findChrPositionAfter(s, 'E', 5)

	if n == 11 {
		t.Log(succeed, "found character at position = ", n)
	} else {
		t.Error(failed, "failed to find correct character position, got = ", n)
	}
}

func TestFindChrPositionBefore(t *testing.T) {
	s := "THIS IS SOME STRING TO DO SOME CHECKING ON!"
	n := findChrPositionBefore(s, 'H', 5)

	if n == 1 {
		t.Log(succeed, "found character at position = ", n)
	} else {
		t.Error(failed, "failed to find correct character position, got = ", n)
	}
}

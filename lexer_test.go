package main

import (
	"strings"
	"testing"
)

const succeed = "*** \u2713 *** "
const failed = "*** \u2717 **** "

func TestTokenSendConsole(t *testing.T) {
	//TODO
}

func TestRealRun(t *testing.T) {
	rdr := strings.NewReader(`
	<shiporder orderid="889923"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:noNamespaceSchemaLocation="shiporder.xsd">
		<orderperson>John Smith</orderperson>
		<shipto>
			<name>Ola Nordmann</name>
		</shipto>
		<cmd name="TakeOff" id="1"></cmd>
	</shiporder>
	<comment
		title="Take off"
		desc="Ask the drone to take off.\n
		On the fixed wings (such as Disco): not used except to cancel a land."
		support="0901;090c;090e"
		result="On the quadcopters: the drone takes off if its [FlyingState](#1-4-1) was landed.\n
		On the fixed wings, the landing process is aborted if the [FlyingState](#1-4-1) walanding.\n
		Then, event [FlyingState](#1-4-1) is triggered."/>
	`)

	tokenChan = make(chan token)

	wg.Add(1)
	go readToken()
	defer wg.Wait()

	lex := newLexer(rdr, tokenOutputType(1))
	lex.lexStart()

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

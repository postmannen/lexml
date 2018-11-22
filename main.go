package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const fileName = "ardrone3.xml"

type lexer struct {
	fileReader      *bufio.Reader //buffer used for reading file lines.
	currentLineNR   int           //the line nr. being read
	currentLine     string        //the current single line we mainly work on.
	nextLine        string        //the line after the current line.
	EOF             bool          //used for proper ending of readLine method.
	workingLine     string        //the line being worked on. Can be a collection of several lines.
	workingPosition int           //the chr position we're currently working at in line
}

//newLexer will return a *lexer type, it takes a pointer to a file as input.
func newLexer(fh *os.File) *lexer {
	return &lexer{
		fileReader:    bufio.NewReader(fh),
		currentLineNR: -1,
	}
}

//stateFunc is the defenition of a state function.
type stateFunc func() stateFunc

//readLine will allways read the next line, and move the previous next line
// into current line on next run. All spaces and carriage returns are removed.
// Since we are working on the line that was read on the prevoius run, we will
// set l.EOF to true as our exit parameter if error==io.EOF, so the whole
// function is called one more time if error=io.EOF so we also get the last line
// of the file moved into l.currentLine.
func (l *lexer) readLine() stateFunc {
	l.workingPosition = 0
	l.currentLineNR++
	l.currentLine = l.nextLine
	line, _, err := l.fileReader.ReadLine()
	if err != nil {
		if l.EOF {
			return nil
		}
		if err == io.EOF {
			l.EOF = true
		}
	}
	l.nextLine = strings.TrimSpace(string(line))
	return l.checkItemInLine
}

//start will start the reading of lines from file, and then kickstart it all
// by running the returned function inside the for loop.
// Since all methods return a new method to be executed on the next run, we
// will check if the current ran method returned nil instead of a new method
// to exit.
func (l *lexer) start() {
	fn := l.readLine()
	for {
		fn = fn()
		if fn == nil {
			log.Println("done with for loop")
			break
		}
	}
}

//print will print the current working line.
func (l *lexer) print() stateFunc {
	fmt.Println(l.currentLineNR, l.currentLine)
	fmt.Println("-------------------------------------------------------------------------")
	return l.readLine()
}

//checkItemInLine will work itselves one character position at a time the string line,
// and do some action based on the type of character found.
// If string is blank, or end of string is reached we exit, and read a new line.
//
// While working on the same line any state functions used  should return back here
// until the end of the line is reached.
//
func (l *lexer) checkItemInLine() stateFunc {
	if len(l.currentLine) == 0 {
		log.Println("NOTE: blank line, reading the next line")
		return l.readLine()
	}

	//Check what kind of line it is. If it is a start tag with close on same line,
	// end tag, or a comment line
	//
	if strings.HasPrefix(l.currentLine, "<") && strings.HasSuffix(l.currentLine, ">") {
		fmt.Println(" ***HAS BOTH START AND END BRACKET, Normal tag line ***")
	}
	if strings.HasPrefix(l.currentLine, "<") && !strings.HasSuffix(l.currentLine, ">") {
		fmt.Println(" ***HAS START, BUT NO END BRACKET, TAG CONTINUES ON NEXT LINE ***")
	}
	if !strings.HasPrefix(l.currentLine, "<") && !strings.HasSuffix(l.currentLine, ">") {
		fmt.Println(" ***HAS NO START, NO END BRACKET, PROBABLY COMMENT, ALSO TAG CONTINUES ON NEXT LINE ***")
	}

	//---NB: Here the line should be complete, and concatenated by others if needed.

	l.workingLine = l.currentLine
	//Check all the individual characters of the string
	//
	for l.workingPosition < len(l.workingLine) {
		switch l.workingLine[l.workingPosition] {
		case '<':
			fmt.Println("------FOUND START BRACKET CHR--------")
		case '>':
			fmt.Println("------FOUND END BRACKET CHR----------")
		case '=':
			fmt.Println("------FOUND EQUAL SIGN CHR----------")
		}

		l.workingPosition++
	}

	return l.print
}

func main() {
	fh, err := os.Open(fileName)
	if err != nil {
		log.Println("Error: opening file: ", err)
	}

	lex := newLexer(fh)
	lex.start()

}

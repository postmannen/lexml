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
	fileReader  *bufio.Reader //buffer used for reading file lines.
	currentLine string        //the current line we mainly work on.
	nextLine    string        //the line next to current line.
	EOF         bool          //used for proper ending of readLine method.
}

//newLexer will return a *lexer type, it takes a pointer to a file as input.
func newLexer(fh *os.File) *lexer {
	return &lexer{
		fileReader: bufio.NewReader(fh),
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
	return l.print
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
	fmt.Println(l.currentLine)
	return l.readLine()
}

func main() {
	fh, err := os.Open(fileName)
	if err != nil {
		log.Println("Error: opening file: ", err)
	}

	lex := newLexer(fh)
	lex.start()

}

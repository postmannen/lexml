/* The idea here is to create a lexer for the parrot.xml file which
holds the definition of the protocol used to control the Parrot Bebop 2 drone.
The lexer will be build't by having one main run function who executes a function,
and get a new function in return, that again will be executed next.
The program will be build't up by many smaller functions who server one single purpose
in the production line, and they know what function to return next based on they're own
simple logic.
*/
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
	startFound      bool          //used to tell if a start was found so we can separate descriptions (which have to tags) from the rest while making the lexer.workingLine.
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

//lexReadLine will allways read the next line, and move the previous next line
// into current line on next run. All spaces and carriage returns are removed.
// Since we are working on the line that was read on the prevoius run, we will
// set l.EOF to true as our exit parameter if error==io.EOF, so the whole
// function is called one more time if error=io.EOF so we also get the last line
// of the file moved into l.currentLine.
func (l *lexer) lexReadLine() stateFunc {
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
	return l.lexCheckLineType
}

//lexStart will start the reading of lines from file, and then kickstart it all
// by running the returned function inside the for loop.
// Since all methods return a new method to be executed on the next run, we
// will check if the current ran method returned nil instead of a new method
// to exit.
func (l *lexer) lexStart() {
	fn := l.lexReadLine()
	for {
		fn = fn()
		if fn == nil {
			log.Println("done with for loop")
			break
		}
	}
}

//lexPrint will print the current working line.
func (l *lexer) lexPrint() stateFunc {
	fmt.Println(l.currentLineNR, l.workingLine)
	fmt.Println("-------------------------------------------------------------------------")
	//We reset the working line here
	l.workingLine = ""
	return l.lexReadLine()
}

//lexChr will work itselves one character position at a time the string line,
// and do some action based on the type of character found.
//
func (l *lexer) lexChr() stateFunc {
	//Check all the individual characters of the string
	//
	for l.workingPosition < len(l.workingLine) {
		switch l.workingLine[l.workingPosition] {
		case '<':
			//fmt.Println("------FOUND START BRACKET CHR--------")
			//TODO: Do something...........................
		case '>':
			//fmt.Println("------FOUND END BRACKET CHR----------")
			//TODO: Do something...........................
		case '=':
			//fmt.Println("------FOUND EQUAL SIGN CHR----------")
			//TODO: Do something...........................
		}

		l.workingPosition++
	}

	return l.lexPrint
}

//lexCheckLineType checks what kind of line we are dealing with. If the line belongs
// together with the line following after, the lines will be combined into a single
// string
// If string is blank, or end of string is reached we exit, and read a new line.
func (l *lexer) lexCheckLineType() stateFunc {
	// If the line is blank, return and read a new line
	if len(l.currentLine) == 0 {
		log.Println("NOTE ", l.currentLineNR, ": blank line, getting out and reading the next line")
		return l.lexReadLine
	}

	start := strings.HasPrefix(l.currentLine, "<")
	end := strings.HasSuffix(l.currentLine, ">")
	nextLineStart := strings.HasPrefix(l.nextLine, "<")

	//set the workingLine = currentLine and go directly to lexing.
	if start && end {
		fmt.Println(" ***TAG ", l.currentLineNR, ": HAS START AND END BRACKET, Normal tag line ***")
		l.startFound = false
		l.workingLine = l.currentLine
		return l.lexChr
	}

	// This indicates this is the first line with a start tag, and the rest are on the following lines.
	// Set initial workingLine=currentLine, and read another line. We set l.startfound to true, to signal
	// that we want to add more lines later to the current working line.
	if start && !end {
		fmt.Println(" ***TAG ", l.currentLineNR, " : start && !end, CONTINUES ON NEXT LINE ***")
		l.startFound = true
		l.workingLine = l.workingLine + " " + l.currentLine
		return l.lexReadLine
	}

	// This indicates we have found a start earlier, and that we need to add this currentLine to the
	// existing content of the workingLine, and read another line
	if !start && !end && l.startFound {
		fmt.Println(" ***TAG ", l.currentLineNR, " : !start && !end && l.startFound, CONTINUES ON NEXT LINE ***")
		l.workingLine = l.workingLine + " " + l.currentLine
		return l.lexReadLine
	}

	//This should indicate that we found the last line of several that have to be combined
	if !start && end && l.startFound {
		fmt.Println(" ***TAG ", l.currentLineNR, " : !start && !end && l.startFound, CONTINUES ON NEXT LINE ***")
		l.workingLine = l.workingLine + " " + l.currentLine
		return l.lexPrint
	}

	//Description line
	if !start && !end && !l.startFound && !nextLineStart {
		fmt.Println(" ***DESC", l.currentLineNR, ": !start && !end && !l.startFound && !nextLineStart, CONTINUES ON NEXT LINE ***")
		l.workingLine = l.workingLine + " " + l.currentLine
		return l.lexReadLine
	}

	//Description line
	if !start && !end && !l.startFound && nextLineStart {
		fmt.Println(" ***DESC", l.currentLineNR, " : !start && !end && !l.startFound && nextLineStart ***")
		l.workingLine = l.workingLine + " " + l.currentLine
		return l.lexPrint
	}

	l.startFound = false
	// The print and return below should optimally never happen.
	// Check it's output to figure what is not detected in the if's above.
	fmt.Println("DEBUG: *uncaught line!! :", l.currentLineNR, l.workingLine)
	return l.lexPrint
}

func main() {
	fh, err := os.Open(fileName)
	if err != nil {
		log.Println("Error: opening file: ", err)
	}

	lex := newLexer(fh)
	lex.lexStart()

}

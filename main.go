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
	"sync"
)

const fileName = "ardrone3.xml"

const (
	chrAngleStart = '<'
	chrAngleEnd   = '>'
	chrEqual      = '='
)

type lexer struct {
	fileReader      *bufio.Reader //buffer used for reading file lines.
	currentLineNR   int           //the line nr. being read
	currentLine     string        //the current single line we mainly work on.
	nextLine        string        //the line after the current line.
	EOF             bool          //used for proper ending of readLine method.
	workingLine     string        //the line being worked on. Can be a collection of several lines.
	workingPosition int           //the chr position we're currently working at in line
	firstLineFound  bool          //used to tell if a start tag was found so we can separate descriptions (which have to tags) from the rest while making the lexer.workingLine.
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

//lexReadFileLine will allways read the next line, and move the previous next line
// into current line on next run. All spaces and carriage returns are removed.
// Since we are working on the line that was read on the prevoius run, we will
// set l.EOF to true as our exit parameter if error==io.EOF, so the whole
// function is called one more time if error=io.EOF so we also get the last line
// of the file moved into l.currentLine.
func (l *lexer) lexReadFileLine() stateFunc {
	l.workingPosition = 0
	l.currentLineNR++
	l.currentLine = l.nextLine
	line, _, err := l.fileReader.ReadLine()
	if err != nil {
		if l.EOF {
			close(tokenChan)
			log.Println("closing channel, EOF of file")
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

	fn := l.lexReadFileLine()
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
	//We reset variables here, since this is the last link in the chain of functions.
	l.workingLine = ""
	return l.lexReadFileLine()
}

//lexLineContent will work itselves one character position at a time the string line,
// and do some action based on the type of character found.
// Checking if a line is done with lexing is done here by checking if working
// position < len. If greater we're done lexing the line, and can continue with
// the next operation.
//
//All functions returned from this function will return back here to continue the iteration
// of the specific line until the end of line, then it will chose the last return at the
// bottom of the function, and get out of the looping back into this function for the specific
// line
//
func (l *lexer) lexLineContent() stateFunc {
	//Check all the individual characters of the string
	//
	for l.workingPosition < len(l.workingLine) {
		switch l.workingLine[l.workingPosition] {
		case chrAngleStart:
			fmt.Println("------FOUND START BRACKET CHR--------")
			return l.lexTagName //find tag name

		case chrAngleEnd:
			fmt.Println("------FOUND END BRACKET CHR----------")
			//TODO: Do something...........................
		case chrEqual:
			fmt.Println("------FOUND EQUAL SIGN CHR----------")
			return l.lexTagArguments
		}

		l.workingPosition++
	}

	return l.lexPrint
}

//lexTagArguments will pick out all the arguments "arg=value".
func (l *lexer) lexTagArguments() stateFunc {
	p1 := findChrPositionBefore(l.workingLine, ' ', l.workingPosition)
	arg := findLettersBetween(l.workingLine, p1, l.workingPosition)
	fmt.Printf("---------------Found argument : %v \n", arg)

	p2 := findChrPositionAfter(l.workingLine, ' ', l.workingPosition)
	value := findLettersBetween(l.workingLine, l.workingPosition+1, p2)
	fmt.Printf("---------------Found argument value : %v \n", value)

	l.workingPosition++
	return l.lexLineContent
}

//findChrPositionBefore .
// Searches backwards in a string from a given positions,
// for the first occurence of a character.
//
func findChrPositionBefore(s string, preChr byte, origChrPosition int) (preChrPosition int) {
	p := origChrPosition
	//find the first space before the preceding word
	for {
		p--
		if p < 0 {
			log.Println("Found no space before the equal sign, reached the beginning of the line")
			break
		}
		if s[p] == preChr {
			preChrPosition = p
			break
		}
	}

	//Will return the position of the prior occurance of the a the character
	return
}

//findChrPositionAfter .
// Searches forward in a string from a given positions,
// for the first occurence of a character after it.
//
func findChrPositionAfter(s string, preChr byte, origChrPosition int) (nextChrPosition int) {
	p := origChrPosition
	//find the first space before the preceding word
	for {
		p++
		if p > len(s)-1 {
			log.Println("Found no space before the equal sign, reached the end of the line")
			break
		}

		//When value is the last thing in a line, it will be followed by '>' and not a space.
		if s[p] == preChr || s[p] == '>' {
			nextChrPosition = p
			break
		}
	}

	//will return the preceding chr's positions found
	return
}

//findLettersBetween
// takes a string, and two positions given as slices as input,
// and returns a slice of string with the words found.
//'
func findLettersBetween(s string, firstPosition int, secondPosition int) (word string) {
	letters := []byte{}
	//as long as first pos'ition is lower than second position....
	for firstPosition < secondPosition {
		letters = append(letters, s[firstPosition])
		firstPosition++
	}
	word = string(letters)
	word = strings.Trim(word, "\"")

	return
}

func (l *lexer) lexTagName() stateFunc {
	var start bool
	var end bool
	var tn []byte

	l.workingPosition++
	//....check if there is a / following the <, then it is an end tag.
	if l.workingLine[l.workingPosition] == '/' {
		fmt.Println("--------FOUND / AFTER <", l.currentLineNR)
		end = true
		l.workingPosition++
	}

	for {
		//look for space, the name ends where the space is.
		if l.workingLine[l.workingPosition] == ' ' {
			fmt.Println("---------FOUND SPACE", l.currentLineNR)
			fmt.Printf("start = %v, end = %v \n", start, end)
			break
		}
		//End tags dont have any attributes, so the '>' will come directly after the tag name.
		// The name ends where the '>' is.
		if l.workingLine[l.workingPosition] == '>' {
			fmt.Println("---------FOUND '>', WHICH INDICATED AN END TAG", l.currentLineNR)
			fmt.Printf("start = %v, end = %v \n", start, end)
			break
		}
		//if none of the above, we can safely add the chr to the slice
		tn = append(tn, l.workingLine[l.workingPosition])
		l.workingPosition++
	}

	fmt.Printf("--- Found tag name '%v'\n", string(tn))
	tokenChan <- token{tokenType: tokenStartTag, tokenText: string(tn)}

	//we return lexLineContent since we know we want to check if there is more to do with the line
	return l.lexLineContent
}

//lexCheckLineType checks what kind of line we are dealing with. If the line belongs
// together with the line following after, the lines will be combined into a single
// string to make it clearer what belongs to a given tag while lexing.
// If string is blank, or end of string is reached we exit, and read a new line.
func (l *lexer) lexCheckLineType() stateFunc {
	// If the line is blank, return and read a new line
	if len(l.currentLine) == 0 {
		log.Println("NOTE ", l.currentLineNR, ": blank line, getting out and reading the next line")
		return l.lexReadFileLine
	}

	start := strings.HasPrefix(l.currentLine, "<")
	end := strings.HasSuffix(l.currentLine, ">")
	nextLineStart := strings.HasPrefix(l.nextLine, "<")

	//TAG: set the workingLine = currentLine and go directly to lexing.
	if start && end {
		fmt.Println(" ***TAG ", l.currentLineNR, ": HAS START AND END BRACKET, Normal tag line ***")
		l.firstLineFound = false
		l.workingLine = l.currentLine
		return l.lexLineContent
	}

	// TAG: This indicates this is the first line with a start tag, and the rest are on the following lines.
	// Set initial workingLine=currentLine, and read another line. We set l.firstLineFound to true, to signal
	// that we want to add more lines later to the current working line.
	if start && !end {
		fmt.Println(" ***TAG ", l.currentLineNR, " : start && !end, CONTINUES ON NEXT LINE ***")
		l.firstLineFound = true
		l.workingLine = l.workingLine + " " + l.currentLine
		return l.lexReadFileLine
	}

	// TAG: This indicates we have found a start earlier, and that we need to add this currentLine to the
	// existing content of the workingLine, and read another line
	if !start && !end && l.firstLineFound {
		fmt.Println(" ***TAG ", l.currentLineNR, " : !start && !end && l.firstLineFound, CONTINUES ON NEXT LINE ***")
		l.workingLine = l.workingLine + " " + l.currentLine
		return l.lexReadFileLine
	}

	//TAG: This should indicate that we found the last line of several that have to be combined
	if !start && end && l.firstLineFound {
		fmt.Println(" ***TAG ", l.currentLineNR, " : !start && !end && l.firstLineFound, DONE COMBINING LINES ***")
		l.workingLine = l.workingLine + " " + l.currentLine
		l.firstLineFound = false //end found, set firstLineFound back to false to be ready for finding new tag.
		return l.lexLineContent
	}

	//Description line. These are lines that have no start or end tag that belong to them.
	// Starts, but continues on the next line.
	if !start && !end && !l.firstLineFound && !nextLineStart {
		fmt.Println(" ***DESC", l.currentLineNR, ": !start && !end && !l.firstLineFound && !nextLineStart, CONTINUES ON NEXT LINE ***")
		l.workingLine = l.workingLine + " " + l.currentLine
		return l.lexReadFileLine
	}

	//Description line. These are lines that have no start or end tag that belong to them.
	// End's here.
	//TODO: This one should not need to lexLineContent since there is just descriptions here,
	// no need to lex inside, and it should be given a token immediately.
	if !start && !end && !l.firstLineFound && nextLineStart {
		fmt.Println(" ***DESC", l.currentLineNR, " : !start && !end && !l.firstLineFound && nextLineStart ***")
		l.workingLine = l.workingLine + " " + l.currentLine
		return l.lexLineContent
	}

	// ---------------------The code should never exexute the code below-----------------------
	// The print and return below should optimally never happen.
	// Check it's output to figure what is not detected in the if's above.
	fmt.Println("DEBUG: *uncaught line!! :", l.currentLineNR, l.workingLine)
	return l.lexPrint
}

//------------------Tokens--------------------------------------------
/*
Example for how to send and receive tokens from the lexer.
When the lexer find something, it creates a token of type token. It will choose the type of token found and put that into tokenType, and the text found will be put into the argument.
Then the token is put on the channel to be received by the parser, and the go struct code will be generated within the switch/case selection.
*/

//tokenType is the type describing a token.
// A <start> start tag will have token start.
// An </start> end tag will have token end.
type tokenType int

const (
	tokenStartTag tokenType = iota
	tokenEndTag
	tokenArgumentName
	tokenArgumentValue
)

type token struct {
	tokenType        //type of token, tokenStart, tokenStop, etc.
	tokenText string //the actual text found in the xml while lexing
}

//readToken will pickup and read all the tokens that is received from the lexer
func readToken() {
	for v := range tokenChan {
		fmt.Println("*readToken*", v)
	}
	wg.Done()
}

//--------------------------------------------------------------------

var tokenChan chan token
var wg sync.WaitGroup

func main() {
	tokenChan = make(chan token)

	fh, err := os.Open(fileName)
	if err != nil {
		log.Println("Error: opening file: ", err)
	}

	wg.Add(1)
	go readToken()
	tokenChan <- token{tokenType: tokenStartTag, tokenText: "hest"}

	lex := newLexer(fh)
	lex.lexStart()

	wg.Wait()
}

/*
 This main.go file will typically be the parser, but here we just print out the tokens that is received.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/postmannen/lexml"
)

func main() {

	a := os.Args
	if len(a) < 2 {
		log.Fatal("Specify an xml file\n")

	}

	fileName := flag.String("fileName", "", "specify the filename to check")

	flag.Parse()

	fh, err := os.Open(*fileName)
	if err != nil {
		log.Fatal("Error: opening file: ", err)
	}

	tCh := lexml.LexStart(fh)

	for v := range tCh {
		fmt.Println("*readToken from channel * ", v.TokenType, ", tokenText = ", v.TokenText)
	}

}

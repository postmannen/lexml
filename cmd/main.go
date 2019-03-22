package main

import (
	"flag"
	"log"
	"os"

	"github.com/postmannen/xmltogo"
)

func main() {

	a := os.Args
	if len(a) < 2 {
		log.Fatal("Specify an xml file\n")

	}

	fileName := flag.String("fileName", "", "specify the filename to check")
	tokenOutput := flag.String("tokenOutput", "console", "specify '0' for console or '1' for channel. If you want to simulate a read locally without a parser who picks up the data from the channel, remember to enable -readChannel=yes.")

	flag.Parse()

	fh, err := os.Open(*fileName)
	if err != nil {
		log.Fatal("Error: opening file: ", err)
	}

	xmltogo.LexStart(fh, *tokenOutput)

}

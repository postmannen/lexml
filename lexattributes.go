package main

import (
	"log"
	"strings"
)

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

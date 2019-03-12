package main

import (
	"log"
)

//lexStartStopTag will check various versions of tags. Like if is a start tag, stop tag
// tag with elements inside.
func (l *lexer) lexStartStopTag() (elementText string) {
	//If no equal was detected in the line, it is most likely a line with a start,
	// and an end tag, but just text inbetween, and we want to pick out that text.
	// Example : <someTag>WE WANT THIS TEXT</someTag>
	if !l.foundEqual {
		pos := 0
		var posTextStart int
		var posTextStop int
		var text []byte

		for {
			if pos >= len(l.workingLine) {
				break
			}
			switch {
			case l.workingLine[pos] == '<' && pos == 0:
				//fmt.Println("-- firstStartAngle found", firstStartAngleFound)
			case l.workingLine[pos] == '>' && pos == len(l.workingLine)-1:
				//fmt.Println("-- secondStopAngle found", secondStopAngleFound)
			case l.workingLine[pos] == '>':
				posTextStart = pos + 1
				//fmt.Println("-- firstStopAngle found", firstStopAngleFound)
			case l.workingLine[pos] == '<':
				posTextStop = pos - 1
				//fmt.Println("-- secondStartAngle found", secondStartAngleFound)
			default:
				//if there are more angle brackets than needed for a start and end tag
				// there is something malformed in xml, and we break out
				if l.workingLine[pos] == '<' || l.workingLine[pos] == '>' {
					log.Println("error: malformed xml with to man angle brackets")
				}
			}

			pos++
		}

		if posTextStart != 0 || posTextStop != 0 {
			for i := posTextStart; i <= posTextStop; i++ {
				text = append(text, l.workingLine[i])
			}
			return string(text)
		}

	}
	return ""
}

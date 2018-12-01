# parrotxml2

Lex and parse the parrot ardrone3.xml file.

## How it works

First we create a lexer struct to hold all the info lexed. We then create methods on that struct for all the various lex actions.

We have one function that starts it all called lexStart(), and all it does is that it calls a method, and put whatever method that method returns into a variable. We then call that returned method on the next loop, and so on.

The for loop in lexStart() will exit if it at any point receives NIL from any of the methods executed. NIL is for example received when the last line of the XML file is read.
The program will then return to main, and terminate.

## The flow of the functions

A flowchart diagram showing the flow of the program.
![alt_text](flow.jpg)
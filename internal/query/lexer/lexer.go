package lexer

import (
	"bytes"

	"github.com/madeinheaven91/black-turtle-go/internal/query/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	literal := l.readWord()

	if token.Lookup(literal) == token.NAME {
		for {
			l.skipWhitespace()
			nextWord := l.peekWord()
			if token.Lookup(nextWord) != token.NAME {
				break
			}
			literal = literal + " " + l.readWord()
		}
		tok.Type = token.NAME
		tok.Literal = literal
	}
	tok.Type = token.Lookup(literal)
	tok.Literal = literal

	return tok
}

func (l *Lexer) readWord() string {
	var out bytes.Buffer
	for l.peekChar() != ' ' && l.peekChar() != 0 {
		cur := l.ch
		out.WriteByte(cur)
		l.readChar()
	}
	out.WriteByte(l.ch)
	l.readChar()
	return out.String()
}

func (l *Lexer) peekWord() string {
	oldPos, oldReadPos, oldCh := l.position, l.readPosition, l.ch
	out := l.readWord()
	l.ch = oldCh
	l.position = oldPos
	l.readPosition = oldReadPos
	return out
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

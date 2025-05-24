package parser

import (
	"fmt"
	"time"

	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/internal/query/ir"
	"github.com/madeinheaven91/black-turtle-go/internal/query/lexer"
	"github.com/madeinheaven91/black-turtle-go/internal/query/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

// Hides lexing part
func FromString(input string) *Parser {
	l := lexer.New(input)
	p := New(l)
	return p
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: make([]string, 0)}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseQuery() *ir.QueryRaw {
	var res ir.QueryRaw
	switch p.curToken.Type {
	case token.LESSONS:
		res = p.parseLessonQuery()
	default:
		logging.Warning("couldn't parse query starting with token %s of type %s", p.curToken.Literal, p.curToken.Type)
		return nil
	}
	logging.Trace(res.String())
	return &res
}

func (p *Parser) parseLessonQuery() *ir.LessonsQueryRaw {
	query := ir.LessonsQueryRaw{}
	for !p.curTokenIs(token.EOF) {
		switch p.curToken.Type {
		case token.LESSONS:
		case token.NAME:
			if query.StudyEntityName == nil {
				name := p.curToken.Literal
				query.StudyEntityName = &name
			} else {
				p.errors = append(p.errors, "more than 1 study entity name provided")
			}
		case token.DAY:
			if query.TimeFrame.Date == nil {
				day := p.curToken.Literal
				query.TimeFrame.Day = &day
			} else {
				p.errors = append(p.errors, "provided day token when date token is present")
			}
		case token.WEEK:
			if query.TimeFrame.Date == nil {
				week := p.curToken.Literal
				query.TimeFrame.Week = &week
			} else {
				p.errors = append(p.errors, "provided week token when date token is present")
			}
		case token.DATE:
			if query.TimeFrame.Day == nil && query.TimeFrame.Week == nil {
				date, err := parseDate(p.curToken.Literal)
				if err != nil {
					p.errors = append(p.errors, "couldn't parse date token to time.Time")
				} else {
					query.TimeFrame.Date = date
				}
			} else {
				p.errors = append(p.errors, "provided date token when day or week token is present")
			}
		default:
			p.errors = append(p.errors, fmt.Sprintf("unexpected token %v", p.curToken))
		}
		p.nextToken()
	}
	if len(p.errors) != 0 {
		return nil
	}
	return &query
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) Errors() []string {
	return p.errors
}

func parseDate(input string) (*time.Time, error) {
	// форматы:
	// 02.01.2006
	// 02.01.06
	// 02.01
	// пн
	// завтра

	var res time.Time
	// TODO: maybe could be written prettier
	res, err := time.Parse("02.01", input)
	if err == nil {
		res = res.AddDate(time.Now().Year(), 0, 0)
		return &res, nil
	}
	res, err = time.Parse("02.01.06", input)
	if err == nil {
		return &res, nil
	}
	res, err = time.Parse("02.01.2006", input)
	if err == nil {
		return &res, nil
	}

	e := fmt.Errorf("couldn't parse date")
	return nil, e
}

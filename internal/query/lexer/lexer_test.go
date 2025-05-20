package lexer

import (
	"encoding/json"
	"testing"

	"github.com/madeinheaven91/black-turtle-go/internal/query/token"
)

func TestLesson(t *testing.T) {
	tests := []struct {
		input  string
		tokens []token.Token
	}{
		{"пары", []token.Token{
			token.New("пары", token.LESSONS),
		}},
		{"пары 921", []token.Token{
			token.New("пары", token.LESSONS),
			token.New("921", token.NAME),
		}},
		{"пары димитриев александр олегович", []token.Token{
			token.New("пары", token.LESSONS),
			token.New("димитриев александр олегович", token.NAME),
		}},
		{"пары завтра", []token.Token{
			token.New("пары", token.LESSONS),
			token.New("завтра", token.DAY),
		}},
		{"пары 02.02.2024", []token.Token{
			token.New("пары", token.LESSONS),
			token.New("02.02.2024", token.DATE),
		}},
		{"пары 02.02.24", []token.Token{
			token.New("пары", token.LESSONS),
			token.New("02.02.24", token.DATE),
		}},
		{"пары 02.02", []token.Token{
			token.New("пары", token.LESSONS),
			token.New("02.02", token.DATE),
		}},
		{"пары 921 завтра", []token.Token{
			token.New("пары", token.LESSONS),
			token.New("921", token.NAME),
			token.New("завтра", token.DAY),
		}},
		{"пары 921 пн след", []token.Token{
			token.New("пары", token.LESSONS),
			token.New("921", token.NAME),
			token.New("пн", token.DAY),
			token.New("след", token.WEEK),
		}},
		{"пары александр пн", []token.Token{
			token.New("пары", token.LESSONS),
			token.New("александр", token.NAME),
			token.New("пн", token.DAY),
		}},
		{"пары александр олегович пн", []token.Token{
			token.New("пары", token.LESSONS),
			token.New("александр олегович", token.NAME),
			token.New("пн", token.DAY),
		}},
		{"пары александр олегович пн след", []token.Token{
			token.New("пары", token.LESSONS),
			token.New("александр олегович", token.NAME),
			token.New("пн", token.DAY),
			token.New("след", token.WEEK),
		}},
		{"пары димитриев александр олегович пн след", []token.Token{
			token.New("пары", token.LESSONS),
			token.New("димитриев александр олегович", token.NAME),
			token.New("пн", token.DAY),
			token.New("след", token.WEEK),
		}},
	}

	for testIndex, testcase := range tests {
		l := New(testcase.input)
		tokens := make([]token.Token, 0)
		for l.ch != 0 {
			tok := l.NextToken()
			tokens = append(tokens, tok)
		}

		if len(tokens) != len(testcase.tokens) {
			tmp, _ := json.MarshalIndent(tokens, "", "  ")
			tokensPretty := string(tmp)
			tmp, _ = json.MarshalIndent(testcase.tokens, "", "  ")
			testcasePretty := string(tmp)
			t.Logf("test %d: len(tokens) (%d) doesnt match len(testcase.tokens) (%d)", testIndex, len(tokens), len(testcase.tokens))
			t.Log("tokens: ", tokensPretty)
			t.Fatal("testcase.tokens: ", testcasePretty)
		}

		for tokenIndex, token := range testcase.tokens {
			if token != tokens[tokenIndex] {
				t.Fatalf("test %d token %d: expected %#v, got=%#v", testIndex, tokenIndex, token, tokens[tokenIndex])
			}
		}
	}
}

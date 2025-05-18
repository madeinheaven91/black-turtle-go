package parser

import (
	"testing"

	"github.com/madeinheaven91/black-turtle-go/internal/parser/ir"
	"github.com/madeinheaven91/black-turtle-go/internal/parser/lexer"
)

func TestLessonQuery(t *testing.T) {
	testcases := []struct {
		input    string
		expected string
	}{
		{
			"пары",
			"пары nil nil nil nil",
		},
		{
			"пары 921",
			"пары 921 nil nil nil",
		},
		{
			"пары 921 завтра",
			"пары 921 завтра nil nil",
		},
		{
			"пары 921 пн след",
			"пары 921 пн след nil",
		},
		{
			"пары 921 02.02",
			"пары 921 nil nil 02.02.2025",
		},
		{
			"пары 921 пн 02.02",
			"",
		},
		{
			"пары 921 след 02.02",
			"",
		},
		{
			"пары 921 пн след 02.02",
			"",
		},
	}

	for testIndex, testcase := range testcases {
		l := lexer.New(testcase.input)
		p := New(l)
		q := p.ParseQuery()
		if testcase.expected != "" {
			checkParserErrors(t, p)
		} else if len(p.Errors()) == 0 {
			t.Fatalf("test %d: should have failed", testIndex)
		} else {
			continue
		}

		lq, ok := q.Command.(*ir.LessonsQuery)
		if ok {
			if lq.String() != testcase.expected {
				t.Fatalf("test %d: expected %s, got %s", testIndex, testcase.expected, lq.String())
			}
		} else {
			t.Fatalf("test %d: parsed query not of type *ir.LessonsQuery. got=%T", testIndex, q)
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser had %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

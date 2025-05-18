package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/madeinheaven91/black-turtle-go/internal/logging"
	"github.com/madeinheaven91/black-turtle-go/internal/parser"
	"github.com/madeinheaven91/black-turtle-go/internal/parser/ir"
	"github.com/madeinheaven91/black-turtle-go/internal/parser/lexer"
	"github.com/madeinheaven91/black-turtle-go/internal/parser/token"
)

func main() {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		logging.Critical("%s\n", err)
		os.Exit(1)
	}
	for k, v := range envFile {
		os.Setenv(k, v)
	}
	logging.InitLogLevel()

	logging.Info("REPL for testing parser")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf(">>> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			logging.Debug("%#v", tok)
		}

		l = lexer.New(line)
		p := parser.New(l)
		query := p.ParseQuery()
		if query == nil {
			logging.Error("unknown command")
			return
		}
		lq, ok := query.Command.(*ir.LessonsQuery)
		if ok {
			logging.Debug("%s", lq)
		}
		if len(p.Errors()) != 0 {
			for _, e := range p.Errors() {
				logging.Error("%s", e)
			}
		}
	}
}

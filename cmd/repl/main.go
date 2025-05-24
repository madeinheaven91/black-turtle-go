package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/madeinheaven91/black-turtle-go/internal/query/ir"
	"github.com/madeinheaven91/black-turtle-go/internal/query/lexer"
	"github.com/madeinheaven91/black-turtle-go/internal/query/parser"
	"github.com/madeinheaven91/black-turtle-go/internal/query/token"
	"github.com/madeinheaven91/black-turtle-go/pkg/config"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
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
	config.InitFromEnv()
	logging.InitLoggers()

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
		raw := p.ParseQuery()
		if raw == nil {
			logging.Error("unknown command")
			return
		}
		lq, ok := (*raw).(*ir.LessonsQueryRaw)
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

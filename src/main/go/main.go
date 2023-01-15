package main

import (
	"bufio"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"gradle-go-generated/parser"
	"os"
)

func main() {
	calc := CreateCalcListener()

	cin := bufio.NewReader(os.Stdin)

	for {
		s, _ := cin.ReadString('\n')
		input := antlr.NewInputStream(s)
		lexer := parser.NewCalcLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parser.NewCalcParser(stream)
		p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
		p.BuildParseTrees = true
		tree := p.Stat()
		calc.instSet.Clear()
		antlr.ParseTreeWalkerDefault.Walk(calc, tree)
		calc.instSet.Evaluate()
	}
}

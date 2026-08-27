package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"simenf05/flc/ast"
	"simenf05/flc/lexer"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: flc [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 2 {
		fmt.Println("Input file is missing")
		os.Exit(1)
	}

	f, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(f)

	tokenizer := lexer.NewTokenizer(*reader)
	tokens, err := tokenizer.Parse()
	if err != nil {
		panic(err)
	}

	buffer := ast.NewBuffer(tokens)

	expr := buffer.ParseExpr()
	fmt.Printf("ret expr %v\n", expr)

	for _, token := range tokens {
		fmt.Printf("type: %d numval: %d strval: %s\n",
			token.Type,
			token.NumValue,
			token.StrValue)
	}
}

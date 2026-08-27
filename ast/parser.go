package ast

import (
	"simenf05/flc/lexer"
)

type Buffer struct {
	vector []lexer.Token
	index  int
}

func (b *Buffer) nextToken() lexer.Token {
	tok := b.vector[b.index]
	b.index++
	return tok
}

func (b *Buffer) peekToken() lexer.Token {
	return b.vector[b.index]
}

func NewBuffer(tokens []lexer.Token) Buffer {
	return Buffer{
		vector: tokens,
	}
}

func (b *Buffer) ParseExpr() Expr {

	t := b.nextToken()

	switch t.Type {
	case lexer.Identifier:
		return b.ParseIdentifierExpr(t)
	case lexer.Number:
		return b.ParseNumberExpr(t)
	case lexer.EOF:
		return nil
	case '(':
		return b.ParseParenExpr(t)
	}

	return nil
}

func (b *Buffer) ParseNumberExpr(t lexer.Token) Expr {
	return NumberExpr{
		Value: t.NumValue,
	}
}

func (b *Buffer) ParseParenExpr(t lexer.Token) Expr {
	if t.Type != '(' {
		return nil
	}

	expr := b.ParseExpr()
	tok := b.nextToken()

	if tok.Type != ')' {
		return nil
	}
	return expr
}

func (b *Buffer) ParseIdentifierExpr(t lexer.Token) Expr {

	tok := b.nextToken()

	if tok.Type != '(' {
		return VariableExpr{
			Name: t.StrValue,
		}
	}

	var args []Expr

	for {
		tok = b.nextToken()
		if tok.Type != ')' {
			break
		}
		arg := b.ParseExpr()
		if arg != nil {
			args = append(args, arg)
		}
		tok = b.nextToken()
		if tok.Type != ',' && tok.Type != ')' {
			return nil
		}
	}

	return CallExpr{
		Callee: t.StrValue,
		Args:   args,
	}
}

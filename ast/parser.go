package ast

import (
	"simenf05/flc/lexer"
)

type Buffer struct {
	vector []lexer.Token
	index  int
}

func (b *Buffer) eat() {
	b.index++
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
	case lexer.Func:
		return b.ParseFunction(t)
	case lexer.TokenType(Addition),
		lexer.TokenType(Subtraction),
		lexer.TokenType(Multiplication),
		lexer.TokenType(Division):
		return b.ParseBinaryExpr(t)
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
	tok := b.peekToken()
	if tok.Type != '(' {
		return VariableExpr{
			Name: t.StrValue,
		}
	}
	b.eat()

	var args []Expr
	for {
		tok = b.peekToken()

		switch tok.Type {
		case ',':
			b.eat()
			continue
		case ')':
			b.eat()
			return CallExpr{
				Callee: t.StrValue,
				Args:   args,
			}
		default:
			arg := b.ParseExpr()
			if arg == nil {
				return nil
			}
			args = append(args, arg)
		}
	}
}

func (b *Buffer) ParseBinaryExpr(t lexer.Token) Expr {
	left := b.ParseExpr()
	right := b.ParseExpr()

	if left == nil || right == nil {
		return nil
	}

	return BinaryExpr{
		Op:    OperatorType(t.Type),
		Left:  left,
		Right: right,
	}
}

func (b *Buffer) ParseFunction(t lexer.Token) Expr {

	if t.Type != lexer.Func {
		return nil
	}

	proto := b.ParsePrototype(b.nextToken())
	if proto == nil {
		return nil
	}

	if tok := b.nextToken(); tok.Type != '{' {
		return nil
	}

	body := b.ParseExpr()
	if body == nil {
		return nil
	}

	if tok := b.nextToken(); tok.Type != '}' {
		return nil
	}

	return FunctionExpr{
		Proto: proto.(PrototypeExpr),
		Body:  body,
	}
}

func (b *Buffer) ParsePrototype(t lexer.Token) Expr {
	if t.Type != lexer.Identifier {
		return nil
	}

	tok := b.nextToken()
	if tok.Type != '(' {
		return nil
	}

	var args []string
	for {
		tok = b.peekToken()

		switch tok.Type {
		case ',':
			b.eat()
			continue
		case ')':
			b.eat()
			return PrototypeExpr{
				Name: t.StrValue,
				Args: args,
			}
		default:
			args = append(args, b.nextToken().StrValue)
		}
	}
}

package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	EOF TokenType = iota * -1
	Identifier
	Number
	Func
)

type Token struct {
	Type     TokenType
	NumValue int
	StrValue string
}

type Tokenizer struct {
	reader bufio.Reader
	tokens []Token
}

func (t *Tokenizer) readRunes(lastChar rune, cond func(lastChar rune) bool) (string, error) {

	var str strings.Builder
	var err error

	for cond(lastChar) {
		str.WriteRune(lastChar)
		lastChar, _, err = t.reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("reading rune: %w", err)
		}
	}
	err = t.reader.UnreadRune()
	if err != nil {
		return "", fmt.Errorf("unreading rune: %w", err)
	}

	return str.String(), nil
}

func (t *Tokenizer) getToken() (*Token, error) {

	var err error
	lastChar := ' '

	for unicode.IsSpace(lastChar) {
		lastChar, _, err = t.reader.ReadRune()
		if err == io.EOF {
			return &Token{
				Type: EOF,
			}, nil
		}
		if err != nil {
			return nil, fmt.Errorf("reading rune: %w", err)
		}
	}

	if unicode.IsLetter(lastChar) {
		word, err := t.readRunes(lastChar, func(char rune) bool {
			return unicode.IsLetter(char) || unicode.IsDigit(char)
		})
		if err != nil {
			return nil, fmt.Errorf("identifier reading runes: %w", err)
		}

		switch word {
		case "func":
			return &Token{
				Type: Func,
			}, nil
		}

		return &Token{
			Type:     Identifier,
			StrValue: word,
		}, nil
	}

	if unicode.IsDigit(lastChar) {
		num, err := t.readRunes(lastChar, func(char rune) bool {
			return unicode.IsDigit(char)
		})
		if err != nil {
			return nil, fmt.Errorf("numeric reading runes: %w", err)
		}
		numVal, err := strconv.Atoi(num)
		if err != nil {
			return nil, fmt.Errorf("failed to convert number: %w", err)
		}

		return &Token{
			Type:     Number,
			NumValue: numVal,
		}, nil
	}

	return &Token{
		Type: TokenType(lastChar),
	}, nil
}

func (t *Tokenizer) Parse() ([]Token, error) {
	for {
		token, err := t.getToken()
		if err != nil {
			return nil, fmt.Errorf("getting tokens: %w", err)
		}
		t.tokens = append(t.tokens, *token)

		if token.Type == EOF {
			break
		}
	}

	return t.tokens, nil
}

func NewTokenizer(reader bufio.Reader) *Tokenizer {
	return &Tokenizer{
		reader: reader,
		tokens: []Token{},
	}
}

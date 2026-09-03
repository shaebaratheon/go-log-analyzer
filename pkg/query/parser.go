package query

import (
	"errors"
	"strconv"
	"strings"
)

type TokenKind int

const (
	TokenSelect TokenKind = iota
	TokenFrom
	TokenWhere
	TokenAnd
	TokenOr
	TokenIdent
	TokenString
	TokenNumber
	TokenOperator
	TokenEOF
)

type Token struct {
	Kind  TokenKind
	Value string
}

type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	if l.pos >= len(l.input) {
		return Token{Kind: TokenEOF, Value: ""}
	}

	ch := l.input[l.pos]
	if ch == ''' || ch == '"' {
		return l.readString(ch)
	}
	if isDigit(ch) {
		return l.readNumber()
	}
	if isLetter(ch) {
		ident := l.readIdent()
		switch strings.ToUpper(ident) {
		case "SELECT":
			return Token{Kind: TokenSelect, Value: ident}
		case "FROM":
			return Token{Kind: TokenFrom, Value: ident}
		case "WHERE":
			return Token{Kind: TokenWhere, Value: ident}
		case "AND":
			return Token{Kind: TokenAnd, Value: ident}
		case "OR":
			return Token{Kind: TokenOr, Value: ident}
		case "CONTAINS":
			return Token{Kind: TokenOperator, Value: "CONTAINS"}
		default:
			return Token{Kind: TokenIdent, Value: ident}
		}
	}
	if ch == '=' || ch == '>' || ch == '<' || ch == '!' {
		op := string(ch)
		l.pos++
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			op += "="
			l.pos++
		}
		return Token{Kind: TokenOperator, Value: op}
	}

	l.pos++
	return Token{Kind: TokenIdent, Value: string(ch)}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && (l.input[l.pos] == ' ' || l.input[l.pos] == '	' || l.input[l.pos] == '
') {
		l.pos++
	}
}

func (l *Lexer) readString(quote byte) Token {
	l.pos++ // consume open
	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != quote {
		l.pos++
	}
	val := l.input[start:l.pos]
	if l.pos < len(l.input) {
		l.pos++ // consume close
	}
	return Token{Kind: TokenString, Value: val}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	for l.pos < len(l.input) && (isDigit(l.input[l.pos]) || l.input[l.pos] == '.') {
		l.pos++
	}
	return Token{Kind: TokenNumber, Value: l.input[start:l.pos]}
}

func (l *Lexer) readIdent() Token {
	start := l.pos
	for l.pos < len(l.input) && (isLetter(l.input[l.pos]) || isDigit(l.input[l.pos]) || l.input[l.pos] == '_') {
		l.pos++
	}
	return Token{Kind: TokenIdent, Value: l.input[start:l.pos]}
}

func isDigit(ch byte) bool  { return ch >= '0' && ch <= '9' }
func isLetter(ch byte) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

type SQLQuery struct {
	Fields []string
	Table  string
	Where  *QueryExpression
}

type Parser struct {
	lexer *Lexer
	curr  Token
}

func NewParser(query string) *Parser {
	l := NewLexer(query)
	return &Parser{lexer: l, curr: l.NextToken()}
}

func (p *Parser) Parse() (*SQLQuery, error) {
	if p.curr.Kind != TokenSelect {
		return nil, errors.New("expected SELECT")
	}
	p.advance()

	var fields []string
	for p.curr.Kind == TokenIdent || p.curr.Kind == TokenString {
		fields = append(fields, p.curr.Value)
		p.advance()
	}

	if p.curr.Kind != TokenFrom {
		return nil, errors.New("expected FROM")
	}
	p.advance()

	table := p.curr.Value
	p.advance()

	var expr *QueryExpression
	if p.curr.Kind == TokenWhere {
		p.advance()
		expr = &QueryExpression{}
		for p.curr.Kind != TokenEOF {
			field := p.curr.Value
			p.advance()
			op := Operator(p.curr.Value)
			p.advance()
			val := p.curr.Value
			p.advance()

			expr.Conditions = append(expr.Conditions, Condition{
				Field: field,
				Op:    op,
				Value: val,
			})
			if p.curr.Kind == TokenAnd {
				p.advance()
			} else {
				break
			}
		}
	}

	return &SQLQuery{Fields: fields, Table: table, Where: expr}, nil
}

func (p *Parser) advance() {
	p.curr = p.lexer.NextToken()
}

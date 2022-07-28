/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package eval

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

// ---- lexer ----

// This lexer is similar to the one described in Chapter 13.
type lexer struct {
	scan  scanner.Scanner
	token rune // current lookahead token
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

// describe returns a string describing the current token, for use in errors.
func (lex *lexer) describe() string {
	switch lex.token {
	case scanner.EOF:
		return "end of file"
	case scanner.Ident:
		return fmt.Sprintf("identifier %s", lex.text())
	case scanner.Int, scanner.Float:
		return fmt.Sprintf("number %s", lex.text())
	}
	return fmt.Sprintf("%q", rune(lex.token)) // any other rune
}

func precedence(op rune) int {
	// 运算符优先级，越高对应树型的越上层
	switch op {
	case '*', '/', '%':
		return 4
	case '+', '-', '>', '<', '≥', '≤', '≡', '≠':
		return 3
	case '∈', '∉':
		return 2
	case '∩', '∪':
		return 1
	}
	return 0
}

// ---- parser ----

// Parse parses the input string as an arithmetic expression.
//
//   expr = num                         a literal number, e.g., 3.14159
//        | id                          a variable name, e.g., x
//        | id '(' expr ',' ... ')'     a function call
//        | '-' expr                    a unary operator (+-)
//        | expr '+' expr               a binary operator (+-*/)
//
func Parse(input string) (_ Expr, err error) {
	lex := new(lexer)
	lex.scan.Init(strings.NewReader(input))
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats
	lex.next() // initial lookahead
	if lex.token == scanner.EOF {
		return nil, fmt.Errorf("unexpected %s", lex.describe())
	}
	return parseExpr(lex)
}

func parseExpr(lex *lexer) (Expr, error) { return parseBinary(lex, 1) }

// binary = unary ('+' binary)*
// parseBinary stops when it encounters an
// operator of lower precedence than prec1.
func parseBinary(lex *lexer, prec1 int) (Expr, error) {
	lhs, err := parseUnary(lex)
	if err != nil {
		return nil, err
	}
	for prec := precedence(lex.token); prec >= prec1; prec-- {
		for precedence(lex.token) == prec {
			op := lex.token
			lex.next() // consume operator
			rhs, err := parseBinary(lex, prec+1)
			if err != nil {
				return nil, err
			}
			if op == '∈' || op == '∉' {
				lhs = section{op, lhs, rhs.(section).y}
			} else {
				lhs = binary{op, lhs, rhs}
			}
		}
	}
	return lhs, nil
}

// unary = '+' expr | primary
func parseUnary(lex *lexer) (Expr, error) {
	if lex.token == '+' || lex.token == '-' {
		op := lex.token
		lex.next() // consume '+' or '-'
		expr, err := parseUnary(lex)
		if err != nil {
			return nil, err
		}
		return unary{op, expr}, nil
	}
	return parsePrimary(lex)
}

// primary = id
//         | id '(' expr ',' ... ',' expr ')'
//         | num
//         | '(' expr ')'
func parsePrimary(lex *lexer) (Expr, error) {
	switch lex.token {
	case scanner.Ident:
		id := lex.text()
		lex.next() // consume Ident
		if lex.token != '(' {
			return Var(id), nil
		}
		lex.next() // consume '('
		var args []Expr
		if lex.token != ')' {
			for {
				expr, err := parseExpr(lex)
				if err != nil {
					return nil, err
				}
				args = append(args, expr)
				if lex.token != ',' {
					break
				}
				lex.next() // consume ','
			}
			if lex.token != ')' {
				return nil, fmt.Errorf("got %s, want ')'", lex.describe())
			}
		}
		lex.next() // consume ')'
		return call{id, args}, nil

	case scanner.Int, scanner.Float:
		f, err := strconv.ParseFloat(lex.text(), 64)
		if err != nil {
			return nil, err
		}
		lex.next() // consume number
		return literal(f), nil

	case '(':
		lex.next() // consume '('
		e, err := parseExpr(lex)
		if err != nil {
			return nil, err
		}
		if lex.token != ')' {
			return nil, fmt.Errorf("got %s, want ')'", lex.describe())
		}
		lex.next() // consume ')'
		return e, nil
	case '{':
		lex.next() // consume '('
		var args []Expr
		if lex.token != '}' {
			for {
				expr, err := parseExpr(lex)
				if err != nil {
					return nil, err
				}
				args = append(args, expr)
				if lex.token != ',' {
					break
				}
				lex.next() // consume ','
			}
			if lex.token != '}' {
				return nil, fmt.Errorf("got %s, want '}'", lex.describe())
			}
		}
		lex.next() // consume ')'
		return section{y: args}, nil
	case '\'':
		id := lex.text()
		lex.next() // consume Ident
		if lex.token != '\'' {
			s := lex.text()
			lex.next()
			if lex.token == '\'' {
				lex.next()
			}
			return Var(id + s + id), nil
		}
		return Var(""), nil
	}
	return nil, fmt.Errorf("unexpected %s", lex.describe())
}

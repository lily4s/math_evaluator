package math_evaluator

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"unicode"
)

type stack struct {
	items []interface{}
}

type tokenType int

const (
	Number tokenType = iota
	Operator
	LeftParen
	RightParen
)

type token struct {
	Kind  tokenType
	Value any
}

type operator_detail struct {
	symbol rune
	prec   int
	assoc  rune
}

var operators = map[rune]operator_detail{
	'+': {'+', 2, 'L'},
	'-': {'-', 2, 'L'},
	'*': {'*', 3, 'L'},
	'/': {'/', 3, 'L'},
	'^': {'^', 4, 'R'},
}

func (this *stack) push(item interface{}) {
	this.items = append(this.items, item)
}

func (this *stack) pop() (interface{}, error) {
	if len(this.items) == 0 {
		return nil, errors.New("stack is empty")
	}

	n := this.items[len(this.items)-1]
	this.items = this.items[:len(this.items)-1]
	return n, nil
}

func (this *stack) peek() (interface{}, error) {
	if len(this.items) == 0 {
		return nil, errors.New("stack is empty")
	}

	n := this.items[len(this.items)-1]
	return n, nil
}

func (this *stack) isEmpty() bool {
	return len(this.items) == 0
}

func newStack() *stack {
	return &stack{make([]interface{}, 0)}
}

func isNumber(n string) bool {
	_, err := strconv.ParseFloat(n, 64)
	return err == nil
}

func isOperator(n rune) bool {
	_, ok := operators[n]
	return ok
}

func tokenize(input string) ([]token, error) {
	var tokens []token
	i := 0

	for i < len(input) {
		switch {
		case unicode.IsDigit(rune(input[i])) || input[i] == '.':
			start := i
			for i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.') {
				i++
			}

			number := input[start:i]
			if isNumber(number) {
				n, err := strconv.ParseFloat(number, 64)
				if err != nil {
					return nil, err
				}
				tokens = append(tokens, token{Number, n})
			} else {
				return nil, errors.New("invalid number")
			}
		case input[i] == '(':
			tokens = append(tokens, token{LeftParen, '('})
			i++
		case input[i] == ')':
			tokens = append(tokens, token{RightParen, ')'})
			i++
		case isOperator(rune(input[i])):
			tokens = append(tokens, token{Operator, rune(input[i])})
			i++
		case unicode.IsSpace(rune(input[i])):
			i++
		default:
			return nil, fmt.Errorf("unexpected character %v at column %d", input[i], i)
		}
	}

	return tokens, nil
}

func precedence(op rune) int {
	if n, ok := operators[op]; ok {
		return n.prec
	} else {
		return -1
	}
}

func isLeftAssoc(op rune) bool {
	if n, ok := operators[op]; ok {
		return n.assoc == 'L'
	}
	return false
}

func shuntingYard(tokens []token) ([]interface{}, error) {
	queue := make([]interface{}, 0)
	stck := newStack()

	for _, t := range tokens {
		switch t.Kind {
		case Number:
			queue = append(queue, t.Value)
		case Operator:
			for len(stck.items) != 0 {
				top, err := stck.peek()
				if err != nil {
					return nil, err
				}

				if top != '(' && (precedence(top.(rune)) >= precedence(t.Value.(rune))) &&
					isLeftAssoc(t.Value.(rune)) {
					n, err := stck.pop()
					if err != nil {
						return nil, err
					}
					queue = append(queue, n)
				} else {
					break
				}
			}
			stck.push(t.Value)
		case LeftParen:
			stck.push(t.Value)
		case RightParen:
			n, err := stck.pop()
			if err != nil {
				return nil, err
			}

			for n != '(' {
				queue = append(queue, n)
				n, err = stck.pop()
				if err != nil {
					return nil, errors.New("mismatched parentheses")
				}

			}
		}
	}

	for len(stck.items) != 0 {
		n, err := stck.pop()
		if err != nil {
			return nil, err
		}

		if n == '(' || n == ')' {
			return nil, fmt.Errorf("mismatched parentheses")
		}

		queue = append(queue, n)
	}

	return queue, nil
}

func evaluatePostfix(tokens []interface{}) (float64, error) {
	output := newStack()

	for _, t := range tokens {
		switch t.(type) {
		case float64:
			output.push(t)
		case rune:

			o1, err := output.pop()
			if err != nil {
				return -1, err
			}
			b, _ := o1.(float64)

			o2, err := output.pop()
			if err != nil {
				return -1, err
			}
			a, _ := o2.(float64)

			switch t {
			case '+':
				output.push(a + b)
			case '-':
				output.push(a - b)
			case '*':
				output.push(a * b)
			case '/':
				if b == 0 {
					return -1, fmt.Errorf("division by zero")
				}
				output.push(a / b)
			case '^':
				output.push(math.Pow(a, b))
			}
		}
	}

	res, err := output.pop()
	if err != nil {
		return -1, fmt.Errorf("invalid opre")
	}

	v, _ := res.(float64)

	return v, nil
}

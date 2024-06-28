package math_evaluator

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	got, err := tokenize("12+167/(1-5)*999^4")
	if err != nil {
		t.Error(err)
	}

	t.Log(got)
}

func TestShuntingYard(t *testing.T) {
	n := []token{
		{0, 12},
		{1, '+'},
		{0, 167},
		{1, '*'},
		{0, 22222},
		{1, '/'},
		{2, '('},
		{0, 1},
		{1, '-'},
		{0, 5},
		{3, ')'},
		{1, '^'},
		{0, 999},
	}

	got, err := shuntingYard(n)
	if err != nil {
		t.Error(err)
	}

	t.Log(got)
}

func TestEvaluatePostfix(t *testing.T) {
	got, err := evaluatePostfix(
		[]any{12.0, 167.0, 22222.0, '*', 1.0, 5.0, '-', 999.0, '^', '/', '+'},
	)
	if err != nil {
		t.Error(err)
	}

	t.Log(got)
}

func TestOverall(t *testing.T) {
	tkn, err := tokenize("12+45-60*(3^6)")
	if err != nil {
		t.Errorf("error tokenize, %v", err)
	}
	t.Logf("the result of tokenize:\n%v", tkn)

	postfix, err := shuntingYard(tkn)
	if err != nil {
		t.Errorf("error postfix, %v", err)
	}
	t.Logf("the result of shuntingyard:\n%v", postfix)

	got, err := evaluatePostfix(postfix)
	if err != nil {
		t.Errorf("error evaluation, %v", err)
	}
	t.Logf("the final result: %3.f", got)
}

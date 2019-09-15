package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

const (
	testSimpleExpression              = "1+2"
	testSimpleExpressionCorrectResult = "3"

	testComplicatedExpression              = "((123+6/2)/2+1*2-(2*3))"
	testComplicatedExpressionCorrectResult = "59"

	testIncorrectInput = "(1++2)"

	testIncorrectBrackets = "(1+2))"
)

func TestSimpleExpression(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(testSimpleExpression))
	out := new(bytes.Buffer)
	err := calc(in, out)
	if err != nil {
		t.Errorf("testSimpleExpression for OK Failed - error")
	}
	result := out.String()
	if result != testSimpleExpressionCorrectResult {
		t.Errorf("testSimpleExpression for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testSimpleExpressionCorrectResult)
	}
}

func TestComplicatedExpresssion(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(testComplicatedExpression))
	out := new(bytes.Buffer)
	err := calc(in, out)
	if err != nil {
		t.Errorf("testComplicatedExpression for OK Failed - error")
	}
	result := out.String()
	if result != testComplicatedExpressionCorrectResult {
		t.Errorf("testComplicatedExpression for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testComplicatedExpressionCorrectResult)
	}
}

func TestIncorrectInput(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(testIncorrectInput))
	out := new(bytes.Buffer)
	err := calc(in, out)
	if err == nil {
		t.Errorf("testIncorrectInput for OK Failed - error")
	}
}

func TestIncorrectBrackets(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(testIncorrectBrackets))
	out := new(bytes.Buffer)
	err := calc(in, out)
	if err == nil {
		t.Errorf("testIncorrectBrackets for OK Failed - error")
	}
}

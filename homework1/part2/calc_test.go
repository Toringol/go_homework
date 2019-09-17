package main

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestSimpleExpression(t *testing.T) {
	const testSimpleExpression = "1+2"
	const testSimpleExpressionCorrectResult = "3"

	in := bufio.NewReader(strings.NewReader(testSimpleExpression))
	out := new(bytes.Buffer)

	err := calc(in, out)

	require.Equal(t, err, nil, "TestSimpleExpression for OK Failed - error")

	require.Equal(t, out.String(), testSimpleExpressionCorrectResult, "TestSimpleExpression for OK Failed - results not match")
}

func TestComplicatedExpresssion(t *testing.T) {
	const testComplicatedExpression = "((123+6/2)/2+1*2-(2*3))"
	const testComplicatedExpressionCorrectResult = "59"

	in := bufio.NewReader(strings.NewReader(testComplicatedExpression))
	out := new(bytes.Buffer)

	err := calc(in, out)

	require.Equal(t, err, nil, "TestComplicatedExpression for OK Failed - error")

	require.Equal(t, out.String(), testComplicatedExpressionCorrectResult, "TestComplicatedExpresssion for OK Failed - results not match")
}

func TestIncorrectInput(t *testing.T) {
	const testIncorrectInput = "(1++2)"
	in := bufio.NewReader(strings.NewReader(testIncorrectInput))
	out := new(bytes.Buffer)
	err := calc(in, out)

	require.NotEqual(t, err, nil, "TestIncorrectInput for OK Failed - error")
}

func TestIncorrectBrackets(t *testing.T) {
	const testIncorrectBrackets = "(1+2))"
	in := bufio.NewReader(strings.NewReader(testIncorrectBrackets))
	out := new(bytes.Buffer)
	err := calc(in, out)
	require.NotEqual(t, err, nil, "TestIncorrectBrackets for OK Failed - error")
}

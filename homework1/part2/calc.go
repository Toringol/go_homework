package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/emirpasic/gods/stacks/arraystack"
	"io"
	"os"
	"strconv"
)

func checkInputData(expression string) error {
	bracketsCounter := 0
	operationFlag := 0

	for i := 0; i < len(expression); i++ {
		if (expression[i] < '0' || expression[i] > '9') &&
			!isOperator(string(expression[i])) &&
			expression[i] != '(' && expression[i] != ')' {
			return errors.New("Incorrect symbols")
		}

		if isOperator(string(expression[i])) && operationFlag == 1 {
			return errors.New("Incorrect sequence of operations")
		}

		if isOperator(string(expression[i])) {
			operationFlag = 1
		} else {
			operationFlag = 0
		}

		if expression[i] == '(' {
			bracketsCounter++
		}

		if expression[i] == ')' {
			bracketsCounter--
		}

		if bracketsCounter < 0 {
			return errors.New("Incorrect sequence of brackets")
		}

	}

	if bracketsCounter != 0 {
		return errors.New("Incorrect sequence of brackets")
	}

	return nil
}

func splitOnOp(expression string) []string {
	var operand string
	operandAndOperations := []string{}

	for _, ch := range expression {
		if ch >= '0' && ch <= '9' {
			operand += string(ch)
			continue
		}
		if len(operand) != 0 {
			operandAndOperations = append(operandAndOperations, operand)
			operand = ""
		}
		operandAndOperations = append(operandAndOperations, string(ch))
	}

	if len(operand) != 0 {
		operandAndOperations = append(operandAndOperations, operand)
	}

	return operandAndOperations
}

func isOperator(operation string) bool {
	if operation == "+" || operation == "-" ||
		operation == "*" || operation == "/" {
		return true
	}
	return false
}

func makeOp(operand1 int, operand2 int, operation string) int {
	if operation == "+" {
		return operand1 + operand2
	}
	if operation == "-" {
		return operand2 - operand1
	}
	if operation == "*" {
		return operand1 * operand2
	}
	if operation == "/" {
		return operand2 / operand1
	}
	return 0
}

func makeOperation(operandsStack *arraystack.Stack, operationsStack *arraystack.Stack) error {

	topOperation, okPopOperation := operationsStack.Peek()
	operand1, okPopOp1 := operandsStack.Pop()
	operand2, okPopOp2 := operandsStack.Pop()

	if !okPopOp1 || !okPopOp2 || !okPopOperation {
		return errors.New("Error in Pop")
	}

	operandsStack.Push(makeOp(operand1.(int), operand2.(int), topOperation.(string)))

	_, ok := operationsStack.Pop()

	if !ok {
		return errors.New("Error in Pop")
	}

	return nil
}

func calc(input io.Reader, output io.Writer) error {

	in := bufio.NewScanner(input)
	in.Scan()

	expression := in.Text()

	err := checkInputData(expression)

	if err != nil {
		return err
	}

	var operandAndOperations []string
	operandsStack := arraystack.New()
	operationsStack := arraystack.New()

	operandAndOperations = splitOnOp(expression)

	for _, value := range operandAndOperations {
		switch value {
		case "+", "-":
			for operationsStack.Size() != 0 {
				topOperation, _ := operationsStack.Peek()
				if isOperator(topOperation.(string)) {
					makeOperation(operandsStack, operationsStack)
				} else {
					break
				}
			}
			operationsStack.Push(value)
		case "/", "*":
			for operationsStack.Size() != 0 {
				topOperation, _ := operationsStack.Peek()
				if topOperation == "*" || topOperation == "/" {
					makeOperation(operandsStack, operationsStack)
				} else {
					break
				}
			}
			operationsStack.Push(value)
		case "(":
			operationsStack.Push("(")
		case ")":
			for {
				topOperation, _ := operationsStack.Peek()
				if topOperation == "(" {
					break
				} else {
					makeOperation(operandsStack, operationsStack)
				}
			}
			operationsStack.Pop()
		default:
			operand, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			operandsStack.Push(operand)
		}

	}

	for operationsStack.Size() != 0 {
		makeOperation(operandsStack, operationsStack)
	}

	result, ok := operandsStack.Pop()

	if !ok {
		return errors.New("No operands in stack")
	}

	fmt.Fprintf(output, "%d", result)

	return nil
}

func main() {
	err := calc(os.Stdin, os.Stdout)

	if err != nil {
		panic(err.Error())
	}
}

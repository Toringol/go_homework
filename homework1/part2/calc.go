package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type OperationsStack struct {
	Slice []string
	Pos   int
}

type OperandsStack struct {
	Slice []int
	Pos   int
}

func NewOperationsStack() *OperationsStack {
	return &OperationsStack{
		Slice: []string{},
		Pos:   0,
	}
}

func NewOperandsStack() *OperandsStack {
	return &OperandsStack{
		Slice: []int{},
		Pos:   0,
	}
}

func (operations *OperationsStack) Push(str string) {
	if operations.Pos < len(operations.Slice) {
		operations.Slice[operations.Pos] = str
	} else {
		operations.Slice = append(operations.Slice, str)
	}
	operations.Pos++
}

func (operands *OperandsStack) Push(number int) {
	if operands.Pos < len(operands.Slice) {
		operands.Slice[operands.Pos] = number
	} else {
		operands.Slice = append(operands.Slice, number)
	}
	operands.Pos++
}

func (operations *OperationsStack) Pop() (string, error) {
	ret, err := operations.Top()
	if err != nil {
		return "", errors.New("Can't pop; stack is empty!")
	}
	operations.Pos--
	return ret, nil
}

func (operands *OperandsStack) Pop() (int, error) {
	ret, err := operands.Top()
	if err != nil {
		return 0, errors.New("Can't pop; stack is empty!")
	}
	operands.Pos--
	return ret, nil
}

func (operations *OperationsStack) Top() (string, error) {
	if operations.Pos < 1 {
		return "", errors.New("No elements in stack")
	}
	return operations.Slice[operations.Pos-1], nil
}

func (operands *OperandsStack) Top() (int, error) {
	if operands.Pos < 1 {
		return 0, errors.New("No elements in stack")
	}
	return operands.Slice[operands.Pos-1], nil
}

func checkInputData(expression string) error {
	bracketsCounter := 0
	operationFlag := 0

	for i := 0; i < len(expression); i++ {
		if (expression[i] < '0' || expression[i] > '9') && expression[i] != '+' &&
			expression[i] != '-' && expression[i] != '/' && expression[i] != '*' &&
			expression[i] != '(' && expression[i] != ')' {
			return errors.New("Incorrect symbols")
		}

		if (expression[i] == '+' || expression[i] == '-' ||
			expression[i] == '/' || expression[i] == '*') && operationFlag == 1 {
			return errors.New("Incorrect sequence of operations")
		}

		if expression[i] == '+' || expression[i] == '-' ||
			expression[i] == '/' || expression[i] == '*' {
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

func makeOperation(operandsStack *OperandsStack, operationsStack *OperationsStack) {
	var (
		result   int
		operand1 int
		operand2 int
	)

	topOperation, _ := operationsStack.Top()

	if topOperation == "+" {
		operand1, _ = operandsStack.Pop()
		operand2, _ = operandsStack.Pop()
		result = operand1 + operand2
		operandsStack.Push(result)
	}
	if topOperation == "-" {
		operand1, _ = operandsStack.Pop()
		operand2, _ = operandsStack.Pop()
		result = operand2 - operand1
		operandsStack.Push(result)
	}
	if topOperation == "*" {
		operand1, _ = operandsStack.Pop()
		operand2, _ = operandsStack.Pop()
		result = operand1 * operand2
		operandsStack.Push(result)
	}
	if topOperation == "/" {
		operand1, _ = operandsStack.Pop()
		operand2, _ = operandsStack.Pop()
		result = operand2 / operand1
		operandsStack.Push(result)
	}

	operationsStack.Pop()
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
	operandsStack := NewOperandsStack()
	operationsStack := NewOperationsStack()

	operandAndOperations = splitOnOp(expression)

	for _, value := range operandAndOperations {
		switch value {
		case "+":
			for operationsStack.Pos != 0 {
				topOperation, _ := operationsStack.Top()
				if topOperation == "+" || topOperation == "-" ||
					topOperation == "*" || topOperation == "/" {
					makeOperation(operandsStack, operationsStack)
				} else {
					break
				}
			}
			operationsStack.Push("+")
		case "-":
			for operationsStack.Pos != 0 {
				topOperation, _ := operationsStack.Top()
				if topOperation == "+" || topOperation == "-" ||
					topOperation == "*" || topOperation == "/" {
					makeOperation(operandsStack, operationsStack)
				} else {
					break
				}
			}
			operationsStack.Push("-")
		case "/":
			for operationsStack.Pos != 0 {
				topOperation, _ := operationsStack.Top()
				if topOperation == "*" || topOperation == "/" {
					makeOperation(operandsStack, operationsStack)
				} else {
					break
				}
			}
			operationsStack.Push("/")
		case "*":
			for operationsStack.Pos != 0 {
				topOperation, _ := operationsStack.Top()
				if topOperation == "*" || topOperation == "/" {
					makeOperation(operandsStack, operationsStack)
				} else {
					break
				}
			}
			operationsStack.Push("*")
		case "(":
			operationsStack.Push("(")
		case ")":
			for {
				topOperation, _ := operationsStack.Top()
				if topOperation == "(" {
					break
				} else {
					makeOperation(operandsStack, operationsStack)
				}
			}
			operationsStack.Pop()
		default:
			operand, _ := strconv.Atoi(value)
			operandsStack.Push(operand)
		}

	}

	for operationsStack.Pos != 0 {
		makeOperation(operandsStack, operationsStack)
	}

	result, _ := operandsStack.Pop()

	fmt.Fprintf(output, "%d", result)

	return nil
}

func main() {
	err := calc(os.Stdin, os.Stdout)

	if err != nil {
		panic(err.Error())
	}
}

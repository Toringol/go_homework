package main

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

const testFile = `data.txt`
const testFileColumn = `column.txt`
const testFileNumber = `numbers.txt`
const testFileOutput = `output.txt`
const testFileForAllFlags = `test.txt`
const testBadFile = `data.txt data.txt`
const testNoExistFile = `IamNotExist.txt`

const testWithoutFlagsCorrectOutput = `Apple
BOOK
Book
Go
Hauptbahnhof
January
January
Napkin
`

const testReverseFlagCorrectOutput = `Napkin
January
January
Hauptbahnhof
Go
Book
BOOK
Apple
`

const testFirstEntryCorrectOutput = `Apple
BOOK
Book
Go
Hauptbahnhof
January
Napkin
`

const testFirstEntryAndLetterCaseCorrectOutput = `Apple
BOOK
Go
Hauptbahnhof
January
Napkin
`

const testColumnSortCorrectOutput = `Napkin Apple
January BOOK
Hauptbahnhof Book
Go Book
Book Go
January Hauptbahnhof
Apple January
BOOK January
`

const testNumbersCorrectOutput = `1
2
3
4
5
`

const testAllFlagsCorrectOutput = `Apple January
January Hauptbahnhof
Book Go
January BOOK
Napkin Apple
`

func TestWithoutFlags(t *testing.T) {
	in := testFile
	out := new(bytes.Buffer)
	err := sortUtil(in, out, false, false, false, "", false, 0)
	if err != nil {
		t.Errorf("TestWithoutFlags for OK Failed - error")
	}
	result := out.String()
	if result != testWithoutFlagsCorrectOutput {
		t.Errorf("TestWithoutFlags for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testWithoutFlagsCorrectOutput)
	}
}

func TestFlagReverse(t *testing.T) {
	in := testFile
	out := new(bytes.Buffer)
	err := sortUtil(in, out, false, false, true, "", false, 0)
	if err != nil {
		t.Errorf("TestFlagReverse for OK Failed - error")
	}
	result := out.String()
	if result != testReverseFlagCorrectOutput {
		t.Errorf("TestFlagReverse for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testReverseFlagCorrectOutput)
	}
}

func TestFlagFirstEntry(t *testing.T) {
	in := testFile
	out := new(bytes.Buffer)
	err := sortUtil(in, out, false, true, false, "", false, 0)
	if err != nil {
		t.Errorf("TestFlagFirstEntry for OK Failed - error")
	}
	result := out.String()
	if result != testFirstEntryCorrectOutput {
		t.Errorf("TestFlagFirstEntry for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testFirstEntryCorrectOutput)
	}
}

func TestFlagFirstEntryAndLetterCase(t *testing.T) {
	in := testFile
	out := new(bytes.Buffer)
	err := sortUtil(in, out, true, true, false, "", false, 0)
	if err != nil {
		t.Errorf("TestFlagFirstEntryAndLetterCase for OK Failed - error")
	}
	result := out.String()
	if result != testFirstEntryAndLetterCaseCorrectOutput {
		t.Errorf("TestFlagFirstEntryAndLetterCase for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testFirstEntryAndLetterCaseCorrectOutput)
	}
}

func TestColumnSort(t *testing.T) {
	in := testFileColumn
	out := new(bytes.Buffer)
	err := sortUtil(in, out, false, false, false, "", false, 1)
	if err != nil {
		t.Errorf("TestColumnSort for OK Failed - error")
	}
	result := out.String()
	if result != testColumnSortCorrectOutput {
		t.Errorf("TestColumnSort for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testColumnSortCorrectOutput)
	}
}

func TestNumbers(t *testing.T) {
	in := testFileNumber
	out := new(bytes.Buffer)
	err := sortUtil(in, out, false, false, false, "", true, 0)
	if err != nil {
		t.Errorf("TestNumbers for OK Failed - error")
	}
	result := out.String()
	if result != testNumbersCorrectOutput {
		t.Errorf("TestNumbers for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testNumbersCorrectOutput)
	}
}

func TestOutputFile(t *testing.T) {
	in := testFile
	out := new(bytes.Buffer)
	err := sortUtil(in, out, false, false, false, testFileOutput, false, 0)
	if err != nil {
		t.Errorf("TestOutputFile for OK Failed - error")
	}
	var result string

	file, err := os.Open(testFileOutput)

	if err != nil {
		t.Errorf("TestOutputFile for OK Failed - error")
	}

	defer file.Close()

	input := bufio.NewScanner(file)
	const delim = "\n"

	for input.Scan() {
		line := input.Text()
		result += line + delim
	}

	if result != testWithoutFlagsCorrectOutput {
		t.Errorf("TestOutputFile for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testWithoutFlagsCorrectOutput)
	}
}

func TestAllFlags(t *testing.T) {
	in := testFileForAllFlags
	out := new(bytes.Buffer)
	err := sortUtil(in, out, true, true, true, "", false, 1)
	if err != nil {
		t.Errorf("TestAllFlags for OK Failed - error")
	}
	result := out.String()
	if result != testAllFlagsCorrectOutput {
		t.Errorf("TestAllFlags for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testAllFlagsCorrectOutput)
	}
}

func TestIncorrectFile(t *testing.T) {
	in := testBadFile
	out := new(bytes.Buffer)
	err := sortUtil(in, out, false, false, false, "", false, 0)
	if err == nil {
		t.Errorf("TestIncorrectFile for OK Failed - error")
	}
}

func TestNoExistFile(t *testing.T) {
	in := testNoExistFile
	out := new(bytes.Buffer)
	err := sortUtil(in, out, false, false, false, "", false, 0)
	if err == nil {
		t.Errorf("TestNoExistFile for OK Failed - error")
	}
}

package main

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestWithoutFlags(t *testing.T) {
	const testWithoutFlags = `Napkin
Apple
January
BOOK
January
Hauptbahnhof
Book
Go
`
	const testWithoutFlagsCorrectOutput = `Apple
BOOK
Book
Go
Hauptbahnhof
January
January
Napkin
`

	in := bufio.NewReader(strings.NewReader(testWithoutFlags))
	out := new(bytes.Buffer)

	lines, err := readInput(in)
	if err != nil {
		t.Errorf("TestWithoutFlags for OK Failed - error")
	}

	opts := Opts{}

	sortedLines, errSort := sortUtil(lines, opts)

	if errSort != nil {
		t.Errorf("TestWithoutFlags for OK Failed - error")
	}

	writeResult(out, sortedLines)

	require.Equal(t, out.String(), testWithoutFlagsCorrectOutput, "TestWithoutFlags for OK Failed - results not match")
}

func TestFlagReverse(t *testing.T) {
	const testReverse = `Napkin
Apple
January
BOOK
January
Hauptbahnhof
Book
Go
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

	in := bufio.NewReader(strings.NewReader(testReverse))
	out := new(bytes.Buffer)

	lines, err := readInput(in)
	if err != nil {
		t.Errorf("TestFlagReverse for OK Failed - error")
	}

	opts := Opts{
		Reverse: true,
	}

	sortedLines, errSort := sortUtil(lines, opts)

	if errSort != nil {
		t.Errorf("TestFlagReverse for OK Failed - error")
	}

	writeResult(out, sortedLines)

	require.Equal(t, out.String(), testReverseFlagCorrectOutput, "TestFlagReverse for OK Failed - results not match")
}

func TestFlagFirstEntry(t *testing.T) {
	const testFirstEntry = `Napkin
Apple
January
BOOK
January
Hauptbahnhof
Book
Go
`
	const testFirstEntryCorrectOutput = `Apple
BOOK
Book
Go
Hauptbahnhof
January
Napkin
`

	in := bufio.NewReader(strings.NewReader(testFirstEntry))
	out := new(bytes.Buffer)

	lines, err := readInput(in)
	if err != nil {
		t.Errorf("TestFlagFirstEntry for OK Failed - error")
	}

	opts := Opts{
		FirstEntry: true,
	}

	sortedLines, errSort := sortUtil(lines, opts)

	if errSort != nil {
		t.Errorf("TestFlagFirstEntry for OK Failed - error")
	}

	writeResult(out, sortedLines)

	require.Equal(t, out.String(), testFirstEntryCorrectOutput, "TestFlagFirstEntry for OK Failed - results not match")

}

func TestFlagFirstEntryAndLetterCase(t *testing.T) {
	const testFirstEntryAndLetterCase = `Napkin
Apple
January
BOOK
January
Hauptbahnhof
Book
Go
`
	const testFirstEntryAndLetterCaseCorrectOutput = `Apple
BOOK
Go
Hauptbahnhof
January
Napkin
`

	in := bufio.NewReader(strings.NewReader(testFirstEntryAndLetterCase))
	out := new(bytes.Buffer)

	lines, err := readInput(in)
	if err != nil {
		t.Errorf("TestFlagFirstEntryAndLetterCase for OK Failed - error")
	}

	opts := Opts{
		FirstEntry: true,
		LetterCase: true,
	}

	sortedLines, errSort := sortUtil(lines, opts)

	if errSort != nil {
		t.Errorf("TestFlagFirstEntryAndLetterCase for OK Failed - error")
	}

	writeResult(out, sortedLines)

	require.Equal(t, out.String(), testFirstEntryAndLetterCaseCorrectOutput, "TestFlagFirstEntryAndLetterCase for OK Failed - results not match")

}

func TestColumnSort(t *testing.T) {
	const testColumnSort = `Napkin Apple
Apple January
January BOOK
BOOK January
January Hauptbahnhof
Hauptbahnhof Book
Book Go
Go Book
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

	in := bufio.NewReader(strings.NewReader(testColumnSort))
	out := new(bytes.Buffer)

	lines, err := readInput(in)
	if err != nil {
		t.Errorf("TestColumnSort for OK Failed - error")
	}

	opts := Opts{
		Column: 1,
	}

	sortedLines, errSort := sortUtil(lines, opts)

	if errSort != nil {
		t.Errorf("TestColumnSort for OK Failed - error")
	}

	writeResult(out, sortedLines)

	require.Equal(t, out.String(), testColumnSortCorrectOutput, "TestColumnSort for OK Failed - results not match")

}

func TestNumbers(t *testing.T) {
	const testSortNumbers = `2
5
3
4
1
`
	const testNumbersCorrectOutput = `1
2
3
4
5
`

	in := bufio.NewReader(strings.NewReader(testSortNumbers))
	out := new(bytes.Buffer)

	lines, err := readInput(in)
	if err != nil {
		t.Errorf("TestNumbers for OK Failed - error")
	}

	opts := Opts{
		SortNumbers: true,
	}

	sortedLines, errSort := sortUtil(lines, opts)

	if errSort != nil {
		t.Errorf("TestNumbers for OK Failed - error")
	}

	writeResult(out, sortedLines)

	require.Equal(t, out.String(), testNumbersCorrectOutput, "TestNumbers for OK Failed - results not match")

}

func TestOutputFile(t *testing.T) {
	const testSortNumbers = `2
5
3
4
1
`
	const testNumbersCorrectOutput = `1
2
3
4
5
`

	in := bufio.NewReader(strings.NewReader(testSortNumbers))

	lines, err := readInput(in)
	if err != nil {
		t.Errorf("TestOutputFile for OK Failed - error")
	}

	opts := Opts{
		SortNumbers:  true,
		DirectOutput: "test.txt",
	}

	sortedLines, errSort := sortUtil(lines, opts)

	if errSort != nil {
		t.Errorf("TestOutputFile for OK Failed - error")
	}

	out, err := os.Create("test.txt")

	if err != nil {
		t.Errorf("TestOutputFile for OK Failed - error")
	}

	defer out.Close()

	writeResult(out, sortedLines)

	var result string

	file, err := os.Open("test.txt")

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

	require.Equal(t, result, testNumbersCorrectOutput, "TestOutputFile for OK Failed - results not match")
}

func TestAllFlags(t *testing.T) {
	const testAllFlags = `Napkin Apple
Apple January
January BOOK
BOOK January
January Hauptbahnhof
Hauptbahnhof Book
Book Go
Go Book
`
	const testAllFlagsCorrectOutput = `Apple January
January Hauptbahnhof
Book Go
January BOOK
Napkin Apple
`

	in := bufio.NewReader(strings.NewReader(testAllFlags))
	out := new(bytes.Buffer)

	lines, err := readInput(in)
	if err != nil {
		t.Errorf("TestAllFlags for OK Failed - error")
	}

	opts := Opts{
		LetterCase: true,
		Column:     1,
		FirstEntry: true,
		Reverse:    true,
	}

	sortedLines, errSort := sortUtil(lines, opts)

	if errSort != nil {
		t.Errorf("TestAllFlags for OK Failed - error")
	}

	writeResult(out, sortedLines)

	require.Equal(t, out.String(), testAllFlagsCorrectOutput, "TestAllFlags for OK Failed - results not match")

}

func TestIncorrectFile(t *testing.T) {
	const testWithoutFlags = `1
Apple
3
2
5
6
8
10
`

	in := bufio.NewReader(strings.NewReader(testWithoutFlags))
	lines, err := readInput(in)
	if err != nil {
		t.Errorf("TestWithoutFlags for OK Failed - error")
	}

	opts := Opts{
		SortNumbers: true,
	}

	_, errSort := sortUtil(lines, opts)

	if errSort == nil {
		t.Errorf("TestWithoutFlags for OK Failed - error")
	}
}

func TestNoExistFile(t *testing.T) {
	const testWithoutFlags = `Napkin
Apple
January
BOOK
January
Hauptbahnhof
Book
Go
`

	in := bufio.NewReader(strings.NewReader(testWithoutFlags))
	lines, err := readInput(in)
	if err != nil {
		t.Errorf("TestWithoutFlags for OK Failed - error")
	}

	opts := Opts{
		Column: 1,
	}

	_, errSort := sortUtil(lines, opts)

	if errSort == nil {
		t.Errorf("TestWithoutFlags for OK Failed - error")
	}
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func readFile(input string) (lines []string, err error) {
	file, err := os.Open(input)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	in := bufio.NewScanner(file)

	for in.Scan() {
		line := in.Text()
		lines = append(lines, line)
	}

	return lines, nil
}

func writeLines(lines []string, output io.Writer) {
	for _, line := range lines {
		fmt.Fprintln(output, line)
	}
}

func writeIntoFile(lines []string, output string) error {
	file, err := os.Create(output)

	if err != nil {
		return err
	}

	defer file.Close()

	const delim = '\n'
	for _, text := range lines {
		file.WriteString(text + string(delim))

	}
	return nil
}

type byLetterCase []string

func (s byLetterCase) Len() int {
	return len(s)
}

func (s byLetterCase) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byLetterCase) Less(i, j int) bool {
	return strings.ToLower(s[i]) < strings.ToLower(s[j])
}

func sortUtil(input string, output io.Writer, letterCase bool, firstEntry bool,
	reverse bool, directOutput string, sortNumbers bool, columnSort int) error {
	lines, err := readFile(input)

	if err != nil {
		return err
	}

	if reverse {
		if letterCase {
			sort.Sort(sort.Reverse(byLetterCase(lines)))
		} else {
			sort.Sort(sort.Reverse(sort.StringSlice(lines)))
		}
	} else {
		if letterCase {
			sort.Slice(lines, func(i, j int) bool {
				return strings.ToLower(lines[i]) < strings.ToLower(lines[j])
			})
		} else {
			sort.Strings(lines)
		}
	}

	if directOutput == "" {
		writeLines(lines, output)
	} else {
		err := writeIntoFile(lines, directOutput)
		if err != nil {
			return err
		}

	}
	return nil
}

func main() {

	letterCasePtr := flag.Bool("f", false, "letterCase")
	firstEntryPtr := flag.Bool("u", false, "firstEntry")
	reversePtr := flag.Bool("r", false, "reverse")
	directOutputPtr := flag.String("o", "", "directOutput")
	sortNumbersPtr := flag.Bool("n", false, "sortNumbers")
	columnSortPtr := flag.Int("k", 0, "columnSort")

	flag.Parse()

	err := sortUtil(flag.Args()[0], os.Stdout, *letterCasePtr, *firstEntryPtr,
		*reversePtr, *directOutputPtr, *sortNumbersPtr, *columnSortPtr)

	if err != nil {
		panic(err.Error())
	}
}

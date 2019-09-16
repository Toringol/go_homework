package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	letterCasePtr   = flag.Bool("f", false, "letterCase")
	firstEntryPtr   = flag.Bool("u", false, "firstEntry")
	reversePtr      = flag.Bool("r", false, "reverse")
	directOutputPtr = flag.String("o", "", "directOutput")
	sortNumbersPtr  = flag.Bool("n", false, "sortNumbers")
	columnSortPtr   = flag.Int("k", 0, "columnSort")
)

type Opts struct {
	Reverse      bool
	Column       int
	LetterCase   bool
	FirstEntry   bool
	DirectOutput string
	SortNumbers  bool
}

func readInput(input io.Reader) (lines []string, err error) {

	in := bufio.NewScanner(input)

	for in.Scan() {
		line := in.Text()
		lines = append(lines, line)
	}

	if err := in.Err(); err != nil {
		emptyLine := []string{}
		return emptyLine, err
	}

	return lines, nil
}

func writeResult(output io.Writer, res []string) {
	const delim = "\n"
	for _, line := range res {
		output.Write([]byte(line))
		output.Write([]byte(delim))
	}
}

func sliceFromString(lines []string) [][]string {
	var result [][]string

	for _, line := range lines {
		result = append(result, strings.Split(line, " "))
	}

	return result
}

func reverseSlice(lines []string) {
	for left, right := 0, len(lines)-1; left < right; left, right = left+1, right-1 {
		lines[left], lines[right] = lines[right], lines[left]
	}
}

func uniq(lines []string, letterCase bool) []string {
	keys := make(map[string]bool)
	list := []string{}
	var checkEntry string
	for _, entry := range lines {
		if letterCase {
			checkEntry = strings.ToLower(entry)
		} else {
			checkEntry = entry
		}
		if _, value := keys[checkEntry]; !value {
			keys[checkEntry] = true
			list = append(list, entry)
		}
	}
	return list
}

func uniqColumn(lines []string, letterCase bool, column int) []string {
	sliceSliceString := sliceFromString(lines)

	keys := make(map[string]bool)
	list := [][]string{}
	var checkEntry string
	for i, entry := range sliceSliceString {
		if letterCase {
			checkEntry = strings.ToLower(entry[column])
		} else {
			checkEntry = entry[column]
		}
		if _, value := keys[checkEntry]; !value {
			keys[checkEntry] = true
			list = append(list, sliceSliceString[i])
		}
	}

	var result []string
	for _, line := range list {
		result = append(result, strings.Join(line[:], " "))
	}

	return result
}

func sortNumbers(lines []string) error {
	values := []int{}
	for _, line := range lines {
		number, err := strconv.Atoi(line)
		if err != nil {
			return err
		}
		values = append(values, number)
	}
	sort.Ints(values)
	for i, line := range values {
		lines[i] = strconv.Itoa(line)
	}
	return nil
}

func sortStrings(lines []string, letterCase bool) {
	if letterCase {
		sort.SliceStable(lines, func(i, j int) bool {
			return strings.ToLower(lines[i]) < strings.ToLower(lines[j])
		})
	} else {
		sort.Strings(lines)
	}
}

func columnSort(lines []string, letterCase bool, column int) error {
	sliceSliceString := sliceFromString(lines)
	for _, value := range sliceSliceString {
		if column > len(value)-1 {
			return errors.New("Out of range")
		}
	}
	if letterCase {
		sort.SliceStable(sliceSliceString, func(i, j int) bool {
			return strings.ToLower(sliceSliceString[i][column]) <
				strings.ToLower(sliceSliceString[j][column])
		})
	} else {
		sort.SliceStable(sliceSliceString, func(i, j int) bool {
			return sliceSliceString[i][column] < sliceSliceString[j][column]
		})
	}
	for index, line := range sliceSliceString {
		lines[index] = strings.Join(line[:], " ")
	}
	return nil
}

func sortUtil(lines []string, opts Opts) ([]string, error) {

	if opts.Column != 0 {
		err := columnSort(lines, opts.LetterCase, opts.Column)
		if err != nil {
			emptyLine := []string{}
			return emptyLine, err
		}
	} else if opts.SortNumbers {
		err := sortNumbers(lines)
		if err != nil {
			emptyLine := []string{}
			return emptyLine, err
		}
	} else {
		sortStrings(lines, opts.LetterCase)
	}

	if opts.FirstEntry {
		if opts.Column != 0 {
			lines = uniqColumn(lines, opts.LetterCase, opts.Column)
		} else {
			lines = uniq(lines, opts.LetterCase)
		}
	}

	if opts.Reverse {
		reverseSlice(lines)
	}

	return lines, nil
}

func main() {

	flag.Parse()

	opts := Opts{
		Reverse:      *reversePtr,
		Column:       *columnSortPtr,
		LetterCase:   *letterCasePtr,
		FirstEntry:   *firstEntryPtr,
		DirectOutput: *directOutputPtr,
		SortNumbers:  *sortNumbersPtr,
	}

	out := os.Stdout

	if len(flag.Args()) > 1 {
		log.Fatal("Too many argument")
	}

	var lines []string

	if len(flag.Args()) == 1 {
		in, err := os.Open(flag.Args()[0])

		if err != nil {
			log.Fatal("Can`t open file")
		}

		defer in.Close()

		lines, err = readInput(in)

		if err != nil {
			log.Fatal("Mistake in read")
		}

	} else {
		var err error
		lines, err = readInput(os.Stdin)

		if err != nil {
			log.Fatal("Mistake in read")
		}
	}

	sortedLines, err := sortUtil(lines, opts)

	if opts.DirectOutput == "" {
		writeResult(out, sortedLines)
	} else {
		file, err := os.Create(opts.DirectOutput)

		if err != nil {
			log.Fatal("Can`t create or open file")
		}

		defer file.Close()

		writeResult(file, lines)
	}

	if err != nil {
		log.Fatal("Error")
	}
}

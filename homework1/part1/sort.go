package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
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

func columnSort(lines []string, letterCase bool, column int) {
	sliceSliceString := sliceFromString(lines)
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
}

func sortUtil(input string, output io.Writer, letterCase bool, firstEntry bool,
	reverse bool, directOutput string, sortNum bool, column int) error {
	lines, err := readFile(input)

	if err != nil {
		return err
	}

	if column != 0 {
		columnSort(lines, letterCase, column)
	} else if sortNum {
		err := sortNumbers(lines)
		if err != nil {
			return err
		}
	} else {
		sortStrings(lines, letterCase)
	}

	if firstEntry {
		if column != 0 {
			lines = uniqColumn(lines, letterCase, column)
		} else {
			lines = uniq(lines, letterCase)
		}
	}

	if reverse {
		reverseSlice(lines)
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

	if len(flag.Args()) > 1 {
		panic("Too many argument")
	}

	err := sortUtil(flag.Args()[0], os.Stdout, *letterCasePtr, *firstEntryPtr,
		*reversePtr, *directOutputPtr, *sortNumbersPtr, *columnSortPtr)

	if err != nil {
		panic(err.Error())
	}
}

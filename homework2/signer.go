package main

import (
	"fmt"
	"sort"
	"strconv"
)

func SingleHash(data int) string {

	var strData = strconv.Itoa(data)

	var crt32Data = DataSignerCrc32(strData)
	var crt32Md5Data = DataSignerCrc32(DataSignerMd5(strData))

	return crt32Data + "~" + crt32Md5Data
}

func MultiHash(data int) string {
	var result string
	for th := 0; th < 6; th++ {
		result += DataSignerCrc32(strconv.Itoa(th) + SingleHash(data))
	}
	return result
}

func CombineResults(input []int) string {

	if len(input) == 0 {
		return ""
	}

	var result string

	var hashArr []string
	for _, item := range input {
		hashArr = append(hashArr, MultiHash(item))
	}
	sort.Strings(hashArr)

	if len(hashArr) > 1 {
		for index, item := range hashArr {
			result += item
			if index < len(hashArr)-1 {
				result += "_"
			}
		}
	} else {
		result = hashArr[0]
	}

	return result
}

func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	fmt.Println(CombineResults(inputData))
}

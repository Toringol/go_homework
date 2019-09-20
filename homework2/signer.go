package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

func SingleHash(in, out chan interface{}) {

	wg := &sync.WaitGroup{}

LOOP:
	for {
		select {
		case data := <-in:
			var strData = strconv.Itoa(data.(int))
			var crt32Data string
			var crt32Md5Data string
			wg.Add(1)

			go func(data string) {
				defer wg.Done()
				crt32Data = DataSignerCrc32(data)
			}(strData)

			wg.Add(1)

			go func(data string) {
				defer wg.Done()
				crt32Md5Data = DataSignerCrc32(DataSignerMd5(data))
			}(strData)

			wg.Wait()
			fmt.Println("SingleHash: " + crt32Data + "~" + crt32Md5Data)
			out <- crt32Data + "~" + crt32Md5Data

		default:
			break LOOP
		}

	}
}

func MultiHash(in, out chan interface{}) {

	wg := &sync.WaitGroup{}

LOOP:
	for {
		select {
		case data := <-in:
			var strData = data.(string)
			var result string
			var futuresArr [6]string

			for th := 0; th < 6; th++ {
				wg.Add(1)

				go func(index int, data string) {
					defer wg.Done()
					futuresArr[index] = DataSignerCrc32(strconv.Itoa(index) + data)
					var testIndex string
					testIndex = strconv.Itoa(index)
					fmt.Println(data + " " + "MultiHash: " + testIndex + " " + futuresArr[index])
				}(th, strData)
			}

			wg.Wait()

			for _, item := range futuresArr {
				result += item
			}

			out <- result

		default:
			break LOOP
		}
	}
}

func CombineResults(in, out chan interface{}) {

	var result string

	var hashArr []string

LOOP:
	for {
		select {
		case data := <-in:
			hashArr = append(hashArr, data.(string))
		default:
			break LOOP
		}
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
	fmt.Println(result)
	out <- result
}

func main() {
	chanIn := make(chan interface{}, 2)
	chanOut := make(chan interface{}, 2)

	chanIn <- 0
	chanIn <- 1

	SingleHash(chanIn, chanOut)

	chanIn = chanOut
	chanOut = make(chan interface{}, 2)

	MultiHash(chanIn, chanOut)

	chanIn = chanOut
	chanOut = make(chan interface{}, 2)

	CombineResults(chanIn, chanOut)
}

package main

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ExecutePipeline(freeFlowJobs ...job) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	in := make(chan interface{})
	out := make(chan interface{})

	for _, freeJob := range freeFlowJobs {
		in = out
		out = make(chan interface{})
		wg.Add(1)

		go func(jobFunc job, input, output chan interface{}) {
			defer wg.Done()
			jobFunc(input, output)
			close(output)
			runtime.Gosched()
		}(freeJob, in, out)
	}
}

func SingleHash(in, out chan interface{}) {

	wgOutput := &sync.WaitGroup{}
	defer wgOutput.Wait()

LOOP:
	for {
		// Выявлено опытным путем
		time.Sleep(time.Millisecond * 11)
		select {
		case dataInterface := <-in:

			data, ok := dataInterface.(int)
			if !ok {
				fmt.Println("SingleHash convert error!")
				break LOOP
			}

			var strData = strconv.Itoa(data)
			var crt32Data string
			var crt32Md5Data string

			wgOutput.Add(1)

			go func(targetData string, wgOut *sync.WaitGroup) {
				defer wgOut.Done()
				wgInput := &sync.WaitGroup{}
				wgInput.Add(2)

				go func(data string, wgIn *sync.WaitGroup) {
					defer wgIn.Done()
					crt32Data = DataSignerCrc32(data)
					runtime.Gosched()
				}(targetData, wgInput)

				go func(data string, wgIn *sync.WaitGroup) {
					defer wgIn.Done()
					crt32Md5Data = DataSignerCrc32(DataSignerMd5(data))
					runtime.Gosched()
				}(targetData, wgInput)

				wgInput.Wait()

				fmt.Println(targetData + " SingleHash result: " + crt32Data + "~" + crt32Md5Data)
				out <- crt32Data + "~" + crt32Md5Data

			}(strData, wgOutput)
		}
	}
}

func MultiHash(in, out chan interface{}) {

	wgOutput := &sync.WaitGroup{}
	defer wgOutput.Wait()

LOOP:
	for {
		select {
		case dataInterface := <-in:
			data, ok := dataInterface.(string)
			if !ok {
				fmt.Println("Multihash convert error!")
				break LOOP
			}

			var strData = data
			var result string
			var futuresArr [6]string

			wgOutput.Add(1)

			go func(targetData string, wgOut *sync.WaitGroup) {
				defer wgOut.Done()
				wgInput := &sync.WaitGroup{}

				for th := 0; th < 6; th++ {
					wgInput.Add(1)

					go func(index int, data string, wgIn *sync.WaitGroup) {
						defer wgIn.Done()
						futuresArr[index] = DataSignerCrc32(strconv.Itoa(index) + data)
						fmt.Println(data + " " + "MultiHash: " + strconv.Itoa(index) + futuresArr[index])
						runtime.Gosched()
					}(th, targetData, wgInput)
				}

				wgInput.Wait()

				for _, item := range futuresArr {
					result += item
				}

				out <- result
			}(strData, wgOutput)
		}
	}
}

func CombineResults(in, out chan interface{}) {

	var result []string

LOOP:
	for {
		select {
		case dataInterface := <-in:
			data, ok := dataInterface.(string)
			if !ok {
				fmt.Println("CombineResults convert error!")
				break LOOP
			}
			result = append(result, data)
		}
	}
	sort.Strings(result)

	fmt.Println("CombineResults: " + strings.Join(result, "_"))

	out <- strings.Join(result, "_")
}

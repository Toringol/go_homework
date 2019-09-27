package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
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

		go runJob(freeJob, in, out, wg)
	}
}

func runJob(jobFunc job, input, output chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	jobFunc(input, output)
	close(output)
}

func SingleHash(in, out chan interface{}) {
	mu := &sync.Mutex{}
	wgOutput := &sync.WaitGroup{}
	defer wgOutput.Wait()

	for dataInterface := range in {
		var (
			strData string
			data    int
			ok      bool
			err     error
		)

		switch dataInterface.(type) {
		case int:
			data, ok = dataInterface.(int)
		case string:
			data, err = strconv.Atoi(dataInterface.(string))
		}

		if !ok || err != nil {
			fmt.Println("SingleHash convert error!")
			break
		}

		strData = strconv.Itoa(data)

		wgOutput.Add(1)
		go processCalculatingSingleHash(out, wgOutput, mu, strData)
	}
}

func processCalculatingSingleHash(out chan interface{}, wgOutput *sync.WaitGroup, mu *sync.Mutex, data string) {
	defer wgOutput.Done()

	mu.Lock()
	md5hash := DataSignerMd5(data)
	mu.Unlock()

	savedData := map[string]string{
		"data":    data,
		"md5hash": md5hash,
	}

	resultData := make(map[string]string, 2)

	wgInput := &sync.WaitGroup{}
	for key := range savedData {
		wgInput.Add(1)
		go calculatingSingleHash(key, savedData, resultData, wgInput, mu)
	}
	wgInput.Wait()

	result := resultData["data"] + "~" + resultData["md5hash"]
	out <- result
}

func calculatingSingleHash(key string, savedData map[string]string, resultData map[string]string, wgInput *sync.WaitGroup, mu *sync.Mutex) {
	defer wgInput.Done()
	hash := DataSignerCrc32(savedData[key])
	mu.Lock()
	resultData[key] = hash
	mu.Unlock()
}

func MultiHash(in, out chan interface{}) {

	wgOutput := &sync.WaitGroup{}
	defer wgOutput.Wait()

	for dataInterface := range in {
		data, ok := dataInterface.(string)
		if !ok {
			fmt.Println("Multihash convert error!")
			break
		}

		wgOutput.Add(1)

		go processCalculatingMultiHash(data, wgOutput, out)
	}
}

func processCalculatingMultiHash(targetData string, wgOut *sync.WaitGroup, out chan interface{}) {

	defer wgOut.Done()
	wgInput := &sync.WaitGroup{}

	var result string
	var futuresArr [6]string

	for th := 0; th < 6; th++ {
		wgInput.Add(1)

		go calculateMultihashByIndex(th, targetData, wgInput, &futuresArr)
	}

	wgInput.Wait()

	for _, item := range futuresArr {
		result += item
	}

	out <- result
}

func calculateMultihashByIndex(index int, data string, wgIn *sync.WaitGroup, futuresArr *[6]string) {
	defer wgIn.Done()
	mu := &sync.Mutex{}
	mu.Lock()
	(*futuresArr)[index] = DataSignerCrc32(strconv.Itoa(index) + data)
	mu.Unlock()
	fmt.Println(data + " " + "MultiHash: " + strconv.Itoa(index) + futuresArr[index])
}

func CombineResults(in, out chan interface{}) {

	var result []string

	for dataInterface := range in {
		data, ok := dataInterface.(string)
		if !ok {
			fmt.Println("CombineResults convert error!")
			break
		}
		result = append(result, data)
	}
	sort.Strings(result)

	fmt.Println("CombineResults: " + strings.Join(result, "_"))

	out <- strings.Join(result, "_")
}

package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type MD5Hasher struct {
	mu   *sync.Mutex
	hash string
}

func (m *MD5Hasher) Hash(data string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.hash = DataSignerMd5(data)

	return m.hash
}

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
	wgOutput := &sync.WaitGroup{}
	defer wgOutput.Wait()

	var md5Hasher = MD5Hasher{
		mu:   &sync.Mutex{},
		hash: "",
	}

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
		go processCalculatingSingleHash(out, wgOutput, strData, md5Hasher)
	}
}

func processCalculatingSingleHash(out chan interface{}, wgOutput *sync.WaitGroup, data string, md5hasher MD5Hasher) {
	defer wgOutput.Done()

	hash := md5hasher.Hash(data)

	md5hasher.mu.Lock()
	savedData := map[string]string{
		"data":    data,
		"md5hash": hash,
	}
	md5hasher.mu.Unlock()

	resultData := make(map[string]string, 2)

	wgInput := &sync.WaitGroup{}
	for key := range savedData {
		wgInput.Add(1)
		go calculatingSingleHash(key, savedData, resultData, wgInput, md5hasher.mu)
	}
	wgInput.Wait()

	md5hasher.mu.Lock()
	resultString := resultData["data"] + "~" + resultData["md5hash"]
	md5hasher.mu.Unlock()

	result := resultString

	out <- result
}

func calculatingSingleHash(key string, savedData map[string]string, resultData map[string]string, wgInput *sync.WaitGroup, mu *sync.Mutex) {
	defer wgInput.Done()

	mu.Lock()
	tempData := savedData[key]
	mu.Unlock()

	hash := DataSignerCrc32(tempData)

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

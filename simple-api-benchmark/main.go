package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var ada = "123456"
var orderNumber = 123456789
var cnt = 10
var goroutineCnt = 20

var gateway = "http://tw-openapi-pt.intranet.local/pt/user-center-privilege/eWallet/v1"
var client = http.DefaultClient
var token = "eyJraWQiOiIyQXFzRWdoWDEzSmxtQVNjaWc5a0loZFN5OW9abk9SSzE0MG9scFFvYWcwPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiIzaGI3Mzd0ZG5sbDJyOG5kczA4aW5oMGhyciIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiZWNvbVwvYXBpLmFsbCIsImF1dGhfdGltZSI6MTY0OTY2NjYxMiwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmFwLW5vcnRoZWFzdC0xLmFtYXpvbmF3cy5jb21cL2FwLW5vcnRoZWFzdC0xX05aWDJUcXNHZCIsImV4cCI6MTY0OTc1MzAxMiwiaWF0IjoxNjQ5NjY2NjEyLCJ2ZXJzaW9uIjoyLCJqdGkiOiJlOTQwOTQyMi0zZjZhLTRiMjktYTczOC0xYWM0ODhiOGNlYTUiLCJjbGllbnRfaWQiOiIzaGI3Mzd0ZG5sbDJyOG5kczA4aW5oMGhyciJ9.C94ciSa4bYAhm6MSv_uAT0Tlz01HZ2htvVOYVYntf3LiveGU-49jNzOq3-35ZY2IwyocO8GR1NQx5ilxnBgdOIpu8XYgQ3Ih5EeMJ4ydlvQdy2H8WMxg4iYMmhjFW05nGD4zEhwe6oGaRqvCXjeMtcQgrPB9lqEVNXQayiYc5-7WETeY4arj_H3yqpYzyfCdV83l4OfYCTO2um7AZTGD8aYGDWlP7ImRuRvUe_jOVZbX8tXl-XjnociYqtAz8oArGkSz-RJ0b9vlr5MFqgFmQQ45tHJk5Du_kCIoOw3gBxMcOxWD1qmpphsbOc1tgWMnBjjVi_XqYGoYfZELZetrCw"

func main() {
	group := &sync.WaitGroup{}

	group.Add(8)
	go add(group, orderNumber)
	go rec(group, orderNumber)

	go consume(group, orderNumber)
	go resume(group, orderNumber)

	go add(group, orderNumber)
	go rec(group, orderNumber)

	go consume(group, orderNumber)
	go resume(group, orderNumber)

	group.Wait()
}

func rec(w *sync.WaitGroup, number int) {
	defer w.Done()

	no := int64(number)
	template := "{\"regionCode\":\"130\",\"bizCode\":\"INT\",\"language\":\"zh-TW\",\"channel\":\"INT\",\"transactionType\":\"COR\",\"aboNumber\":\"#{ada}\",\"referenceOrderNumber\":\"#{on}\"}"

	waitGroup := sync.WaitGroup{}
	for i := 0; i < cnt; i++ {
		for j := 0; j < goroutineCnt; j++ {
			waitGroup.Add(1)
			go func() {
				body := strings.Replace(strings.Replace(template, "#{ada}", ada, -1), "#{on}", strconv.Itoa(int(atomic.AddInt64(&no, 1))), -1)
				for {
					err := doRequest("recoverAccountBalance", body)
					if err == nil {
						break
					}
					log.Print("rec fail ", err)
				}
				waitGroup.Done()
			}()
		}
	}
	waitGroup.Wait()
}

func add(w *sync.WaitGroup, number int) {
	defer w.Done()

	no := int64(number)
	template := "{\"regionCode\":\"130\",\"bizCode\":\"INT\",\"language\":\"zh-TW\",\"channel\":\"INT\",\"accounts\":[{\"type\":1,\"amount\":100.01},{\"type\":2,\"amount\":100.02}],\"transactionType\":\"GOR\",\"aboNumber\":\"#{ada}\",\"referenceOrderNumber\":\"#{on}\"}"

	waitGroup := sync.WaitGroup{}
	for i := 0; i < cnt*2; i++ {
		for j := 0; j < goroutineCnt; j++ {
			waitGroup.Add(1)
			go func() {
				body := strings.Replace(strings.Replace(template, "#{ada}", ada, -1), "#{on}", strconv.Itoa(int(atomic.AddInt64(&no, 1))), -1)
				for {
					err := doRequest("addAccountBalance", body)
					if err == nil {
						break
					}
					log.Print("add fail ", err)
				}
				waitGroup.Done()
			}()
		}
	}
	waitGroup.Wait()
}

func resume(w *sync.WaitGroup, number int) {
	defer w.Done()

	no := int64(number)
	template := "{\"regionCode\":\"130\",\"bizCode\":\"INT\",\"language\":\"zh-TW\",\"channel\":\"INT\",\"transactionType\":\"CCR\",\"aboNumber\":\"#{ada}\",\"referenceOrderNumber\":\"#{on}\"}"

	waitGroup := sync.WaitGroup{}
	for i := 0; i < cnt; i++ {
		for j := 0; j < goroutineCnt; j++ {
			waitGroup.Add(1)
			go func() {
				body := strings.Replace(strings.Replace(template, "#{ada}", ada, -1), "#{on}", strconv.Itoa(int(atomic.AddInt64(&no, 1))), -1)
				for {
					err := doRequest("resumeAccountBalance", body)
					if err == nil {
						break
					}
					log.Print("resume fail ", err)
				}
				waitGroup.Done()
			}()
		}
	}
	waitGroup.Wait()
}

func consume(w *sync.WaitGroup, number int) {
	defer w.Done()

	no := int64(number)
	template := "{\"regionCode\":\"130\",\"bizCode\":\"INT\",\"language\":\"zh-TW\",\"channel\":\"INT\",\"accounts\":[{\"type\":1,\"amount\":100.01},{\"type\":2,\"amount\":100.02}],\"transactionType\":\"UCR\",\"aboNumber\":\"#{ada}\",\"referenceOrderNumber\":\"#{on}\"}"

	waitGroup := sync.WaitGroup{}
	for i := 0; i < cnt; i++ {
		for j := 0; j < goroutineCnt; j++ {
			waitGroup.Add(1)
			go func() {
				body := strings.Replace(strings.Replace(template, "#{ada}", ada, -1), "#{on}", strconv.Itoa(int(atomic.AddInt64(&no, 1))), -1)
				for {
					err := doRequest("consumeAccountBalance", body)
					if err == nil {
						break
					}
					log.Print("consume fail ", err)
				}
				waitGroup.Done()
			}()
		}
	}
	waitGroup.Wait()
}

func doRequest(url, requestBody string) error {
	request, _ := http.NewRequest("POST", gateway+"/"+url, strings.NewReader(requestBody))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-api-key", "E8K7PMSAR217STANCBK2")
	request.Header.Add("x-apigw-api-id", "3a8r77t89b")
	request.Header.Add("Authorization", "Bearer "+token)
	do, err := client.Do(request)
	if err != nil {
		return err
	}
	defer do.Body.Close()
	b, _ := ioutil.ReadAll(do.Body)
	res := make(map[string]string)
	_ = json.Unmarshal(b, &res)
	if res["code"] != "0" && res["code"] != "UC.B.TRANSACTION_DUPLICATE.01" {
		return errors.New(res["code"])
	}
	return nil
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var ada = "123456"
var orderNumber = 123456789
var cnt = 10
var goroutineCnt = 20

var gateway = "http://tw-openapi-pt.intranet.local/pt/user-center-lifecycle/userInfo/v1/queryMemberInfo"
var client = http.DefaultClient
var token = "eyJraWQiOiIyQXFzRWdoWDEzSmxtQVNjaWc5a0loZFN5OW9abk9SSzE0MG9scFFvYWcwPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiIzaGI3Mzd0ZG5sbDJyOG5kczA4aW5oMGhyciIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiZWNvbVwvYXBpLmFsbCIsImF1dGhfdGltZSI6MTY1MjM0MTc3NCwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmFwLW5vcnRoZWFzdC0xLmFtYXpvbmF3cy5jb21cL2FwLW5vcnRoZWFzdC0xX05aWDJUcXNHZCIsImV4cCI6MTY1MjQyODE3NCwiaWF0IjoxNjUyMzQxNzc0LCJ2ZXJzaW9uIjoyLCJqdGkiOiI4OWI3ZGFlYS1hMDZhLTQyYzAtODc3YS0zYTQwZjM3ZjVjZmMiLCJjbGllbnRfaWQiOiIzaGI3Mzd0ZG5sbDJyOG5kczA4aW5oMGhyciJ9.eY4koN9p1BwS_EWygxWOPwqN0wzPJ25i-yQv1HuVt8E6XkSNTaDj7lP0EyTsYkqAARFRQpz33zpbxOnzkalRJI1ydtdd244peVUcE1C0RTr6qf6a1noj-SHFM7614bT77ReFMnZkjJI6DR1O3tgBe-Jy8HfZsV3SZ9lkh_y_4-KTQWYyP2_inuNUYHVj0_4SEokGrekZh5yLe--BsOdbfiIZlr9x2JZe3OQ5Xretefaz3g8V1_BeUeoIMncoUoLj-yOkrVEamm-Oc9g6nYz2_xd_E2SjGz2Wf_ls073JiuqSIkuIPkIn-_IifHUF_Qj8Cz5WCi4w8pUpjwgwCWVmIQ"

func main() {
	add(orderNumber)
}

func add(number int) {
	no := 0

	template := "{\n  \"searchKeys\": [\n    {\n      \"type\": \"aboNumber\",\n      \"value\": \"#{ada}\"\n    }\n  ],\n  \"scope\": [\n    \"basic\",\n    \"applicant\",\n    \"extend\",\n    \"phone\",\n    \"email\"\n  ],\n  \"regionCode\": \"130\",\n  \"channel\": \"INT\",\n  \"bizCode\": \"INT\",\n  \"language\": \"zh-CN\"\n}"
	//open a file and read the content
	file, err := ioutil.ReadFile("benchmark.txt")
	if err != nil {
		panic(err)
	}
	//split the content by \n
	lines := strings.Split(string(file), "\n")
	//create a goroutine for each line
	for _, line := range lines {
		split := strings.Split(line, "\t")

		i := strings.Split(split[0], "_")
		//split the line by ,
		body := strings.Replace(template, "#{ada}", split[1], -1)
		body = strings.Replace(body, "#{c}", i[0], -1)
		if len(i) > 1 {
			body = strings.Replace(body, "#{d}", i[1], -1)
		} else {
			body = strings.Replace(body, "#{d}", "", -1)
		}
		err = doRequest("", body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(no)
		no++
	}
}

func doRequest(url, requestBody string) error {
	request, _ := http.NewRequest("POST", gateway+"/"+url, strings.NewReader(requestBody))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-api-key", "1")
	request.Header.Add("x-apigw-api-id", "2")
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

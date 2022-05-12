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

var gateway = "http://3/pt/"
var client = http.DefaultClient
var token = "1"

func main() {
	add(orderNumber)
}

func add(number int) {
	no := 0

	template := "{\n  \"regionCode\": \"130\",\n  \"bizCode\": \"INT\",\n  \"language\": \"zh-TW\",\n  \"channel\": \"INT\",\n  \"items\": [\n    {\n      \"cityCode\": \"#{c}\",\n      \"districtCode\": \"#{d}\",\n      \"ada\": \"#{ada}\",\n      \"action\": 1,\n      \"memberFlag\": \"1\",\n      \"channel\": \"test\",\n      \"regionCode\": \"130\"\n    }\n  ]\n}"

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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//RT
const (
	url         = "https://g-db257d1d75.grafana-workspace.ap-northeast-1.amazonaws.com/api/ds/query"
	cookie      = "grafana_session=fec967c99f83e3962784707fca309e2c"
	contentType = "application/json"
	queries     = `
{
    "queries": [
        {
            "exemplar": true,
            "expr": "sum(irate(http_server_requests_seconds_sum{environment=\"pd\", app=~\"(user-center-lifecycle|user-center-privilege)\"}[2m])) by (uri) / sum(irate(http_server_requests_seconds_count{environment=\"pd\", app=~\"(user-center-lifecycle|user-center-privilege)\"}[2m])) by (uri)",
            "interval": "1m",
            "legendFormat": "{{kubernetes_pod_name}}  - {{uri}}",
            "refId": "A",
            "datasource": {
                "uid": "yPdqCed7z",
                "type": "prometheus"
            },
            "queryType": "timeSeriesQuery",
            "requestId": "36A",
            "utcOffsetSec": 28800,
            "datasourceId": 34,
            "intervalMs": 30000,
            "maxDataPoints": 580
        }
    ]
}
					`
)

var client = &http.Client{}

type query struct {
	expr          string
	requestId     string
	maxDataPoints int
}

var rt = &query{
	expr:          `sum(irate(http_server_requests_seconds_sum{environment=\"pd\", app=~\"(user-center-lifecycle|user-center-privilege)\"}[2m])) by (uri) / sum(irate(http_server_requests_seconds_count{environment=\"pd\", app=~\"(user-center-lifecycle|user-center-privilege)\"}[2m])) by (uri)`,
	requestId:     "36A",
	maxDataPoints: 511,
}

func main() {
	queryRT()
}

func queryRT() {
	var in map[string]interface{}
	err := json.Unmarshal([]byte(queries), &in)
	if err != nil {
		fmt.Println(err)
	}
	in["queries"].([]interface{})[0].(map[string]interface{})["expr"] = rt.expr
	in["queries"].([]interface{})[0].(map[string]interface{})["requestId"] = rt.requestId
	in["queries"].([]interface{})[0].(map[string]interface{})["maxDataPoints"] = rt.maxDataPoints
	in["from"] = "1652760856566"
	in["to"] = "1652761156566"

	var out map[string]interface{}
	call(in, &out)
	fmt.Println(out)
}

func cpu() {
	var out map[string]interface{}
	call(queries, &out)
	fmt.Println(out)
}

func call(in interface{}, out interface{}) {
	inStr, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(inStr))
	req, err := http.NewRequest("POST", url, strings.NewReader(string(inStr)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, out)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.Printf("Starting UC")
	uc := center{
		pod:       "(user-center-lifecycle|user-center-privilege)",
		nameSpace: "ecommerce-user-center",
		db:        "pd-ecom-uc-common-auroramysql",
		redis: `"pd-ecom-uc-lifecycle-redis-0001-001",
             "pd-ecom-uc-lifecycle-redis-0001-002",
             "pd-ecom-uc-privilege-redis-0001-001",
             "pd-ecom-uc-privilege-redis-0001-002"`,
		es: "test",
	}
	uc.call()
	log.Printf("\n")
	log.Printf("Starting MGC")
	mgc := center{
		pod:       "message-center",
		nameSpace: "ecommerce-message-center",
		db:        "pd-ecom-mgc-common-auroramysql",
		redis: `"pd-ecom-mgc-app-redis-0001-001",
			 "pd-ecom-mgc-app-redis-0001-002"`,
	}
	mgc.call()
}

const (
	url         = "https://g-db257d1d75.grafana-workspace.ap-northeast-1.amazonaws.com/api/ds/query"
	cookie      = "grafana_session=183085299e3ffada0b760b2a62d87072"
	contentType = "application/json"
	queries     = `
{
    "queries": [
        {
            "exemplar": true,
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
	rdsQueries = `
{
    "queries": [
        {
            "intervalMs": 500,
            "maxDataPoints": 511,
            "alias": "{{DBClusterIdentifier}}",
            "dimensions": {
                "DBClusterIdentifier": [
                    "{{environment}}"
                ]
            },
            "expression": "",
            "id": "",
            "matchExact": true,
            "metricEditorMode": 0,
            "metricName": "CPUUtilization",
            "metricQueryType": 0,
            "namespace": "AWS/RDS",
            "period": "",
            "refId": "A",
            "region": "ap-northeast-1",
            "statistic": "Maximum",
            "datasource": {
                "type": "cloudwatch",
                "uid": "9z4hi1Onz"
            },
            "sqlExpression": "",
            "type": "timeSeriesQuery"
        }
    ]
}`
	rdsMemoryQueries = `
{
    "queries": [
        {
            "alias": "{{InstanceId}}",
            "dimensions": {
                "InstanceId": "*"
            },
            "expression": "",
            "hide": true,
            "id": "",
            "matchExact": false,
            "metricEditorMode": 0,
            "metricName": "RDSBuffersMemory",
            "metricQueryType": 0,
            "namespace": "CustomRDSMetrics",
            "period": "",
            "queryMode": "Metrics",
            "refId": "A",
            "region": "ap-northeast-1",
            "sqlExpression": "",
            "statistic": "Average",
            "datasource": {
                "uid": "9z4hi1Onz",
                "type": "cloudwatch"
            },
            "datasourceId": 8,
            "intervalMs": 30000,
            "maxDataPoints": 536
        },
        {
            "alias": "{{InstanceId}}",
            "dimensions": {
                "InstanceId": "*"
            },
            "expression": "",
            "hide": true,
            "id": "",
            "matchExact": false,
            "metricEditorMode": 0,
            "metricName": "RDSCachedMemory",
            "metricQueryType": 0,
            "namespace": "CustomRDSMetrics",
            "period": "",
            "queryMode": "Metrics",
            "refId": "B",
            "region": "ap-northeast-1",
            "sqlExpression": "",
            "statistic": "Average",
            "datasource": {
                "uid": "9z4hi1Onz",
                "type": "cloudwatch"
            },
            "datasourceId": 8,
            "intervalMs": 30000,
            "maxDataPoints": 536
        },
        {
            "alias": "{{InstanceId}}",
            "dimensions": {
                "InstanceId": "*"
            },
            "expression": "",
            "hide": true,
            "id": "",
            "matchExact": false,
            "metricEditorMode": 0,
            "metricName": "RdsFreeMemory",
            "metricQueryType": 0,
            "namespace": "CustomRDSMetrics",
            "period": "",
            "queryMode": "Metrics",
            "refId": "C",
            "region": "ap-northeast-1",
            "sqlExpression": "",
            "statistic": "Average",
            "datasource": {
                "uid": "9z4hi1Onz",
                "type": "cloudwatch"
            },
            "datasourceId": 8,
            "intervalMs": 30000,
            "maxDataPoints": 536
        },
        {
            "alias": "{{InstanceId}}",
            "dimensions": {
                "InstanceId": "*"
            },
            "expression": "",
            "hide": true,
            "id": "",
            "matchExact": false,
            "metricEditorMode": 0,
            "metricName": "RDSTotalMemory",
            "metricQueryType": 0,
            "namespace": "CustomRDSMetrics",
            "period": "",
            "queryMode": "Metrics",
            "refId": "D",
            "region": "ap-northeast-1",
            "sqlExpression": "",
            "statistic": "Average",
            "datasource": {
                "uid": "9z4hi1Onz",
                "type": "cloudwatch"
            },
            "datasourceId": 8,
            "intervalMs": 30000,
            "maxDataPoints": 536
        },
        {
            "datasource": {
                "type": "__expr__",
                "uid": "__expr__"
            },
            "expression": "($A + $C)/$D",
            "hide": false,
            "refId": "E",
            "type": "math"
        }
    ]
}
`
	redisQueries = `
{
    "queries": [
        {
            "intervalMs": 30000,
            "maxDataPoints": 536,
            "alias": "{{CacheClusterId}}",
            "dimensions": {
                "CacheClusterId": [
                    {{env}}
                ]
            },
            "expression": "",
            "id": "",
            "matchExact": true,
            "metricEditorMode": 0,
            "metricName": "CPUUtilization",
            "metricQueryType": 0,
            "namespace": "AWS/ElastiCache",
            "period": "",
            "queryType": "randomWalk",
            "refId": "A",
            "region": "ap-northeast-1",
            "statistic": "Average",
            "datasource": {
                "type": "cloudwatch",
                "uid": "9z4hi1Onz"
            },
            "sqlExpression": "",
            "type": "timeSeriesQuery"
        }
    ]
}`
	mqQueries = `
{
    "queries": [
        {
            "intervalMs": 30000,
            "maxDataPoints": 518,
            "alias": "{{Broker}}",
            "dimensions": {
                "Broker": [
                    "pd-ecom-shared-common-rabbitmq"
                ]
            },
            "expression": "",
            "id": "",
            "matchExact": true,
            "metricEditorMode": 0,
            "metricName": "MessageUnacknowledgedCount",
            "metricQueryType": 0,
            "namespace": "AWS/AmazonMQ",
            "period": "",
            "queryType": "randomWalk",
            "refId": "A",
            "region": "ap-northeast-1",
            "statistic": "Sum",
            "datasource": {
                "type": "cloudwatch",
                "uid": "9z4hi1Onz"
            },
            "sqlExpression": "",
            "type": "timeSeriesQuery"
        }
    ]
}`
	esQueries = `
{
    "queries": [
        {
            "intervalMs": 500,
            "maxDataPoints": 518,
            "alias": "{{DomainName}} - {{NodeId}}",
            "dimensions": {
                "DomainName": [
                    "pd-ecom-uc-common-es"
                ],
                "NodeId": [
                    "*"
                ]
            },
            "expression": "",
            "id": "",
            "matchExact": false,
            "metricEditorMode": 0,
            "metricName": "CPUUtilization",
            "metricQueryType": 0,
            "namespace": "AWS/ES",
            "period": "",
            "queryType": "randomWalk",
            "refId": "A",
            "region": "ap-northeast-1",
            "statistic": "Average",
            "datasource": {
                "type": "cloudwatch",
                "uid": "9z4hi1Onz"
            },
            "sqlExpression": "",
            "type": "timeSeriesQuery"
        }
    ]
}`
)

type center struct {
	pod       string
	nameSpace string
	db        string
	mq        string
	es        string
	redis     string
}

func (c *center) call() {
	rt := &Query{
		query:         queries,
		expr:          `sum(irate(http_server_requests_seconds_sum{environment="pd", app=~"({{environment}})"}[2m])) by (uri) / sum(irate(http_server_requests_seconds_count{environment="pd", app=~"({{environment}})"}[2m])) by (uri)`,
		requestId:     "36A",
		maxDataPoints: 511,
		legendFormat:  "{{uri}}",
		out: func(urlMax string, rtMax float64) {
			log.Printf("[RT]: %s %vms", urlMax, int(rtMax*1000))
		},
	}
	rt.expr = strings.Replace(rt.expr, "{{environment}}", c.pod, -1)
	rt.call()
	cpu := &Query{
		query:         queries,
		expr:          `sum(rate(container_cpu_usage_seconds_total{environment="pd", pod =~".*{{pod}}.*", namespace=~"{{environment}}", container=~".*center.*"}[2m]) )by (pod)/ sum(kube_pod_container_resource_limits{app_kubernetes_io_instance="amp-prometheus",environment="pd", resource="cpu", namespace=~"{{environment}}", container=~".*center.*"}) by (pod)`,
		requestId:     "143A",
		maxDataPoints: 511,
		legendFormat:  "{{pod}}",
		out: func(urlMax string, rtMax float64) {
			log.Printf("[CPU]: %s %%%v", urlMax, int(rtMax*100))
		},
	}
	cpu.expr = strings.Replace(cpu.expr, "{{environment}}", c.nameSpace, -1)
	cpu.expr = strings.Replace(cpu.expr, "{{pod}}", c.pod, -1)
	cpu.call()
	mem := &Query{
		query:         queries,
		expr:          `sum(container_memory_working_set_bytes{environment="pd", image!="",name=~"^k8s_.*",pod =~".*{{pod}}.*", namespace=~"{{environment}}"})by (pod) / sum(kube_pod_container_resource_limits{app_kubernetes_io_instance="amp-prometheus",environment="pd", pod =~".*{{pod}}.*", namespace=~"{{environment}}",resource="memory"})by (pod) * 100`,
		requestId:     "145A",
		maxDataPoints: 511,
		legendFormat:  "{{pod}}",
		out: func(urlMax string, rtMax float64) {
			log.Printf("[MEM]: %s %%%v", urlMax, int(rtMax))
		},
	}
	mem.expr = strings.Replace(mem.expr, "{{environment}}", c.nameSpace, -1)
	mem.expr = strings.Replace(mem.expr, "{{pod}}", c.pod, -1)
	mem.call()

	rdsCpu := &Query{
		query: strings.Replace(rdsQueries, "{{environment}}", c.db, -1),
		out: func(urlMax string, rtMax float64) {
			log.Printf("[RDS-CPU]: %s %%%v", urlMax, int(rtMax))
		},
	}
	rdsCpu.call()
	rdsMemory := &Query{
		query:        strings.Replace(rdsMemoryQueries, "{{environment}}", c.db, -1),
		resultsLabel: "E",
		out: func(urlMax string, rtMax float64) {
			float, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", rtMax*100), 64)
			log.Printf("[RDS-MEM]: %s %%%v", urlMax, float)
		},
	}
	rdsMemory.call()
	redisCpu := &Query{
		query: strings.Replace(redisQueries, "{{env}}", c.redis, -1),
		out: func(urlMax string, rtMax float64) {
			log.Printf("[REDIS-CPU]: %s %%%v", urlMax, int(rtMax))
		},
	}
	redisCpu.call()
	redisMemory := &Query{
		query: strings.Replace(redisQueries, "{{env}}", c.redis, -1),
		others: map[string]interface{}{
			"matchExact": true,
			"metricName": "DatabaseMemoryUsagePercentage",
			"queryType":  "randomWalk",
		},
		out: func(urlMax string, rtMax float64) {
			float, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", rtMax), 64)
			log.Printf("[REDIS-MEM]: %s %%%v", urlMax, float)
		},
	}
	redisMemory.call()
	mq := &Query{
		query:    mqQueries,
		duration: 5 * time.Minute,
		others: map[string]interface{}{
			"intervalMs": 500,
		},
		outFuncCustom: func(out map[string]interface{}) (string, float64) {
			resultLabel := "A"
			values := out["results"].(map[string]interface{})[resultLabel].(map[string]interface{})["frames"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["values"].([]interface{})[1].([]interface{})
			v := values[len(values)-1].(float64)
			return "", v
		},
		out: func(urlMax string, rtMax float64) {
			log.Printf("[mq-CNT]: %v", int(rtMax))
		},
	}
	mq.call()
	if len(c.es) > 0 {
		es := &Query{
			query: esQueries,
			out: func(urlMax string, rtMax float64) {
				log.Printf("[es-CPU]: %v", int(rtMax))
			},
		}
		es.call()

		esMemory := &Query{
			query: esQueries,
			others: map[string]interface{}{
				"metricName": "MasterJVMMemoryPressure",
				"statistic":  "Maximum",
			},
			out: func(urlMax string, rtMax float64) {
				log.Printf("[es-MEM]: %v", int(rtMax))
			},
		}
		esMemory.call()
	}
}

type Query struct {
	query         string
	expr          string
	requestId     string
	maxDataPoints int
	outFuncCustom func(out map[string]interface{}) (string, float64)
	out           func(urlMax string, rtMax float64)
	resultsLabel  string
	legendFormat  string
	others        map[string]interface{}
	duration      time.Duration
	now           time.Time
}

var (
	client = &http.Client{}
)

func (q *Query) outFunc(out map[string]interface{}) (string, float64) {
	resultLabel := "A"
	if len(q.resultsLabel) > 0 {
		resultLabel = q.resultsLabel
	}
	frames := out["results"].(map[string]interface{})[resultLabel].(map[string]interface{})["frames"].([]interface{})
	var rtMax float64
	var urlMax string
	for _, frame := range frames {
		frame := frame.(map[string]interface{})
		frameName, ok := frame["schema"].(map[string]interface{})["name"].(string)
		if !ok {
			frameName = ""
		}
		values := frame["data"].(map[string]interface{})["values"].([]interface{})[1]
		var valueMax float64
		for _, value := range values.([]interface{}) {
			if value == nil {
				continue
			}
			value := value.(float64)
			if value > valueMax {
				valueMax = value
			}
		}
		if valueMax > rtMax {
			rtMax = valueMax
			urlMax = frameName
		}
	}
	return urlMax, rtMax
}

func (q *Query) call() {
	var in map[string]interface{}
	err := json.Unmarshal([]byte(q.query), &in)
	if err != nil {
		fmt.Println(err)
	}
	if len(q.expr) > 0 {
		in["queries"].([]interface{})[0].(map[string]interface{})["expr"] = q.expr
	}
	if len(q.requestId) > 0 {
		in["queries"].([]interface{})[0].(map[string]interface{})["requestId"] = q.requestId
	}
	if q.maxDataPoints > 0 {
		in["queries"].([]interface{})[0].(map[string]interface{})["maxDataPoints"] = q.maxDataPoints
	}
	if len(q.legendFormat) > 0 {
		in["queries"].([]interface{})[0].(map[string]interface{})["legendFormat"] = q.legendFormat
	}

	for k, v := range q.others {
		in["queries"].([]interface{})[0].(map[string]interface{})[k] = v
	}

	now := q.now
	if now.IsZero() {
		now = time.Now()
	}

	duration := -q.duration
	if duration == 0 {
		duration = -time.Hour * 6
	}
	to := now.UnixNano()
	from := now.Add(duration).UnixNano()
	in["from"] = strconv.FormatInt(from/int64(1e6), 10)
	in["to"] = fmt.Sprintf("%v", to/int64(1e6))

	var out map[string]interface{}
	call(in, &out)

	var urlMax string
	var rtMax float64
	if q.outFuncCustom != nil {
		urlMax, rtMax = q.outFuncCustom(out)
	} else {
		urlMax, rtMax = q.outFunc(out)
	}
	q.out(urlMax, rtMax)
}

func call(in interface{}, out interface{}) {
	inStr, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
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

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
	Request(time.Unix(1595669112000, 0), time.Unix(1595669114000, 0))
}

func Request(start_time time.Time, end_time time.Time) {
	url := "http://localhost:9200/order/_search"
	method := "GET"

	payload := strings.NewReader("{\n  \"size\":0,\n  \"_source\": [\"order_id\",\"receiver_phone\",\"receiver_province\",\"delivery_date\"],\n  \"query\": {\n    \"bool\": {\n      \"must\": { \"match_all\": {} },\n      \"filter\": {\n        \"range\": {\n          \"delivery_date\": {\n            \"gte\": " + strconv.FormatInt(start_time.Unix(), 10) + ",\n            \"lte\":" + strconv.FormatInt(end_time.Unix(), 10) + "\n          }\n        }\n      }\n    }\n  },\n  \"aggs\":{\n    \"group_by_phone\":{\n      \"terms\":{\n        \"field\":\"receiver_phone.keyword\"\n      },\n      \"aggs\":{\n        \"group_by_province\":{\n          \"terms\": {\n            \"field\": \"receiver_province\"\n          },\n          \"aggs\":{\n            \"group_by_product_type\":{\n              \"terms\": {\n                \"field\": \"product_type.keyword\"\n              },\n              \"aggs\":{\n                \"group_by_order_status\":{\n                  \"terms\": {\n                    \"field\": \"order_status\"\n                  },\n                  \"aggs\":{\n                    \"sum_of_money\":{\n                      \"sum\":{\n                        \"field\": \"money_collection\"  \n                      } \n                    }\n                  }\n                }\n              }\n            }\n          }\n        }\n      }\n    }\n  }\n}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	GetProperty(body, start_time, end_time)
}

// GetProperty :
func GetProperty(bytes []byte, start_time time.Time, end_time time.Time) []Model {
	var myMap map[string]interface{}
	err := json.Unmarshal(bytes, &myMap)
	if err != nil {
		log.Println(err)
	}
	models := make([]Model, 0)
	aggregations := myMap["aggregations"]
	receiverPhoneBuckets := aggregations.(map[string]interface{})["group_by_phone"]
	receiverPhoneInterfaceArray := (receiverPhoneBuckets.(map[string]interface{})["buckets"]).([]interface{})
	for _, rPhone := range receiverPhoneInterfaceArray {
		receiverProvinceBuckets := rPhone.(map[string]interface{})["group_by_province"]
		receiverProvinceInterfaceArray := (receiverProvinceBuckets.(map[string]interface{})["buckets"]).([]interface{})
		for _, rProvince := range receiverProvinceInterfaceArray {
			productTypeBuckets := rProvince.(map[string]interface{})["group_by_product_type"]
			productTypeInterfaceArray := (productTypeBuckets.(map[string]interface{})["buckets"]).([]interface{})
			for _, pType := range productTypeInterfaceArray {
				receiver_phone := rPhone.(map[string]interface{})["key"].(string)
				recetver_province := int(rProvince.(map[string]interface{})["key"].(float64))
				product_type := pType.(map[string]interface{})["key"].(string)
				order := make(map[string]Value)
				orderStatusBuckets := pType.(map[string]interface{})["group_by_order_status"]
				orderStatusInterfaceArray := (orderStatusBuckets.(map[string]interface{})["buckets"]).([]interface{})

				for _, oStatus := range orderStatusInterfaceArray {
					order_status := oStatus.(map[string]interface{})["key"]
					money_collection := oStatus.(map[string]interface{})["sum_of_money"].(map[string]interface{})["value"]
					order_count := oStatus.(map[string]interface{})["doc_count"]
					order[fmt.Sprintf("%f", order_status.(float64))] = Value{int(order_count.(float64)), int(money_collection.(float64))}
				}
				model := NewModel(receiver_phone, order, product_type, recetver_province, start_time, end_time)
				model.Print()
				models = append(models, model)
			}
		}
	}
	return models
}

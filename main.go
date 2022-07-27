package main

import (
	"encoding/json"
	"fmt"
)

type Miner struct {
	MinerID    string
	City       string
	CountyCode string
}

func main() {
	responseJson := `{
    "responseId": "ACYDBNgZW2V0gIzE0yhv16UuqRxQbyqRhrN2Pv8z-3FeXAWNwCfbkfg7Hg5OBNje2bh7rV0",
    "timestamp": "2022-07-22T21:25:44.093103Z",
    "0_name": "Magik",
    "1_minerid": "f02620",
    "1_city": "Krakow",
    "1_country_code": "PL",
    "1_do_you_have_another_minerid_to": "No",
    "3_anything_else_you_want_us_to": "Nothing"
  }`
	var responseFields map[string]interface{}
	json.Unmarshal([]byte(responseJson), &responseFields)

	for key, value := range responseFields {
		fmt.Println(key, value.(string))
	}
}

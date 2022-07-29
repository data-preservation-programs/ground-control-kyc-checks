package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Miner struct {
	MinerID    string
	City       string
	CountyCode string
}

func main() {
	responsesFile := "testdata/responses-1.json"
	responsesJsonBytes, err := os.ReadFile(responsesFile)
	if err != nil {
		log.Fatalln("Could not load responses from", responsesFile)
	}
	var responses []map[string]interface{}
	json.Unmarshal(responsesJsonBytes, &responses)

	for _, response := range responses {
		for key, value := range response {
			fmt.Println(key, value.(string))
		}
		fmt.Println()
	}
}

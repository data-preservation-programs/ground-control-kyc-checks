package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
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
	// var responses []map[string]interface{}
	var responses []map[string]string
	json.Unmarshal(responsesJsonBytes, &responses)

	for _, response := range responses {
		for key, value := range response {
			fmt.Println(key, value)
		}
		i := 1
		for {
			minerID, ok := response[fmt.Sprintf("%d_minerid", i)]
			if !ok {
				break
			}
			city := response[fmt.Sprintf("%d_city", i)]
			countrycode := response[fmt.Sprintf("%d_country_code", i)]
			fmt.Println(i, minerID, city, countrycode)
			miner := Miner{minerID, city, countrycode}
			result, err := test_miner(context.Background(), miner)
			fmt.Printf("Result %v %v\n", result, err)
			i++
		}
		fmt.Println()
	}
}

func test_miner(ctx context.Context, miner Miner) (bool, error) {
	cmd := exec.CommandContext(ctx, "go", "test", "./minpower", "-json")
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	fmt.Println(string(out))
	return true, nil
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Miner struct {
	MinerID    string
	City       string
	CountyCode string
}

type MinerCheckResult struct {
	Miner   Miner
	Success bool

	OutputLines []TestOutput
}

type ResponseResult struct {
	ResponseFields    map[string]string
	MinerCheckResults []MinerCheckResult
}

func main() {
	responsesFile := os.Args[1]
	responsesJsonBytes, err := os.ReadFile(responsesFile)
	if err != nil {
		log.Fatalln("Could not load responses from", responsesFile)
	}
	var responses []map[string]string
	json.Unmarshal(responsesJsonBytes, &responses)

	var responseResults []ResponseResult

	for num, response := range responses {
		log.Printf("Response %d:\n", num)
		for key, value := range response {
			log.Println(" ", key, value)
		}
		i := 1
		var minerChecks []MinerCheckResult
		for {
			minerID, ok := response[fmt.Sprintf("%d_minerid", i)]
			if !ok {
				break
			}
			city := response[fmt.Sprintf("%d_city", i)]
			countrycode := response[fmt.Sprintf("%d_country_code", i)]
			log.Printf("Miner %d: %s - %s, %s\n", i, minerID, city, countrycode)
			miner := Miner{minerID, city, countrycode}
			success, testOutput, err := test_miner(context.Background(), miner)
			log.Printf("Result: %v\n", success)
			if err != nil {
				log.Printf("Error: %v\n", err)
			}
			minerCheck := MinerCheckResult{miner, success, testOutput}
			minerChecks = append(minerChecks, minerCheck)
			i++
		}
		responseResult := ResponseResult{response, minerChecks}
		responseResults = append(responseResults, responseResult)
		/*
			jsonData, err := json.MarshalIndent(minerChecks, "", "  ")
			if err != nil {
				log.Fatalln("Json marshal error", err)
			}
			fmt.Printf("JSON: %v\n", string(jsonData))
		*/
		log.Println()
	}
	if len(responseResults) == 0 {
		fmt.Println("[]")
	} else {
		jsonData, err := json.MarshalIndent(responseResults, "", "  ")
		if err != nil {
			log.Fatalln("Json marshal error", err)
		}
		fmt.Println(string(jsonData))
	}
}

type TestOutput struct {
	Time    string
	Action  *string `json:",omitempty"`
	Package *string `json:",omitempty"`
	Test    *string `json:",omitempty"`
	Output  *string `json:",omitempty"`
}

func test_miner(ctx context.Context, miner Miner) (bool, []TestOutput, error) {
	cmd := exec.CommandContext(ctx, "go", "test", "./minpower", "-json")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("MINER_ID=%s", miner.MinerID),
		fmt.Sprintf("CITY=%s", miner.City),
		fmt.Sprintf("COUNTY_CODE=%s", miner.CountyCode),
	)
	out, err := cmd.Output()
	lines := strings.Split(string(out), "\n")
	var outputLines []TestOutput
	for _, line := range lines {
		log.Println(line)
		outputLine := TestOutput{}
		json.Unmarshal([]byte(line), &outputLine)
		if outputLine.Time != "" {
			outputLines = append(outputLines, outputLine)
		}
	}
	if err != nil {
		return false, outputLines, err
	}
	return true, outputLines, nil
}

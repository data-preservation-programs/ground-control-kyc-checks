package testrig

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runTestFile(testFile string) []ResponseResult {
	fmt.Printf("Testing %s\n", testFile)
	file := fmt.Sprintf("testdata/%s.json", testFile)
	result, err := RunChecksForFormResponses(context.Background(), file, true)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(result)
	var results []ResponseResult
	json.Unmarshal([]byte(result), &results)
	return results
}

func TestResponses(t *testing.T) {

	if os.Getenv("GOOGLE_MAPS_API_KEY") == "skip" {
		log.Printf("Warning: Skipping tests because GOOGLE_MAPS_API_KEY set to 'skip'")
		return
	}
	os.Chdir("../..")

	results := runTestFile("responses-1-pass")
	assert.Equal(t, len(results), 1)
	for _, result := range results {
		for _, minerCheckResult := range result.MinerCheckResults {
			assert.True(t, minerCheckResult.Success)
		}
	}

	results = runTestFile("responses-2-fail-min-power")
	assert.Equal(t, len(results), 1)
	for _, result := range results {
		for _, minerCheckResult := range result.MinerCheckResults {
			assert.False(t, minerCheckResult.Success)
		}
	}

	results = runTestFile("responses-3-empty")
	assert.Empty(t, results)

	/*
		results = runTestFile("responses-4-china")
		assert.Equal(t, len(results), 1)
		for _, result := range results {
			for _, minerCheckResult := range result.MinerCheckResults {
				assert.True(t, minerCheckResult.Success)
			}
		}
	*/

}

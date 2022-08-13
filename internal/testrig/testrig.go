package testrig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Miner struct {
	MinerID     string
	City        string
	CountryCode string
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

func RunChecksForFormResponses(ctx context.Context, filename string,
	forceEpoch bool) (string, error) {
	responsesJsonBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	var responses []map[string]string
	json.Unmarshal(responsesJsonBytes, &responses)

	if len(responses) == 0 {
		return "[]", nil
	}

	var responseResults []ResponseResult

	// Download location data
	err = download_location_data(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed downloading location data: %w", err)
	}

	// Loop over responses and miners and run checks
	for num, response := range responses {
		log.Printf("Response %d:\n", num+1)
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
			countrycode := response[fmt.Sprintf("%d_country", i)]
			log.Printf("Miner %d: %s - %s, %s\n", i, minerID, city, countrycode)
			miner := Miner{minerID, city, countrycode}
			success, testOutput, err := test_miner(ctx, miner, forceEpoch)
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
		log.Println()
	}
	jsonData, err := json.MarshalIndent(responseResults, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON marshal error: %w", err)
	}
	return string(jsonData), nil
}

func download_location_data(ctx context.Context) error {
	downloadsDir := "downloads"
	if _, err := os.Stat(downloadsDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(downloadsDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	urls := []string{
		"https://multiaddrs-ips.feeds.provider.quest/multiaddrs-ips-latest.json",
		"https://geoip.feeds.provider.quest/ips-geolite2-latest.json",
		"https://geoip.feeds.provider.quest/ips-baidu-latest.json",
	}

	for _, dataUrl := range urls {
		u, err := url.Parse(dataUrl)
		if err != nil {
			return err
		}
		base := path.Base(u.Path)
		dest := path.Join(downloadsDir, base)

		if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
			log.Printf("Downloading %s ...\n", base)
			resp, err := http.Get(dataUrl)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			out, err := os.Create(dest)
			if err != nil {
				return err
			}
			defer out.Close()
			io.Copy(out, resp.Body)
		}
	}
	return nil
}

type TestOutput struct {
	Time    string
	Action  *string `json:",omitempty"`
	Package *string `json:",omitempty"`
	Test    *string `json:",omitempty"`
	Output  *string `json:",omitempty"`
}

func test_miner(ctx context.Context, miner Miner, forceEpoch bool) (bool,
	[]TestOutput, error) {
	cmd := exec.CommandContext(
		ctx,
		"go",
		"test",
		"./checks/minpower",
		"./checks/geoip",
		"-json",
	)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("MINER_ID=%s", miner.MinerID),
		fmt.Sprintf("CITY=%s", miner.City),
		fmt.Sprintf("COUNTRY_CODE=%s", miner.CountryCode),
		"MULTIADDRS_IPS=../../downloads/multiaddrs-ips-latest.json",
		"IPS_GEOLITE2=../../downloads/ips-geolite2-latest.json",
		"IPS_BAIDU=../../downloads/ips-baidu-latest.json",
	)
	if forceEpoch {
		cmd.Env = append(cmd.Env, "EPOCH=205500") // To match JSON files for testing
	}
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

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jimpick/sp-kyc-checks/internal/testrig"
)

func main() {
	// Load Google Form responses from data file
	responsesFile := os.Args[1]
	json, err := testrig.RunChecksForFormResponses(context.Background(), responsesFile)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(json)
}

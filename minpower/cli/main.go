package main

import (
	"context"
	"log"

	"github.com/jimpick/sp-kyc-checks/minpower"
)

func main() {
	power, err := minpower.LookupPower(context.Background(), "f02620")
	if err != nil {
		log.Fatalf("lookup power error: %s", err)
	}
	log.Println(power)
}

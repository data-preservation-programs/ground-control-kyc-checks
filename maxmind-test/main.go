package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/savaki/geoip2"
)

func main() {
	fmt.Println("MAXMIND_USER_ID", os.Getenv("MAXMIND_USER_ID"))
	fmt.Println("MAXMIND_LICENSE_KEY", os.Getenv("MAXMIND_LICENSE_KEY"))
	api := geoip2.New(os.Getenv("MAXMIND_USER_ID"), os.Getenv("MAXMIND_LICENSE_KEY"))
	resp, _ := api.Insights(context.Background(), "38.140.198.26")
	// resp, _ := api.City(context.Background(), "1.1.1.1")
	json.NewEncoder(os.Stdout).Encode(resp)
}

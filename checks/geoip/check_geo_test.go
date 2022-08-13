package geoip

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	minerID     string
	city        string
	countryCode string
	want        bool
}

func TestGeoMatchExists(t *testing.T) {
	minerID := os.Getenv("MINER_ID")
	city := os.Getenv("CITY")
	countryCode := os.Getenv("COUNTRY_CODE")

	var currentEpoch int64

	cases := make([]TestCase, 0)
	if minerID == "" {
		currentEpoch = 2055000 // To match JSON files
		cases = append(
			cases,
			TestCase{
				minerID:     "f02620",
				city:        "Warsaw",
				countryCode: "PL",
				want:        true,
			},
			TestCase{
				minerID:     "f02620",
				city:        "Toronto",
				countryCode: "CA",
				want:        false,
			},
			/*
				TestCase{ // No GeoLite2 city
					minerID:     "f01660837",
					city:        "Las Vegas",
					countryCode: "US",
					want:        true,
				},
			*/
		)
		if os.Getenv("GOOGLE_MAPS_API_KEY") == "skip" {
			log.Println("Warning: Skipping tests as GOOGLE_MAP_API_KEY set to 'skip'")
		} else {
			cases = append(
				cases,
				TestCase{ // Distance match, 500km
					minerID:     "f01558688",
					city:        "Montreal",
					countryCode: "CA",
					want:        true,
				},
				TestCase{ // More than 500 km
					minerID:     "f01558688",
					city:        "Toronto",
					countryCode: "CA",
					want:        false,
				},
				TestCase{ // China - City Name match
					minerID:     "f01012",
					city:        "Hangzhou",
					countryCode: "CN",
					want:        true,
				},
				TestCase{ // China - Distance match
					minerID:     "f01012",
					city:        "Jiaxing",
					countryCode: "CN",
					want:        true,
				},
			)
		}
	} else {
		var err error
		currentEpoch, err = strconv.ParseInt(os.Getenv("EPOCH"), 10, 64)
		if currentEpoch == 0 || err != nil {
			currentEpoch, err = GetCurrentEpoch(context.Background())
			if err != nil {
				log.Fatalf("Error getting current epoch: %v\n", err)
			}
		}

		cases = append(cases, TestCase{minerID, city, countryCode, true})
	}

	geodata, err := LoadGeoData()
	assert.Nil(t, err)

	// FIXME: Filter by date

	for _, c := range cases {
		ok := GeoMatchExists(geodata, currentEpoch, c.minerID, c.city, c.countryCode)
		assert.Equal(t, c.want, ok)
	}
}

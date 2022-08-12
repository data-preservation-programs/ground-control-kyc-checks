package geoip

import (
	"os"
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

	cases := make([]TestCase, 0)
	if minerID == "" {
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
				TestCase{
					minerID:     "f0478563",
					city:        "Hangzhou",
					countryCode: "CN",
					want:        true,
				},
			*/
		)
	} else {
		cases = append(cases, TestCase{minerID, city, countryCode, true})
	}

	geodata, err := LoadGeoData()
	assert.Nil(t, err)

	// FIXME: Filter by date

	for _, c := range cases {
		ok := GeoMatchExists(geodata, c.minerID, c.city, c.countryCode)
		assert.Equal(t, c.want, ok)
	}
}

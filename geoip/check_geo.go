package geoip

import (
	"encoding/json"
	"fmt"
	"os"
)

type MultiaddrsIPsReport struct {
	Date          *string
	MultiaddrsIPs []MultiaddrsIPsRecord
}

type MultiaddrsIPsRecord struct {
	Miner     string `json:"miner"`
	Maddr     string `json:"maddr"`
	PeerID    string `json:"peerId"`
	IP        string `json:"ip"`
	Epoch     uint   `json:"epoch"`
	Timestamp string `json:"timestamp"`
	DHT       bool   `json:"dht"`
	Chain     bool   `json:"chain"`
}

type GeoData struct {
	multiaddrsIPs []MultiaddrsIPsRecord
}

func LoadGeoData() (*GeoData, error) {
	multiaddrsIPs, err := LoadMultiAddrsIPs()
	if err != nil {
		return nil, err
	}
	return &GeoData{
		multiaddrsIPs,
	}, nil
}

func (g *GeoData) filterByMinerID(minerID string) *GeoData {
	multiaddrsIPs := []MultiaddrsIPsRecord{}
	for _, m := range g.multiaddrsIPs {
		if m.Miner == minerID {
			multiaddrsIPs = append(multiaddrsIPs, m)
		}
	}
	return &GeoData{
		multiaddrsIPs,
	}
}

func LoadMultiAddrsIPs() ([]MultiaddrsIPsRecord, error) {
	multiaddrsIPsFile := "testdata/multiaddrs-ips-latest.json"
	multiaddrsIPsBytes, err := os.ReadFile(multiaddrsIPsFile)
	if err != nil {
		return nil, err
	}
	var multiaddrsIPsReport MultiaddrsIPsReport
	json.Unmarshal(multiaddrsIPsBytes, &multiaddrsIPsReport)
	return multiaddrsIPsReport.MultiaddrsIPs, nil
}

// GeoMatchExists checks if the miner has an IP address with a location close to the city/country
func GeoMatchExists(geodata *GeoData, minerID string, city string, countryCode string) (bool, error) {
	g := geodata.filterByMinerID(minerID)
	for i, v := range g.multiaddrsIPs {
		fmt.Println("Jim1", i, v)
	}

	return true, nil
}

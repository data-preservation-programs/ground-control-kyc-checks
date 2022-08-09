package geoip

import (
	"fmt"
	"log"
)

type GeoData struct {
	MultiaddrsIPs []MultiaddrsIPsRecord
	IPsGeolite2   map[string]IPsGeolite2Record
}

func LoadGeoData() (*GeoData, error) {
	multiaddrsIPs, err := LoadMultiAddrsIPs()
	if err != nil {
		return nil, err
	}

	ipsGeolite2, err := LoadIPsGeolite2()
	if err != nil {
		return nil, err
	}

	return &GeoData{
		multiaddrsIPs,
		ipsGeolite2,
	}, nil
}

func (g *GeoData) filterByMinerID(minerID string) *GeoData {
	multiaddrsIPs := []MultiaddrsIPsRecord{}
	ipsGeoLite2 := make(map[string]IPsGeolite2Record)
	for _, m := range g.MultiaddrsIPs {
		if m.Miner == minerID {
			multiaddrsIPs = append(multiaddrsIPs, m)
			if r, ok := g.IPsGeolite2[m.IP]; ok {
				ipsGeoLite2[m.IP] = r
			}
		}
	}

	return &GeoData{
		multiaddrsIPs,
		ipsGeoLite2,
	}
}

// GeoMatchExists checks if the miner has an IP address with a location close to the city/country
func GeoMatchExists(geodata *GeoData, minerID string, city string, countryCode string) bool {
	g := geodata.filterByMinerID(minerID)
	for i, v := range g.MultiaddrsIPs {
		fmt.Println("Jim1", i, v)
	}
	fmt.Println("Jim2", g.IPsGeolite2)

	var match_found bool
	for ip, geolite2 := range g.IPsGeolite2 {
		if geolite2.Country == countryCode {
			log.Printf("Matching country for %s (%s, %s) found, IP: %s (%s)\n",
				minerID, city, countryCode, ip, geolite2.Country)
			match_found = true
		}
	}

	return match_found
}

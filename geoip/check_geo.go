package geoip

import (
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

	var match_found bool = false
	for ip, geolite2 := range g.IPsGeolite2 {
		// Match country
		if geolite2.Country != countryCode {
			continue
		}
		log.Printf("Matching country for %s (%s) found, IP: %s\n",
			minerID, countryCode, ip)

		// Try to match city
		if geolite2.City == city {
			log.Printf("Matching city for %s (%s) found, IP: %s\n",
				minerID, city, ip)
			match_found = true
			continue
		}

		// Distance based matching
	}

	return match_found
}

package geoip

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jftuga/geodist"
	"googlemaps.github.io/maps"
)

const MAX_DISTANCE = 500

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

func geocodeAddress(address string) ([]geodist.Coord, error) {
	key := os.Getenv("GOOGLE_MAPS_API_KEY")
	if key == "" {
		log.Println("Warning! Missing GOOGLE_MAPS_API_KEY")
		return []geodist.Coord{}, nil
	}
	c, err := maps.NewClient(maps.WithAPIKey(key))
	if err != nil {
		return []geodist.Coord{}, err
	}

	r := &maps.GeocodingRequest{
		Address: address,
	}
	resp, err := c.Geocode(context.Background(), r)
	if err != nil {
		return []geodist.Coord{}, err
	}

	var locations []geodist.Coord
	for _, r := range resp {
		location := geodist.Coord{
			Lat: r.Geometry.Location.Lat,
			Lon: r.Geometry.Location.Lng,
		}
		locations = append(locations, location)
	}

	return locations, nil
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
		log.Printf("No city match for %s (%s != GeoLite2:%s), IP: %s\n",
			minerID, city, geolite2.City, ip)

		locations, err := geocodeAddress(fmt.Sprintf("%s, %s", city, countryCode))
		if err != nil {
			log.Fatalf("Geocode error: %s", err)
		}

		l := geolite2.Geolite2["location"].(map[string]interface{})
		geolite2Location := geodist.Coord{
			Lat: l["latitude"].(float64),
			Lon: l["longitude"].(float64),
		}
		// log.Printf("Geolite2: %v\n", geolite2.Geolite2["location"])
		log.Printf("Geolite2 Lat/Lng: %v for IP %s\n", geolite2Location, ip)
		// Distance based matching
		for i, location := range locations {
			log.Printf("Geocoded via Google %s, %s #%d Lat/Long %v", city,
				countryCode, i+1, location)
			_, distance, err := geodist.VincentyDistance(location, geolite2Location)
			if err != nil {
				log.Println("Unable to compute Vincenty Distance.")
				continue
			} else {
				if distance <= MAX_DISTANCE {
					log.Printf("Match found! Distance %f km\n", distance)
					match_found = true
					continue
				}
				log.Printf("No match, distance %f km > %d km\n", distance, MAX_DISTANCE)
			}
		}
	}

	return match_found
}

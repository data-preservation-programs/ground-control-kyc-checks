package geoip

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jftuga/geodist"
	"googlemaps.github.io/maps"
)

const MAX_DISTANCE = 500

type GeoData struct {
	MultiaddrsIPs []MultiaddrsIPsRecord
	IPsGeolite2   map[string]IPsGeolite2Record
	IPsBaidu      map[string]IPsBaiduRecord
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

	ipsBaidu, err := LoadIPsBaidu()
	if err != nil {
		return nil, err
	}

	return &GeoData{
		multiaddrsIPs,
		ipsGeolite2,
		ipsBaidu,
	}, nil
}

func geocodeAddress(address string) ([]geodist.Coord, error) {
	key := os.Getenv("GOOGLE_MAPS_API_KEY")
	if key == "" {
		log.Fatalf("Missing GOOGLE_MAPS_API_KEY")
	}
	if key == "skip" {
		log.Println("Warning: GOOGLE_MAPS_API_KEY set to 'skip'")
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

func (g *GeoData) filterByMinerID(minerID string, currentEpoch int64) *GeoData {
	minEpoch := currentEpoch - 14*24*60*2 // 2 weeks
	multiaddrsIPs := []MultiaddrsIPsRecord{}
	ipsGeoLite2 := make(map[string]IPsGeolite2Record)
	ipsBaidu := make(map[string]IPsBaiduRecord)
	for _, m := range g.MultiaddrsIPs {
		if m.Miner == minerID {
			if int64(m.Epoch) < minEpoch {
				log.Printf("IP address %s rejected, too old: %d < %d\n",
					m.IP, m.Epoch, minEpoch)
			} else {
				multiaddrsIPs = append(multiaddrsIPs, m)
				if r, ok := g.IPsGeolite2[m.IP]; ok {
					ipsGeoLite2[m.IP] = r
				}
				if r, ok := g.IPsBaidu[m.IP]; ok {
					ipsBaidu[m.IP] = r
				}
			}
		}
	}

	return &GeoData{
		multiaddrsIPs,
		ipsGeoLite2,
		ipsBaidu,
	}
}

// GeoMatchExists checks if the miner has an IP address with a location close to the city/country
func GeoMatchExists(geodata *GeoData, currentEpoch int64, minerID string,
	city string, countryCode string) bool {

	log.Printf("Searching for geo matches for %s (%s, %s)", minerID,
		city, countryCode)
	g := geodata.filterByMinerID(minerID, currentEpoch)

	var match_found bool = false
	if countryCode != "CN" {
		for ip, geolite2 := range g.IPsGeolite2 {
			// Match country
			if geolite2.Country != countryCode {
				log.Printf("No country match for %s (%s != GeoLite2:%s), IP: %s\n",
					minerID, countryCode, geolite2.Country, ip)
				continue
			}
			log.Printf("Matching country for %s (%s) found, IP: %s\n",
				minerID, countryCode, ip)

			// Try to match city
			if geolite2.City == city {
				log.Printf("Match found! %s matches city name (%s), IP: %s\n",
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
	} else {
		for ip, baidu := range g.IPsBaidu {
			// Try to match city
			if baidu.City == city {
				log.Printf("Match found! %s matches city name (%s), IP: %s\n",
					minerID, city, ip)
				match_found = true
				continue
			}
			log.Printf("No city match for %s (%s != Baidu:%s), IP: %s\n",
				minerID, city, baidu.City, ip)

			locations, err := geocodeAddress(fmt.Sprintf("%s, %s", city, countryCode))
			if err != nil {
				log.Fatalf("Geocode error: %s", err)
			}

			baiduContent := baidu.Baidu["content"].(map[string]interface{})
			baiduPoint := baiduContent["point"].(map[string]interface{})
			lon, err := strconv.ParseFloat(baiduPoint["x"].(string), 64)
			if err != nil {
				log.Println("Error parsing baidu longitude (x)", err)
				continue
			}
			lat, err := strconv.ParseFloat(baiduPoint["y"].(string), 64)
			if err != nil {
				log.Println("Error parsing baidu latitude (y)", err)
				continue
			}
			baiduLocation := geodist.Coord{
				Lat: lat,
				Lon: lon,
			}
			log.Printf("Baidu Lat/Lng: %v for IP %s\n", baiduLocation, ip)
			// Distance based matching
			for i, location := range locations {
				log.Printf("Geocoded via Google %s, %s #%d Lat/Long %v", city,
					countryCode, i+1, location)
				_, distance, err := geodist.VincentyDistance(location, baiduLocation)
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
	}

	if !match_found {
		log.Println("No match found.")
	}
	return match_found
}

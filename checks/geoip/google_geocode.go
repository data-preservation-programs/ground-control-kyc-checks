package geoip

import (
	"context"
	"log"
	"os"

	"github.com/jftuga/geodist"
	"googlemaps.github.io/maps"
)

func getGeocodeClient() (*maps.Client, error) {
	key := os.Getenv("GOOGLE_MAPS_API_KEY")
	if key == "" {
		log.Fatalf("Missing GOOGLE_MAPS_API_KEY")
	}
	if key == "skip" {
		log.Println("Warning: GOOGLE_MAPS_API_KEY set to 'skip'")
		return nil, nil
	}
	return maps.NewClient(maps.WithAPIKey(key))
}

func geocodeAddress(ctx context.Context, client *maps.Client,
	address string) ([]geodist.Coord, []maps.GeocodingResult, error) {
	if client == nil {
		return []geodist.Coord{}, []maps.GeocodingResult{}, nil
	}

	r := &maps.GeocodingRequest{
		Address: address,
	}
	resp, err := client.Geocode(ctx, r)
	if err != nil {
		return []geodist.Coord{}, resp, err
	}

	var locations []geodist.Coord
	for _, r := range resp {
		location := geodist.Coord{
			Lat: r.Geometry.Location.Lat,
			Lon: r.Geometry.Location.Lng,
		}
		locations = append(locations, location)
	}

	return locations, resp, nil
}

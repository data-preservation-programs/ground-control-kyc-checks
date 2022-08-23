sp-kyc-checks
===

Implements "Know Your Customer" checks for Filecoin Storage Providers.

It is implemented as a GitHub Action.

It takes a JSON file of submitted data with claims about particular
Miner IDs and their locations (City, Country).

## Implemented Checks

### minpower

Uses the Lotus API to check that the miner ID has a minimum of 10 TiB of Quality
Adjusted Power on the network.

### geoip

Uses a combination of data sources to determine IP addresses for the Miner ID,
and checks to see if they match the claimed location.

## Inputs

### `google-form-responses`

**Required** The name of the file with the JSON input data.

### `google-maps-api-key`

**Required** The API key used with the Google Maps API for Geocoding submitted locations.

### `maxmind-user-id`

**Required** The user ID for MaxMind's GeoIP2 service.

### `maxmind-license-key`

**Required** The license key for MaxMind's GeoIP2 service.

## Example usage

Used from a GitHub Actions workflow file:

```
      - name: Run checks against responses
        uses: jimpick/sp-kyc-checks@v1.2
        with:
          google-form-responses: /github/workspace/output/google-form-responses.json
          google-maps-api-key: ${{ secrets.GOOGLE_MAPS_API_KEY }}
          maxmind-user-id: ${{ secrets.MAXMIND_USER_ID }}
          maxmind-license-key: ${{ secrets.MAXMIND_LICENSE_KEY }}
```      

## License

MIT/Apache-2 (Permissive License Stack)

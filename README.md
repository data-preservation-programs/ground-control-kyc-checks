sp-kyc-checks
===

Desired checks:

* The adjusted power of the SP is > 10 TiB
* The Peer ID associated with the miner has a valid multi address
  * The provider does not have a “delegate” in [this json](https://geoip.feeds.provider.quest/synthetic-country-state-province-locations-latest.json)
* The ISP is a non-cloud provider
* The IP address is public
* The IP address is a likely static IP
* The location of the IP address is within a 50 km radius of the city self-reported by the SP in the application
  * If the SP is using a national service that cannot be pinned to a specific city, a PL human’s discretion is needed here

## License

MIT/Apache-2 (Permissive License Stack)

package main

import "net/http"

func NewRouter(client *OpenF1Client) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/api/car-data/raw", carDataRawHandler(client))
	mux.Handle("/api/car-data/speed-summary", carDataSpeedSummaryHandler(client))
	mux.Handle("/api/car-data/gear-usage", carDataGearUsageHandler(client))

	mux.Handle("/api/drivers", driversHandler(client))
	mux.Handle("/api/drivers/speed-summary", driverSpeedSummaryHandler(client))

	return mux
}

package main

import "net/http"

func NewRouter(client *OpenF1Client) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/api/car-data/raw", carDataRawHandler(client))
	mux.Handle("/api/car-data/speed-summary", carDataSpeedSummaryHandler(client))
	mux.Handle("/api/car-data/gear-usage", carDataGearUsageHandler(client))

	mux.Handle("/api/drivers", driversHandler(client))
	mux.Handle("/api/drivers/speed-summary", driverSpeedSummaryHandler(client))
	mux.Handle("/api/meetings", meetingsHandler(client))
	mux.Handle("/api/sessions", sessionsHandler(client))
	mux.Handle("/api/session-results", sessionResultsHandler(client))

	return withCORS(mux)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

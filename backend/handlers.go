package main

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func parseCarDataFilter(r *http.Request) (CarDataFilter, error) {
	q := r.URL.Query()
	var f CarDataFilter

	if dnStr := q.Get("driver_number"); dnStr != "" {
		if dn, err := strconv.Atoi(dnStr); err == nil {
			f.DriverNumber = &dn
		} else {
			return f, err
		}
	}

	if skStr := q.Get("session_key"); skStr != "" {
		if sk, err := strconv.Atoi(skStr); err == nil {
			f.SessionKey = &sk
		} else {
			return f, err
		}
	}

	if msStr := q.Get("min_speed"); msStr != "" {
		if ms, err := strconv.ParseFloat(msStr, 64); err == nil {
			f.MinSpeed = &ms
		} else {
			return f, err
		}
	}

	return f, nil
}

// GET /api/car-data/raw?driver_number=55&session_key=9159&min_speed=315
// returns raw telemetry array (directly usable in charts)
func carDataRawHandler(client *OpenF1Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, err := parseCarDataFilter(r)
		if err != nil {
			http.Error(w, "invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		data, err := client.GetCarData(r.Context(), filter)
		if err != nil {
			http.Error(w, "failed to fetch car data: "+err.Error(), http.StatusBadGateway)
			return
		}

		writeJSON(w, http.StatusOK, data)
	}
}

// GET /api/car-data/speed-summary?driver_number=55&session_key=9159&min_speed=315
// returns min/max/avg for speed
func carDataSpeedSummaryHandler(client *OpenF1Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, err := parseCarDataFilter(r)
		if err != nil {
			http.Error(w, "invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		data, err := client.GetCarData(r.Context(), filter)
		if err != nil {
			http.Error(w, "failed to fetch car data: "+err.Error(), http.StatusBadGateway)
			return
		}

		if len(data) == 0 {
			writeJSON(w, http.StatusOK, SpeedSummary{
				Samples: 0,
				Min:     0,
				Max:     0,
				Avg:     0,
			})
			return
		}

		min := math.MaxFloat64
		max := -math.MaxFloat64
		sum := 0.0

		for _, d := range data {
			if d.Speed < min {
				min = d.Speed
			}
			if d.Speed > max {
				max = d.Speed
			}
			sum += d.Speed
		}

		summary := SpeedSummary{
			Samples: len(data),
			Min:     min,
			Max:     max,
			Avg:     sum / float64(len(data)),
		}

		writeJSON(w, http.StatusOK, summary)
	}
}

// GET /api/car-data/gear-usage?driver_number=55&session_key=9159&min_speed=100
// returns counts per gear, e.g. [{gear:1,count:10},...]
func carDataGearUsageHandler(client *OpenF1Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, err := parseCarDataFilter(r)
		if err != nil {
			http.Error(w, "invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		data, err := client.GetCarData(r.Context(), filter)
		if err != nil {
			http.Error(w, "failed to fetch car data: "+err.Error(), http.StatusBadGateway)
			return
		}

		counts := map[int]int{}
		for _, d := range data {
			counts[d.Gear]++
		}

		var result []GearUsageItem
		for gear, count := range counts {
			result = append(result, GearUsageItem{
				Gear:  gear,
				Count: count,
			})
		}

		writeJSON(w, http.StatusOK, result)
	}
}

func parseDriverFilter(r *http.Request) (DriverFilter, error) {
	q := r.URL.Query()
	var f DriverFilter

	if skStr := q.Get("session_key"); skStr != "" {
		if sk, err := strconv.Atoi(skStr); err == nil {
			f.SessionKey = &sk
		} else {
			return f, err
		}
	}

	if mkStr := q.Get("meeting_key"); mkStr != "" {
		if mk, err := strconv.Atoi(mkStr); err == nil {
			f.MeetingKey = &mk
		} else {
			return f, err
		}
	}

	if dnStr := q.Get("driver_number"); dnStr != "" {
		if dn, err := strconv.Atoi(dnStr); err == nil {
			f.DriverNumber = &dn
		} else {
			return f, err
		}
	}

	return f, nil
}

func driversHandler(client *OpenF1Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, err := parseDriverFilter(r)
		if err != nil {
			http.Error(w, "invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		drivers, err := client.GetDrivers(r.Context(), filter)
		if err != nil {
			http.Error(w, "failed to fetch drivers: "+err.Error(), http.StatusBadGateway)
			return
		}

		writeJSON(w, http.StatusOK, drivers)
	}
}

func parseMeetingFilter(r *http.Request) (MeetingFilter, error) {
	q := r.URL.Query()
	var f MeetingFilter

	if mkStr := q.Get("meeting_key"); mkStr != "" {
		if mk, err := strconv.Atoi(mkStr); err == nil {
			f.MeetingKey = &mk
		} else {
			return f, err
		}
	}

	if yearStr := q.Get("year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			f.Year = &year
		} else {
			return f, err
		}
	}

	if name := q.Get("meeting_name"); name != "" {
		f.MeetingName = &name
	}
	if location := q.Get("location"); location != "" {
		f.Location = &location
	}
	if country := q.Get("country_name"); country != "" {
		f.CountryName = &country
	}

	return f, nil
}

func meetingsHandler(client *OpenF1Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, err := parseMeetingFilter(r)
		if err != nil {
			http.Error(w, "invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		meetings, err := client.GetMeetings(r.Context(), filter)
		if err != nil {
			http.Error(w, "failed to fetch meetings: "+err.Error(), http.StatusBadGateway)
			return
		}

		writeJSON(w, http.StatusOK, meetings)
	}
}

func parseSessionFilter(r *http.Request) (SessionFilter, error) {
	q := r.URL.Query()
	var f SessionFilter

	if mkStr := q.Get("meeting_key"); mkStr != "" {
		if mk, err := strconv.Atoi(mkStr); err == nil {
			f.MeetingKey = &mk
		} else {
			return f, err
		}
	}

	if skStr := q.Get("session_key"); skStr != "" {
		if sk, err := strconv.Atoi(skStr); err == nil {
			f.SessionKey = &sk
		} else {
			return f, err
		}
	}

	if sessionType := q.Get("session_type"); sessionType != "" {
		f.SessionType = &sessionType
	}
	if sessionName := q.Get("session_name"); sessionName != "" {
		f.SessionName = &sessionName
	}

	return f, nil
}

func sessionsHandler(client *OpenF1Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, err := parseSessionFilter(r)
		if err != nil {
			http.Error(w, "invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		sessions, err := client.GetSessions(r.Context(), filter)
		if err != nil {
			http.Error(w, "failed to fetch sessions: "+err.Error(), http.StatusBadGateway)
			return
		}

		writeJSON(w, http.StatusOK, sessions)
	}
}

func parseSessionResultFilter(r *http.Request) (SessionResultFilter, error) {
	q := r.URL.Query()
	var f SessionResultFilter

	if skStr := q.Get("session_key"); skStr != "" {
		if sk, err := strconv.Atoi(skStr); err == nil {
			f.SessionKey = &sk
		} else {
			return f, err
		}
	}

	if mpStr := q.Get("max_position"); mpStr != "" {
		if mp, err := strconv.Atoi(mpStr); err == nil {
			f.MaxPosition = &mp
		} else {
			return f, err
		}
	}

	return f, nil
}

func sessionResultsHandler(client *OpenF1Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, err := parseSessionResultFilter(r)
		if err != nil {
			http.Error(w, "invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		if filter.SessionKey == nil {
			http.Error(w, "session_key is required", http.StatusBadRequest)
			return
		}

		results, err := client.GetSessionResults(r.Context(), filter)
		if err != nil {
			http.Error(w, "failed to fetch session results: "+err.Error(), http.StatusBadGateway)
			return
		}

		writeJSON(w, http.StatusOK, results)
	}
}

// GET /api/drivers/speed-summary?driver_number=55&session_key=9159&min_speed=315
func driverSpeedSummaryHandler(client *OpenF1Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// reuse car data filter parsing (it already has driver_number, session_key, min_speed)
		filter, err := parseCarDataFilter(r)
		if err != nil {
			http.Error(w, "invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		// enforce driver_number + session_key for this endpoint
		if filter.DriverNumber == nil || filter.SessionKey == nil {
			http.Error(w, "driver_number and session_key are required", http.StatusBadRequest)
			return
		}

		data, err := client.GetCarData(r.Context(), filter)
		if err != nil {
			http.Error(w, "failed to fetch car data: "+err.Error(), http.StatusBadGateway)
			return
		}

		if len(data) == 0 {
			writeJSON(w, http.StatusOK, SpeedSummary{
				Samples: 0,
				Min:     0,
				Max:     0,
				Avg:     0,
			})
			return
		}

		min := math.MaxFloat64
		max := -math.MaxFloat64
		sum := 0.0

		for _, d := range data {
			if d.Speed < min {
				min = d.Speed
			}
			if d.Speed > max {
				max = d.Speed
			}
			sum += d.Speed
		}

		summary := SpeedSummary{
			Samples: len(data),
			Min:     min,
			Max:     max,
			Avg:     sum / float64(len(data)),
		}

		writeJSON(w, http.StatusOK, summary)
	}
}

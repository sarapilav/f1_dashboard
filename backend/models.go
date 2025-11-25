package main

type CarData struct {
	Brake        float64 `json:"brake"`
	Date         string  `json:"date"`
	DriverNumber int     `json:"driver_number"`
	DRS          int     `json:"drs"`
	MeetingKey   int     `json:"meeting_key"`
	Gear         int     `json:"n_gear"`
	RPM          float64 `json:"rpm"`
	SessionKey   int     `json:"session_key"`
	Speed        float64 `json:"speed"`
	Throttle     float64 `json:"throttle"`
}

type CarDataFilter struct {
	DriverNumber *int
	SessionKey   *int
	MinSpeed     *float64
}

type SpeedSummary struct {
	Samples int     `json:"samples"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Avg     float64 `json:"avg"`
}

type GearUsageItem struct {
	Gear  int `json:"gear"`
	Count int `json:"count"`
}

type Driver struct {
	DriverNumber int    `json:"driver_number"`
	FullName     string `json:"full_name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	NameAcronym  string `json:"name_acronym"`

	CountryCode string `json:"country_code"`
	TeamName    string `json:"team_name"`
	TeamColour  string `json:"team_colour"`

	HeadshotURL string `json:"headshot_url"`
	MeetingKey  int    `json:"meeting_key"`
	SessionKey  int    `json:"session_key"`
}

type DriverFilter struct {
	SessionKey   *int
	MeetingKey   *int
	DriverNumber *int
}

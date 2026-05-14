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

type Meeting struct {
	MeetingKey  int    `json:"meeting_key"`
	MeetingName string `json:"meeting_name"`
	Location    string `json:"location"`
	CountryName string `json:"country_name"`
	Year        int    `json:"year"`
}

type MeetingFilter struct {
	MeetingKey  *int
	Year        *int
	MeetingName *string
	Location    *string
	CountryName *string
}

type Session struct {
	SessionKey  int    `json:"session_key"`
	MeetingKey  int    `json:"meeting_key"`
	SessionName string `json:"session_name"`
	SessionType string `json:"session_type"`
	DateStart   string `json:"date_start"`
	DateEnd     string `json:"date_end"`
}

type SessionFilter struct {
	MeetingKey  *int
	SessionKey  *int
	SessionType *string
	SessionName *string
}

type SessionResult struct {
	DNF          bool `json:"dnf"`
	DNS          bool `json:"dns"`
	DSQ          bool `json:"dsq"`
	DriverNumber int  `json:"driver_number"`
	Duration     any  `json:"duration"`
	GapToLeader  any  `json:"gap_to_leader"`
	NumberOfLaps int  `json:"number_of_laps"`
	MeetingKey   int  `json:"meeting_key"`
	Position     int  `json:"position"`
	SessionKey   int  `json:"session_key"`
}

type SessionResultFilter struct {
	SessionKey  *int
	MaxPosition *int
}

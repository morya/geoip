package model

type RspBase struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CityResult struct {
	Country  string `json:"country,omitempty"`
	Province string `json:"province,omitempty"`
	City     string `json:"city,omitempty"`
}

type RspLookup struct {
	RspBase
	Data CityResult `json:"data"`
}

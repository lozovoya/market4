package model

type Shop struct {
	ID           int    `json:"id,string,omitempty"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	WorkingHours string `json:"workingHours"`
	LON          string `json:"lon"`
	LAT          string `json:"lat"`
}

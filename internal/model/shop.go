package model

type Shop struct {
	Id           int    `json:"id,string,omitempty"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	WorkingHours string `json:"workingHours"`
	Lon          string `json:"lon"`
	Lat          string `json:"lat"`
}

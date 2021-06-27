package model

type Shop struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	WorkingHours string `json:"working_hours"`
	Lon          string `json:"lon"`
	Lat          string `json:"lat"`
}

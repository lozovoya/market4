package model

type Category struct {
	ID       int    `json:"id,string"`
	Name     string `json:"name"`
	URI_name string `json:"uri_name"`
}

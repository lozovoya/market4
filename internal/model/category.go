package model

type Category struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Uri_name string `json:"uri_name"`
}

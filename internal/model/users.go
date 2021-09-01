package model

type User struct {
	ID       int    `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

type Roles struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
}

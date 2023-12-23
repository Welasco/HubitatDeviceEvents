package model

type Device struct {
	Id    string `json:"Id"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Type  string `json:"type"`
	Room  string `json:"room"`
}

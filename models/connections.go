package models

type Connection struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Feeling     string `json:"feeling"`
	PhotoURL    string `json:"photo_url"`
	Type        string `json:"type"`
}

type Connections []*Connection

package models

type Campaign struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Campaigns struct {
	Campaigns []Campaign `json:"compaigns"`
}

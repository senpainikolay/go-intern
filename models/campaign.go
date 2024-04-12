package models

type Campaign struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Campaigns struct {
	Campaigns []Campaign `json:"campaigns"`
}

type CampaignWithDomains struct {
	Campaign
	Domains Domains
}

type CampaignsWithDomain struct {
	Campaigns []CampaignWithDomains `json:"campaigns"`
}

type Domains struct {
	Type string              `json:"type_filter"`
	Data map[string]struct{} `json:"list"`
}

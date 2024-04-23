package models

type Campaign struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Campaigns struct {
	Campaigns []Campaign `json:"campaigns"`
}

type CampaignsWithSelectedMinPriceCampaign struct {
	Campaigns
	MinPriceCampaign
}

type MinPriceCampaign struct {
	Name  string `json:"min_price_cam_name"`
	Price int    `json:"price"`
}

type CampaignWithDomains struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Domains Domains `json:"domains"`
}

type CampaignsWithDomain struct {
	Campaigns []CampaignWithDomains `json:"campaigns"`
}

type Domains struct {
	Type string              `json:"type_filter"`
	Data map[string]struct{} `json:"list"`
}

package utils

import (
	"encoding/json"
	"io"
	"math/rand"

	"github.com/senpainikolay/go-tasks/models"

	"os"
)

type RandomDomainsSubdomains struct {
	Domains    []string `json:"domains"`
	Subdomains []string `json:"subdomains"`
}

func GetRandomDomainsJsonByteArr() ([]byte, bool) {
	domains := readRandomDomainsJSON()

	var domainsPerCampaign models.CampaignDomains

	randomNumIterator := rand.Intn(len(domains.Domains))

	for i := 0; i < randomNumIterator; i++ {

		randomIndex := rand.Intn(len(domains.Domains))

		// Random goes to WhiteList or Blocked
		if rand.Intn(2)%2 == 0 {
			// added for blocked
			domainsPerCampaign.Blocked = append(domainsPerCampaign.Blocked, domains.Domains[randomIndex])
			// Random To Inlcude Subdomain or not
			if rand.Intn(2)%2 == 0 {
				domainsPerCampaign.Blocked = append(domainsPerCampaign.Blocked, domains.Subdomains[randomIndex])
			}
		} else {
			// added for white listed
			domainsPerCampaign.WhiteListed = append(domainsPerCampaign.WhiteListed, domains.Domains[randomIndex])
			// Random To Inlcude Subdomain or not
			if rand.Intn(2)%2 == 0 {
				domainsPerCampaign.WhiteListed = append(domainsPerCampaign.WhiteListed, domains.Subdomains[randomIndex])
			}
		}

		// shrink the domains slice
		domains.Domains = append(domains.Domains[:randomIndex], domains.Domains[randomIndex+1:]...)
		domains.Subdomains = append(domains.Subdomains[:randomIndex], domains.Subdomains[randomIndex+1:]...)

	}

	if len(domainsPerCampaign.Blocked) == 0 && len(domainsPerCampaign.WhiteListed) == 0 {
		return nil, false
	}

	jsonData, err := json.Marshal(domainsPerCampaign)
	if err != nil {
		panic(jsonData)
	}

	return jsonData, true

}

func readRandomDomainsJSON() RandomDomainsSubdomains {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file, err := os.Open(cwd + "/utils/domains.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var domains RandomDomainsSubdomains
	err = json.Unmarshal(jsonData, &domains)
	if err != nil {
		panic(err)
	}
	return domains
}

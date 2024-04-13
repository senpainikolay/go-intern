package utils

import (
	"encoding/json"
	"io"
	"math/rand"

	"os"

	"github.com/senpainikolay/go-tasks/models"
)

type RandomDomainsSubdomains struct {
	Domains []string `json:"domains"`
}

func GetJsonByteArr() ([]byte, bool) {
	domains := readRandomDomainsJSON()

	domainsPerCampaign := make(map[string]struct{}, len(domains.Domains))

	randomNumIterator := rand.Intn(len(domains.Domains))

	for i := 0; i < randomNumIterator; i++ {

		randomIndex := rand.Intn(len(domains.Domains))

		domainsPerCampaign[domains.Domains[randomIndex]] = struct{}{}

		// shrinking
		domains.Domains = append(domains.Domains[:randomIndex], domains.Domains[randomIndex+1:]...)

	}

	if len(domainsPerCampaign) == 0 {
		return nil, false
	}

	type_filter := "white"
	if rand.Intn(2)%2 == 0 {
		type_filter = "black"
	}

	jsonData, err := json.Marshal(models.Domains{Type: type_filter, Data: domainsPerCampaign})
	if err != nil {
		panic(err)
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

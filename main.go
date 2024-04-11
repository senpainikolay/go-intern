package main

import (
	my_sql_db "github.com/senpainikolay/go-tasks/pkg"
	"github.com/senpainikolay/go-tasks/repository"
)

func main() {
	db := my_sql_db.NewDbConnection()

	defer db.Close()

	generalRepository := repository.NewGeneralRepository(db)

	err := generalRepository.TryCreate()
	if err != nil {
		panic(err)
	}

	err = generalRepository.PopulateRandomSources()
	if err != nil {
		panic(err)
	}

	err = generalRepository.PopulateRandomCampaigns()
	if err != nil {
		panic(err)
	}

	err = generalRepository.PopulateRandomSourcesCampaigns()
	if err != nil {
		panic(err)
	}
}

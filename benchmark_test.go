package main

import (
	"strconv"
	"testing"

	controller "github.com/senpainikolay/go-tasks/http-controller"
	"github.com/senpainikolay/go-tasks/models"

	my_sql_db "github.com/senpainikolay/go-tasks/pkg"

	"github.com/senpainikolay/go-tasks/repository"
	"github.com/valyala/fasthttp"
)

func BenchmarkGetCampaignsPerMultipleSourceIDs(b *testing.B) {

	databaseTestConfig := models.DatabaseConfig{
		DbName: "nikolayinternbenchmarktestdb",
		DbUser: "root",
		DbPass: "password",
		DbHost: "localhost",
		DbPort: "3306",
	}

	db := my_sql_db.NewDbConnection(databaseTestConfig)
	defer db.Close()

	repo := repository.NewGeneralRepository(db)
	err := repo.TryCreate()
	if err != nil {
		panic(err)
	}
	err = repo.PopulateRandomDB()
	if err != nil {
		panic(err)
	}
	controller := controller.NewController(repo)

	randomSourceIDs := []int{1, 2, 3, 10, 11, 12, 50, 51, 52, 60, 61, 62, 63, 70, 75, 90, 100}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, id := range randomSourceIDs {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.SetRequestURI("/campaignsBySource?id=" + strconv.Itoa(id))
			controller.GetCampaginsPerSource(ctx)

		}
	}

}

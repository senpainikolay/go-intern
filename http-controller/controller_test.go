package controller

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/senpainikolay/go-tasks/models"
	"github.com/valyala/fasthttp"

	my_sql_db "github.com/senpainikolay/go-tasks/pkg"

	"github.com/senpainikolay/go-tasks/repository"
)

func BenchmarkGetCampaginsPerSource(b *testing.B) {
	db := newTestDBConnection()
	defer db.Close()

	repo := repository.NewGeneralRepository(db)
	err := repo.TryCreate()
	if err != nil {
		panic(err)
	}

	sourceIDinDB, cleanUpDB := insertDataInDB(repo)
	defer cleanUpDB()

	controller := NewController(repo)

	b.ResetTimer()
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/campaignsBySource?id=" + strconv.Itoa(sourceIDinDB))
	controller.GetCampaginsPerSource(ctx)

}

func TestGetCampaignsWithDomainsPerSourceIdAndFilterByType(t *testing.T) {
	db := newTestDBConnection()
	defer db.Close()

	repo := repository.NewGeneralRepository(db)
	err := repo.TryCreate()
	if err != nil {
		panic(err)
	}

	sourceIDinDB, cleanUpDB := insertDataInDB(repo)
	defer cleanUpDB()

	controller := NewController(repo)

	tcs := getTestCasesBasedOnDataForDBFunc()

	for i, tc := range tcs {
		ctx := &fasthttp.RequestCtx{}
		ctx.Request.SetRequestURI("/campaignsBySource?id=" + strconv.Itoa(sourceIDinDB) + "&domain=" + tc.Domain)
		controller.GetCampaignsWithDomainsPerSourceIdAndFilterByType(ctx)

		var res models.Campaigns

		err = json.Unmarshal(ctx.Response.Body(), &res)
		if err != nil {
			panic(err)
		}

		if len(res.Campaigns) != len(tc.TrueNameRes) {
			t.Fatalf("wrong number of entities returned. Returned: %d Expected:  %d ", len(res.Campaigns), len(tc.TrueNameRes))
		}

		for _, r := range res.Campaigns {

			if _, ok := tc.TrueNameRes[r.Name]; !ok {
				t.Fatalf("wrong companies name returned: %s for test case nr %d ", r.Name, i)
			}
		}

	}

}

type TestCase struct {
	Domain      string
	TrueNameRes map[string]struct{}
}

func getTestCasesBasedOnDataForDBFunc() []TestCase {

	return []TestCase{
		TestCase{
			Domain: "google.com",
			TrueNameRes: map[string]struct{}{
				"cam2": struct{}{},
				"cam4": struct{}{},
			},
		},
		TestCase{
			Domain: "testsub.google.com",
			TrueNameRes: map[string]struct{}{
				"cam2": struct{}{},
				"cam4": struct{}{},
			},
		},
		TestCase{
			Domain: "noninlist.com",
			TrueNameRes: map[string]struct{}{
				"cam1": struct{}{},
				"cam3": struct{}{},
			},
		},
		TestCase{
			Domain: "yandex.com",
			TrueNameRes: map[string]struct{}{
				"cam2": struct{}{},
				"cam3": struct{}{},
			},
		},
		TestCase{
			Domain: "yahoo.com",
			TrueNameRes: map[string]struct{}{
				"cam2": struct{}{},
				"cam3": struct{}{},
				"cam4": struct{}{},
			},
		},
	}
}

func insertDataInDB(repo *repository.GeneralRepository) (int, func()) {
	campaigns, sourceName := dataForDB()

	sourceId, err := repo.InsertSource(sourceName)
	if err != nil {
		panic(err)
	}
	for i, val := range campaigns.Campaigns {

		id, err := repo.InsertCampaignWithDomains(val)
		if err != nil {
			panic(err)
		}
		err = repo.InsertSourceCampaign(sourceId, id)
		if err != nil {
			panic(err)
		}
		campaigns.Campaigns[i].ID = id
	}

	return sourceId, func() {
		for _, val := range campaigns.Campaigns {
			err := repo.DeleteCampaignByID(val.ID)
			if err != nil {
				panic(err)
			}
		}

		err := repo.DeleteSourceByID(sourceId)
		if err != nil {
			panic(err)
		}

	}
}

func dataForDB() (models.CampaignsWithDomain, string) {

	return models.CampaignsWithDomain{
			Campaigns: []models.CampaignWithDomains{
				models.CampaignWithDomains{
					Name: "cam1",
					Domains: models.Domains{
						Type: "black",
						Data: map[string]struct{}{
							"google.com": struct{}{},
							"yahoo.com":  struct{}{},
							"yandex.com": struct{}{},
						},
					},
				},

				models.CampaignWithDomains{
					Name: "cam2",
					Domains: models.Domains{
						Type: "white",
						Data: map[string]struct{}{
							"google.com": struct{}{},
							"yahoo.com":  struct{}{},
							"yandex.com": struct{}{},
						},
					},
				},
				models.CampaignWithDomains{
					Name: "cam3",
					Domains: models.Domains{
						Type: "black",
						Data: map[string]struct{}{
							"google.com": struct{}{},
						},
					},
				},
				models.CampaignWithDomains{
					Name: "cam4",
					Domains: models.Domains{
						Type: "white",
						Data: map[string]struct{}{
							"google.com": struct{}{},
							"yahoo.com":  struct{}{},
						},
					},
				},
			},
		},
		"source1"
}

func newTestDBConnection() *sql.DB {
	return my_sql_db.NewDbConnection(models.DatabaseConfig{
		DbName: "nikolayinterntestdb",
		DbUser: "root",
		DbPass: "password",
		DbHost: "localhost",
		DbPort: "3306",
	})
}

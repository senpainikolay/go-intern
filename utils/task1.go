package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

var (
	rowsToGeneratePerColumn = 100
	maxCampaignsPerSource   = 10
)

func GetRandomPopulateColumnSqlString(lastId int, tableName, entityName string) string {

	var sqlStringBuilder strings.Builder

	recalculated_max_rows := rowsToGeneratePerColumn + lastId

	fmt.Fprintf(&sqlStringBuilder, "INSERT INTO %s (name) VALUES ", tableName)

	for i := lastId; i < recalculated_max_rows; i++ {
		fmt.Fprintf(&sqlStringBuilder, "('%s%v'),", entityName, i)
	}

	fmt.Fprintf(&sqlStringBuilder, "('%s%v');", entityName, recalculated_max_rows)

	return sqlStringBuilder.String()

}

func GetRandomPopulateCampaignSqlString(lastId int) string {

	var sqlStringBuilder strings.Builder

	recalculated_max_rows := rowsToGeneratePerColumn + lastId

	fmt.Fprintf(&sqlStringBuilder, "INSERT INTO campaigns (name,domains) VALUES ")

	for i := lastId; i < recalculated_max_rows; i++ {

		if jsonByteArr, ok := GetRandomDomainsJsonByteArr(); ok {
			fmt.Fprintf(&sqlStringBuilder, "('campaign%v', '%s'),", i, string(jsonByteArr))
		} else {
			fmt.Fprintf(&sqlStringBuilder, "('campaign%v', NULL ),", i)
		}

	}

	fmt.Fprintf(&sqlStringBuilder, "('campaign%v', NULL);", recalculated_max_rows)

	return sqlStringBuilder.String()

}

func GetRandomPopulateJunctionTableSqlString(minSourceID, maxSouceID, minCampaignID, maxCampaignID int) string {

	var sqlStringBuilder strings.Builder

	sqlStringBuilder.WriteString("INSERT INTO sources_campaigns (source_id,campaign_id) VALUES ")

	completeCampaignEliminated := map[int]bool{}
	for i := minSourceID; i <= maxSouceID; i++ {

		campaignsPerCurrentSource := rand.Intn(maxCampaignsPerSource)

		if campaignsPerCurrentSource <= 7 {
			completeCampaignEliminated[rand.Intn(maxSouceID)+minSourceID] = true
		}

		if campaignsPerCurrentSource == 0 {
			continue
		}

		j := 0
		generationSet := map[int]bool{}
		for j < campaignsPerCurrentSource {

			randomCompaignId := rand.Intn(maxCampaignID) + minCampaignID

			if !generationSet[randomCompaignId] && !completeCampaignEliminated[randomCompaignId] {
				fmt.Fprintf(&sqlStringBuilder, "(%v,%v),", i, randomCompaignId)
				generationSet[randomCompaignId] = true
			}
			j++
		}
	}

	sqlString := sqlStringBuilder.String()

	return sqlString[:len(sqlString)-1] + ";"

}

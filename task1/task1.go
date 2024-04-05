package task1

import (
	"fmt"
	"math/rand"
	"strings"
)

var (
	rowsToGeneratePerColumn = 100
)

func Run() {

	InitDB()

	db, truncateTablesClosure := NewDbConnection()
	defer db.Close()
	truncateTablesClosure()
	//defer truncateTablesClosure()

	_, err := db.Exec(populateColumnSqlString("sources", "source"))
	checkError(err)

	_, err = db.Exec(populateColumnSqlString("campaigns", "campaign"))
	checkError(err)

	_, err = db.Exec(populateJunctionTableSqlString())
	checkError(err)

}

func populateColumnSqlString(tableName, entityName string) string {

	var sqlStringBuilder strings.Builder

	fmt.Fprintf(&sqlStringBuilder, "INSERT INTO %s (name) VALUES ", tableName)

	for i := 1; i < rowsToGeneratePerColumn; i++ {
		fmt.Fprintf(&sqlStringBuilder, "('%s%v'),", entityName, i)
	}

	fmt.Fprintf(&sqlStringBuilder, "('%s%v');", entityName, rowsToGeneratePerColumn)

	return sqlStringBuilder.String()

}

func populateJunctionTableSqlString() string {

	var sqlStringBuilder strings.Builder

	sqlStringBuilder.WriteString("INSERT INTO sources_campaigns (source_id,campaign_id) VALUES ")

	completeCampaignEliminated := map[int]bool{}
	for sourceId := 1; sourceId <= rowsToGeneratePerColumn; sourceId++ {

		campaignsPerCurrentSource := rand.Intn(10)

		if campaignsPerCurrentSource <= 7 {
			completeCampaignEliminated[rand.Intn(rowsToGeneratePerColumn)+1] = true
		}

		if campaignsPerCurrentSource == 0 {
			continue
		}

		j := 0
		generationSet := map[int]bool{}
		for j < campaignsPerCurrentSource {

			randomCompaignId := rand.Intn(rowsToGeneratePerColumn) + 1

			if !generationSet[randomCompaignId] && !completeCampaignEliminated[randomCompaignId] {
				fmt.Fprintf(&sqlStringBuilder, "(%v,%v),", sourceId, randomCompaignId)
				generationSet[randomCompaignId] = true
			}
			j++
		}
	}

	sqlString := sqlStringBuilder.String()

	return sqlString[:len(sqlString)-1] + ";"

}

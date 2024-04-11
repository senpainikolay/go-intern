package repository

import (
	"database/sql"
	"fmt"

	"github.com/senpainikolay/go-tasks/utils"

	_ "github.com/go-sql-driver/mysql"
)

type GeneralRepository struct {
	dbClient *sql.DB
}

func NewGeneralRepository(dbClient *sql.DB) *GeneralRepository {
	return &GeneralRepository{
		dbClient: dbClient,
	}
}

func (repo *GeneralRepository) PopulateRandomSources() error {

	tableName := "sources"

	_, lastSourceId, err := repo.getIDRangePerTable(tableName)
	if err != nil {
		return err
	}

	_, err = repo.dbClient.Exec(utils.GetRandomPopulateColumnSqlString(lastSourceId+1, tableName, "source"))
	if err != nil {
		return err
	}
	return nil
}
func (repo *GeneralRepository) PopulateRandomCampaigns() error {

	tableName := "campaigns"

	_, lastCampaignId, err := repo.getIDRangePerTable(tableName)
	if err != nil {
		return err
	}

	_, err = repo.dbClient.Exec(utils.GetRandomPopulateColumnSqlString(lastCampaignId+1, tableName, "campaign"))
	if err != nil {
		return err
	}
	return nil
}

func (repo *GeneralRepository) PopulateRandomSourcesCampaigns() error {

	err := repo.truncateSoucesCampaignsTable()
	if err != nil {
		return err
	}

	firstSourceId, lastSourceId, err := repo.getIDRangePerTable("sources")
	if err != nil {
		return err
	}
	firstCampaignId, lastCampaignId, err := repo.getIDRangePerTable("campaigns")
	if err != nil {
		return err
	}

	_, err = repo.dbClient.Exec(utils.GetRandomPopulateJunctionTableSqlString(firstSourceId, lastSourceId, firstCampaignId, lastCampaignId))
	if err != nil {
		return err
	}
	return nil
}

func (repo *GeneralRepository) TryCreate() error {

	_, err := repo.dbClient.Exec(`
	    CREATE TABLE IF NOT EXISTS sources (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL);
		`)
	if err != nil {
		return err
	}

	_, err = repo.dbClient.Exec(`
	    CREATE TABLE IF NOT EXISTS campaigns (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL);
		`)
	if err != nil {
		return err
	}

	_, err = repo.dbClient.Exec(`
	    CREATE TABLE IF NOT EXISTS sources_campaigns (
		source_id INT NOT NULL, 
		campaign_id INT NOT NULL, 
		PRIMARY KEY (source_id, campaign_id),
		FOREIGN KEY (source_id) REFERENCES sources(id),
		FOREIGN KEY (campaign_id) REFERENCES campaigns(id) );
		`)
	if err != nil {
		return err
	}
	return nil
}

func (repo *GeneralRepository) getIDRangePerTable(tableName string) (int, int, error) {
	var firstID, lastID int

	err := repo.dbClient.QueryRow(fmt.Sprintf("select COALESCE(MIN(ID),0), COALESCE(MAX(ID),0) from %s", tableName)).Scan(&firstID, &lastID)
	if err != nil {
		return 0, 0, err
	}

	return firstID, lastID, nil
}

func (repo *GeneralRepository) truncateSoucesCampaignsTable() error {
	_, err := repo.dbClient.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	if err != nil {
		return err
	}
	_, err = repo.dbClient.Exec("truncate table sources_campaigns;")
	if err != nil {
		return err
	}
	_, err = repo.dbClient.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	if err != nil {
		return err
	}
	return nil
}

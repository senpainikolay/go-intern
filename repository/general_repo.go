package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/senpainikolay/go-tasks/models"
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
		name VARCHAR(100) NOT NULL,
		domains JSON
	    );`)
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

func (repo *GeneralRepository) GetCampaignsPerSourceId(id int) (models.Campaigns, error) {

	rows, err := repo.dbClient.Query(`SELECT id, name
									FROM campaigns c
									JOIN sources_campaigns sc ON sc.campaign_id = c.id
									WHERE sc.source_id = ?
								   `, id)
	if err != nil {
		return models.Campaigns{}, err
	}
	defer rows.Close()

	var campaigns models.Campaigns

	for rows.Next() {

		var campaign models.Campaign

		err := rows.Scan(&campaign.ID, &campaign.Name)
		if err != nil {
			return models.Campaigns{}, err
		}

		campaigns.Campaigns = append(campaigns.Campaigns, campaign)
	}

	if err := rows.Err(); err != nil {
		return models.Campaigns{}, err
	}

	if len(campaigns.Campaigns) == 0 {
		campaigns.Campaigns = make([]models.Campaign, 0)
		return campaigns, nil
	}

	return campaigns, nil
}
func (repo *GeneralRepository) GetCampaignsWithDomainsPerSourceIdAndFilterByType(id int, domain string) (models.Campaigns, error) {

	compaignsBySourceIdWithDomains, err := repo.getCampaignsWithDomainsPerSourceId(id)
	if err != nil {
		return models.Campaigns{}, err
	}

	res := models.Campaigns{Campaigns: make([]models.Campaign, 0)}

	for i := 0; i < len(compaignsBySourceIdWithDomains.Campaigns); i++ {

		if compaignsBySourceIdWithDomains.Campaigns[i].Domains.Type == "black" {

			if _, ok := compaignsBySourceIdWithDomains.Campaigns[i].Domains.Data[domain]; !ok {
				res.Campaigns = append(res.Campaigns, models.Campaign{ID: compaignsBySourceIdWithDomains.Campaigns[i].ID, Name: compaignsBySourceIdWithDomains.Campaigns[i].Name})

			}

		} else { // Type == "white"
			if _, ok := compaignsBySourceIdWithDomains.Campaigns[i].Domains.Data[domain]; ok {
				res.Campaigns = append(res.Campaigns, models.Campaign{ID: compaignsBySourceIdWithDomains.Campaigns[i].ID, Name: compaignsBySourceIdWithDomains.Campaigns[i].Name})
			}

		}

	}

	return res, nil
}

func (repo *GeneralRepository) getCampaignsWithDomainsPerSourceId(id int) (models.CampaignsWithDomain, error) {

	rows, err := repo.dbClient.Query(`SELECT id, name, domains
									FROM campaigns c
									JOIN sources_campaigns sc ON sc.campaign_id = c.id
									WHERE sc.source_id = ? and domains is not NULL
								   `, id)
	if err != nil {
		return models.CampaignsWithDomain{}, err
	}
	defer rows.Close()

	var campaigns models.CampaignsWithDomain

	for rows.Next() {

		var campaign models.CampaignWithDomains

		var byteJSONData []byte

		err := rows.Scan(&campaign.ID, &campaign.Name, &byteJSONData)
		if err != nil {
			return models.CampaignsWithDomain{}, err
		}
		err = json.Unmarshal(byteJSONData, &campaign.Domains)
		if err != nil {
			panic(err)
		}

		campaigns.Campaigns = append(campaigns.Campaigns, campaign)
	}

	if err := rows.Err(); err != nil {
		return models.CampaignsWithDomain{}, err
	}

	if len(campaigns.Campaigns) == 0 {
		campaigns.Campaigns = make([]models.CampaignWithDomains, 0)
		return campaigns, nil
	}

	return campaigns, nil
}

func (repo *GeneralRepository) PopulateRandomDB() error {

	err := repo.populateRandomSources()
	if err != nil {
		return err
	}

	err = repo.populateRandomCampaigns()
	if err != nil {
		return err
	}

	err = repo.populateRandomSourcesCampaigns()
	if err != nil {
		return err
	}
	return nil
}

func (repo *GeneralRepository) populateRandomSources() error {

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
func (repo *GeneralRepository) populateRandomCampaigns() error {

	tableName := "campaigns"

	_, lastCampaignId, err := repo.getIDRangePerTable(tableName)
	if err != nil {
		return err
	}

	_, err = repo.dbClient.Exec(utils.GetRandomPopulateCampaignSqlString(lastCampaignId + 1))
	if err != nil {
		return err
	}
	return nil
}

func (repo *GeneralRepository) populateRandomSourcesCampaigns() error {

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

func (repo *GeneralRepository) getIDRangePerTable(tableName string) (int, int, error) {
	var firstID, lastID int

	err := repo.dbClient.QueryRow(fmt.Sprintf("select COALESCE(MIN(ID),0), COALESCE(MAX(ID),0) from %s", tableName)).Scan(&firstID, &lastID)
	if err != nil {
		return 0, 0, err
	}

	return firstID, lastID, nil
}

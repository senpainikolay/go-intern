USE nikolayinterndb;

SELECT id,name
FROM campaigns
LEFT JOIN sources_campaigns ON sources_campaigns.campaign_id = campaigns.id
WHERE campaign_id IS NULL
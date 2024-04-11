USE nikolayinterndb;

SELECT source_id, s.name, COUNT(campaign_id) as nr_of_campaigns
FROM sources_campaigns sc
JOIN sources s ON sc.source_id = s.id
GROUP BY source_id
ORDER BY  nr_of_campaigns DESC
LIMIT 5

USE nikolayinterndb;


SELECT  source_id, name, nr_of_campaigns  FROM 
(
SELECT source_id, COUNT(campaign_id) as nr_of_campaigns
FROM sources_campaigns 
GROUP BY source_id
ORDER BY  nr_of_campaigns DESC
LIMIT 5
) as inner_select_res
JOIN sources ON sources.id = inner_select_res.source_id

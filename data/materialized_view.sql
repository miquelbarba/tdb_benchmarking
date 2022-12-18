CREATE MATERIALIZED VIEW cpu_usage_summary_minute
WITH (timescaledb.continuous) AS
SELECT host,
       time_bucket(INTERVAL '1 minute', ts) AS bucket,
       MAX(usage),
       MIN(usage)
FROM cpu_usage
GROUP BY host, bucket;

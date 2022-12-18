docker-compose -p timescale-bench_compose -f docker-compose.yml cp data/cpu_usage.csv timescale:/tmp/cpu_usage.csv
docker-compose -p timescale-bench_compose -f docker-compose.yml exec -T timescale psql -U postgres < data/cpu_usage.sql
docker-compose -p timescale-bench_compose -f docker-compose.yml exec -T timescale psql -U postgres -d homework < data/materialized_view.sql
docker-compose -p timescale-bench_compose -f docker-compose.yml exec -T timescale psql -U postgres -d homework -c "\COPY cpu_usage FROM /tmp/cpu_usage.csv CSV HEADER"

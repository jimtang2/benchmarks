# Notes

## PostgreSQL Service Configuration

- postgres `-c max_connections`
- postgres `-c shared_buffers`
- postgres `-c effective_cache_size`
- postgres `-c work_mem`
- postgres `-c maintenance_work_mem`
- postgres `-c wal_buffers`
- postgres `-c max_parallel_workers_per_gather`
- postgres `-c max_worker_processes`
- postgres `-c max_parallel_workers`
- docker `shm_size`
- docker `mem_limit`
- docker `mem_reservation`

## PostgreSQL Driver Configuration

- jackc/pgx `config.MaxConns`
- jackc/pgx `config.MinConns`
- jackc/pgx `config.MaxConnLifetime`
- jackc/pgx `config.MaxConnIdleTime`
- jackc/pgx `config.HealthCheckPeriod`
- lib/pq `db.SetMaxOpenConns`
- lib/pq `db.SetMaxIdleConns`
- lib/pq `db.SetMaxLifetime`
- lib/pq `db.SetMaxIdleTime`
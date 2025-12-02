# Benchmarks of Web Servers and SQL Drivers

This is a load test of different web servers, observation of relative performance and identification of bottlenecks.

## Description

### Web Servers

- go net/http
- fastapi
- express.js
- next.js
- spring boot + webflux
- rust actix-web

### SQL Drivers

- github.com/lib/pq
- github.com/jackc/pgx
- asyncpg
- psycopg
- node-postgres
- porsager/postgres
- postgresql r2dbc
- rust sqlx

### Details

- hardware: Apple Silicon M4 (10 cores) 32GB RAM
- runtime: Docker Desktop 4.53 (8 cores, 24GB RAM, 4GB swap)
- database engine: Postgres 17.5 (Docker)
- database max connection: 2,000
- database pooler: pgcat 0.2.5
- load tester: grafana/k6
- load test duration: 1m sustained
- load success threshold: error rate < 0.01

# Results

| Web Server             | SQL Driver          | VUs | RPS | avg-ms | min-ms | max-ms | p90-ms | p95-ms |
|------------------------|------------------------|------|-------|-------|------|------|-------|------|
| go stdlib              | jackc/pgx              |  150 |  6696 | 17.3  | 0.69 | 558  | 40.7  | 58.3 |
|                        |                        |  600 |  6163 | 92.1  | 47.9 | 1370 | 136   | 170  |
|                        |                        | 1200 |  5856 | 199   | 8.95 | 1220 | 253   | 287  |
|                        |                        | 2400 |  5521 | 429   | 43.2 | 1670 | 498   | 534  |
|                        |                        | 4800 |  6369 | 737   | 77.0 | 1830 | 816   | 847  |
|                        |                        | 9600 |  6219 | 1480  | 319  | 2900 | 1590  | 1620 |
|                        | lib/pq                 |  150 |  7095 | 16.1  | 0.90 | 378  | 36.9  | 52.4 |
|                        |                        |  600 |  6806 | 82.9  | 40.7 | 852  | 123   | 152  |
|                        |                        | 1200 |  6819 | 170   | 39.5 | 1370 | 305   | 375  |
|                        |                        | 2400 |  6676 | 353   | 42.5 | 3660 | 717   | 909  |
|                        |                        | 4800 |  6604 | 714   | 45.3 | 9160 | 1540  | 1980 |
|                        |                        | 9600 |  6371 | 1470  | 45.2 | 15530 | 3260 | 4200 |
| rust actix-web         | sqlx 0.8               |  150 |  7528 | 14.8  | 0.90 | 319  | 33.7  | 47.7 |
|                        |                        |  600 |  7045 | 85.0  | 2.67 | 965  | 117   | 145  |
|                        |                        | 1200 |  6975 | 171   | 72.5 | 1100 | 206   | 235  |
|                        |                        | 2400 |  6848 | 342   | 67.0 | 1260 | 387   | 413  |
|                        |                        | 4800 |  6955 | 676   | 67.3 | 1960 | 726   | 754  |
|                        |                        | 9600 |  6091 | 1510  | 184  | 2980 | 1630  | 1720 |
| spring boot + webflux  | postgresql r2dbc       |  150 |  6437 | 18.2  | 0.97 | 658  | 39.6  | 55.2 |
|                        |                        |  600 |  6117 | 92.8  | 18.7 | 945  | 137   | 170  |
|                        |                        | 1200 |  6103 | 190   | 11.2 | 1190 | 235   | 268  |
|                        |                        | 2400 |  6353 | 369   | 72.6 | 1250 | 418   | 448  |
|                        |                        | 4800 |  6075 | 771   | 230  | 1960 | 856   | 1040 |
|                        |                        | 9600 |  5892 | 1570  | 137  | 3230 | 1680  | 1730 |
| fastapi                | psycopg                |  150 |  3403 | 38.8  | 6.27 | 588  | 78.0  | 106  |
|                        |                        |  600 |  2726 | 214   | 4.65 | 2480 | 308   | 383  |
|                        |                        | 1200 |  2423 | 486   | 149  | 2660 | 594   | 679  |
|                        | asyncpg                |  150 |  4164 | 30.9  | 11.5 | 1380 | 35.8  | 38.7 |
|                        |                        |  600 |  3392 | 171   | 5.77 | 8170 | 183   | 209  |
| express.js             | node-postgres          |  150 |  7002 | 16.3  | 1.42 | 547  | 21.2  | 22.8 |
|                        |                        |  600 |  6520 | 86.8  | 1.62 | 1140 | 109   | 115  |
|                        | porsager/postgres      |  150 |  5940 | 20.2  | 2.93 | 641  | 23.2  | 24.6 |
|                        |                        |  600 |  5302 | 108   | 5.75 | 8100 | 124   | 131  |
| next.js                | porsager/postgres      |  150 |  2448 | 56.1  | 14.7 | 1650 | 67.1  | 72.3 |
|                        |                        |  600 |  2352 | 255   | 9.71 | 1810 | 278   | 297  |

### Observations

- pgcat does not support asyncpg
- pgcat does not support node-postgres
- pgcat does not support porsager/postgres
- pgcat does not support nextjs
- pgcat unsupported web servers are bottlenecked by the number of database connections
- driver pool implementations do not release connections well enough to significantly impact active_connections / max_connections ratio; hence pgcat or pgbouncer or similar is much needed at scale
- overall results were expected
- Go jackc/pgx dropped 7% in RPS between 150 and 9600 VUs (i.e., concurrent request rate)
- Go lib/pq dropped 10% in RPS between 150 and 9600 VUs 
- Rust Actix-web + Sqlx dropped 19% in RPS between 150 and 9600 VUs 
- Spring Boot Webflux + R2DBC dropped 8.5% in RPS between 150 and 9600 VUs
- All four average response time is 1.5s when using 9600 VUs
- express.js with node-postgres performance is on par with system languages at low concurrent request rate


## Setup

### Testing

1. run `docker compose build` to build all services (you may need to adjust Docker and Postgres resources for your hardware)
2. run `docker compose up -d postgres` to start database
3. run `docker compose up service_name` (check compose.yml for list of services)
4. configure `benchmark.js` options (vux, duration, or stages)
5. run `k6 run benchmark.js --env URL=http://localhost:8080` (for next.js, use `http://localhost:8080/api/tx`)

### Workload

The workload is for the web server to call the following SQL prepared statement with RNG values.

```sql
CREATE OR REPLACE FUNCTION public.pgbench_tx(
    p_aid   bigint,
    p_tid   bigint,
    p_bid   bigint,
    p_delta int
)
RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE pgbench_accounts
       SET abalance = abalance + p_delta
     WHERE aid = p_aid;

    UPDATE pgbench_tellers
       SET tbalance = tbalance + p_delta
     WHERE tid = p_tid;

    UPDATE pgbench_branches
       SET bbalance = bbalance + p_delta
     WHERE bid = p_bid;

    INSERT INTO pgbench_history (tid, bid, aid, delta, mtime, filler)
    VALUES (p_tid, p_bid, p_aid, p_delta, now(), repeat('x', 22));
END;
$$;
```

### SQL Schema

#### pgbench_accounts
| | |
|---------|----------------|
| aid      | integer       |
| bid      | integer       |
| abalance | integer       |
| filler   | character(84) |

#### pgbench_branches
| | |
----------|----------------|
| bid      | integer       |
| bbalance | integer       |
| filler   | character(88) |

#### pgbench_history
| | |
--------|------------------------------|
| tid    | integer                     |
| bid    | integer                     |
| aid    | integer                     |
| delta  | integer                     |
| mtime  | timestamp without time zone |
| filler | character(22)               |

#### pgbench_tellers
| | |
----------|----------------|
| tid      | integer       |
| bid      | integer       |
| tbalance | integer       |
| filler   | character(84) |

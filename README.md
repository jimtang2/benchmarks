# Benchmarks of Web Servers and SQL Drivers

The word benchmark may be a bit too scientific for this exercise, as it is too big of an endeavour to define the right criteria and system states to run benchmarks against. As such, this exercise is a reference whereby I have set up multiple web servers including basic "out-of-the-box" optimizations and load tested them, in order to create a data-driven comparison of different web servers.

## Description

### Web Servers

- go net/http
- fastapi
- express.js
- next.js
- spring boot + webflux
- rust actix-web
- rust axum
- rust hyper

### SQL Drivers

- github.com/lib/pq
- github.com/jackc/pgx
- asyncpg
- psycopg
- node-postgres
- porsager/postgres
- postgresql r2dbc
- rust sqlx
- rust tokio-postgres

### Details

- hardware: Apple Silicon M4 (10 cores) 32GB RAM
- runtime: Docker Desktop 4.53 (8 cores, 24GB RAM, 4GB swap)
- database engine: Postgres 17.5 (Docker)
- database pooler: pgcat 0.2.5
- load tester: grafana/k6
- load test duration: 3m sustained
- load success threshold: error rate < 0.01

# Top Results

| Rank | Web Server + SQL Driver                    |  VUs | RPS[K] | avg[ms] | p90[ms] |
|------|--------------------------------------------|------|--------|---------|---------|
|  1   | go net/http + jackc/pgx                    |  130 |   9.1  |      4  |      8  |
|  2   | rust actix-web + sqlx                      |  100 |   8.1  |    7.2  |     16  |
|  3   | spring boot + webflux + postgresql r2dbc   |  100 |   7.7  |      3  |      5  |
|  4   | express.js + node-postgres                 |  100 |   5.8  |      7  |      9  |
|  5   | express.js + porsager/postgres             |  100 |   5.2  |      9  |     13  |
|  6   | fastapi + asyncpg                          |  130 |   3.7  |     24  |     31  |
|  7   | fastapi + psycopg                          |  100 |   3.0  |     22  |     36  |
|  8   | fastapi + psycopg + pgcat                  |  100 |   2.8  |     35  |     51  |
|  9   | next.js + porsager/postgres                |  100 |   2.6  |     37  |     43  |
| 10   | go net/http + lib/pq w/o p. stmt           |  100 |   1.8  |     43  |     87  |


## Results (Full)

| Web Server             | SQL Driver             |  VUs | RPS[K] | min[ms] | max[ms] | avg[ms] | p90[ms] | p95[ms] |
|------------------------|------------------------|------|--------|---------|---------|---------|---------|---------|
| go std lib             | jackc/pgx              |  100 |   8.6  |   0.21  |    128  |      1  |      2  |     13  |
|                        |                        |  130 |   9.1  |   0.30  |    114  |      4  |      8  |     11  |
|                        |                        |  500 |   6.3  |   0.67  |    2.4s |     68  |    166  |    252  |
|                        |                        |   1K |   5.7  |    57   |    3.7s |    176  |    273  |    362  |
|                        | jackc/pgx + pgcat      |  100 |   7.9  |    10   |    243  |     12  |     14  |     15  |
|                        |                        |  130 |   8.4  |   0.69  |     77  |      5  |     10  |     15  |
|                        |                        |  250 |   7.0  |    10   |     .9s |     35  |     72  |     99  |
|                        |                        |  500 |   6.6  |    14   |      1s |     64  |    105  |    135  |
|                        |                        |   2K |   6.5  |   233   |    1.3s |    305  |    357  |    388  |
|                        |                        |   4K |   6.5  |    61   |    2.1s |     .6s |    .7s  |    .7s  |
|                        |                        |   6K |   6.9  |   151   |    2.3s |     .8s |     1s  |     1s  |
|                        |                        |  10K |   5.7  |   451   |    3.7s |    1.7s |   2.1s  |   2.4s  |
|                        | jackc/pgx w/o p. stmt  |  100 |   7.1  |   0.79  |    169  |      3  |      6  |      9  |
|                        | lib/pq w/o p. stmt     |  100 |   1.8  |   1.2   |    510  |     43  |     87  |    115  |
| rust actix-web         | sqlx 0.8               |  100 |   8.1  |  0.58   |   156   |   7.2   |   15.7  |  22     |
|                        |                        |  130 |   7.7  |  0.7    |   304   |   11.7  |   27    |  38     |
|                        |                        |  500 |   5.6  |   5.8   |  3.1s   |   88    |    211  |  319    |
|                        |                        |   2K |   2.0  |   1.1   |  39.7s  |   964.8 |    2.5s  | 3.9s   |
|                        | sqlx + pgcat           |  100 |   6.9  |  0.73   |   254   |   9.2   |   19.4   |   27   |
|                        |                        |  200 |   6.4  |   5.9   |   819   |   31    |    66    | 93     |
|                        |                        |   2K |   5.9  |    70   |   1.6s  |   331   |    421   | 472    |
|                        |                        |   5K |   6.4  |    96   |   2.2s  |   775   |   845    | 891    |
|                        |                        |  10K |   2.0  |   1.1   |  39.7s  |   964.8 |    2.5s  | 3.9s   |
| spring boot + webflux  | postgresql r2dbc       |  100 |   7.7  |   0.29  |     46  |      3  |      5  |      7  |
|                        |                        |  130 |   7.5  |   1.9   |     81  |      7  |     10  |     13  |
|                        |                        |  250 |   7.7  |   3.2   |     83  |     33  |     26  |     27  |
|                        |                        |  500 |   7.0  |    27   |    192  |     60  |     70  |     77  |
|                        |                        | 1000 |   6.8  |    31   |    250  |    136  |    153  |    160  |
|                     | postgresql r2dbc + pgcat | 10000 |   4.4  |    59   |    3.6s |    2.2s |    2.7s |    2.8s |
| fastapi                | asyncpg                |  100 |   3.7  |    11   |    156  |     16  |     20  |     23  |
|                        |                        |  130 |   3.7  |     2   |    356  |     24  |     31  |     35  |
|                        |                        |  200 |   3.2  |    16   |     .8s |     51  |     66  |     75  |
|                        | psycopg                |  100 |   3.0  |    2.3  |    272  |     22  |     36  |     45  |
|                        |                        |  200 |   2.7  |    14   |    1.6s |     62  |    141  |    191  |
|                        | psycopg + pgcat        |  100 |   2.8  |    14   |    328  |     35  |     51  |     62  |
|                        |                        |  200 |   2.4  |    15   |    1.5s |     69  |    157  |    203  |
| express.js             | node-postgres          |  100 |   5.8  |   0.83  |    291  |      7  |      9  |     10  |
|                        |                        |  250 |   4.8  |   4.0   |    1.2s |     51  |     61  |     66  |
|                        | porsager/postgres      |  100 |   5.2  |   0.43  |    329  |      9  |     13  |     15  |
|                        |                        |  250 |   4.9  |   4.3   |    1.1s |     41  |     48  |     51  |
| next.js                | porsager/postgres      |  100 |   2.6  |   4.1   |    543  |     37  |     43  |     47  |
|                        |                        |  250 |   2.4  |   8.5   |    3.3s |     94  |    112  |    119  |
| rust axum 0.7          | sqlx                   |  |  |  |  |  |  |  |
| rust hyper 1.8.1       | tokio-postgres 0.7.15  |  |  |  |  |  |  |  |

## Setup

### Testing

1. run `docker compose build` to build all services (you may need to adjust Docker and Postgres resources for your hardware)
2. run `docker compose up -d postgres` to start database
3. run `docker compose up service_name` (check compose.yml for list of services)
4. configure `benchmark.js` options (vux, duration, or stages)
5. run `k6 run benchmark.js --env URL=http://localhost:8080` (for next.js, use `http://localhost:8080/api/tx`)

### Workload

The workload consists in performing the following SQL transactions. We quickly conclude the execution of prepared statement by database engine is most efficient. Performance difference is shown in results for Go net/http + jackc/pgx.

#### SQL

```sql
BEGIN;
  UPDATE pgbench_accounts 
    SET abalance = abalance + :delta 
    WHERE aid = :aid; 
  SELECT abalance 
    FROM pgbench_accounts 
    WHERE aid = :aid;
  UPDATE pgbench_tellers 
    SET tbalance = tbalance + :delta 
    WHERE tid = :tid; 
  UPDATE pgbench_branches 
    SET bbalance = bbalance + :delta 
    WHERE bid = :bid;
  INSERT INTO pgbench_history (tid, bid, aid, delta, mtime, filler) 
    VALUES (:tid, :bid, :aid, :delta, now(), repeat('x',22));
END;
```

### Schema

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

### Notes

- pgcat does not support asyncpg
- pgcat does not support node-postgres
- pgcat does not support porsager/postgres
- issues with axum and hyper setups

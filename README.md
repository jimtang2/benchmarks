# Comparison of Benchmarks for 7 Different Web Servers with 9 Different SQL Drivers

Benchmarks of different web servers + PostgreSQL i/o.

# Conclusion



# Results 

- Apple Silicon M4 32GB RAM
- Docker Desktop 4.53
- Postgres 17.5 (Docker)
- grafana/k6
- threshold: error rate < 0.01
- duration: 1m sustained

| Web Server                  | SQL Driver             | RPS    | VUs   | min[ms]  | max[ms]  | avg[ms]  | p90[ms]  | p95[ms]  |
|-----------------------------|------------------------|--------|-------|----------|----------|----------|----------|----------|
| Go net/http                 | lib/pq                 | 1,865  |   100 |   1.29   |   510    |    43    |    87    |    115   |
|                             | jackc/pgx              | 7,171  |   100 |   0.79   |   169    |     3    |     6    |     9    |
|                             | jackc/pgx (prep. stmt) | 8,640  |   100 |   0.21   |   128    |     1    |     2    |    13    |
|                             |                        | 9,107  |   130 |   0.30   |   114    |     4    |     8    |    11    |
|                             |                        | 6,366  |   500 |   0.67   |  2,470   |    68    |    166   |    252   |
|                     | jackc/pgx (prep. stmt + pgcat) | 7,985  |   100 |    10    |   243    |    12    |    14    |    15    |
|                             |                        | 8,412  |   130 |   0.69   |    77    |     5    |    10    |    15    |
|                             |                        | 7,013  |   250 |    10    |   882    |    35    |    72    |    99    |
|                             |                        | 6,874  |   500 |    29    |   721    |    62    |    100   |    132   |
|                             |                        | 6,905  | 2,000 |   193    |   950    |   277    |    330   |    361   |
|                             |                        | 6,550  | 4,000 |    61    |  2,050   |   591    |    675   |    715   |
|                             |                        | 6,916  | 6,000 |   151    |  2,280   |   846    |    928   |    977   |
|                             |                        | 5,723  |10,000 |   451    |  3,690   |  1,700   |   2,110  |  2,420   |
| FastAPI                     | asyncpg                |        |       |          |          |          |          |          |
| FastAPI                     | psycopg                |        |       |          |          |          |          |          |
| Express.js                  | node-postgres          |        |       |          |          |          |          |          |
| Express.js                  | porsager/postgres      |        |       |          |          |          |          |          |
| Next.js                     | node-postgres          |        |       |          |          |          |          |          |
| Next.js                     | porsager/postgres      |        |       |          |          |          |          |          |
| Spring Boot + Undertow      | —                      |        |       |          |          |          |          |          |
| Spring Boot + Tomcat        | —                      |        |       |          |          |          |          |          |
| Spring Boot + Netty         | —                      |        |       |          |          |          |          |          |

## Setup

### Workload

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

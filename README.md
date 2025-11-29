Benchmarks of different web servers + PostgreSQL i/o performed on Macbook Air M4 32GB RAM.

# Results 

| Web Server                  | SQL Driver             | min      | max       | ave      | p(90)    | p(95)    | req/s       | VUs |
|-----------------------------|------------------------|----------|-----------|----------|----------|----------|-------------|-----|
| Express.js                  | node-postgres          |          |           |          |          |          |             |     |
| Express.js                  | porsager/postgres      |          |           |          |          |          |             |     |
| FastAPI Setup               | asyncpg                |          |           |          |          |          |             |     |
| FastAPI Setup               | psycopg                |          |           |          |          |          |          |             |     |
| Go net/http                 | lib/pq                 | 1.21ms   | 131.57ms  | 6.62ms   | 12.37ms  | 16.78ms  | 5976.87   | 100 |
| Go net/http                 | jackc/pgx              | 717µs    | 115.46ms  | 3.87ms   | 6.86ms   | 9.09ms   | 7158.35   | 100 |
| Next.js                     | node-postgres          |          |           |          |          |          |             |     |
| Next.js                     | porsager/postgres      |          |           |          |          |          |             |     |
| Spring Boot (Undertow)      | —                      |          |           |          |          |          |             |     |
| Spring Boot (Tomcat)        | —                      |          |           |          |          |          |             |     |
| Spring Boot (Netty)         | —                      |          |           |          |          |          |             |     |

# Configuration

## Web Servers

1. All web servers to use 1,200 SQL pool size (10MB/connection => 12GB RAM)
2. All web servers to use prepared SQL statement

### Express.js
- run with `--experimental-worker --max-old-space-size=8192` 

### FastAPI
- `asyncpg` vs sync `psycopg`
- run with `--uvicorn --workers 1`

### Go net/http

### Next.js
- use server actions or route handlers
- `@vercel/postgres or neon.tech/serverless`
- `export const runtime = 'edge'`

### Spring Boot (Undertow)
- `server.undertow.threads.worker=64`
- `spring.r2dbc.pool.enabled=true`
- `spring-boot-starter-webflux + r2dbc-postgresql`

## Test Script

For each web server, we will implement an http handler to execute the following SQL statement (using pgbench OLTP tables):

```sql
BEGIN;
-- 1. udpate account balance (:aid = random account ID (1 .. N×100000))
UPDATE pgbench_accounts
   SET abalance = abalance + :delta
 WHERE aid = :aid; -- 
-- 2. read account balance
SELECT abalance
  FROM pgbench_accounts
 WHERE aid = :aid;
-- 3. update teller total (:tid = random teller (1 .. N×10))
UPDATE pgbench_tellers
   SET tbalance = tbalance + :delta
 WHERE tid = :tid; 
-- 4. update branch total (:bid = random branch (1 .. N))
UPDATE pgbench_branches
   SET bbalance = bbalance + :delta
 WHERE bid = :bid;
-- 5. insert history
INSERT INTO pgbench_history (tid, bid, aid, delta, mtime, filler)
VALUES (:tid, :bid, :aid, :delta, now(), repeat('x',22));
END;
```

We use `grafana/k6` (install with `brew install k6`):

```bash
docker compose up -d
k6 run benchmark.js --env URL=http://localhost:$port
```

## SQL

### OLTP Schema

```
demo=# \d pgbench*
              Table "public.pgbench_accounts"
  Column  |     Type      | Collation | Nullable | Default
----------+---------------+-----------+----------+---------
 aid      | integer       |           | not null |
 bid      | integer       |           |          |
 abalance | integer       |           |          |
 filler   | character(84) |           |          |
Indexes:
    "pgbench_accounts_pkey" PRIMARY KEY, btree (aid)

              Table "public.pgbench_branches"
  Column  |     Type      | Collation | Nullable | Default
----------+---------------+-----------+----------+---------
 bid      | integer       |           | not null |
 bbalance | integer       |           |          |
 filler   | character(88) |           |          |
Indexes:
    "pgbench_branches_pkey" PRIMARY KEY, btree (bid)

                    Table "public.pgbench_history"
 Column |            Type             | Collation | Nullable | Default
--------+-----------------------------+-----------+----------+---------
 tid    | integer                     |           |          |
 bid    | integer                     |           |          |
 aid    | integer                     |           |          |
 delta  | integer                     |           |          |
 mtime  | timestamp without time zone |           |          |
 filler | character(22)               |           |          |

              Table "public.pgbench_tellers"
  Column  |     Type      | Collation | Nullable | Default
----------+---------------+-----------+----------+---------
 tid      | integer       |           | not null |
 bid      | integer       |           |          |
 tbalance | integer       |           |          |
 filler   | character(84) |           |          |
Indexes:
    "pgbench_tellers_pkey" PRIMARY KEY, btree (tid)

```

import os
import os.path
import random
import asyncpg
import uvicorn
from fastapi import FastAPI, Request, Response
from dynaconf import Dynaconf

settings = Dynaconf(
    settings_files=["config.yaml", "../../config.yaml", "/config/config.yaml"],
    envvar_prefix=False
)

pool: asyncpg.Pool

app = FastAPI()

@app.on_event("startup")
async def startup():
    global pool
    pool = await asyncpg.create_pool(
        dsn=settings.db,
        min_size=5,
        max_size=settings.db_max_connections
    )

@app.on_event("shutdown")
async def shutdown():
    await pool.close()

@app.get("/")
async def tx_handler(request: Request):
    scale = 10
    aid = random.randint(1, scale * 100_000)
    tid = random.randint(1, scale * 10)
    bid = random.randint(1, scale)
    delta = random.randint(-5000, 4999)

    async with pool.acquire() as conn:
        await conn.execute(
            "SELECT pgbench_tx($1, $2, $3, $4)",
            aid, tid, bid, delta
        )
    return Response(content="OK", media_type="text/plain")
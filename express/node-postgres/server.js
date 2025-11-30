import express from 'express';
import { Pool } from 'pg';
import yaml from 'js-yaml';
import { readFileSync } from 'fs';
import { resolve } from 'path';

function loadConfig() {
  const paths = [
    './config.yaml',
    '../../config.yaml',
    '/config/config.yaml'
  ];
  for (const p of paths) {
    try {
      return yaml.load(readFileSync(p, 'utf8'));
    } catch (_) {}
  }
  throw new Error('config.yaml not found');
}

const config = loadConfig();

const pool = new Pool({
  connectionString: config.db,
  max: config.db_max_connections || 400,
  idleTimeoutMillis: 0,
  connectionTimeoutMillis: 5000
});

const app = express();

const SCALE = 10;

app.get('/', async (req, res) => {
  const aid = Math.floor(Math.random() * SCALE * 100000) + 1;
  const tid = Math.floor(Math.random() * SCALE * 10) + 1;
  const bid = Math.floor(Math.random() * SCALE) + 1;
  const delta = Math.floor(Math.random() * 10000) - 5000;

  const client = await pool.connect();
  try {
    await client.query('SELECT pgbench_tx($1, $2, $3, $4)', [aid, tid, bid, delta]);
    res.type('text/plain').send('OK');
  } catch (err) {
    console.error(err);
    res.status(500).send('KO');
  } finally {
    client.release();
  }
});

app.listen(8080, () => console.log('listening on :8080'));
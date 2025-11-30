import express from 'express';
import postgres from 'postgres';
import yaml from 'js-yaml';
import { readFileSync } from 'fs';

function loadConfig() {
  const paths = ['./config.yaml', '../../config.yaml', '/config/config.yaml'];
  for (const p of paths) {
    try { return yaml.load(readFileSync(p, 'utf8')); }
    catch (_) {}
  }
  throw new Error('config.yaml not found');
}

const config = loadConfig();

const sql = postgres(config.db, {
  max: config.db_max_connections || 400,
  idle_timeout: null,
  connect_timeout: 5,
  prepare: false,
  transform: { undefined: null }
});

const app = express();

const SCALE = 10;

app.get('/', async (req, res) => {
  const aid = Math.floor(Math.random() * SCALE * 100000) + 1;
  const tid = Math.floor(Math.random() * SCALE * 10) + 1;
  const bid = Math.floor(Math.random() * SCALE) + 1;
  const delta = Math.floor(Math.random() * 10000) - 5000;

  try {
    await sql`SELECT pgbench_tx(${aid}, ${tid}, ${bid}, ${delta})`;
    res.type('text/plain').send('OK');
  } catch (err) {
    console.error(err);
    res.status(500).send('KO');
  }
});

const server = app.listen(8080, () => console.log('listening on :8080'));

process.on('SIGTERM', () => sql.end().then(() => server.close()));
process.on('SIGINT', () => sql.end().then(() => server.close()));
import postgres from 'postgres';
import yaml from 'js-yaml';
import { readFileSync } from 'fs';
import { join } from 'path';

let sql;

function loadConfig() {
  const paths = [
    join(process.cwd(), 'config.yaml'),
    join(process.cwd(), '..', '..', 'config.yaml'),
    '/config/config.yaml'
  ];

  for (const p of paths) {
    try {
      const content = readFileSync(p, 'utf8');
      return yaml.load(content);
    } catch (_) {}
  }
  throw new Error('config.yaml not found in any expected location');
}

export function getSql() {
  if (!sql) {
    const config = loadConfig();
    sql = postgres(config.db, {
      max: config.db_max_connections || 400,
      idle_timeout: null,
      connect_timeout: 5,
      prepare: false,
      transform: { undefined: null },
      // Graceful shutdown
      onclose: () => { sql = null; }
    });
  }
  return sql;
}
import { Hono } from 'hono';
import postgres from 'postgres';

const app = new Hono();

const configPaths = ['./config.yaml', '../../config.yaml', '/config/config.yaml'];
let dbUrl = '';

for (const path of configPaths) {
  const file = Bun.file(path);
  if (await file.exists()) {
    const text = await file.text();
    const match = text.match(/db:\s*["']?([^"'\n]+)["']?/);
    if (match) {
      dbUrl = match[1];
      break;
    }
  }
}

if (!dbUrl) {
  console.error('Could not find db connection string in config.yaml');
  process.exit(1);
}

const maxConns = Number(Bun.env.DB_MAX_CONNECTIONS) || 10;
const sql = postgres(dbUrl, { max: maxConns });

await sql`SELECT 1`.catch(err => {
  console.error('Database connection failed:', err);
  process.exit(1);
});

const SCALE = 10;

app.get('/*', async (c) => {
  const aid = Math.floor(Math.random() * SCALE * 100_000) + 1;
  const tid = Math.floor(Math.random() * SCALE * 10) + 1;
  const bid = Math.floor(Math.random() * SCALE) + 1;
  const delta = Math.floor(Math.random() * 10_000) - 5_000;

  try {
    await sql`SELECT pgbench_tx(${aid}, ${tid}, ${bid}, ${delta})`;
    return c.text('OK', 200);
  } catch (err: any) {
    console.error('tx failed:', err);
    return c.text('tx failed: ' + err.message, 500);
  }
});

const server = Bun.serve({
  port: 8080,
  fetch: app.fetch,
  reusePort: true,
});

console.log(`listening on http://0.0.0.0:${server.port}`);

const shutdown = async () => {
  console.log('\nShutting down...');
  server.stop(true);
  await sql.end();
  process.exit(0);
};

process.on('SIGINT', shutdown);
process.on('SIGTERM', shutdown);
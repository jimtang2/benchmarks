const SCALE = 10;

export const dynamic = 'force-dynamic'; 

export async function GET() {
  const { getSql } = await import('@/server/db');
  const sql = getSql();

  const aid = Math.floor(Math.random() * SCALE * 100000) + 1;
  const tid = Math.floor(Math.random() * SCALE * 10) + 1;
  const bid = Math.floor(Math.random() * SCALE) + 1;
  const delta = Math.floor(Math.random() * 10000) - 5000;

  try {
    await sql`SELECT pgbench_tx(${aid}, ${tid}, ${bid}, ${delta})`;
    return new Response('OK', {
      status: 200,
      headers: { 'Content-Type': 'text/plain' }
    });
  } catch (err) {
    console.error('tx error:', err);
    return new Response('KO', { status: 500 });
  }
}
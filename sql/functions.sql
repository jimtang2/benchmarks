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
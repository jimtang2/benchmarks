package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")
	viper.AddConfigPath("/config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	srv := http.Server{
		Addr:    ":8080",
		Handler: initHandler(),
	}
	go srv.ListenAndServe()
	defer srv.Close()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	log.Println("exiting:", <-c)
}

type Handler struct {
	pool *pgxpool.Pool
}

func initHandler() *Handler {
	handler := &Handler{}
	config, err := pgxpool.ParseConfig(viper.GetString("db"))
	if err != nil {
		log.Fatal("parse config:", err)
	}
	config.MaxConns = viper.GetInt32("db_max_connections")
	config.MinConns = 20
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 15 * time.Minute
	config.HealthCheckPeriod = time.Minute

	handler.pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("create pool:", err)
	}
	if err := handler.pool.Ping(context.Background()); err != nil {
		log.Fatal("ping:", err)
	}
	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const scale = 10
	aid := rand.Int63n(int64(scale*100000)) + 1
	tid := rand.Int63n(int64(scale*10)) + 1
	bid := rand.Int63n(int64(scale)) + 1
	delta := rand.Int63n(10000) - 5000
	tx, err := h.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		http.Error(w, "begin tx failed", 500)
		return
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx,
		`UPDATE pgbench_accounts SET abalance = abalance + $1 WHERE aid = $2`,
		delta, aid)
	if err != nil {
		http.Error(w, "accounts", 500)
		return
	}
	_, err = tx.Exec(ctx,
		`UPDATE pgbench_tellers SET tbalance = tbalance + $1 WHERE tid = $2`,
		delta, tid)
	if err != nil {
		http.Error(w, "tellers", 500)
		return
	}
	_, err = tx.Exec(ctx,
		`UPDATE pgbench_branches SET bbalance = bbalance + $1 WHERE bid = $2`,
		delta, bid)
	if err != nil {
		http.Error(w, "branches", 500)
		return
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO pgbench_history (tid, bid, aid, delta, mtime, filler)
		 VALUES ($1, $2, $3, $4, now(), repeat('x',22))`,
		tid, bid, aid, delta)
	if err != nil {
		http.Error(w, "history", 500)
		return
	}
	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "commit failed", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

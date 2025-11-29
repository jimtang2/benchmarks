package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

var (
	pool *pgxpool.Pool
	srv  http.Server
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
	config, err := pgxpool.ParseConfig(viper.GetString("db"))
	if err != nil {
		log.Fatal("parse config:", err)
	}
	config.MaxConns = viper.GetInt32("db_max_connections")
	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("create pool:", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("ping:", err)
	}
	srv = http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(txHandleFunc),
	}
}

func main() {
	go srv.ListenAndServe()
	defer srv.Close()
	log.Println("listening")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	log.Println("exiting:", <-c)
}

func txHandleFunc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const scale = 10
	aid := rand.Int63n(int64(scale*100000)) + 1
	tid := rand.Int63n(int64(scale*10)) + 1
	bid := rand.Int63n(int64(scale)) + 1
	delta := rand.Int63n(10000) - 5000
	if _, err := pool.Exec(ctx,
		`SELECT pgbench_tx($1, $2, $3, $4)`, aid, tid, bid, int32(delta)); err != nil {
		log.Println(err)
		http.Error(w, "tx failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

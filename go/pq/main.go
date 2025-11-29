package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
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
	db *sql.DB
}

func initHandler() *Handler {
	handler := &Handler{}
	var err error
	handler.db, err = sql.Open("postgres", viper.GetString("db"))
	if err != nil {
		log.Fatal(err)
	}
	handler.db.SetMaxOpenConns(viper.GetInt("db_max_connections"))
	handler.db.SetMaxIdleConns(50)
	handler.db.SetConnMaxLifetime(time.Hour)
	handler.db.SetConnMaxIdleTime(15 * time.Minute)
	if err = handler.db.Ping(); err != nil {
		log.Fatal("ping:", err)
	}
	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const scale = 10
	aid := rand.Int63n(int64(scale*100000)) + 1
	tid := rand.Int63n(int64(scale*10)) + 1
	bid := rand.Int63n(int64(scale)) + 1
	delta := rand.Int63n(10000) - 5000
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, "begin tx failed", 500)
		return
	}
	defer tx.Rollback()
	_, err = tx.Exec(`UPDATE pgbench_accounts SET abalance = abalance + $1 WHERE aid = $2`, delta, aid)
	if err != nil {
		http.Error(w, "accounts", 500)
		return
	}
	_, err = tx.Exec(`UPDATE pgbench_tellers SET tbalance = tbalance + $1 WHERE tid = $2`, delta, tid)
	if err != nil {
		http.Error(w, "tellers", 500)
		return
	}
	_, err = tx.Exec(`UPDATE pgbench_branches SET bbalance = bbalance + $1 WHERE bid = $2`, delta, bid)
	if err != nil {
		http.Error(w, "branches", 500)
		return
	}
	_, err = tx.Exec(`INSERT INTO pgbench_history (tid, bid, aid, delta, mtime, filler) VALUES ($1, $2, $3, $4, now(), repeat('x',22))`,
		tid, bid, aid, delta)
	if err != nil {
		http.Error(w, "history", 500)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "commit failed", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

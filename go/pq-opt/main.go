package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var (
	db  *sql.DB
	srv http.Server
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
	var err error
	db, err = sql.Open("postgres", viper.GetString("db"))
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(viper.GetInt("db_max_connections"))
	if err = db.Ping(); err != nil {
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
	c := make(chan os.Signal, 1)
	log.Println("listening")
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	log.Println("exiting:", <-c)
}

func txHandleFunc(w http.ResponseWriter, r *http.Request) {
	const scale = 10
	aid := rand.Int63n(int64(scale*100000)) + 1
	tid := rand.Int63n(int64(scale*10)) + 1
	bid := rand.Int63n(int64(scale)) + 1
	delta := rand.Int63n(10000) - 5000
	if _, err := db.Exec(`SELECT pgbench_tx($1, $2, $3, $4)`, aid, tid, bid, int32(delta)); err != nil {
		log.Println(err)
		http.Error(w, "tx failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

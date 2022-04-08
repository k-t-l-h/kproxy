package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"kproxy/m/internal"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)

	conn := "postgres://postgres:postgres@127.0.0.1:5432/postgres?pool_max_conns=1000"

	pool, err := pgxpool.Connect(context.Background(), conn)
	if err != nil {
		log.Println(err)
	}


	//tcp part
	p := internal.Proxy{
		pool,
	}
	p.Run()

	muxRoute := mux.NewRouter()
	api := muxRoute.PathPrefix("/api/v1").Subrouter()
	{
		api.HandleFunc("/requests", p.GetList).Methods(http.MethodGet)
		api.HandleFunc("/requests/{id}", p.GetOne).Methods(http.MethodGet)
		api.HandleFunc("/repeat/{id}", p.RepeateOne).Methods(http.MethodGet)
		api.HandleFunc("/scan/{id}", p.RepeatSQLInj).Methods(http.MethodGet)
	}

	//graceful shutdown
	log.Print("Signal received: ", <-signals)


}

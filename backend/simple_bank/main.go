package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/litmus-zhang/simple_bank/api"
	"github.com/litmus-zhang/simple_bank/bank"
	"github.com/litmus-zhang/simple_bank/util"
)

func main() {

	var err error
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBURL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := bank.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

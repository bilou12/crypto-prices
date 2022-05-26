package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	sdk "github.com/adshao/go-binance/v2"

	"github.com/rs/zerolog"
	"gitlab-v2.assetik.com/benoit/crypto-prices/pkg/binance"
	"gitlab-v2.assetik.com/benoit/crypto-prices/pkg/persistance"
)

func main() {
	// read config from file
	file, _ := ioutil.ReadFile("config.json")
	var m map[string]string
	json.Unmarshal([]byte(file), &m)
	fmt.Println(m)

	ctx := context.Background()
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	host := m["host"]
	port, err := strconv.Atoi(m["port"])
	if err != nil {
		panic(err)
	}
	user := m["user"]
	password := m["password"]
	dbmktbar := m["dbmktbar"]
	dbinfra := m["dbinfra"]

	// create dbMktBar
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbmktbar)
	dbMktBar, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer dbMktBar.Close()

	// create dbInfra
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbinfra)
	dbInfra, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer dbInfra.Close()

	myPostgres := persistance.NewPostgres(dbMktBar, dbInfra)
	binance := binance.NewBinance(myPostgres, &log, &ctx, sdk.NewClient("", ""))

	interval := "1h"
	sqlStatement := `select distinct (symbol)
	from signal_preprocessing.cryptos_ison_universe_liq ciul
	order by symbol;`

	rows, err := binance.MyPostgres.DbInfra.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var (
		pair string
	)
	pairs := make(map[string]string)
	for rows.Next() {
		err := rows.Scan(&pair)
		if err != nil {
			panic(err)
		}
		fmt.Println(pair)
		pairs[pair] = interval
	}

	binance.GetPriceRt(pairs, interval)
}

package persistance

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type MyPostgres struct {
	DbMktBar *sql.DB
	DbInfra  *sql.DB
}

func NewPostgres(dbMktBar *sql.DB, dbInfra *sql.DB) *MyPostgres {
	return &MyPostgres{
		DbMktBar: dbMktBar,
		DbInfra:  dbInfra,
	}
}

func (me *MyPostgres) CreateTablePrice(symbol string, interval string) {
	sqlStatement := `CREATE TABLE if not exists binance_` + interval + `_rt."` + symbol + `_` + interval + `" (
		open_time numeric NULL,
		open_date_time timestamptz NOT NULL,
		"open" numeric NULL,
		"high" numeric NULL,
		"low" numeric NULL,
		"close" numeric NULL,
		volume numeric NULL,
		close_time numeric NULL,
		close_date_time timestamptz NOT NULL,
		quote_asset_volume numeric NULL,
		number_of_trade numeric NULL,
		taker_buy_base_asset_volume numeric NULL,
		taker_buy_quote_asset_volume numeric NULL,
		log_timestamp timestamptz NULL DEFAULT now(),
		CONSTRAINT "` + symbol + `_` + interval + `_rt_pk" PRIMARY KEY (open_date_time, close_date_time)
	);`

	_, err := me.DbMktBar.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v2"
)

func dbinit() {
	/*
		テーブル構成
		主キー 記事URL(ドメインも) : String
		いいね数 : int
	*/
	type likeAPITable struct {
		URL                 string "db:url, primarykey notnull" //SQLiteでは歴史的経緯で主キーがnullで登録できてしまうらしいのでnotnull
		Like                int    "db:like notnull"
		LastUpdatedUnixTime int64  "db:last_updated_unixtime notnull"
	}
	db, err := sql.Open("sqlite3", "./likeapiserver.db")
	checkErr(err, "sql.Open Failed")
	defer db.Close()

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}} //gorpの初期化

	//データベースのテーブルを作成
	dbmap.AddTableWithName(likeAPITable{}, baseurl[0]) //TODO baseurl[0]は仮
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create table failed")
}

package db

import (
	"database/sql"
	"log"
)

func dbinit() {
	db, err := sql.Open("sqlite3", "./likeapiserver.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}

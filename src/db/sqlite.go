package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v2"
	"log"
	"net/http"
	"time"
)

/*
テーブル構成
主キー 記事URL(ドメインも) : String
いいね数 : int
*/
type like struct {
	URL                 string        `db:"url, primarykey, notnull"` //SQLiteでは歴史的経緯で主キーがnullで登録できてしまうらしいのでnotnull必須
	Like                int           `db:"like, notnull"`
	LastUpdatedUnixTime sql.NullInt64 `db:"last_updated_unixtime, notnull"`
}

// Init DB初期化
func Init(baseurl string) *gorp.DbMap {
	db, err := sql.Open("sqlite3", "./likeapiserver.db")
	if err != nil {
		log.Fatal("sql.Open Failed")
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}} //gorpの初期化
	//データベースのテーブルを作成
	dbmap.AddTableWithName(like{}, baseurl)
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatal("Create table failed : %s", err)
	}

	return dbmap
}

// GetLike intが0を返したらエラーです
func GetLike(url string, dbmap *gorp.DbMap) (int, error) {
	resultobj, err := dbmap.Get(like{}, url)
	if err != nil {
		return 0, err
	}
	return resultobj.(like).Like, nil
}

// LikeInc いいね数を増やす
func LikeInc(url string, inc int, dbmap *gorp.DbMap) error {
	//select
	resultobj, err := dbmap.Get(like{}, url)
	if err != nil {
		log.Printf("dbmap.Get failed : %s", err)
		return err
	}
	//すでにDBに存在しない場合は新規登録
	if resultobj == nil {
		//記事の存在をHTTPリクエストで検証
		if req, err := http.Get(url); err != nil || req.StatusCode != http.StatusOK {
			if err != nil {
				log.Printf("http.Get failed : %s", err)
				return err
			}
		}
		//挿入
		err := dbmap.Insert(&like{URL: url, Like: 0 + inc, LastUpdatedUnixTime: sql.NullInt64{time.Now().Unix(), true}})
		if err != nil {
			log.Printf("dbmap.Insert failed : %s", err)
			return err
		}
		return nil
	}
	//すでにDBに存在する場合はupdate
	_, err = dbmap.Update(&like{URL: url, Like: resultobj.(like).Like + inc, LastUpdatedUnixTime: sql.NullInt64{time.Now().Unix(), true}})
	if err != nil {
		log.Printf("dbmap.Update failed : %s", err)
		return err
	}
	return nil
}

package main

import (
	"database/sql"
	"errors"
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

// DBInit DB初期化
func DBInit(baseurl string) *gorp.DbMap {
	db, err := sql.Open("sqlite3", "./likeapiserver.db")
	if err != nil {
		log.Fatal("sql.Open Failed")
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}} //gorpの初期化
	//データベースのテーブルを作成
	dbmap.AddTableWithName(like{}, baseurl)
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalf("Create table failed : %s", err)
	}

	return dbmap
}

// GetLike
func GetLike(url string, dbmap *gorp.DbMap) (int, error) {
	resultobj, err := dbmap.Get(like{}, url)
	if err != nil {
		return -1, err
	}
	//存在しない場合0を返す
	if resultobj == nil {
		return 0, nil
	}
	return resultobj.(*like).Like, nil
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
		//本当はmain関数内でやって、未登録の場合BadRequestのほうがいいが、未登録か新規登録かはDBにアクセスしないとわからないので渋々ここでやる(2回クエリ叩くのはちょっと...)
		//記事の存在をHTTPリクエストで検証
		if req, err := http.Get(url); err != nil || req.StatusCode != http.StatusOK {
			if err != nil {
				log.Printf("http.Get failed : %s", err)
				return err
			}
			return errors.New("http.Get not succeed : " + url + " " + req.Status)
		}
		//新規挿入
		err := dbmap.Insert(&like{URL: url, Like: 0 + inc, LastUpdatedUnixTime: sql.NullInt64{Int64: time.Now().Unix(), Valid: true}})
		if err != nil {
			log.Printf("dbmap.Insert failed : %s", err)
			return err
		}
		return nil
	}
	t, ok := resultobj.(*like)
	if !ok {
		log.Print("resultobj.(like) failed")
		return errors.New("resultobj.(like) failed")
	}
	//すでにDBに存在する場合はupdate
	_, err = dbmap.Update(&like{URL: url, Like: t.Like + inc, LastUpdatedUnixTime: sql.NullInt64{time.Now().Unix(), true}})
	if err != nil {
		log.Printf("dbmap.Update failed : %s", err)
		return err
	}
	return nil
}

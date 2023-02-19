package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v2"
	"log"
	"net/http"
	"strings"
	"time"
)

// TODO  仮実装(configから読むべき)
var baseurl = []string{"https://powerfulfamily.net/"}

/*
テーブル構成
主キー 記事URL(ドメインも) : String
いいね数 : int
*/
type like struct {
	URL                 string        "db:url, primarykey notnull" //SQLiteでは歴史的経緯で主キーがnullで登録できてしまうらしいのでnotnull必須
	Like                int           "db:like notnull"
	LastUpdatedUnixTime sql.NullInt64 "db:last_updated_unixtime"
}

var dbmap *gorp.DbMap

// Init DB初期化
func Init() {
	db, err := sql.Open("sqlite3", "./likeapiserver.db")
	checkErr(err, "sql.Open Failed")
	defer db.Close()

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}} //gorpの初期化

	//データベースのテーブルを作成
	dbmap.AddTableWithName(like{}, baseurl[0]) //TODO baseurl[0]は仮
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create table failed")
}

// GetLike intが0を返したらエラーです
func GetLike(url string) (error, int) {
	if notContainBaseURL(url) {
		return errors.New("URL is not in baseurl(not valid domain)"), 0
	}
	selectResult, err := dbmap.Select(like{}, "select * from "+baseurl[0]+" where url=?", url)
	checkErr(err, "select failed")
	if len(selectResult) > 2 {
		return errors.New("select result more than 2 (DataBase mismatch error)"), 0
	}
	selectResult
}

// LikeInc いいね数を増やす
func LikeInc(url string, inc int) error {
	//URLが正しいか検証
	if notContainBaseURL(url) {
		return errors.New("URL is not in baseurl(not valid domain)")
	}
	//select
	selectResult, err := dbmap.SelectOne(like{}, "select * from "+baseurl[0]+" where url=?", url)
	checkErr(err, "select failed")
	////検索結果が2つ以上出た場合
	//if len(selectResult) > 2 {
	//	return errors.New("select result more than 2 (DataBase mismatch error)")
	//}
	//すでにDBに存在しない場合は新規登録
	if len(selectResult) == 0 {
		//記事の存在を検証
		if req, err := http.Get(url); err != nil || req.StatusCode != http.StatusOK {
			checkErr(err, "http.Get failed")
		}
		//挿入
		err := dbmap.Insert(&like{URL: url, Like: 0 + inc, LastUpdatedUnixTime: sql.NullInt64{}})
		checkErr(err, "Insert failed")
	}
	//すでにDBに存在する場合はupdate
	_, err = dbmap.Update(&like{URL: url, Like: selectResult[0].(like).Like + inc, LastUpdatedUnixTime: sql.NullInt64{time.Now().Unix(), false}})
	checkErr(err, "likeIncrement : update Failed")
	return nil
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln("Error : "+msg, err)
	}
}
func notContainBaseURL(url string) bool {
	if strings.Contains(url, baseurl[0]) {
		return true
	}
	return false
}

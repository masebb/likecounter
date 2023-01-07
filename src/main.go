package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

// 仮実装//

var (
	like    int
	baseurl = []string{"https://powerfulfamily.net/"}
)

//仮実装//

func main() {
	//DB初期化
	dbinit()

	//Echo初期化
	e := echo.New()
	//レートリミッター(5req/sec)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	e.POST("/v1/like", likeIncrement)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "It works!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

type Request struct {
	Increment int `json:"increment"`
}
type Response struct {
	Like  int    `json:"like"`
	Error string `json:"error"`
}

func likeIncrement(c echo.Context) error {
	post := new(Request)
	resp := new(Response)
	//リクエストのパース
	if err := c.Bind(post); err != nil {
		resp.Like = like
		resp.Error = "Bad Request"
		return c.JSON(http.StatusBadRequest, resp)
	}
	//0はありえない
	if post.Increment == 0 {
		resp.Like = like
		resp.Error = "Bad Request"
		return c.JSON(http.StatusBadRequest, resp)
	}
	//いいねが多すぎる場合はシャットアウト
	if post.Increment > 31 {
		resp.Like = like
		resp.Error = "increment is too large"
		return c.JSON(http.StatusBadRequest, resp)
	}
	////リファラを見て代入
	rurl := c.Request().Header.Get("Referer")

	//いいねの数を増やす
	//TODO DB化
	//URLチェック後にDB投入
	like = like + post.Increment
	//レスポンスにlikeの数を入れる
	resp.Like = like
	return c.JSON(http.StatusOK, resp)
}

// TODO
// 初めていいねのリクエストが送られるものに本当にその記事は存在するのかチェックする
func isUrlExist(url string) bool {
	req, _ := http.Get(url)
	if req.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln("Error : "+msg, err)
	}
}

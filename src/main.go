package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"likeapiserver/src/db"
	"log"
	"net/http"
	"strings"
)

// TODO  仮実装(configから読むべき)
var baseurl = []string{"https://powerfulfamily.net/"}

// Request なぜリクエストにURLパラメータを含めないのか : クライアント側で自身のURLを組み立てるのは面倒だし、別にリファラを見ることで事足りるため
type Request struct {
	Increment int `json:"increment"`
}
type PostResponse struct {
	Error string `json:"error"`
}
type GetResponse struct {
	Like int `json:like`
}

func main() {
	if len(baseurl) == 0 {
		// エラーハンドリング
		log.Fatal("baseurl is empty")
	}
	dbmap := db.Init(baseurl[0])
	defer dbmap.Db.Close()
	//Echo初期化
	e := echo.New()
	//レートリミッター(5req/sec)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	e.Use(middleware.Recover())
	//Likeを取得する
	e.GET("/v1/like", func(c echo.Context) error {
		post := new(Request)
		resp := new(GetResponse)
		//リクエストのパース
		if err := c.Bind(post); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		refurl := c.Request().Header.Get("Referer")
		resp.Like, _ = db.GetLike(refurl, dbmap)
		err := c.JSON(http.StatusOK, resp)
		if err != nil {
			return err
		}
		return nil
	})
	//Likeをインクリメントする
	e.POST("/v1/like", func(c echo.Context) error {
		post := new(Request)
		resp := new(PostResponse)
		refurl := c.Request().Header.Get("Referer")
		//リクエストのパース
		if err := c.Bind(post); err != nil {
			resp.Error = "Bad Request"
			return c.JSON(http.StatusBadRequest, resp)
		}
		//0以下はありえない
		if post.Increment < 0 {
			resp.Error = "Bad Request"
			return c.JSON(http.StatusBadRequest, resp)
		}
		//いいねが多すぎる場合は拒否
		if post.Increment > 31 {
			resp.Error = "increment is too large"
			return c.JSON(http.StatusBadRequest, resp)
		}
		if notContainBaseURL(refurl) {
			resp.Error = "Bad Request"
			return c.JSON(http.StatusBadRequest, resp)
		}
		//いいねの数を増やす
		err := db.LikeInc(refurl, post.Increment, dbmap)
		if err != nil {
			//DBエラー
			log.Printf("DBError: %s", err)
			resp.Error = "Internal Error"
			return c.JSON(http.StatusInternalServerError, resp)
		}
		resp.Error = ""
		return c.JSON(http.StatusOK, resp)
	})
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "It works!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
func notContainBaseURL(url string) bool {
	return !strings.Contains(url, baseurl[0])
}

package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"strings"
)

// TODO  仮実装(configから読むべき)
var baseurl = []string{"http://powerfulfamily.net/"}

type GetRequest struct {
	URL string `query:"url"`
}
type PostRequest struct {
	//URL       string `query:"url"`
	Increment int `json:"increment"`
}
type GetResponse struct {
	Like int `json:like`
}

func main() {
	if len(baseurl) < 0 {
		// エラーハンドリング
		log.Fatal("baseurl is empty")
	}
	dbmap := DBInit(baseurl[0])
	defer dbmap.Db.Close()
	//Echo初期化
	e := echo.New()
	//レートリミッター(5req/sec)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	e.Use(middleware.Recover())
	//Likeを取得する
	e.GET("/v1/like", func(c echo.Context) error {
		log.Printf("GET /v1/like %s", c)
		resp := new(GetResponse)
		var req GetRequest
		//リクエストのパース
		if err := c.Bind(&req); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		//DBから取得
		var err error
		if resp.Like, err = GetLike(req.URL, dbmap); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		//レスポンス
		if err := c.JSON(http.StatusOK, resp); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return nil
	})
	//Likeをインクリメントする
	e.POST("/v1/like", func(c echo.Context) error {
		post := new(PostRequest)
		//リクエストのパース
		if err := c.Bind(&post); err != nil {
			log.Printf("BindError: %s", err)
			return c.NoContent(http.StatusBadRequest)
		}
		var url = c.QueryParam("url")
		//likeが多すぎる場合は拒否
		if post.Increment > 31 || post.Increment < 0 {
			log.Printf("IncrementError: %s", c)
			return c.NoContent(http.StatusBadRequest)
		}
		//ベースURLが入っていないやつは拒否
		if notContainBaseURL(url) {
			log.Printf("NotContainBaseURL: %s", c)
			return c.NoContent(http.StatusBadRequest)
		}
		//いいねの数を増やす
		if err := LikeInc(url, post.Increment, dbmap); err != nil {
			//DBエラー
			log.Printf("DBError: %s", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	})
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "It works!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
func notContainBaseURL(url string) bool {
	return !strings.Contains(url, baseurl[0])
}

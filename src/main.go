package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"likeapiserver/src/db"
	"net/http"
)

func main() {
	//DB初期化
	db.Init()
	//Echo初期化
	e := echo.New()
	//レートリミッター(5req/sec)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	e.POST("/v1/like", likeIncrement)
	e.GET("/v1/like", likeGet)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "It works!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

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

func likeGet(c echo.Context) error {
	post := new(Request)
	resp := new(GetResponse)
	//リクエストのパース
	if err := c.Bind(post); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	refurl := c.Request().Header.Get("Referer")
	resp.Like, _ = db.GetLike(refurl)
	err := c.JSON(http.StatusOK, resp)
	if err != nil {
		return err
	}
	return nil
}
func likeIncrement(c echo.Context) error {
	post := new(Request)
	resp := new(PostResponse)
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
	//リファラを見て代入
	refurl := c.Request().Header.Get("Referer")

	//いいねの数を増やす
	err := db.LikeInc(refurl, post.Increment)
	if err != nil {
		//DBエラー
		resp.Error = "Internal Error"
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Error = ""
	return c.JSON(http.StatusOK, resp)
}

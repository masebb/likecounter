package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

var like int

func main() {
	e := echo.New()
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	e.POST("/v1/like", likeIncrement)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "It works!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

type request struct {
	Increment int `json:"increment"`
}
type response struct {
	Like int `json:"like"`
}

func likeIncrement(c echo.Context) error {
	post := new(request)
	resp := new(response)
	//リクエストのパース
	if err := c.Bind(post); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	//いいねが多すぎる場合はシャットアウト
	if post.Increment > 31 {
		resp.Like = like
		return c.NoContent(http.StatusBadRequest)
	}
	//いいねの数を増やす
	//TODO DB化
	like = like + post.Increment
	//レスポンスにlikeの数を入れる
	resp.Like = like
	return c.JSON(http.StatusOK, resp)
}

// 初めていいねのリクエストが送られるものに本当にその記事は存在するのかチェックする
func isUrlExist(url string) bool {
	req, _ := http.Get(url)
	if req.StatusCode != http.StatusOK {
		return false
	}
	return true
}

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type lotteryController struct {
	Ctx iris.Context
}

func NewApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}
func main() {
	app := NewApp()
	app.Run(iris.Addr(":8080"))
}

// 即开即得型
func (c *lotteryController) Get() string {
	var prize string
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Intn(10)
	switch {
	case code == 1:
		prize = "一等奖"
	case code == 2:
		prize = "二等奖"
	case code == 3:
		prize = "三等奖"
	default:
		return fmt.Sprintf("未中奖，%v", code)
	}
	return fmt.Sprintf("%v,%v", code, prize)
}

// 自选型
func (c *lotteryController) GetPrize() [7]int {
	var codes [7]int
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < 6; i++ {
		code := r.Intn(33) + 1
		codes[i] = code
	}
	codes[6] = r.Intn(16) + 1
	return codes
}

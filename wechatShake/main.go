package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

const (
	giftTypeCoin = iota
	giftTypeCoupon
	giftTypeCouponFix
	giftTypeRealSmall
	giftTypeRealLarge
)

type gift struct {
	id       int
	name     string
	pic      string
	link     string
	gtype    int
	data     string
	datalist []string //奖品数据集合
	total    int      // 0 不限量； 1只有1个
	left     int
	inuse    bool
	rate     int //中奖概率万分之n
	rateMin  int
	rateMax  int
}

const rateMax = 10000

var logger *log.Logger
var giftList []*gift

type lotteryController struct {
	Ctx iris.Context
}

func initLog() {
	f, _ := os.Create("./lottery.log")
	logger = log.New(f, "", log.Ldate|log.Lmicroseconds)
}

func initGift() {
	giftList = make([]*gift, 5)
	g1 := gift{
		id:       1,
		name:     "手机大奖",
		pic:      "",
		link:     "",
		gtype:    giftTypeRealLarge,
		data:     "",
		datalist: nil,  //奖品数据集合
		total:    1000, // 0 不限量； 1只有1个
		left:     1000,
		inuse:    true,
		rate:     10000, //中奖概率万分之n
		rateMin:  0,
		rateMax:  0,
	}
	g2 := gift{
		id:       2,
		name:     "充电器",
		pic:      "",
		link:     "",
		gtype:    giftTypeRealLarge,
		data:     "",
		datalist: nil, //奖品数据集合
		total:    5,   // 0 不限量； 1只有1个
		left:     5,
		inuse:    true,
		rate:     100, //中奖概率万分之n
		rateMin:  0,
		rateMax:  0,
	}
	g3 := gift{
		id:       3,
		name:     "满50减20",
		pic:      "",
		link:     "",
		gtype:    giftTypeCouponFix,
		data:     "mall-coupon-2019",
		datalist: nil, //奖品数据集合
		total:    5,   // 0 不限量； 1只有1个
		left:     5,
		inuse:    true,
		rate:     5000, //中奖概率万分之n
		rateMin:  0,
		rateMax:  0,
	}
	g4 := gift{
		id:       4,
		name:     "直降50元优惠券",
		pic:      "",
		link:     "",
		gtype:    giftTypeCoupon,
		data:     "",
		datalist: []string{"c1", "c2", "c3", "c4", "c5"}, //奖品数据集合
		total:    5,                                      // 0 不限量； 1只有1个
		left:     5,
		inuse:    true,
		rate:     5000, //中奖概率万分之n
		rateMin:  0,
		rateMax:  0,
	}
	g5 := gift{
		id:       5,
		name:     "金币",
		pic:      "",
		link:     "",
		gtype:    giftTypeCoin,
		data:     "",
		datalist: nil, //奖品数据集合
		total:    5,   // 0 不限量； 1只有1个
		left:     5,
		inuse:    true,
		rate:     5000, //中奖概率万分之n
		rateMin:  0,
		rateMax:  0,
	}
	giftList[0] = &g1
	giftList[1] = &g2
	giftList[2] = &g3
	giftList[3] = &g4
	giftList[4] = &g5

}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	initLog()
	initGift()
	return app
}
func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}
func (c *lotteryController) Get() string {
	return fmt.Sprintf("current total gift:%v\n", len(giftList))
}

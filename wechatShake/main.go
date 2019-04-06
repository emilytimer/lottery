package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

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
		rate:     1000, //中奖概率万分之n
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
		rate:     1000, //中奖概率万分之n
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
		rate:     1000, //中奖概率万分之n
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
		rate:     1000, //中奖概率万分之n
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
		rate:     1000, //中奖概率万分之n
		rateMin:  0,
		rateMax:  0,
	}
	giftList[0] = &g1
	giftList[1] = &g2
	giftList[2] = &g3
	giftList[3] = &g4
	giftList[4] = &g5
	rateStart := 0
	for _, data := range giftList {
		if !data.inuse {
			continue
		}
		data.rateMin = rateStart
		data.rateMax = rateStart + data.rate
		if data.rateMax > rateMax {
			data.rateMax = rateMax
			rateStart = 0 // ? 该有重复了啊
		} else {
			rateStart = rateStart + data.rate
		}
	}

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
func (c *lotteryController) GetLucky() map[string]interface{} {
	code := luckyCode()
	ok := false
	result := make(map[string]interface{})

	for _, data := range giftList {
		if !data.inuse || (data.total > 0 && data.left <= 0) {
			continue
		}

		if data.rateMin <= int(code) && data.rateMax > int(code) {
			sendData := ""

			switch data.gtype {
			case giftTypeCoin:
				ok, sendData = sendCoin(data)
			case giftTypeCoupon:
				ok, sendData = sendCoupon(data)

			case giftTypeCouponFix:
				ok, sendData = sendCouponFix(data)

			case giftTypeRealLarge:
				ok, sendData = sendRealLarge(data)

			case giftTypeRealSmall:
				ok, sendData = sendRealSmall(data)
			}
			if ok {
				saveLucky(code, data.id, data.name, data.link, sendData, data.left)
				result["success"] = ok
				result["id"] = data.id
				result["name"] = data.name
				result["data"] = sendData
			}
		}
	}
	return result
}
func saveLucky(code int32, id int, name, link, SendData string, left int) {
	logger.Printf("lucky, code=%d,gift=%d,name=%s,link=%s,data=%s,left=%d", code, id, name, link, SendData, left)
}
func sendCoin(data *gift) (bool, string) {
	fmt.Println("data.total:", data.total)
	if data.total == 0 {
		return true, data.data
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

// 不同面额得优惠券
func sendCoupon(data *gift) (bool, string) {
	if data.left > 0 {
		left := data.left - 1
		data.left = left
		return true, data.datalist[left]
	} else {
		return false, "奖品已发完"
	}
}

// 一样的优惠券
func sendCouponFix(data *gift) (bool, string) {
	if data.total == 0 {
		return true, data.data
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

// 小的实物奖
func sendRealSmall(data *gift) (bool, string) {
	if data.total == 0 {
		return true, data.data
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

// 大得实物奖
func sendRealLarge(data *gift) (bool, string) {
	if data.total == 0 {
		return true, data.data
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}
func luckyCode() int32 {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed)).Int31n(int32(rateMax))
	return r
}

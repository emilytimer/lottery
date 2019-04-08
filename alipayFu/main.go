package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type gift struct {
	Id      int
	Name    string
	Inuse   bool
	Rate    int //中奖概率万分之n
	RateMin int
	RateMax int
}

const rateMax = 10000

type lotteryController struct {
	Ctx iris.Context
}

var logger *log.Logger

func initLogger() {
	f, _ := os.Create("./lottery.log")
	logger = log.New(f, "", log.Ldate|log.Lmicroseconds)
}
func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	initLogger()
	return app
}
func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}
func newGift() *[5]gift {
	giftList := new([5]gift)
	g1 := gift{
		Id:      1,
		Name:    "富强福",
		Inuse:   true,
		Rate:    0, //中奖概率万分之n
		RateMin: 0,
		RateMax: 0,
	}
	g2 := gift{
		Id:      2,
		Name:    "和谐福",
		Inuse:   true,
		Rate:    0, //中奖概率万分之n
		RateMin: 0,
		RateMax: 0,
	}
	g3 := gift{
		Id:      3,
		Name:    "友善福",
		Inuse:   true,
		Rate:    0, //中奖概率万分之n
		RateMin: 0,
		RateMax: 0,
	}
	g4 := gift{
		Id:      4,
		Name:    "爱国福",
		Inuse:   true,
		Rate:    0, //中奖概率万分之n
		RateMin: 0,
		RateMax: 0,
	}
	g5 := gift{
		Id:      5,
		Name:    "敬业福",
		Inuse:   true,
		Rate:    0, //中奖概率万分之n
		RateMin: 0,
		RateMax: 0,
	}
	giftList[0] = g1
	giftList[1] = g2
	giftList[2] = g3
	giftList[3] = g4
	giftList[4] = g5
	return giftList
}
func giftRate(rate string) *[5]gift {
	giftList := newGift()
	rates := strings.Split(rate, ",")
	ratesLen := len(rates)
	rateStart := 0
	for i, data := range giftList {
		if !data.Inuse {
			continue
		}
		grate := 0
		if i < ratesLen {
			grate, _ = strconv.Atoi(rates[i])
		}
		giftList[i].Rate = grate
		giftList[i].RateMin = rateStart
		giftList[i].RateMax = rateStart + grate

		if data.RateMax >= rateMax {
			rateStart = 0
		} else {
			rateStart += grate
		}
	}
	fmt.Println("giftList=%v\n", giftList)

	return giftList
}
func (c *lotteryController) Get() string {
	rate := c.Ctx.URLParamDefault("rate", "4,3,2,1,0")
	giftList := giftRate(rate)
	v, _ := json.Marshal(giftList)
	return string(v)
}
func luckyCode() int32 {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed)).Int31n(int32(rateMax))
	return r
}
func (c *lotteryController) GetLucky() map[string]interface{} {
	uid, _ := c.Ctx.URLParamInt("uid")
	rate := c.Ctx.URLParamDefault("rate", "4000,3000,2000,1000,0")
	giftList := giftRate(rate)
	code := luckyCode()
	result := make(map[string]interface{})
	result["success"] = false
	for _, data := range giftList {
		if !data.Inuse {
			continue
		}
		result["number"] = code
		if data.RateMin < int(code) && data.RateMax > int(code) {
			result["uid"] = uid
			result["success"] = true
			result["gift"] = data.Name
			saveLucky(code, data.Id, data.Name)
		}
	}
	return result
}
func saveLucky(code int32, id int, name string) {
	logger.Printf("lucky, code=%d,id=%d,name=%s", code, id, name)
}

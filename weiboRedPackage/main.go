package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

var packageList map[uint32][]uint = make(map[uint32][]uint)

type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}
func (c *lotteryController) Get() map[uint32][2]int {
	rs := make(map[uint32][2]int)
	for id, list := range packageList {
		var money int
		for _, m := range list {
			money += int(m)
		}
		rs[id] = [2]int{len(list), money}
	}
	return rs
}

// http://localhost:8080/set?uid=1&money=100&number=100
func (c *lotteryController) GetSet() string {
	uid, uidErr := c.Ctx.URLParamInt("uid")
	money, moneyErr := c.Ctx.URLParamInt("money")
	number, numberErr := c.Ctx.URLParamInt("number")
	if uidErr != nil || moneyErr != nil || numberErr != nil {
		return fmt.Sprintf("参数格式异常,uidErr=%d,moneyErr=%d,numberErr=%d", uid, money, number)
	}
	money = int(money * 100)
	if uid < 1 || money < 1 || number < 1 {
		return fmt.Sprintf("参数格式异常,uidErr=%d,moneyErr=%d,numberErr=%d", uid, money, number)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rMax := 0.55
	list := make([]uint, number)
	leftNumber := number
	leftMoney := money
	for leftNumber > 0 {
		if leftNumber == 1 {
			list[number-1] = uint(leftMoney)
			break
		}
		if leftMoney == leftNumber {
			for i := number - leftNumber; i < number; i++ {
				list[i] = 1
			}
			break
		}
		rMoney := int(float64(leftMoney-number) * rMax)
		m := r.Intn(rMoney)
		if m < 1 {
			m = 1
		}
		list[number-leftNumber] = uint(m)
		leftNumber--
		leftMoney -= m
	}
	redId := r.Uint32()
	packageList[redId] = list
	return fmt.Sprintf("/get?id=%v&money=%v&number=%v", uid, money, number)
}

package main

/*
curl http://localhost:8080
curl http://localhost:8080/lucky
curl --data "users=a,b,c" http://localhost:8080/import
*/
import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

var userList []string

type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}
func (c *lotteryController) Get() string {
	count := len(userList)
	return fmt.Sprintf("total: %d\n", count)
}
func (c *lotteryController) PostImport() string {
	strUsers := c.Ctx.FormValue("users")
	users := strings.Split(strUsers, ",")
	count1 := len(users)
	for _, u := range users {
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			userList = append(userList, u)
		}
	}
	count2 := len(userList)
	return fmt.Sprintf("total: %d, success: %d\n", count1, count2)
}
func (c *lotteryController) GetLucky() string {
	count := len(userList)
	if count > 1 {
		seed := time.Now().UnixNano()
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))
		user := userList[index]
		userList = append(userList[0:index], userList[index+1:]...)
		return fmt.Sprintf("lucky user:%s, still have: %d\n", user, count-1)
	} else if count == 1 {
		user := userList[0]
		userList = userList[:0]
		return fmt.Sprintf("lucky user:%s, still have: %d\n", user, count-1)
	} else {
		return fmt.Sprintf("total user is 0\n")
	}
}
func main() {
	app := newApp()
	userList = []string{}
	app.Run(iris.Addr(":8080"))
}

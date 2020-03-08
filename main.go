package main

import (
	"ds-yibasuo/models"
	_ "ds-yibasuo/routers"
	"ds-yibasuo/utils/blotdb"
	"ds-yibasuo/utils/common"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
)

func main() {
	// 跨域cors
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "content-type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		AllowOrigins:     []string{"http://10.*.*.*:*", "http://localhost:*", "http://127.0.0.1:*"},
	}))

	//日志
	//err := beego.SetLogger("file", `{"filename":"/Users/finup/logs/dig-gene2.log"}`)
	//if err != nil {
	//  logs.Error("init log err: ", err)
	//  os.Exit(0)
	//}

	// 初始化
	blotdb.BlotInit()
	common.AnsibleInit()
	models.UserInit()

	// 启动
	logs.Info("The ds-yibasuo web v1.0.0")
	beego.Run()
}

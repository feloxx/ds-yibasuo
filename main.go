package main

import (
  _ "ds-yibasuo/routers"
  "ds-yibasuo/utils"
  "github.com/astaxie/beego"
  "github.com/astaxie/beego/logs"
)



func main() {
  //日志
  //err := beego.SetLogger("file", `{"filename":"/Users/finup/logs/dig-gene2.log"}`)
  //if err != nil {
  //  logs.Error("init log err: ", err)
  //  os.Exit(0)
  //}

  // 初始化
  utils.BlotInit()

  // 启动
  logs.Info("The ds-yibasuo web v1.0.0")
  beego.Run()
}
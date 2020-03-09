package controllers

import (
	"ds-yibasuo/models"
	"github.com/astaxie/beego/logs"
)

type IndexController struct {
	BaseController
}

func (c *IndexController) Get() {
	c.Data["Website"] = "ds-yibasuo"
	c.TplName = "index.html"
	_ = c.Render()
}

func (c *IndexController) CurrentState() {
	logs.Info("controller select cluster current state")
	res, err := models.SelectClusterList(-1)
	if err == nil {
		for _, value := range res.Data {
			if value.Status == "startsuccess" {
				c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: value}
				break
			}
		}
	} else if err.Error() == "null" {
		c.Data["json"] = models.Response{Code: 200, Message: "没有查到", Result: nil}
	} else {
		c.Data["json"] = models.Response{Code: 500, Message: err.Error(), Result: nil}
	}

	c.ServeJSON()
}

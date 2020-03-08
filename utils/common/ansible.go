package common

import (
	"ds-yibasuo/utils/black"
	"ds-yibasuo/utils/blotdb"
	"github.com/astaxie/beego/logs"
	"os"
	"os/exec"
	"strings"
)

func AnsibleInit() {
	logs.Info("Init ansible")
	res, err := blotdb.Db.SelectVal("init", black.String2Byte("init"))
	if err != nil {
		logs.Error("select init info err: ", err)
	}
	if len(res) < 1 {
		logs.Info("Initializing, may take several minutes, please wait...")
		init := exec.Command("/bin/bash", "-c", "sh init.sh > init.log 2>&1")
		init.Dir = "./offline"
		init.Run()

		check := exec.Command("/bin/bash", "-c", "tail -1 init.log")
		check.Dir = "./offline"
		checkOut, err := check.Output()
		if err != nil {
			logs.Error("select init check err: ", err)
		}
		checkRes := strings.Replace(black.Byte2String(checkOut), "\n", "", -1)
		if checkRes == "ok" {
			logs.Info("Init ansible ok")
			if err := blotdb.Db.Add("init", black.String2Byte("init"), black.String2Byte(Now())); err != nil {
				logs.Error("ansible init add data err: ", err)
				os.Exit(0)
			}
		} else {
			os.Exit(0)
		}
	} else {
		logs.Info("Init ansible ok")
	}
}

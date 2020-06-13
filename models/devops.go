package models

import (
	"bufio"
	"ds-yibasuo/utils/black"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	ANSIBLE_LOG = "./devops/log/ansible"
)

type DevopsInfo struct {
	ExecuteType ExecuteType `json:"executeType"`
	ClusterId   string      `json:"clusterId"`
	ExecTime    string      `json:"execTime"`
}

func (m *DevopsInfo) BackupLog(typ ExecuteType) {
	startLog := ""
	switch typ {
	case Stop:
		startLog = "stop"
	case Start:
		startLog = "start"
	case DeployUpdate:
		startLog = "deployupdate"
	}
	logs.Info("clean log 3 days ago")
	exec.Command("/bin/bash", "-c", `find ./devops/log -type f -name "ansible.*.log" -mtime +3 -exec rm -f {} \;`).Run()
	logs.Info("ansible backup log ")
	exec.Command("/bin/bash", "-c", fmt.Sprintf("mv %s.log %s.log", ANSIBLE_LOG, ANSIBLE_LOG+"."+m.ExecTime)).Run()
	logs.Info("ansible new log")
	exec.Command("/bin/bash", "-c", fmt.Sprintf(`echo "%s" > %s`, startLog, ANSIBLE_LOG+".log")).Run()
}

func (m *DevopsInfo) DeployUpdate() {
	cmd := exec.Command("ansible-playbook", "deploy.yml")
	//cmd := exec.Command("ansible-playbook", "deploy.yml", "-t", "dancer,music,show")
	cmd.Dir = "./devops"
	cmd.Start()
}

// 这样的异步操作，显得易懂又新手。
func (m *DevopsInfo) UpdateStatus(exe ExecuteType, cluster *ClusterInfo) {
	for {
		// 检查ansible是否执行结束
		ansible, _ := exec.Command("/bin/bash", "-c", "ps -ef | grep ansible-playbook | grep -v grep | wc -l").Output()
		ansibleCount, _ := strconv.Atoi(strings.Replace(black.Byte2String(ansible), "\n", "", -1))

		// 检查执行中是否出现了错误
		failed, _ := exec.Command("/bin/bash", "-c", fmt.Sprintf(`grep -E "failed=[1-9]" %s.log | wc -l`, ANSIBLE_LOG)).Output()
		failedCount, _ := strconv.Atoi(strings.Replace(black.Byte2String(failed), "\n", "", -1))

		if ansibleCount == 0 {
			// ansible 进程为0 证明结束了
			switch exe {
			case Start:
				if failedCount == 0 {
					cluster.Status = "startsuccess"
				} else {
					cluster.Status = "startfailure"
				}
			case Stop:
				if failedCount == 0 {
					cluster.Status = "stopsuccess"
				} else {
					cluster.Status = "stopfailure"
				}
			case DeployUpdate:
				if failedCount == 0 {
					cluster.Status = "deploysuccess"
				} else {
					cluster.Status = "deployfailure"
				}
			}
			if err := cluster.UpdateCluster(); err != nil {
				logs.Error("status change err: ", err)
			}
			break
		} else {
			// 没结束，暂停1秒再来一次
			time.Sleep(1 * time.Second)
		}
	}
}

func (m *DevopsInfo) ReadLog(start, end int) ([]string, error) {
	filePath := "./devops/log/ansible.log"
	var result []string

	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return nil, err
	}
	defer file.Close()

	buf := bufio.NewReader(file)

	i := 1
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if i >= start && i <= end {
			result = append(result, line)
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return nil, err
			}
		}

		if i == end {
			return result, nil
		}
		i += 1
	}
	return result, nil
}

func (m *DevopsInfo) GetLogRows() (int, error) {
	fuck := fmt.Sprintf("wc -l %s | awk '{print $1}'", ANSIBLE_LOG)

	cmd := exec.Command("/bin/bash", "-c", fuck)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	fuckRows, _ := strconv.Atoi(strings.Replace(black.Byte2String(out), "\n", "", -1))
	return fuckRows, nil
}

type DevopsLogResult struct {
	CurrentPage int      `json:"currentPage"`
	Rows        int      `json:"rows"`
	Data        []string `json:"data"`
}

type SignalResult struct {
	SingnalType string `json:"singnalType"`
	Success     bool   `json:"success"`
	Over        bool   `json:"over"`
}

// TODO 很多重复逻辑 隐患
func (m *DevopsInfo) GetSignal() (*SignalResult, error) {
	// 获取信号相关的数据
	out, err := exec.Command("/bin/bash", "-c", fmt.Sprintf(`head -1 %s.log`, ANSIBLE_LOG)).Output()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("/bin/bash", "-c", "ps -ef | grep ansible-playbook | grep -v grep | wc -l")
	bytes, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	ansibleCount, err := strconv.Atoi(strings.Replace(black.Byte2String(bytes), "\n", "", -1))
	if err != nil {
		return nil, err
	}
	fuckSuccess, err := exec.Command("/bin/bash", "-c", fmt.Sprintf(`grep -E "failed=[1-9]" %s.log | wc -l`, ANSIBLE_LOG)).Output()
	if err != nil {
		return nil, err
	}
	fuckSuccessCount, err := strconv.Atoi(strings.Replace(black.Byte2String(fuckSuccess), "\n", "", -1))
	if err != nil {
		return nil, err
	}
	// 判断数据
	if ansibleCount == 0 && fuckSuccessCount == 0 {
		return &SignalResult{
			SingnalType: strings.Replace(black.Byte2String(out), "\n", "", -1),
			Success:     true,
			Over:        true,
		}, nil
	} else if ansibleCount == 0 && fuckSuccessCount != 0 {
		return &SignalResult{
			SingnalType: strings.Replace(black.Byte2String(out), "\n", "", -1),
			Success:     false,
			Over:        true,
		}, nil
	} else {
		return &SignalResult{
			SingnalType: strings.Replace(black.Byte2String(out), "\n", "", -1),
			Success:     false,
			Over:        false,
		}, nil
	}
}

func (m *DevopsInfo) Start() {
	cmd := exec.Command("ansible-playbook", "start.yml")
	cmd.Dir = "./devops"
	cmd.Start()
}

func (m *DevopsInfo) Stop() {
	cmd := exec.Command("ansible-playbook", "stop.yml")
	cmd.Dir = "./devops"
	cmd.Start()
}

func (m *DevopsInfo) RefreshHost(pwd string) error {
	logs.Info("ansible refresh host")
	refresh := "ansible-playbook -i hosts.ini create_users.yml -u root"
	cmd := exec.Command("/bin/bash", "-c", refresh)
	cmd.Dir = "./devops"
	out, err := cmd.Output()
	if strings.Contains(black.Byte2String(out), "Permission denied") {
		return errors.New("密码错误！！")
	}
	if err != nil {
		return err
	}
	return nil
}

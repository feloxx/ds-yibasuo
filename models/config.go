package models

import (
	"ds-yibasuo/utils/black"
	"ds-yibasuo/utils/blotdb"
	"ds-yibasuo/utils/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
)

var (
	ConfigList  = []string{"frontend", "backend", "alert", "master", "worker", "database", "zookeeper", "hadoop", "common"}
	FrontendYml = "./devops/conf/nginx.yml"
	BackendYml  = "./devops/conf/api.yml"
	AlertYml    = "./devops/conf/alert.yml"
	MasterYml   = "./devops/conf/master.yml"
	WorkerYml   = "./devops/conf/worker.yml"
	HadoopYml   = "./devops/conf/hadoop.yml"
	CommonYml   = "./devops/conf/common.yml"
)

// 配置内容接口
type ConfigBody interface{ ConfigToString() string }

// 前端配置
type ConfigFrontend struct {
	Frontend *Frontends `yaml:"nginx" json:"frontend"`
}

func (m *ConfigFrontend) ConfigToString() string { return "" }

type Frontends struct {
	EscProxyPort int `yaml:"esc_proxy_port" json:"escProxyPort"`
}

func (m *Frontends) ConfigToString() string { return "" }

// 后端配置
type ConfigBackend struct {
	Backend *Backends `yaml:"api" json:"backend"`
}

func (m *ConfigBackend) ConfigToString() string { return "" }

type Backends struct {
	ServerPort int `yaml:"server.port" json:"serverPort"`
}

func (m *Backends) ConfigToString() string { return "" }

// 报警配置
//type AlertType int
//
//const (
//	Mail = iota
//	Weixin
//)

type AlertBody interface{ AlertToString() string }

type AlertMail struct {
	AlertType    string `yaml:"alert.type" json:"alertType"`
	MailProtocol string `yaml:"mail.protocol" json:"mailProtocol"`
	MailSender   string `yaml:"mail.sender" json:"mailSender"`
	MailUser     string `yaml:"mail.user" json:"mailUser"`
	MailPasswd   string `yaml:"mail.passwd" json:"mailPasswd"`
	MailSmtpHost string `yaml:"mail.server.host" json:"mailSmtpHost"`
	MailSmtpPort int    `yaml:"mail.server.port" json:"mailSmtpPort"`
	MailSsl      bool   `yaml:"mail.smtp.ssl.enable" json:"mailSsl"`
	MailTls      bool   `yaml:"mail.smtp.starttls.enable" json:"mailTls"`
}

func (m *AlertMail) AlertToString() string {
	return fmt.Sprintf(`AlertMail:{}`)
}

type AlertWeixin struct {
	//TODO 后面在写
}

func (m *AlertWeixin) AlertToString() string {
	return fmt.Sprintf(`AlertWeixin:{}`)
}

//序列化的时候为什么不认识这个类型呢AlertBody
//错误 * 'Alert' expected type 'models.AlertBody', got 'map[string]interface {}'
type ConfigAlert struct {
	Alert interface{} `yaml:"alert" json:"alert"`
}

func (m *ConfigAlert) ConfigToString() string { return "" }

// Master配置
type ConfigMaster struct {
	Master *Masters `yaml:"master" json:"master"`
}

func (m *ConfigMaster) ConfigToString() string { return fmt.Sprintf(`ConfigMaster:{}`) }

type Masters struct {
	MasterReservedMemory float32 `yaml:"master.reserved.memory" json:"masterReservedMemory"`
}

func (m *Masters) ConfigToString() string { return "" }

// Worker配置
type ConfigWorker struct {
	Worker *Workers `yaml:"worker" json:"worker"`
}

func (m *ConfigWorker) ConfigToString() string { return fmt.Sprintf(`ConfigWorker:{}`) }

type Workers struct {
	WorkerReservedMemory float32 `yaml:"worker.reserved.memory" json:"workerReservedMemory"`
}

func (m *Workers) ConfigToString() string { return "" }

// zookeeper配置
// TODO 第一版不用配zk参数
type ConfigZookeeper struct {
	Zookeeper []string `json:"zookeeper"`
}

func (m *ConfigZookeeper) ConfigToString() string { return "" }

// 资源中心配置
type Hadoops struct {
	FsDefaultFS string `yaml:"fs.defaultFS" json:"fsDefaultFS"`
	//TODO 其他字段后面写
}

func (m *Hadoops) ConfigToString() string { return "" }

type ConfigHadoop struct {
	Hadoop *Hadoops `yaml:"hadoop" json:"hadoop"`
}

func (m *ConfigHadoop) ConfigToString() string { return "" }

type Commons struct {
	DataStore2hdfsBasepath string `yaml:"data.store2hdfs.basepath" json:"dataStore2hdfsBasepath"`
	//TODO 其他字段后面写
}

func (m *Commons) ConfigToString() string { return "" }

type ConfigCommon struct {
	Common *Commons `yaml:"common" json:"common"`
}

func (m *ConfigCommon) ConfigToString() string { return "" }

// database配置
type ConfigDatabase struct {
	DatabaseType string `json:"databaseType"`
	DatabaseName string `json:"databaseName"`
	Account      string `json:"account"`
	Password     string `json:"password"`
}

func (m *ConfigDatabase) ConfigToString() string {
	return fmt.Sprintf(`ConfigDB:{}`)
}

// 配置信息对接前端
// TODO 没办法了,先这样凑合了
type ConfigInfoPuppet struct {
	Id                     string  `json:"id"`
	Name                   string  `json:"name"`
	Typ                    string  `json:"typ"`
	EscProxyPort           int     `json:"escProxyPort,omitempty"`
	ServerPort             int     `json:"serverPort,omitempty"`
	AlertType              string  `json:"alertType,omitempty"`
	MailProtocol           string  `json:"mailProtocol,omitempty"`
	MailSender             string  `json:"mailSender,omitempty"`
	MailUser               string  `json:"mailUser,omitempty"`
	MailPasswd             string  `json:"mailPasswd,omitempty"`
	MailSmtpHost           string  `json:"mailServerHost,omitempty"`
	MailSmtpPort           int     `json:"mailServerPort,omitempty"`
	MailSsl                bool    `json:"mailSmtpSslEnable,omitempty"`
	MailTls                bool    `json:"mailSmtpStarttlsEnable,omitempty"`
	MasterReservedMemory   float32 `json:"masterReservedMemory,omitempty"`
	WorkerReservedMemory   float32 `json:"workerReservedMemory,omitempty"`
	FsDefaultFS            string  `json:"fsDefaultFS,omitempty"`
	DataStore2hdfsBasepath string  `json:"dataStore2hdfsBasepath,omitempty"`
	DatabaseType           string  `json:"databaseType,omitempty"`
	DatabaseName           string  `json:"databaseName,omitempty"`
	Account                string  `json:"account,omitempty"`
	Password               string  `json:"password,omitempty"`
	Remark                 string  `json:"remark"`
}

// 配置信息
type ConfigType int

const (
	Frontend = iota
	Backend
	Alert
	Master
	Worker
	Zookeeper
	Database
	Hadoop
	Common
)

type ConfigInfo struct {
	Id     string      `json:"id"`   // 配置id
	Name   string      `json:"name"` // 配置别名
	Typ    string      `json:"typ"`  // 配置类型
	Conf   interface{} `json:"conf"` // 配置内容 //TODO interface 隐患
	Remark string      `json:"remark"`
}

func (m *ConfigInfo) CreateConfig() error {
	m.Id = common.MakeUuid(m.Name + m.Typ + common.Now())
	hostBody, _ := json.Marshal(m)
	return blotdb.Db.Add("config", black.String2Byte(m.Id), hostBody)
}

func (m *ConfigInfo) DeleteConfig() error {
	return blotdb.Db.RemoveID("config", black.String2Byte(m.Id))
}

func (m *ConfigInfo) UpdateConfig() error {
	hostBody, _ := json.Marshal(m)
	return blotdb.Db.Update("config", black.String2Byte(m.Id), hostBody)
}

func (m *ConfigInfo) SelectConfig() (*ConfigInfoPuppet, error) {
	res, err := blotdb.Db.SelectVal("config", black.String2Byte(m.Id))
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, errors.New("null")
	}

	var object json.RawMessage
	h := ConfigInfo{Conf: &object}
	if err := json.Unmarshal(black.String2Byte(res[0]), &h); err != nil {
		logs.Error(err)
		return nil, err
	}
	var result ConfigInfoPuppet
	result.Id = h.Id
	result.Name = h.Name
	result.Typ = h.Typ
	result.Remark = h.Remark
	switch h.Typ {
	case "frontend":
		var n ConfigFrontend
		if err := json.Unmarshal(object, &n); err != nil {
			logs.Error(err)
			return nil, err
		}
		result.EscProxyPort = n.Frontend.EscProxyPort
	case "backend":
		var n ConfigBackend
		if err := json.Unmarshal(object, &n); err != nil {
			logs.Error(err)
			return nil, err
		}
		result.ServerPort = n.Backend.ServerPort
	case "alert":
		var insideObject json.RawMessage
		n := ConfigAlert{Alert: &insideObject}
		if err := json.Unmarshal(object, &n); err != nil {
			logs.Error(err)
			return nil, err
		}
		var alert AlertMail
		if err := json.Unmarshal(insideObject, &alert); err != nil {
			logs.Error(err)
			return nil, err
		}
		result.AlertType = alert.AlertType
		result.MailProtocol = alert.MailProtocol
		result.MailSender = alert.MailSender
		result.MailUser = alert.MailUser
		result.MailPasswd = alert.MailPasswd
		result.MailSmtpHost = alert.MailSmtpHost
		result.MailSmtpPort = alert.MailSmtpPort
		result.MailSsl = alert.MailSsl
		result.MailTls = alert.MailTls
	case "master":
		var n ConfigMaster
		if err := json.Unmarshal(object, &n); err != nil {
			logs.Error(err)
			return nil, err
		}
		result.MasterReservedMemory = n.Master.MasterReservedMemory
	case "worker":
		var n ConfigWorker
		if err := json.Unmarshal(object, &n); err != nil {
			logs.Error(err)
			return nil, err
		}
		result.WorkerReservedMemory = n.Worker.WorkerReservedMemory
	case "hadoop":
		var n ConfigHadoop
		if err := json.Unmarshal(object, &n); err != nil {
			logs.Error(err)
			return nil, err
		}
		result.FsDefaultFS = n.Hadoop.FsDefaultFS
	case "common":
		var n ConfigCommon
		if err := json.Unmarshal(object, &n); err != nil {
			logs.Error(err)
			return nil, err
		}
		result.DataStore2hdfsBasepath = n.Common.DataStore2hdfsBasepath
	case "database":
		var n ConfigDatabase
		if err := json.Unmarshal(object, &n); err != nil {
			logs.Error(err)
			return nil, err
		}
		result.DatabaseType = n.DatabaseType
		result.DatabaseName = n.DatabaseName
		result.Account = n.Account
		result.Password = n.Password
	}

	return &result, err
}

type ConfigInfoResult struct {
	CurrentPage int                 `json:"currentPage"`
	Total       int                 `json:"total"`
	Data        []*ConfigInfoPuppet `json:"data"`
}

func SelectConfigList(page int, typ string) (*ConfigInfoResult, error) {
	res, err := blotdb.Db.SelectValues("config")
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, errors.New("null")
	}

	var fuck []*ConfigInfoPuppet
	for _, value := range res {
		var object json.RawMessage
		h := ConfigInfo{Conf: &object}
		if err := json.Unmarshal(value, &h); err != nil {
			logs.Error(err)
			return nil, err
		}
		if h.Typ == typ {
			var result ConfigInfoPuppet
			result.Id = h.Id
			result.Name = h.Name
			result.Typ = h.Typ
			result.Remark = h.Remark
			switch h.Typ {
			case "frontend":
				var n ConfigFrontend
				if err := json.Unmarshal(object, &n); err != nil {
					logs.Error(err)
					return nil, err
				}
				result.EscProxyPort = n.Frontend.EscProxyPort
			case "backend":
				var n ConfigBackend
				if err := json.Unmarshal(object, &n); err != nil {
					logs.Error(err)
					return nil, err
				}
				result.ServerPort = n.Backend.ServerPort
			case "alert":
				var insideObject json.RawMessage
				n := ConfigAlert{Alert: &insideObject}
				if err := json.Unmarshal(object, &n); err != nil {
					logs.Error(err)
					return nil, err
				}
				var alert AlertMail
				if err := json.Unmarshal(insideObject, &alert); err != nil {
					logs.Error(err)
					return nil, err
				}
				result.AlertType = alert.AlertType
				result.MailProtocol = alert.MailProtocol
				result.MailSender = alert.MailSender
				result.MailUser = alert.MailUser
				result.MailPasswd = alert.MailPasswd
				result.MailSmtpHost = alert.MailSmtpHost
				result.MailSmtpPort = alert.MailSmtpPort
				result.MailSsl = alert.MailSsl
				result.MailTls = alert.MailTls
			case "master":
				var n ConfigMaster
				if err := json.Unmarshal(object, &n); err != nil {
					logs.Error(err)
					return nil, err
				}
				result.MasterReservedMemory = n.Master.MasterReservedMemory
			case "worker":
				var n ConfigWorker
				if err := json.Unmarshal(object, &n); err != nil {
					logs.Error(err)
					return nil, err
				}
				result.WorkerReservedMemory = n.Worker.WorkerReservedMemory
			case "hadoop":
				var n ConfigHadoop
				if err := json.Unmarshal(object, &n); err != nil {
					logs.Error(err)
					return nil, err
				}
				result.FsDefaultFS = n.Hadoop.FsDefaultFS
			case "common":
				var n ConfigCommon
				if err := json.Unmarshal(object, &n); err != nil {
					logs.Error(err)
					return nil, err
				}
				result.DataStore2hdfsBasepath = n.Common.DataStore2hdfsBasepath
			case "database":
				var n ConfigDatabase
				if err := json.Unmarshal(object, &n); err != nil {
					logs.Error(err)
					return nil, err
				}
				result.DatabaseType = n.DatabaseType
				result.DatabaseName = n.DatabaseName
				result.Account = n.Account
				result.Password = n.Password
			}
			fuck = append(fuck, &result)
		}
	}
	if len(fuck) == 0 {
		return nil, errors.New("null")
	}

	fucks := slidingConfig(fuck, 10)

	var fuckOff []*ConfigInfoPuppet
	if len(fucks) <= page {
		fuckOff = fucks[len(fucks)-1]
	} else {
		fuckOff = fucks[page-1]
	}

	result := &ConfigInfoResult{
		CurrentPage: page,
		Total:       len(fuck),
		Data:        fuckOff,
	}

	return result, nil
}

type AllConfig struct {
	Frontend []*ConfigInfoPuppet `json:"frontend"`
	Backend  []*ConfigInfoPuppet `json:"backend"`
	Alert    []*ConfigInfoPuppet `json:"alert"`
	Master   []*ConfigInfoPuppet `json:"master"`
	Worker   []*ConfigInfoPuppet `json:"worker"`
	Database []*ConfigInfoPuppet `json:"database"`
	Hadoop   []*ConfigInfoPuppet `json:"hadoop"`
	Common   []*ConfigInfoPuppet `json:"common"`
}

func SelectAllConfig() (*AllConfig, error) {
	res, err := blotdb.Db.SelectValues("config")
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, errors.New("null")
	}

	all := &AllConfig{}
	for _, value := range res {
		var object json.RawMessage
		h := ConfigInfo{Conf: &object}
		if err := json.Unmarshal(value, &h); err != nil {
			logs.Error(err)
			return nil, err
		}
		var result ConfigInfoPuppet
		result.Id = h.Id
		result.Name = h.Name
		result.Typ = h.Typ
		result.Remark = h.Remark
		switch h.Typ {
		case "frontend":
			var n ConfigFrontend
			if err := json.Unmarshal(object, &n); err != nil {
				logs.Error(err)
				return nil, err
			}
			result.EscProxyPort = n.Frontend.EscProxyPort
		case "backend":
			var n ConfigBackend
			if err := json.Unmarshal(object, &n); err != nil {
				logs.Error(err)
				return nil, err
			}
			result.ServerPort = n.Backend.ServerPort
		case "alert":
			var insideObject json.RawMessage
			n := ConfigAlert{Alert: &insideObject}
			if err := json.Unmarshal(object, &n); err != nil {
				logs.Error(err)
				return nil, err
			}
			var alert AlertMail
			if err := json.Unmarshal(insideObject, &alert); err != nil {
				logs.Error(err)
				return nil, err
			}
			result.AlertType = alert.AlertType
			result.MailProtocol = alert.MailProtocol
			result.MailSender = alert.MailSender
			result.MailUser = alert.MailUser
			result.MailPasswd = alert.MailPasswd
			result.MailSmtpHost = alert.MailSmtpHost
			result.MailSmtpPort = alert.MailSmtpPort
			result.MailSsl = alert.MailSsl
			result.MailTls = alert.MailTls
		case "master":
			var n ConfigMaster
			if err := json.Unmarshal(object, &n); err != nil {
				logs.Error(err)
				return nil, err
			}
			result.MasterReservedMemory = n.Master.MasterReservedMemory
		case "worker":
			var n ConfigWorker
			if err := json.Unmarshal(object, &n); err != nil {
				logs.Error(err)
				return nil, err
			}
			result.WorkerReservedMemory = n.Worker.WorkerReservedMemory
		case "hadoop":
			var n ConfigHadoop
			if err := json.Unmarshal(object, &n); err != nil {
				logs.Error(err)
				return nil, err
			}
			result.FsDefaultFS = n.Hadoop.FsDefaultFS
		case "common":
			var n ConfigCommon
			if err := json.Unmarshal(object, &n); err != nil {
				logs.Error(err)
				return nil, err
			}
			result.DataStore2hdfsBasepath = n.Common.DataStore2hdfsBasepath
		case "database":
			var n ConfigDatabase
			if err := json.Unmarshal(object, &n); err != nil {
				logs.Error(err)
				return nil, err
			}
			result.DatabaseType = n.DatabaseType
			result.DatabaseName = n.DatabaseName
			result.Account = n.Account
			result.Password = n.Password
		}
		switch h.Typ {
		case "frontend":
			all.Frontend = append(all.Frontend, &result)
		case "backend":
			all.Backend = append(all.Backend, &result)
		case "alert":
			all.Alert = append(all.Alert, &result)
		case "master":
			all.Master = append(all.Master, &result)
		case "worker":
			all.Worker = append(all.Worker, &result)
		case "database":
			all.Database = append(all.Database, &result)
		case "hadoop":
			all.Hadoop = append(all.Hadoop, &result)
		case "common":
			all.Common = append(all.Common, &result)
		}
	}

	return all, nil
}

func (m *ConfigInfo) CheckName() (bool, error) {
	res, err := blotdb.Db.SelectValues("config")
	if err != nil || len(res) < 1 {
		return false, errors.New("查询错误 或者 没有内容！")
	}

	for _, value := range res {
		h := ConfigInfo{}
		if err := json.Unmarshal(value, &h); err != nil {
			logs.Error(err)
			return false, err
		}
		if h.Typ == m.Typ && h.Name == m.Name && h.Id != m.Id {
			return true, nil
		}
	}

	return false, nil
}

func slidingConfig(list []*ConfigInfoPuppet, step int) (res [][]*ConfigInfoPuppet) {
	start, end := 0, 0
	for {
		if len(list) <= 0 {
			break
		}
		if (start + step) > len(list) {
			end = len(list)
		} else {
			end += step
		}
		res = append(res, list[start:end])
		start += step
		if start > len(list) {
			break
		}
	}
	return
}

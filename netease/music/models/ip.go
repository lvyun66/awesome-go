package models

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/lunny/log"
	"github.com/lvyun66/awesome-go/netease/conf"
	"time"
)

type Ip struct {
	Id    int64  `xorm:"pk autoincr BIGINT(20)"`
	Data  string `xorm:"not null VARCHAR(255)"`
	Type1 string `xorm:"not null VARCHAR(255)"`
	Type2 string `xorm:"VARCHAR(255)"`
	Speed int64  `xorm:"not null BIGINT(20)"`
}

func defaultProxyPoolX() *xorm.Engine {
	var my = conf.DefaultConf.Services.Mysql
	var dataSourceName = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8", my.User, my.Password, my.Host, my.Port, "proxy_pool")
	var x, _ = xorm.NewEngine("mysql", dataSourceName)
	return x
}

func DeleteIP(ip *Ip) bool {
	if _, err := defaultProxyPoolX().Delete(ip); err != nil {
		return false
	}
	log.Println("[proxy_pool][info]", time.Now().Format("2006-01-02T15:04:05Z07:00"), "delete invalid proxy: ", ip.Data)
	return true
}

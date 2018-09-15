package music

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/lvyun66/awesome-go/netease/conf"
	"github.com/lvyun66/awesome-go/netease/music/basetool"
	"github.com/lvyun66/awesome-go/netease/music/models"
	"log"
	"strconv"
	"sync"
	"time"
)

type FanRequest struct {
	UserId    string `json:"userId"`
	Offset    string `json:"offset"`
	Limit     int    `json:"limit"`
	CsrfToken string `json:"csrf_token"`
}

type FanResponse struct {
	Code      int  `json:"code"`
	More      bool `json:"more"`
	Followeds []struct {
		Py            string `json:"py"`
		Time          int64  `json:"time"`
		UserType      int    `json:"userType"`
		ExpertTags    string `json:"expertTags"`
		AuthStatus    int    `json:"authStatus"`
		Followed      bool   `json:"followed"`
		Experts       string `json:"experts"`
		Followeds     int    `json:"followeds"`
		VipType       int    `json:"vipType"`
		Gender        int    `json:"gender"`
		AccountStatus int    `json:"accountStatus"`
		AvatarURL     string `json:"avatarUrl"`
		Nickname      string `json:"nickname"`
		RemarkName    string `json:"remarkName"`
		Follows       int    `json:"follows"`
		Mutual        bool   `json:"mutual"`
		UserID        int    `json:"userId"`
		Signature     string `json:"signature"`
		EventCount    int    `json:"eventCount"`
		PlaylistCount int    `json:"playlistCount"`
	} `json:"followeds"`
}

func Fans() {
	// init xorm
	my := conf.DefaultConf.Services.Mysql
	var dataSource = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8", my.User, my.Password, my.Host, my.Port, "netease")
	engine, err := xorm.NewEngine("mysql", dataSource)
	if err != nil {
		log.Fatalln("Connect mysql error:", err)
	}
	if err := engine.Ping(); err != nil {
		log.Fatalln("Mysql ping error:", err)
	}
	engine.ShowSQL(true)
	engine.SetTableMapper(core.SnakeMapper{})

	var userId = 48353
	var processCount = 10

	wg := &sync.WaitGroup{}
	for i := 0; i < processCount; i++ {
		wg.Add(1)
		time.Sleep(time.Second * 10)
		go func(userId, c, i int) {
			var limit = 20
			var offset = limit * i
			for {
				fanRequest := &FanRequest{
					UserId:    strconv.Itoa(userId),
					Offset:    strconv.Itoa(offset),
					Limit:     limit,
					CsrfToken: "",
				}
				_params, _ := json.Marshal(fanRequest)
				params, encSecKey, encErr := basetool.EncryptParams(string(_params))
				if encErr != nil {
					log.Fatal(encErr)
				}
				var response string
				var err error
				var retryCount = 1
				for retryCount <= 5 {
					url := "https://music.163.com/weapi/user/getfolloweds?csrf_token="
					response, err = basetool.Post(url, params, encSecKey)
					if err != nil {
						fmt.Printf("[CC][retry %d] Get user followed error: %s\n", retryCount, err)
					} else {
						break
					}
					retryCount += 1
					time.Sleep(time.Second * 2)
				}

				fans := &FanResponse{}
				json.Unmarshal([]byte(response), fans)
				for _, value := range fans.Followeds {
					musicUserFan := &models.MusicUserFans{}
					if isExist, _ := engine.Id(value.UserID).Get(musicUserFan); !isExist {
						musicUserFan.UserId = value.UserID
						musicUserFan.UserType = value.UserType
						musicUserFan.NikeName = value.Nickname
						musicUserFan.Time = value.Time
						musicUserFan.Py = value.Py
						musicUserFan.ExpertTags = value.ExpertTags
						musicUserFan.AuthStatus = value.AuthStatus
						if value.Followed {
							musicUserFan.Followed = 1
						} else {
							musicUserFan.Followed = 0
						}
						musicUserFan.VipType = value.VipType
						musicUserFan.Gender = value.Gender
						musicUserFan.AccountStatus = value.AccountStatus
						musicUserFan.AvatarUrl = value.AvatarURL
						musicUserFan.RemarkName = value.RemarkName
						musicUserFan.Follows = value.Follows
						if value.Followed {
							musicUserFan.Mutual = 1
						} else {
							musicUserFan.Mutual = 0
						}
						musicUserFan.Signature = value.Signature
						musicUserFan.EventCount = value.EventCount
						musicUserFan.PlaylistCount = value.PlaylistCount
						_, err := engine.Insert(musicUserFan)
						if err != nil {
							fmt.Println("Insert row error:", err, value)
						}
					}
				}
				offset += limit * processCount
				if fans.More == false {
					break
				}
				time.Sleep(time.Second * 3)
			}
			wg.Done()
		}(userId, processCount, i)
	}
	wg.Wait()
}

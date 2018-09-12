package music

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/lvyun66/awesome-go/netease/music/models"
	"log"
	"strconv"
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
	engine, err := xorm.NewEngine("mysql", "root:firely0506@("+conf.Services.Mysql.Host+":3306)/netease?charset=utf8")
	if err != nil {
		log.Fatalln("Connect mysql error:", err)
	}
	if err := engine.Ping(); err != nil {
		log.Fatalln("Mysql ping error:", err)
	}
	engine.ShowSQL(false)
	engine.SetTableMapper(core.SnakeMapper{})

	var offset = 0
	var limit = 20
	var userId = 48353
	for {
		fanRequest := &FanRequest{
			UserId:    strconv.Itoa(userId),
			Offset:    strconv.Itoa(offset),
			Limit:     limit,
			CsrfToken: "",
		}
		_params, _ := json.Marshal(fanRequest)
		params, encSecKey, encErr := EncryptParams(string(_params))
		if encErr != nil {
			log.Fatal(encErr)
		}
		var response string
		var err error
		var retryCount = 1
		for retryCount <= 5 {
			response, err = Post("https://music.163.com/weapi/user/getfolloweds?csrf_token=", params, encSecKey)
			if err != nil {
				fmt.Printf("[retry %d] Get user followed error: %s\n", retryCount, err)
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
				//fmt.Println("Insert result:", affected)
				if err != nil {
					fmt.Println("Insert row error:", err, value)
				}
			}
		}
		offset += limit
		if fans.More == false {
			break
		}
		time.Sleep(time.Second)
	}
}

package music

import (
	"encoding/json"
	"fmt"
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
		Py            string      `json:"py"`
		Time          int64       `json:"time"`
		UserType      int         `json:"userType"`
		ExpertTags    interface{} `json:"expertTags"`
		AuthStatus    int         `json:"authStatus"`
		Followed      bool        `json:"followed"`
		Experts       interface{} `json:"experts"`
		Followeds     int         `json:"followeds"`
		VipType       int         `json:"vipType"`
		Gender        int         `json:"gender"`
		AccountStatus int         `json:"accountStatus"`
		AvatarURL     string      `json:"avatarUrl"`
		Nickname      string      `json:"nickname"`
		RemarkName    interface{} `json:"remarkName"`
		Follows       int         `json:"follows"`
		Mutual        bool        `json:"mutual"`
		UserID        int         `json:"userId"`
		Signature     interface{} `json:"signature"`
		EventCount    int         `json:"eventCount"`
		PlaylistCount int         `json:"playlistCount"`
	} `json:"followeds"`
}

func Fans() {
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
		for key, value := range fans.Followeds {
			fmt.Println(key, value)
		}
		offset += limit
		if fans.More == false {
			break
		}
		time.Sleep(time.Second)
	}
}

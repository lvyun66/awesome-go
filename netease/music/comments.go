package music

import (
	"encoding/json"
	"fmt"
	"github.com/lvyun66/awesome-go/netease/conf"
	"github.com/lvyun66/awesome-go/netease/music/basetool"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"sync"
	"time"
)

type SongCommentsParams struct {
	Rid       string `json:"rid"`
	Offset    string `json:"offset"`
	Limit     int    `json:"limit"`
	Total     bool   `json:"total"`
	CsrfToken string `json:"csrf_token"`
}

type SongComments struct {
	IsMusician  bool          `json:"isMusician"`
	UserID      int           `json:"userId"`
	TopComments []interface{} `json:"topComments"`
	MoreHot     bool          `json:"moreHot"`
	HotComments []struct {
		User struct {
			LocationInfo interface{} `json:"locationInfo"`
			ExpertTags   interface{} `json:"expertTags"`
			UserID       int         `json:"userId"`
			RemarkName   interface{} `json:"remarkName"`
			AvatarURL    string      `json:"avatarUrl"`
			Experts      interface{} `json:"experts"`
			VipType      int         `json:"vipType"`
			Nickname     string      `json:"nickname"`
			VipRights    interface{} `json:"vipRights"`
			UserType     int         `json:"userType"`
			AuthStatus   int         `json:"authStatus"`
		} `json:"user"`
		SongId        int           `json:"song_id"`
		BeReplied     []interface{} `json:"beReplied"`
		PendantData   interface{}   `json:"pendantData"`
		ExpressionURL interface{}   `json:"expressionUrl"`
		Liked         bool          `json:"liked"`
		LikedCount    int           `json:"likedCount"`
		CommentID     int           `json:"commentId"`
		Time          int64         `json:"time"`
		Content       string        `json:"content"`
	} `json:"hotComments"`
	Code     int       `json:"code"`
	Comments []Comment `json:"comments"`
	Total    int       `json:"total"`
	More     bool      `json:"more"`
}

type Comment struct {
	User               User          `json:"user"`
	SongId             int           `json:"song_id"`
	BeReplied          []interface{} `json:"beReplied"`
	PendantData        interface{}   `json:"pendantData"`
	ExpressionURL      interface{}   `json:"expressionUrl"`
	Liked              bool          `json:"liked"`
	LikedCount         int           `json:"likedCount"`
	CommentID          int           `json:"commentId"`
	Time               int64         `json:"time"`
	Content            string        `json:"content"`
	IsRemoveHotComment bool          `json:"isRemoveHotComment"`
}

type User struct {
	LocationInfo interface{} `json:"locationInfo"`
	ExpertTags   interface{} `json:"expertTags"`
	UserID       int         `json:"userId"`
	RemarkName   interface{} `json:"remarkName"`
	AvatarURL    string      `json:"avatarUrl"`
	Experts      interface{} `json:"experts"`
	VipType      int         `json:"vipType"`
	Nickname     string      `json:"nickname"`
	VipRights    interface{} `json:"vipRights"`
	UserType     int         `json:"userType"`
	AuthStatus   int         `json:"authStatus"`
}

func GetSongComments(songId int) bool {
	var limit = 100
	var rid = strconv.Itoa(songId)
	var processCount = 5

	go func() {
		for {
			mc := NewMongoComments()
			c, err := mc.Find(nil).Count()
			if err != nil {
				log.Println("[COUNT] get comments sum err:", err)
			}
			log.Println("[COUNT] total:", c)
			time.Sleep(time.Second * 30)
		}
	}()

	mongoComments := NewMongoComments()
	wg := &sync.WaitGroup{}
	for i := 0; i < processCount; i++ {
		wg.Add(1)
		go func(i int) {
			log.Println("goroutine", i, "is started")
			var offset = limit * i + 55000
			var total = false
			if offset == 0 {
				total = true
			}
			for {
				params := SongCommentsParams{
					Rid:       rid,
					Limit:     limit,
					Offset:    strconv.Itoa(offset),
					Total:     total,
					CsrfToken: "",
				}
				_params, _ := json.Marshal(params)
				encParams, encSecKey, err := basetool.EncryptParams(string(_params))
				if err != nil {
					log.Fatal(err)
				}

				var str string
				var retryCount = 1
				for retryCount <= 5 {
					_url := "https://music.163.com/weapi/v1/resource/comments/R_SO_4_" + rid + "?csrf_token="
					str, err = basetool.Post(_url, encParams, encSecKey)
					if err != nil {
						fmt.Printf("[CC][retry %d] Get user followed error: %s\n", retryCount, err)
					} else {
						break
					}
					retryCount += 1
					time.Sleep(time.Second)
				}

				response := &SongComments{}
				json.Unmarshal([]byte(str), response)
				if !response.More {
					break
				}
				for _, comment := range response.Comments {
					var c Comment
					if err := mongoComments.Find(bson.M{"commentid": comment.CommentID}).One(&c); err != nil {
						if err.Error() == mgo.ErrNotFound.Error() {
							comment.SongId = songId
							mongoComments.Insert(comment)
							log.Println("[INSERT] a new comment insert mongo:", comment.CommentID)
						}
					} else {
						log.Println("[EXIST] comment is exist in mongo:", comment.CommentID)
					}
				}
				offset += limit * processCount
				time.Sleep(time.Second * 2)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return true
}

func NewMongoComments() *mgo.Collection {
	dialInfo := &mgo.DialInfo{
		Addrs:     []string{conf.DefaultConf.Services.Mongo.Url},
		Source:    conf.DefaultConf.Services.Mongo.Source,
		Username:  conf.DefaultConf.Services.Mongo.Username,
		Password:  conf.DefaultConf.Services.Mongo.Password,
		Direct:    false,
		Timeout:   time.Second * 2,
		PoolLimit: 4096,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal("Connection mongo error:", err)
	}
	database := session.DB("netease")
	return database.C("song_comments")
}

func FlushMongoComments() {
	mongo := NewMongoComments()
	mongo.RemoveAll(nil)
}

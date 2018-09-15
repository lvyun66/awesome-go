package music

import (
	"fmt"
	"github.com/lvyun66/awesome-go/netease/conf"
	"gopkg.in/mgo.v2"
	"log"
	"testing"
	"time"
)

func TestGetSongComments(t *testing.T) {
	GetSongComments(536099160)
}

func TestMongo(t *testing.T) {
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
	c := database.C("song_comments")
	query := c.Find(nil)
	comment := &Comment{}
	query.One(comment)
	fmt.Println(comment.User.Nickname)
	fmt.Println(comment.Content)
	fmt.Println(time.Unix(comment.Time, comment.Time%10000).Format("2016-01-02 15:04:05"))
	c.RemoveAll(nil)
	defer session.Close()
}

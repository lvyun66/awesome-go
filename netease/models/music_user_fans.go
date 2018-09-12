package models

type MusicUserFans struct {
	UserId        int    `xorm:"not null pk INT(11)"`
	UserType      int    `xorm:"INT(11)"`
	NikeName      string `xorm:"VARCHAR(64)"`
	Time          int64  `xorm:"BIGINT(32)"`
	Py            string `xorm:"VARCHAR(256)"`
	ExpertTags    string `xorm:"default '' VARCHAR(128)"`
	AuthStatus    int    `xorm:"INT(11)"`
	Followed      int    `xorm:"TINYINT(1)"`
	Experts       string `xorm:"VARCHAR(128)"`
	Followeds     int    `xorm:"INT(11)"`
	VipType       int    `xorm:"INT(11)"`
	Gender        int    `xorm:"INT(11)"`
	AccountStatus int    `xorm:"INT(11)"`
	AvatarUrl     string `xorm:"default '' VARCHAR(256)"`
	RemarkName    string `xorm:"VARCHAR(128)"`
	Follows       int    `xorm:"INT(11)"`
	Mutual        int    `xorm:"TINYINT(1)"`
	Signature     []byte `xorm:"BLOB"`
	EventCount    int    `xorm:"INT(11)"`
	PlaylistCount int    `xorm:"INT(11)"`
}
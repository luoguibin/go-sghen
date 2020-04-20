package models

import "time"

// UserV2 ...
type UserV2 struct {
	ID int64 `gorm:"primary_key;AUTO_INCREMENT:id;"`

	UserAccount string `gorm:"size:50;comment:'用户个性化账号'"`
	Mobile      string `gorm:"size:50;comment:'手机账号'"`
	UserPWD     string `grom:"column:user_pwd;not null;comment:'登陆密码'"`
	Token       string `gorm:"-" json:"token,omitempty"`

	UserName string `gorm:"size:100;comment:'用户昵称'"`
	Avatar   string `gorm:"comment:'用户头像'" json:"iconUrl"`
	Mood     string `gorm:"size:300;comment:'用户心情'"`
	Level    int    `gorm:"comment:'用户角色等级'" json:"-"`

	TimeCreate time.Time `gorm:"comment:'创建时间'" json:"timeCreate"`
	TimeUpdate time.Time `gorm:"comment:'更新时间'" json:"timeUpdate"`
}

// TableName 数据库用户表名
func (u UserV2) TableName() string {
	return "user_v2"
}

// 初始化系统用户列表
func initSystemUserV2() {
	// 复制手机账号的数据
	// sqlStr := "INSERT INTO user_v2(mobile, user_pwd, user_name, avatar, level, time_create, time_update) SELECT id, password, name, icon_url, level, time_create, time_create FROM user WHERE id > 10000000000;"

	// // 复制个性化账号的数据
	// sqlStr := "INSERT INTO user_v2(user_account, user_pwd, user_name, avatar, level, time_create, time_update) SELECT id, password, name, icon_url, level, time_create, time_create FROM user WHERE id < 30000000000;"

	// dbOrmDefault.Model(&UserV2{}).Exec(sqlStr)

	// usersJSON, err := ioutil.ReadFile("data/sys-account.json")
	// if err != nil {
	// 	fmt.Println("read sys-account.json err by v2", err)
	// 	return
	// }

	// re := gjson.ParseBytes(usersJSON)
	// re.ForEach(func(key, value gjson.Result) bool {
	// 	mobile := value.Get("id").String()
	// 	userPWD := value.Get("password").String()
	// 	userName := value.Get("name").String()
	// 	level := value.Get("level").Int()

	// 	CreateUserV2(mobile, userPWD, userName, int(level))
	// 	return true
	// })
}

// CreateUserV2 ...
func CreateUserV2(mobile string, userPWD string, name string, level int) (*UserV2, error) {
	user := &UserV2{
		Mobile:   mobile,
		UserPWD:  userPWD,
		UserName: name,
	}

	err := dbOrmDefault.Model(&UserV2{}).Create(user).Error
	return user, err
}

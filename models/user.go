package models

import (
	"errors"
	"strings"
	"time"
)

// User ...
type User struct {
	ID int64 `gorm:"primary_key;AUTO_INCREMENT:id;" json:"id"`

	UserAccount string `gorm:"index;size:50;comment:'用户个性化账号'" json:"account"`
	Mobile      string `gorm:"index;size:50;comment:'手机账号'" json:"phone"`
	UserPWD     string `grom:"column:user_pwd;not null;comment:'登陆密码'" json:"-"`
	Token       string `gorm:"-" json:"token,omitempty"`

	UserName string `gorm:"size:100;comment:'用户昵称'" json:"username"`
	Avatar   string `gorm:"comment:'用户头像'" json:"avatar"`
	Mood     string `gorm:"size:300;comment:'用户心情'" json:"mood"`
	Level    int    `gorm:"comment:'用户角色等级'" json:"-"`

	TimeCreate time.Time `gorm:"comment:'创建时间'" json:"timeCreate"`
	TimeUpdate time.Time `gorm:"comment:'更新时间'" json:"timeUpdate"`
}

// TableName 数据库用户表名
func (u User) TableName() string {
	return "user"
}

// 初始化系统用户列表
func initSystemUser() {
	// 复制手机账号的数据
	// sqlStr := "INSERT INTO user(mobile, user_pwd, user_name, avatar, level, time_create, time_update) SELECT id, password, name, icon_url, level, time_create, time_create FROM user WHERE id > 10000000000;"

	// // 复制个性化账号的数据
	// sqlStr := "INSERT INTO user(user_account, user_pwd, user_name, avatar, level, time_create, time_update) SELECT id, password, name, icon_url, level, time_create, time_create FROM user WHERE id < 300000000;"

	// // 更改其他表的外键id
	// UPDATE peotry, user SET peotry.user_id = user.id  WHERE peotry.user_id = user.mobile;
	// UPDATE peotry_set, user SET peotry_set.user_id = user.id  WHERE peotry_set.user_id = user.mobile;
	// UPDATE comment, user SET comment.from_id = user.id  WHERE comment.from_id = user.mobile;
	// UPDATE comment, user SET comment.to_id = user.id  WHERE comment.to_id = user.mobile;
	// UPDATE comment, user SET comment.from_id = user.id  WHERE comment.from_id = user.user_account;
	// UPDATE comment, user SET comment.to_id = user.id  WHERE comment.to_id = user.user_account;
	// UPDATE game_spear, user SET game_spear.id = user.id  WHERE game_spear.id = user.mobile;
	// UPDATE game_shield, user SET game_shield.id = user.id  WHERE game_shield.id = user.mobile;
	// UPDATE game_data, user SET game_data.id = user.id  WHERE game_data.id = user.mobile;
	// UPDATE dynamic_api, user SET dynamic_api.user_id = user.id  WHERE dynamic_api.user_id = user.mobile;

	// dbOrmDefault.Model(&User{}).Exec(sqlStr)

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

	// 	CreateUser(mobile, userPWD, userName, int(level))
	// 	return true
	// })
}

// CreateUser ...
func CreateUser(userAccount, mobile, userPWD, userName, avatar, mood string, level int) (*User, error) {
	user := &User{
		UserAccount: userAccount,
		Mobile:      mobile,
		UserPWD:     userPWD,
		UserName:    userName,
		Avatar:      avatar,
		Mood:        mood,
		Level:       level,
		TimeCreate:  time.Now(),
	}
	user.TimeUpdate = user.TimeCreate

	err := dbOrmDefault.Model(&User{}).Create(user).Error
	return user, err
}

// UpdateUser ...
func UpdateUser(id int64, userPWD, userName, avatar, mood string) (*User, error) {
	user := &User{
		ID: id,
	}

	err := dbOrmDefault.Model(&User{}).Find(user).Error
	if err != nil {
		return user, err
	}

	if len(strings.TrimSpace(userPWD)) > 0 {
		user.UserPWD = userPWD
	}
	if len(strings.TrimSpace(userName)) > 0 {
		user.UserName = userName
	}
	if len(strings.TrimSpace(avatar)) > 0 {
		user.Avatar = avatar
	}
	if len(strings.TrimSpace(mood)) > 0 {
		user.Mood = mood
	}
	user.TimeUpdate = time.Now()

	err = dbOrmDefault.Model(&User{}).Save(user).Error
	return user, err
}

// UpdateUserAccount ...
func UpdateUserAccount(id int64, account, mobile string) (*User, error) {
	user := &User{
		ID: id,
	}

	err := dbOrmDefault.Model(&User{}).Find(user).Error
	if err != nil {
		return user, err
	}

	if len(strings.TrimSpace(account)) > 0 {
		user.UserAccount = account
	} else if len(strings.TrimSpace(mobile)) > 0 {
		user.Mobile = mobile
	} else {
		return user, errors.New("参数不完整")
	}

	user.TimeUpdate = time.Now()
	err = dbOrmDefault.Model(&User{}).Save(user).Error

	return user, err
}

// DeleteUser ...
func DeleteUser(id int64) error {
	user := &User{
		ID: id,
	}

	err := dbOrmDefault.Model(&User{}).Delete(user).Error
	return err
}

// QueryUser ...
func QueryUser(account, mobile string) (*User, error) {
	user := &User{
		UserAccount: account,
		Mobile:      mobile,
	}
	err := dbOrmDefault.Model(&User{}).Where(user).Find(user).Error
	return user, err
}

// QueryUsers 移至动态api

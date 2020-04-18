package models

import (
	"fmt"
	"go-sghen/helper"
	"time"
)

// DynamicAPI 自定义脚本制作接口
type DynamicAPI struct {
	ID         int64  `gorm:"primary_key" json:"id,omitempty"`
	SuffixPath string `gorm:"column:suffix_path;type:varchar(100);not null;unique" json:"suffixPath"`

	Name    string `gorm:"column:name;type:varchar(50)" json:"name,omitempty"`
	Comment string `gorm:"column:comment;type:varchar(200)" json:"comment,omitempty"`
	Content string `gorm:"column:content;type:mediumtext" json:"content"`

	Status int `gorm:"column:status" json:"status"`
	Count  int  `gorm:"column:count" json:"count"`

	TimeCreate time.Time `gorm:"column:time_create" json:"timeCreate"`
	TimeUpdate time.Time `gorm:"column:time_update" json:"timeUpdate"`

	UserID int64 `gorm:"column:user_id" json:"-"`
	User   *User `gorm:"foreignkey:user_id;" json:"user"`
}

// TableName ...
func (u DynamicAPI) TableName() string {
	return "dynamic_api"
}

// CreateDynamicAPI 创建一个接口
func CreateDynamicAPI(suffixPath string, name string, comment string, content string, status int, userID int64) (*DynamicAPI, error) {
	id := helper.NewUinqueID()
	timeNow := time.Now()

	dynamicAPI := &DynamicAPI{
		ID:         id,
		SuffixPath: suffixPath,
		Name:       name,
		Comment:    comment,
		Content:    content,
		Status:     status,
		Count:		0,
		UserID:     userID,
		TimeCreate: timeNow,
		TimeUpdate: timeNow,
	}

	err := dbOrmDefault.Model(&DynamicAPI{}).Create(dynamicAPI).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else if status == 1 {
		initDynamicAPIMap()
	}
	return dynamicAPI, nil
}

// UpdateDynamicAPI 更新一个接口
func UpdateDynamicAPI(id int64, suffixPath string, name string, comment string, content string, status int, count int) (*DynamicAPI, error) {
	dynamicAPI := &DynamicAPI{
		ID:         id,
		SuffixPath: suffixPath,
		Name:       name,
		Comment:    comment,
		Content:    content,
		Status:     status,
		TimeUpdate: time.Now(),
	}

	if count > 0 {
		dynamicAPI.Count = count
	}

	err := dbOrmDefault.Model(&DynamicAPI{}).Update(dynamicAPI).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	initDynamicAPIMap()
	return dynamicAPI, nil
}

// QueryDynamicAPI 查询接口列表
func QueryDynamicAPI(id int64, suffixPath string, name string, comment string, status int, userID int64, limit int, page int) ([]*DynamicAPI, int, int, int, int, error) {
	list := make([]*DynamicAPI, 0)
	totalPage := 0
	count := 0
	curPage := page
	pageIsEnd := 0

	if limit <= 0 {
		limit = 10
	}

	db := dbOrmDefault.Model(&DynamicAPI{})
	query := &DynamicAPI{}
	if id > 0 {
		query.ID = id
	}
	if status > 0 {
		query.Status = status
	}
	if userID > 0 {
		query.UserID = userID
	}
	db = db.Where(query)

	if len(suffixPath) > 1 {
		db = db.Where("suffixPath LIKE ?", "%"+suffixPath+"%")
	}
	if len(name) > 1 {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if len(comment) > 1 {
		db = db.Where("comment LIKE ?", "%"+comment+"%")
	}
	db = db.Order("time_create desc")

	db.Count(&count)
	db = db.Preload("User")
	err := db.Limit(limit).Offset(helper.PageOffset(limit, page)).Find(&list).Error

	if err == nil {
		totalPage, pageIsEnd = helper.PageTotal(limit, page, int64(count))
		return list, count, totalPage, curPage, pageIsEnd, nil
	}
	return nil, 0, 0, 0, 0, err
}

// DeleteDynamicAPI 删除接口
func DeleteDynamicAPI(id int64) error {
	dynamicAPI := &DynamicAPI{
		ID: id,
	}

	err := dbOrmDefault.Model(&DynamicAPI{}).Delete(&dynamicAPI).Error

	if err == nil {
		initDynamicAPIMap()
	}
	return err
}

// GetDynamicData 获取数据
func GetDynamicData(sqlStr string) ([]interface{}, error) {
	rows, err := dbOrmDefault.Raw(sqlStr).Rows()
	if err != nil {
		return nil, err
	}

	//读出查询出的列字段名
	cols, _ := rows.Columns()
	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(cols))
	//rows.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(cols))
	//让每一行数据都填充到[][]byte里面
	for i := range values {
		scans[i] = &values[i]
	}

	var list []interface{}
	for rows.Next() { //循环，让游标往下推
		if err := rows.Scan(scans...); err != nil { //rows.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			MConfig.MLogger.Error(err.Error())
			continue
		}

		row := make(map[string]string) //每行数据

		for k, v := range values { //每行数据是放在values里面，现在把它挪到row里
			key := cols[k]
			row[key] = string(v)
		}
		list = append(list, row)
	}

	return list, nil
}

// // PostDynamicData 更改数据
// func PostDynamicData() {

// }

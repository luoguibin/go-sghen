package models

import (
	"time"
)

// Common ...
type Common struct {
}

// TableDesc ...
type TableDesc struct {
	Name       string    `gorm:"column:Name" json:"name"`
	Comment    string    `gorm:"column:Comment" json:"comment"`
	CreateTime time.Time `gorm:"column:Create_time" json:"createTime"`
	UpdateTime time.Time `gorm:"column:Update_time" json:"updateTime"`
}

// FieldDesc ...
type FieldDesc struct {
	Name     string `gorm:"column:COLUMN_NAME" json:"name"`
	DataType string `gorm:"column:DATA_TYPE" json:"type"`
	Comment  string `gorm:"column:COLUMN_COMMENT" json:"comment"`
}

// TableName ...
func (c Common) TableName() string {
	return "common"
}

// GetTables ...
func GetTables() ([]TableDesc, error) {
	rows, err := dbOrmDynamic.Raw("SHOW TABLE STATUS").Rows()
	if err != nil {
		return nil, err
	}

	var list []TableDesc
	for rows.Next() {
		var tableDesc TableDesc
		dbOrmDynamic.ScanRows(rows, &tableDesc)
		list = append(list, tableDesc)
	}

	return list, nil
}

// GetFieldData ...
func GetFieldData(tableName string) ([]FieldDesc, error) {
	sqlStr := "SELECT * FROM information_schema.COLUMNS WHERE TABLE_NAME=? AND TABLE_SCHEMA=?"
	rows, err := dbOrmDynamic.Raw(sqlStr, tableName, MConfig.dBName0).Rows()
	if err != nil {
		return nil, err
	}

	var list []FieldDesc
	for rows.Next() {
		var tableDesc FieldDesc
		dbOrmDynamic.ScanRows(rows, &tableDesc)
		list = append(list, tableDesc)
	}

	return list, nil
}

// GetTableData ...
func GetTableData(tableName string) ([]interface{}, error) {
	rows, err := dbOrmDynamic.Table(tableName).Rows()
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

// GetSQLData ...
func GetSQLData(sqlStr string) ([]interface{}, error) {
	rows, err := dbOrmDynamic.Raw(sqlStr).Rows()
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

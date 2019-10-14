package models

import "go-sghen/helper"

// SmsCode ... todo: 防攻击
type SmsCode struct {
	ID         int64  `gorm:"primary_key" json:"id,omitempty"`
	Code       string `gorm:"code" json:"code"`
	TimeLife   int64  `gorm:"time_life" json:"timeLife"`
	TimeCreate int64  `gorm:"column:time_create" json:"timeCreate"`
}

// SaveSmsCode ...
func SaveSmsCode(id int64, code string, timeLife int64) (*SmsCode, error) {
	smsCode := &SmsCode{
		ID:         id,
		Code:       code,
		TimeLife:   timeLife,
		TimeCreate: helper.GetMillisecond(),
	}

	err := dbOrmDefault.Model(&SmsCode{}).Save(smsCode).Error
	if err != nil {
		return nil, err
	}
	return smsCode, nil
}

// QuerySmsCode ...
func QuerySmsCode(id int64) (*SmsCode, error) {
	smsCode := &SmsCode{
		ID: id,
	}

	err := dbOrmDefault.Model(&SmsCode{}).Find(smsCode).Error
	if err == nil {
		return smsCode, nil
	}
	return nil, err
}

// DeleteSmsCode ...
func DeleteSmsCode(id int64) error {
	smsCode := &SmsCode{
		ID: id,
	}

	err := dbOrmDefault.Model(&SmsCode{}).Delete(smsCode).Error
	return err
}
